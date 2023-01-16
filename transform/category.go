package transform

import (
	"github.com/kampanosg/go-lsi/types/domain"
	"github.com/kampanosg/go-lsi/types/response"
)

func FromCategoryRespToDomain(resp response.CategoryResponse) domain.Category {
	return domain.Category{
		Id:   resp.Id,
		Name: resp.Name,
	}
}

func FromCategoriesRespToDomain(resp []response.CategoryResponse) (categories []domain.Category) {
	for _, c := range resp {
		categories = append(categories, FromCategoryRespToDomain(c))
	}
	return categories
}
