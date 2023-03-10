package sync

import (
	"fmt"
	"strings"
	"time"

	"github.com/kampanosg/go-lsi/clients/linnworks"
	"github.com/kampanosg/go-lsi/types"
)

const (
	ProductBatchSize = 500
)

var (
	sleepDuration = time.Duration(90 * time.Second)
)

type upsertProduct struct {
	product   types.Product
	isDeleted bool
}

func (s *SyncTool) SyncProducts() error {
	s.logger.Infow("will start syncing products")

	categories, err := s.Db.GetCategories()
	if err != nil {
		s.logger.Errorw("cannot get existing categories", reasonKey, msgDbErr, errKey, err.Error())
		return err
	}

	mappedCatergoriesById := buildMappedCategoriesById(categories)

	oldProducts, err := s.Db.GetProducts()
	if err != nil {
		s.logger.Errorw("cannot get existing products", reasonKey, msgDbErr, errKey, err.Error())
		return err
	}

	lwProduts, err := s.LinnworksClient.GetProducts()
	if err != nil {
		s.logger.Errorw("cannot get new products", reasonKey, msgLwErr, errKey, err.Error())
		return err
	}

	newProducts := fromProductLinnworksResponsesToDomain(lwProduts)
	productsUpsertMap := buildUpsertProductMap(oldProducts)
	productsToUpsert := make([]types.Product, 0)
	productsSquareIdMapping := make(map[string]types.Product, 0)

	s.logger.Infow("will attempt to upsert products", "total", len(newProducts))

	for _, newProduct := range newProducts {

		if isService(newProduct) {
			continue
		}

		upsert, ok := productsUpsertMap[newProduct.LinnworksID]
		if !ok {
			newProduct.SquareID = fmt.Sprintf("#%s", newProduct.LinnworksID)
			newProduct.SquareVarID = fmt.Sprintf("#%s-var", newProduct.LinnworksID)
		} else {

			if isProductUnchanged(newProduct, upsert.product) {
				productsUpsertMap[newProduct.LinnworksID] = upsertProduct{
					product:   newProduct,
					isDeleted: false,
				}
				continue
			}

			newProduct.SquareID = upsert.product.SquareID
			newProduct.SquareVarID = upsert.product.SquareVarID
			newProduct.SquareCategoryID = upsert.product.SquareCategoryID
			newerItemVersion, newerVarVersion, err := s.SquareClient.GetItemVersion(newProduct.SquareID)
			if err != nil {
				s.logger.Debugw("using existing version", "id", newProduct.SquareID, "version", newProduct.Version)
				newProduct.Version = upsert.product.Version
				newProduct.VariationVersion = upsert.product.VariationVersion
			} else {
				s.logger.Debugw("found newer version", "id", newProduct.SquareID, "itemVersion", newerItemVersion, "itemVarVersion", newerVarVersion)
				newProduct.Version = newerItemVersion
				newProduct.VariationVersion = newerVarVersion
			}
		}

		category := mappedCatergoriesById[newProduct.LinnworksCategoryID]
		newProduct.SquareCategoryID = category.SquareID
		newProduct.CategoryID = category.ID

		productsToUpsert = append(productsToUpsert, newProduct)
		productsSquareIdMapping[newProduct.SquareID] = newProduct
		productsUpsertMap[newProduct.LinnworksID] = upsertProduct{
			product:   newProduct,
			isDeleted: false,
		}
	}

	s.logger.Infow("found updated or new products", "total", len(productsToUpsert))
	if len(productsToUpsert) > 0 {

		ptus := make([]types.Product, 0)

		for index, ptu := range productsToUpsert {

			ptus = append(ptus, ptu)
			if len(ptus) < ProductBatchSize && index < len(productsToUpsert)-1 {
				continue
			}

			s.logger.Debugw("sending product batch to square", "length", len(ptus), "items", ptus)

			resp, err := s.SquareClient.UpsertProducts(ptus)
			if err != nil {
				s.logger.Errorw("unable to upsert products", reasonKey, msgSqErr, errKey, err.Error())
				return err
			}

			if len(resp.IDMappings) > 0 {
				s.logger.Debugw("found new product mappings", "total", len(resp.IDMappings))
				for _, idMapping := range resp.IDMappings {
					if !strings.HasSuffix(idMapping.ClientObjectID, "-var") {
						product := productsSquareIdMapping[idMapping.ClientObjectID]
						product.SquareID = idMapping.ObjectID
						for _, varIdMapping := range resp.IDMappings {
							clientObjectId := varIdMapping.ClientObjectID
							clientObjectIdLen := len(clientObjectId)
							productId := clientObjectId[1 : clientObjectIdLen-4]
							if strings.HasSuffix(clientObjectId, "-var") && productId == product.LinnworksID {
								product.SquareVarID = varIdMapping.ObjectID
								break
							}
						}
						productsSquareIdMapping[product.SquareID] = product
					}
				}
			}

			for _, object := range resp.Objects {
				product := productsSquareIdMapping[object.ID]
				product.Version = object.Version
				s.Db.UpsertProduct(product)
			}

			if index >= len(productsToUpsert)-1 {
				break
			}

			ptus = make([]types.Product, 0)
			s.logger.Infow("will sleep to avoid rate limit", "duration", sleepDuration)
			time.Sleep(sleepDuration)
		}
	}

	productsToBeDeleted := getProductsToBeDeleted(productsUpsertMap)
	if len(productsToBeDeleted) > 0 {
		s.logger.Infow("found products to be deleted", "total", len(productsToBeDeleted))

		if err := s.Db.DeleteProductsBySquareIds(productsToBeDeleted); err != nil {
			s.logger.Errorw("unabled to delete products", reasonKey, msgDbErr, errKey, err.Error())
		}

		if err := s.SquareClient.BatchDeleteItems(productsToBeDeleted); err != nil {
			s.logger.Errorw("unable to delete products", reasonKey, msgSqErr, errKey, err.Error())
			return err
		}
	}
	return nil
}

