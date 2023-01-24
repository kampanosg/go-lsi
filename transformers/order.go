package transformers

import (
	"time"

	"github.com/kampanosg/go-lsi/types"
)

func FromOrderDbRowToDomain(id int, squareId, locationId, state string, totalMoney float64, createdAt int64) types.Order {
	return types.Order{
		Id:         id,
		SquareId:   squareId,
		LocationId: locationId,
		State:      state,
		TotalMoney: totalMoney,
		CreatedAt:  time.Unix(createdAt, 0),
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
