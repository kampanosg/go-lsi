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

func main() {
	appId := getEnv("APP_ID")
	secret := getEnv("APP_SECRET")
	token := getEnv("APP_TOKEN")
	dbPath := getEnv("DB")

	c := client.NewLinnworksClient(appId, secret, token)
	// newCategories, _ := c.GetCategories()
	newProducts, _ := c.GetProducts()

	log.Printf("%v\n", newProducts)

	sqliteDb := db.NewSqliteDB(dbPath)
	defer sqliteDb.Connection.Close()

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
	// Assume that all categories in the database are to be deleted
	// Start the merge map, with all the values from the db
	// Go over the new categories
	//    For every new category upsert it to the merged map, and set to_delete=false
	//        New categories are appended
	//        Existing categories are updated
	//        Categories to be deleted, have the flag deleted=true
	// Wipe the database and add the entries again

	// oldCats, _ := sqliteDb.GetCategories()
	// mergedCategories := buildCategoryMap(oldCats)

	// for _, newCat := range newCategories {
	// 	mergedCategories[newCat.Id] = upsertCategory{category: newCat, shouldDelete: false}
	// }

	// sqliteDb.ClearCategories()

	// log.Printf("will merge %d categories", len(mergedCategories))
	// categories := []domain.Category{}
	// for _, entry := range mergedCategories {
	// 	log.Printf("%s - %s - should_delete=%v\n", entry.category.Id, entry.category.Name, entry.shouldDelete)
	// 	if !entry.shouldDelete {
	// 		categories = append(categories, entry.category)
	// 	}
	// }

	// sqliteDb.InsertCategories(categories)
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

func getEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
