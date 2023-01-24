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
		CreatedAt:  time.Unix(int64(createdAt), 0),
	}
}
