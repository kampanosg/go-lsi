package sync

import (
	"fmt"
	"log"
	"strings"

	"github.com/kampanosg/go-lsi/transformers"
	"github.com/kampanosg/go-lsi/types"
)

type upsertProduct struct {
	product   types.Product
	isDeleted bool
}

func (s *SyncTool) SyncProducts() {

	categories, _ := s.Db.GetCategories()
	mappedCatergoriesById := buildMappedCategoriesById(categories)

	oldProducts, _ := s.Db.GetProducts()
	lwProduts, _ := s.LinnworksClient.GetProducts()
	newProducts := transformers.FromProductLinnworksResponsesToDomain(lwProduts)

	// newProducts := []types.Product{
	// 	{Id: "product-1", Title: "Very Good Coffee Beans", CategoryId: "category-1", Price: 69.42, Barcode: "999999999999", SKU: "SKU0420"},
	// }

	log.Printf("will process %d new products\n", len(newProducts))

	productsUpsertMap := buildUpsertProductMap(oldProducts)
	productsToUpsert := make([]types.Product, 0)
	productsSquareIdMapping := make(map[string]types.Product, 0)

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

		category := mappedCatergoriesById[newProduct.CategoryId]
		newProduct.SquareCategoryId = category.SquareId

		productsToUpsert = append(productsToUpsert, newProduct)
		productsSquareIdMapping[newProduct.SquareId] = newProduct
		productsUpsertMap[newProduct.Id] = upsertProduct{
			product:   newProduct,
			isDeleted: false,
		}
	}

	resp, _ := s.SquareClient.UpsertProducts(productsToUpsert)

	if len(resp.IDMappings) > 0 {
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
		log.Printf("adding product: %v\n", product)
		products = append(products, product)
	}

	s.Db.ClearProducts()
	if len(products) > 0 {
		s.Db.InsertProducts(products)
	}

	productsToBeDeleted := getProductsToBeDeleted(productsUpsertMap)
	if len(productsToBeDeleted) > 0 {
		s.SquareClient.BatchDeleteItems(productsToBeDeleted)
	}
}

func buildMappedCategoriesById(categories []types.Category) map[string]types.Category {
	m := make(map[string]types.Category, 0)
	for _, category := range categories {
		m[category.Id] = category
	}
	return m
}

// / Takes a list of Products, converts them to UpsertProduct type and then returns a LinnworksId -> UpsertProduct mapping
// / Assumes that all products in the mapping are to be deleted
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
