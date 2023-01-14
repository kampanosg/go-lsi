package main

import (
	"log"
	"os"

	"kev/client"
	"kev/db"
	"kev/types/domain"

	"github.com/joho/godotenv"
)

type upsertCategory struct {
	category     domain.Category
	shouldDelete bool
}

type upsertProduct struct {
	product      domain.Product
	shouldDelete bool
}

func main() {
	appId := getEnv("APP_ID")
	secret := getEnv("APP_SECRET")
	token := getEnv("APP_TOKEN")
	dbPath := getEnv("DB")

	c := client.NewLinnworksClient(appId, secret, token)
	newCategories, _ := c.GetCategories()
	newProducts, _ := c.GetProducts()

	sqliteDb := db.NewSqliteDB(dbPath)
	defer sqliteDb.Connection.Close()

	// newProducts := []domain.Product{
	// 	{Id: "id-2", CategoryId: "id-1", Title: "Test product 2", Barcode: "012345679", Price: 169.420},
	// 	{Id: "id-3", CategoryId: "id-1", Title: "Test product 3", Barcode: "012345677", Price: 269.420},
	// 	{Id: "id-4", CategoryId: "id-4", Title: "Test product 4", Barcode: "012345676", Price: 369.420},
	// }

	// newCategories := []domain.Category{
	// 	{Id: "id-1", Name: "Category 1"},
	// 	{Id: "id-10", Name: "Category 10"},
	// 	{Id: "id-4", Name: "Category 5"},
	// 	{Id: "id-4", Name: "Category 5"},
	// 	{Id: "id-7", Name: "Category 7"},
	// 	{Id: "id-8", Name: "Category 8"},
	// 	{Id: "id-9", Name: "Category 9"},
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
		mergedCategories[newCat.Id] = upsertCategory{category: newCat, shouldDelete: false}
	}

	sqliteDb.ClearCategories()

	// log.Printf("will merge %d categories", len(mergedCategories))
	categories := []domain.Category{}
	for _, entry := range mergedCategories {
		// log.Printf("%s - %s - should_delete=%v\n", entry.category.Id, entry.category.Name, entry.shouldDelete)
		if !entry.shouldDelete {
			categories = append(categories, entry.category)
		}
	}

	sqliteDb.InsertCategories(categories)

	oldProducts, _ := sqliteDb.GetProducts()
	mergedProducts := buildProductMap(oldProducts)

	for _, newProduct := range newProducts {
		mergedProducts[newProduct.Id] = upsertProduct{product: newProduct, shouldDelete: false}
	}

	sqliteDb.ClearProducts()

	products := []domain.Product{}
	for _, entry := range mergedProducts {
		if !entry.shouldDelete {
			log.Printf("%s - %s - should_delete=%v\n", entry.product.Id, entry.product.Title, entry.shouldDelete)
			products = append(products, entry.product)
		}
	}

	sqliteDb.InsertProducts(products)
}

func buildCategoryMap(categories []domain.Category) map[string]upsertCategory {
	m := map[string]upsertCategory{}
	for _, c := range categories {
		m[c.Id] = upsertCategory{
			category:     c,
			shouldDelete: true,
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
