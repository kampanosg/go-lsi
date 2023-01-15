package request

type BatchDeleteCategoriesRequest struct {
	ObjectIds []string `json:"object_ids"`
}
