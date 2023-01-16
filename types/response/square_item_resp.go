package response

import "time"

type SquareUpsertItemResponse struct {
	Objects []struct {
		Type                  string    `json:"type"`
		ID                    string    `json:"id"`
		UpdatedAt             time.Time `json:"updated_at"`
		CreatedAt             time.Time `json:"created_at"`
		Version               int64     `json:"version"`
		IsDeleted             bool      `json:"is_deleted"`
		PresentAtAllLocations bool      `json:"present_at_all_locations"`
		CategoryData          struct {
			Name       string `json:"name"`
			IsTopLevel bool   `json:"is_top_level"`
		} `json:"category_data"`
	} `json:"objects"`
	IDMappings []SquareItemIdMapping `json:"id_mappings"`
}

type SquareItemIdMapping struct {
	ClientObjectID string `json:"client_object_id"`
	ObjectID       string `json:"object_id"`
}
