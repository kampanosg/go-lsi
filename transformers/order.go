package transformers

import (
	"time"

	"github.com/kampanosg/go-lsi/types"
)

func FromOrderDbRowToDomain(id int, squareId, locationId, state string, totalMoney float64, createdAt int) types.Order {
	return types.Order{
		Id:         id,
		SquareId:   squareId,
		LocationId: locationId,
		State:      state,
		TotalMoney: totalMoney,
		CreatedAt:  time.UnixMilli(int64(createdAt)),
	}
}

func FromOrderProductDbRowToDomain(id int, squareOrderId, squareVarId, qty string) types.OrderProduct {
	return types.OrderProduct{
		Id:            id,
		SquareOrderId: squareOrderId,
		SquareVarId:   squareVarId,
		Quantity:      qty,
	}
}