func isService(product types.Product) bool {
	return strings.HasPrefix(product.SKU, "GTR-")
}

func buildMappedCategoriesById(categories []types.Category) map[string]types.Category {
	m := make(map[string]types.Category, 0)
	for _, category := range categories {
		m[category.LinnworksID] = category
	}
	return m
}

// Takes a list of Products, converts them to UpsertProduct type and then returns a LinnworksId -> UpsertProduct mapping
// Assumes that all products in the mapping are to be deleted
func buildUpsertProductMap(products []types.Product) map[string]upsertProduct {
	m := map[string]upsertProduct{}
	for _, p := range products {
		m[p.LinnworksID] = upsertProduct{
			product:   p,
			isDeleted: true,
		}
	}
	return m
}

func getProductsToBeDeleted(productsUpsertMap map[string]upsertProduct) []string {
	productsToBeDeleted := make([]string, 0)
	for _, v := range productsUpsertMap {
		if v.isDeleted {
			productsToBeDeleted = append(productsToBeDeleted, v.product.SquareID)
		}
	}
	return productsToBeDeleted
}

func fromProductLinnworksResponsesToDomain(lwProducts []linnworks.LinnworksProductResponse) (products []types.Product) {
	for _, p := range lwProducts {
		products = append(products, fromProductLinnworksResponseToDomain(p))
	}
	return products
}

func fromProductLinnworksResponseToDomain(lwProduct linnworks.LinnworksProductResponse) types.Product {
	return types.Product{
		LinnworksID:         lwProduct.StockItemID,
		LinnworksCategoryID: lwProduct.CategoryID,
		Title:               lwProduct.ItemTitle,
		Barcode:             lwProduct.BarcodeNumber,
		Price:               lwProduct.RetailPrice,
		SKU:                 lwProduct.SKU,
	}
}

func isProductUnchanged(newProduct, oldProduct types.Product) bool {
	return newProduct.Title == oldProduct.Title &&
		newProduct.Barcode == oldProduct.Barcode &&
		newProduct.SKU == oldProduct.SKU &&
		newProduct.Price == oldProduct.Price &&
		newProduct.LinnworksCategoryID == oldProduct.LinnworksCategoryID
}
