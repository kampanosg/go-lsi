package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"kev/client"
	"kev/db"
	"kev/types/domain"

	"github.com/joho/godotenv"
)

func main() {
	// appId := getEnv("APP_ID")
	// secret := getEnv("APP_SECRET")
	// token := getEnv("APP_TOKEN")
	squareToken := getEnv("SQ_ACCESS_TOKEN")
	squareHost := getEnv("SQ_HOST")
	dbPath := getEnv("DB")

	// c := client.NewLinnworksClient(appId, secret, token)
	// newCategories, _ := c.GetCategories()
	// newProducts, _ := c.GetProducts()

	sq := client.NewSquareClient(squareToken, squareHost)
	log.Printf("%v", sq)

	sqliteDb := db.NewSqliteDB(dbPath)
	defer sqliteDb.Connection.Close()

	newProducts := []domain.Product{
		{Id: "id-2", CategoryId: "test-cat-7", Title: "Coffee Beans", Barcode: "012345679", SKU: "0x12999", Price: 169.420},
	}
	newCategories := []domain.Category{
		{Id: "test-cat-7", Name: "Test Category"},
	}

	// Strategy:
	// Assume that all entries in the database are to be deleted
	// Start the merge map, with all the values from the db
	// Go over the new entries (from the client)
	//    For every entry upsert it to the merged map, and set to_delete=false
	//        New entries are appended
	//        Existing entries are updated
	//        Entries to be deleted, have the flag deleted=true
	// Wipe the database and add the entries again

	oldCats, _ := sqliteDb.GetCategories()
	mergedCategories := buildCategoryMap(oldCats)

	for _, newCat := range newCategories {
		upsert, ok := mergedCategories[newCat.Id]
		if !ok {
			newCat.SquareId = fmt.Sprintf("#%s", newCat.Id)
		} else {
			newCat.SquareId = upsert.Category.SquareId
			newCat.Version = upsert.Category.Version
		}
		mergedCategories[newCat.Id] = domain.UpsertCategory{
			Category:  newCat,
			IsDeleted: false,
		}
	}

	categoriesToUpsert := []domain.Category{}
	categoriesToDelete := []domain.Category{}

	for _, entry := range mergedCategories {
		if entry.IsDeleted {
			categoriesToDelete = append(categoriesToDelete, entry.Category)
		} else {
			categoriesToUpsert = append(categoriesToUpsert, entry.Category)
		}
	}

	resp := sq.UpsertCategories(categoriesToUpsert)

	categories := []domain.Category{}
	for _, upserted := range categoriesToUpsert {

		var version int64
		var squareId string

		if upserted.SquareId[0] == '#' {
			for _, mapping := range resp.IDMappings {
				if mapping.ClientObjectID == upserted.SquareId {
					squareId = mapping.ObjectID
					break
				}
			}
		} else {
			squareId = upserted.SquareId
		}

		for _, object := range resp.Objects {
			if object.ID == squareId {
				version = object.Version
				break
			}
		}

		upserted.Version = version
		upserted.SquareId = squareId
		categories = append(categories, upserted)
	}

	if len(categoriesToDelete) > 0 {
		sq.DeleteCategories(categoriesToDelete)
	}

	sqliteDb.ClearCategories()
	sqliteDb.InsertCategories(categories)

	oldProducts, _ := sqliteDb.GetProducts()
	mergedProducts := buildProductMap(oldProducts)

	for _, newProd := range newProducts {
		upsert, ok := mergedProducts[newProd.Id]
		if !ok {
			newProd.SquareId = fmt.Sprintf("#%s", newProd.Id)
			newProd.SquareVarId = fmt.Sprintf("#%s-var", newProd.Id)
		} else {
			newProd.SquareId = upsert.Product.SquareId
			newProd.SquareVarId = upsert.Product.SquareVarId
			newProd.Version = upsert.Product.Version
		}

		for _, category := range categories {
			if category.Id == newProd.CategoryId {
				newProd.SquareCategoryId = category.SquareId
				break
			}
		}

		mergedProducts[newProd.Id] = domain.UpsertProduct{
			Product:   newProd,
			IsDeleted: false,
		}
	}

	productsToUpsert := []domain.Product{}
	productsToDelete := []domain.Product{}

	for _, entry := range mergedProducts {
		if entry.IsDeleted {
			productsToDelete = append(productsToDelete, entry.Product)
		} else {
			productsToUpsert = append(productsToUpsert, entry.Product)
		}
	}

	prodResp := sq.UpsertProducts(productsToUpsert)

	products := []domain.Product{}
	for _, upserted := range productsToUpsert {

		var version int64
		var squareId, squareVarId string

		// new product
		if upserted.SquareId[0] == '#' {
			for _, mapping := range prodResp.IDMappings {
				if strings.HasSuffix(mapping.ClientObjectID, "-var") {
					// if the id ends with -var, it means its a variation
					id := mapping.ClientObjectID[1 : len(mapping.ClientObjectID)-4] // need to sanitize the id from #product-id-var to product-id
					if id == upserted.Id {
						squareVarId = mapping.ObjectID
					}
				} else {
					id := mapping.ClientObjectID[1:]
					if id == upserted.Id {
						squareId = mapping.ObjectID
					}
				}
			}

		} else {
			squareId = upserted.SquareId
			squareVarId = upserted.SquareVarId
		}

		for _, object := range prodResp.Objects {
            log.Printf("objkect=%v, squareId=%s\n", object, squareId)
			if object.ID == squareId {
				version = object.Version
				break
			}
		}

		upserted.SquareId = squareId
		upserted.SquareVarId = squareVarId
		upserted.Version = version

		products = append(products, upserted)
	}

    if len(productsToDelete) > 0 {
		sq.DeleteProducts(productsToDelete)
	}

	sqliteDb.ClearProducts()
	sqliteDb.InsertProducts(products)
}

func buildCategoryMap(categories []domain.Category) map[string]domain.UpsertCategory {
	m := map[string]domain.UpsertCategory{}
	for _, c := range categories {
		m[c.Id] = domain.UpsertCategory{
			Category:  c,
			IsDeleted: true,
		}
	}
	return m
}

func buildProductMap(products []domain.Product) map[string]domain.UpsertProduct {
	m := map[string]domain.UpsertProduct{}
	for _, p := range products {
		m[p.Id] = domain.UpsertProduct{
			Product:   p,
			IsDeleted: true,
		}
	}
	return m
}

func getEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
