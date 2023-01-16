package domain

type UpsertProduct struct {
	Product   Product
	IsDeleted bool
}
