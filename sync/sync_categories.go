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
    s.logger.Infow("will start syncing categories")

	oldCategories, err := s.Db.GetCategories()
    if err != nil {
        s.logger.Errorw("unable to sync categories", reasonKey, msgDbErr, errKey, err.Error())
        return err
    }

	lwCategories, err := s.LinnworksClient.GetCategories()
    if err != nil {
        s.logger.Errorw("unable to sync categories", reasonKey, msgLwErr, errKey, err.Error())
        return err
    }

	newCategories := transformers.FromCategoryLinnworksResponsesToDomain(lwCategories)
    s.logger.Infow("found categories from linnworks", "total", len(newCategories))

	categoriesUpsertMap := buildUpsertCategoryMap(oldCategories)
	categoriesToBeUpserted := make([]types.Category, 0)
	categoriesSquareIdMapping := make(map[string]types.Category)

	for _, newCategory := range newCategories {

		upsert, ok := categoriesUpsertMap[newCategory.Id]

		if !ok {
			newCategory.SquareId = fmt.Sprintf("#%s", newCategory.Id)
		} else {
			newCategory.SquareId = upsert.category.SquareId
			newCategory.Version = upsert.category.Version
		}
        s.logger.Debugw("assigned square id and version to category", "id", newCategory.SquareId, "version", newCategory.Version)

		categoriesUpsertMap[newCategory.Id] = upsertCategory{
			category:  newCategory,
			isDeleted: false,
		}
		categoriesToBeUpserted = append(categoriesToBeUpserted, newCategory)
		categoriesSquareIdMapping[newCategory.SquareId] = newCategory
	}

    s.logger.Infow("upserting categories to square", "total", len(categoriesToBeUpserted))
	resp, err := s.SquareClient.UpsertCategories(categoriesToBeUpserted)
	if err != nil {
        s.logger.Errorw("unable to upsert categories", reasonKey, msgSqErr, errKey, err.Error())
		return err
	}

	if len(resp.IDMappings) > 0 {
        s.logger.Debugw("found new category mappings", "total", len(resp.IDMappings))
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
        s.logger.Errorw("unable to clear database from categories", reasonKey, msgDbErr, errKey, err.Error())
		return err
	}

	if len(categories) > 0 {
        s.logger.Infow("inserting categories to db", "total", len(categories))
		if err := s.Db.InsertCategories(categories); err != nil {
            s.logger.Errorw("unable to insert categories", reasonKey, msgDbErr, errKey, err.Error())
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
