package sync

import (
	"fmt"

	"github.com/kampanosg/go-lsi/transformers"
	"github.com/kampanosg/go-lsi/types"
)

type upsertCategory struct {
	category  types.Category
	isDeleted bool
}

func (s *SyncTool) SyncCategories() error {

	oldCategories, _ := s.Db.GetCategories()
	lwCategories, _ := s.LinnworksClient.GetCategories()
	newCategories := transformers.FromCategoryLinnworksResponsesToDomain(lwCategories)

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

	resp, err := s.SquareClient.UpsertCategories(categoriesToBeUpserted)
	if err != nil {
		return err
	}

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

	if err := s.Db.ClearCategories(); err != nil {
		return err
	}

	if len(categories) > 0 {
		if err := s.Db.InsertCategories(categories); err != nil {
			return err
		}
	}

	categoriesToBeDeleted := getCategoriesToBeDeleted(categoriesUpsertMap)
	if len(categoriesToBeDeleted) > 0 {
		if err := s.SquareClient.BatchDeleteItems(categoriesToBeDeleted); err != nil {
			return err
		}
	}
	return nil
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
