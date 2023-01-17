package sync

import (
	"fmt"
	"log"

	"github.com/kampanosg/go-lsi/clients/db"
	"github.com/kampanosg/go-lsi/clients/linnworks"
	"github.com/kampanosg/go-lsi/clients/square"
	"github.com/kampanosg/go-lsi/types"
)

type SyncTool struct {
	LinnworksClient *linnworks.LinnworksClient
	SquareClient    *square.SquareClient
	Db              db.DB
}

type upsertCategory struct {
	category  types.Category
	isDeleted bool
}

func NewSyncTool(lwClient *linnworks.LinnworksClient, sqClient *square.SquareClient, db db.DB) *SyncTool {
	return &SyncTool{
		LinnworksClient: lwClient,
		SquareClient:    sqClient,
		Db:              db,
	}
}

func (s *SyncTool) SyncCategories() {

    oldCategories, _ := s.Db.GetCategories()
	//  newCategories, _ := s.LinnworksClient.GetCategories()

	newCategories := []types.Category{
		{Id: "category-1", Name: "Test Category 1"},
	}

	categoriesUpsertMap := buildUpsertCategoryMap(oldCategories)
	categoriesToBeUpserted := make([]types.Category, 0)
    categoriesSquareIdMapping := make(map[string]types.Category)

	for _, newCategory := range newCategories {

		upsert, ok := categoriesUpsertMap[newCategory.Id]

		if !ok { // new category, need to format square id to specification
			newCategory.SquareId = fmt.Sprintf("#%s", newCategory.Id)
		} else {
			newCategory.SquareId = upsert.category.SquareId
			newCategory.Version = upsert.category.Version
		}

		categoriesUpsertMap[newCategory.Id] = upsertCategory{
			category:  newCategory,
			isDeleted: false,
		}
		categoriesToBeUpserted = append(categoriesToBeUpserted, newCategory)
        categoriesSquareIdMapping[newCategory.SquareId] = newCategory
	}

	categoriesToBeDeleted := getCategoriesToBeDeleted(categoriesUpsertMap)
	if len(categoriesToBeDeleted) > 0 {
		// TODO: dispatch thread to delete categories
	}

	log.Printf("to upsert: %v\n", categoriesToBeUpserted)

    resp, _ := s.SquareClient.UpsertCategories(categoriesToBeUpserted)

    if len(resp.IDMappings) > 0 {
        for _, idMapping := range resp.IDMappings {
            category := categoriesSquareIdMapping[idMapping.ClientObjectID]
            category.SquareId = idMapping.ObjectID
            categoriesSquareIdMapping[category.SquareId] = category 
        }
    }

    categories := make([]types.Category, 0)
    for _, object := range resp.Objects {
        category := categoriesSquareIdMapping[object.ID]
        category.Version = object.Version
        categories = append(categories, category)
    }

    s.Db.ClearCategories()
    s.Db.InsertCategories(categories)
}

func getCategoriesToBeDeleted(categoriesUpsertMap map[string]upsertCategory) []string {
	categoriesToBeDeleted := make([]string, 0)
	for _, v := range categoriesUpsertMap {
		if v.isDeleted {
			categoriesToBeDeleted = append(categoriesToBeDeleted, v.category.SquareId)
		}
	}
	return categoriesToBeDeleted
}

// / Takes a list of Categories, converts them to UpsertCategory type and then returns a LinnworksId -> UpsertCategory mapping
// / Assumes that all categories in the mapping are to be deleted
func buildUpsertCategoryMap(categories []types.Category) map[string]upsertCategory {
	m := map[string]upsertCategory{}
	for _, c := range categories {
		m[c.Id] = upsertCategory{
			category:  c,
			isDeleted: true,
		}
	}
	return m
}
