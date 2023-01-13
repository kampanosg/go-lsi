package transform

import (
	"kev/types/domain"
	"kev/types/response"
)

func FromCategoryRespToDomain(resp response.CategoryResponse) domain.Category {
	return domain.Category{
		Id:   resp.Id,
		Name: resp.Name,
	}
}

func FromArrCategoryRespToDomain(resp []response.CategoryResponse) (categories []domain.Category) {
    for _, c := range resp {
        categories = append(categories, FromCategoryRespToDomain(c))
    }
    return categories
}
