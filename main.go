package main

import (
	"fmt"
	"log"
	"os"

	"kev/client"
	"kev/db"
	"kev/types/domain"

	"github.com/joho/godotenv"
)

type upsertProduct struct {
	product      domain.Product
	shouldDelete bool
}

func main() {
	appId := getEnv("APP_ID")
	secret := getEnv("APP_SECRET")
	token := getEnv("APP_TOKEN")
	squareToken := getEnv("SQ_ACCESS_TOKEN")
	dbPath := getEnv("DB")

	c := client.NewLinnworksClient(appId, secret, token)
	newCategories, _ := c.GetCategories()
	// newProducts, _ := c.GetProducts()

	sq := client.NewSquareClient(squareToken)

	sqliteDb := db.NewSqliteDB(dbPath)
	defer sqliteDb.Connection.Close()

	// newProducts := []domain.Product{
	// 	{Id: "id-2", CategoryId: "id-1", Title: "Test product 2", Barcode: "012345679", Price: 169.420},
	// }
	// newCategories := []domain.Category{
        // { Id: "test-cat-7", Name: "Test Category 7" },
        // { Id: "test-cat-8", Name: "Test Category 8" },
	// }

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

	// log.Printf("%v", mergedCategories)

	resp := sq.UpsertCategories(categoriesToUpsert)

	categories := []domain.Category{}
	for _, mapping := range resp.IDMappings {
		clientId := mapping.ClientObjectID[1:]
		entry := mergedCategories[clientId]
		entry.Category.SquareId = mapping.ObjectID
		categories = append(categories, entry.Category)
	}

	if len(categoriesToDelete) > 0 {
		sq.DeleteCategories(categoriesToDelete)
	}

	sqliteDb.ClearCategories()
	sqliteDb.InsertCategories(categories)

	// oldProducts, _ := sqliteDb.GetProducts()
	// mergedProducts := buildProductMap(oldProducts)

	// for _, newProduct := range newProducts {
	// 	mergedProducts[newProduct.Id] = upsertProduct{product: newProduct, shouldDelete: false}
	// }

	// sqliteDb.ClearProducts()

	// products := []domain.Product{}
	// for _, entry := range mergedProducts {
	// 	if !entry.shouldDelete {
	// 		log.Printf("%s - %s - should_delete=%v\n", entry.product.Id, entry.product.Title, entry.shouldDelete)
	// 		products = append(products, entry.product)
	// 	}
	// }

	// sqliteDb.InsertProducts(products)
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

func buildProductMap(products []domain.Product) map[string]upsertProduct {
	m := map[string]upsertProduct{}
	for _, p := range products {
		m[p.Id] = upsertProduct{
			product:      p,
			shouldDelete: true,
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
