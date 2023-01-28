package sync

import (
	"fmt"
	"strings"

	"github.com/kampanosg/go-lsi/transformers"
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

	newProducts := transformers.FromProductLinnworksResponsesToDomain(lwProduts)

	productsUpsertMap := buildUpsertProductMap(oldProducts)
	productsToUpsert := make([]types.Product, 0)
	productsSquareIdMapping := make(map[string]types.Product, 0)

	s.logger.Infow("will attempt to upsert products", "total", len(newProducts))

	for _, newProduct := range newProducts {
		upsert, ok := productsUpsertMap[newProduct.Id]
		if !ok {
			newProduct.SquareId = fmt.Sprintf("#%s", newProduct.Id)
			newProduct.SquareVarId = fmt.Sprintf("#%s-var", newProduct.Id)
		} else {
			newProduct.SquareId = upsert.product.SquareId
			newProduct.SquareVarId = upsert.product.SquareVarId
			newProduct.Version = upsert.product.Version
		}
		s.logger.Debugw("assigned ids and version to product",
			"squareId", newProduct.SquareId,
			"var id", newProduct.SquareVarId,
			"version", newProduct.Version,
		)

		category := mappedCatergoriesById[newProduct.CategoryId]
		newProduct.SquareCategoryId = category.SquareId

		productsToUpsert = append(productsToUpsert, newProduct)
		productsSquareIdMapping[newProduct.SquareId] = newProduct
		productsUpsertMap[newProduct.Id] = upsertProduct{
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
				product.SquareId = idMapping.ObjectID
				for _, varIdMapping := range resp.IDMappings {
					clientObjectId := varIdMapping.ClientObjectID
					clientObjectIdLen := len(clientObjectId)
					productId := clientObjectId[1 : clientObjectIdLen-4]
					if strings.HasSuffix(clientObjectId, "-var") && productId == product.Id {
						product.SquareVarId = varIdMapping.ObjectID
						break
					}
				}
				productsSquareIdMapping[product.SquareId] = product
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
		s.Db.InsertProducts(products)
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
		m[category.Id] = category
	}
	return m
}

// Takes a list of Products, converts them to UpsertProduct type and then returns a LinnworksId -> UpsertProduct mapping
// Assumes that all products in the mapping are to be deleted
func buildUpsertProductMap(products []types.Product) map[string]upsertProduct {
	m := map[string]upsertProduct{}
	for _, p := range products {
		m[p.Id] = upsertProduct{
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
			productsToBeDeleted = append(productsToBeDeleted, v.product.SquareId)
		}
	}
	return productsToBeDeleted
}
