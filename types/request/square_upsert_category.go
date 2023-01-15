package request

type SquareBatchUpsertCategoryRequest struct {
	IdempotencyKey string        `json:"idempotency_key"`
	Batches        []SquareBatch `json:"batches"`
}

type SquareBatch struct {
	Objects []SquareUpsertCategoryRequest `json:"objects"`
}

type SquareUpsertCategoryRequest struct {
	Type         string       `json:"type"`
	Id           string       `json:"id"`
	IsDeleted    bool         `json:"is_deleted"`
	CategoryData CategoryData `json:"category_data"`
}

type CategoryData struct {
	Name string `json:"name"`
}
