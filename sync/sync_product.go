package sync

import (
	"fmt"
	"strings"

	"github.com/kampanosg/go-lsi/clients/linnworks"
	"github.com/kampanosg/go-lsi/types"
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
		upsert, ok := productsUpsertMap[newProduct.LinnworksID]
		if !ok {
			newProduct.SquareID = fmt.Sprintf("#%s", newProduct.LinnworksID)
			newProduct.SquareVarID = fmt.Sprintf("#%s-var", newProduct.LinnworksID)
		} else {
			newProduct.SquareID = upsert.product.SquareID
			newProduct.SquareVarID = upsert.product.SquareVarID
			newProduct.Version = upsert.product.Version
		}
		s.logger.Debugw("assigned ids and version to product",
			"squareId", newProduct.SquareID,
			"var id", newProduct.SquareVarID,
			"version", newProduct.Version,
		)

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

	resp, err := s.SquareClient.UpsertProducts(productsToUpsert)
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

	products := make([]types.Product, 0)
	for _, object := range resp.Objects {
		product := productsSquareIdMapping[object.ID]
		product.Version = object.Version
		products = append(products, product)
	}

	if err := s.Db.ClearProducts(); err != nil {
		s.logger.Errorw("unable to delete products", reasonKey, msgDbErr, errKey, err.Error())
		return err
	}

	if len(products) > 0 {
		for _, product := range products {
			s.Db.InsertProduct(product)
		}
	}

	productsToBeDeleted := getProductsToBeDeleted(productsUpsertMap)
	if len(productsToBeDeleted) > 0 {
		s.logger.Infow("found products to be deleted", "total", len(productsToBeDeleted))
		if err := s.SquareClient.BatchDeleteItems(productsToBeDeleted); err != nil {
			s.logger.Errorw("unable to delete products", reasonKey, msgSqErr, errKey, err.Error())
			return err
		}
	}
	return nil
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
