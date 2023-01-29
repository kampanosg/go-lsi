package sync

import (
	"fmt"

	"github.com/kampanosg/go-lsi/clients/linnworks"
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

	newCategories := fromCategoryLinnworksResponsesToDomain(lwCategories)
	s.logger.Infow("found categories from linnworks", "total", len(newCategories))

	categoriesUpsertMap := buildUpsertCategoryMap(oldCategories)
	categoriesToBeUpserted := make([]types.Category, 0)
	categoriesSquareIdMapping := make(map[string]types.Category)

	for _, newCategory := range newCategories {

		upsert, ok := categoriesUpsertMap[newCategory.LinnworksID]

		if !ok {
			newCategory.SquareID = fmt.Sprintf("#%s", newCategory.LinnworksID)
		} else {
			newCategory.SquareID = upsert.category.SquareID
			newCategory.Version = upsert.category.Version
		}
		s.logger.Debugw("assigned square id and version to category", "id", newCategory.SquareID, "version", newCategory.Version)

		categoriesUpsertMap[newCategory.LinnworksID] = upsertCategory{
			category:  newCategory,
			isDeleted: false,
		}
		categoriesToBeUpserted = append(categoriesToBeUpserted, newCategory)
		categoriesSquareIdMapping[newCategory.SquareID] = newCategory
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
			category.SquareID = idMapping.ObjectID
			categoriesSquareIdMapping[category.SquareID] = category
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
		s.logger.Infow("found categories to be deleted", "total", len(categoriesToBeDeleted))
		if err := s.SquareClient.BatchDeleteItems(categoriesToBeDeleted); err != nil {
			s.logger.Errorw("unable to delete categories", reasonKey, msgSqErr, errKey, err.Error())
			return err
		}
	}
	return nil
}

func getCategoriesToBeDeleted(categoriesUpsertMap map[string]upsertCategory) []string {
	categoriesToBeDeleted := make([]string, 0)
	for _, v := range categoriesUpsertMap {
		if v.isDeleted {
			categoriesToBeDeleted = append(categoriesToBeDeleted, v.category.SquareID)
		}
	}
	return categoriesToBeDeleted
}

// / Takes a list of Categories, converts them to UpsertCategory type and then returns a LinnworksId -> UpsertCategory mapping
// / Assumes that all categories in the mapping are to be deleted
func buildUpsertCategoryMap(categories []types.Category) map[string]upsertCategory {
	m := map[string]upsertCategory{}
	for _, c := range categories {
		m[c.LinnworksID] = upsertCategory{
			category:  c,
			isDeleted: true,
		}
	}
	return m
}

func fromCategoryLinnworksResponsesToDomain(lwCategories []linnworks.LinnworksCategoryResponse) (categories []types.Category) {
	for _, lwCategory := range lwCategories {
		categories = append(categories, fromCategoryLinnworksResponseToDomain(lwCategory))
	}
	return categories
}

func fromCategoryLinnworksResponseToDomain(lwCategory linnworks.LinnworksCategoryResponse) types.Category {
	return types.Category{
		LinnworksID: lwCategory.Id,
		Name:        lwCategory.Name,
	}
}
