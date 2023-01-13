package response

type CategoryResponse struct {
	Id                string `json:"CategoryId"`
	Name              string `json:"CategoryName"`
	StructureId       int    `json:"StructureCategoryId"`
	ProductCategoryId int    `json:"ProductCategoryId"`
}
