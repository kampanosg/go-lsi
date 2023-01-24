package transformers

import (
	"time"

	"github.com/kampanosg/go-lsi/clients/square"
	"github.com/kampanosg/go-lsi/types"
)

const (
	Pence = 100.0
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

func FromSquareOrderToDomain(order square.SquareOrder) types.Order {
	return types.Order{
		SquareId:   order.ID,
		LocationId: order.LocationID,
		CreatedAt:  order.CreatedAt,
		State:      order.State,
		Version:    order.Version,
		TotalMoney: float64(order.TotalMoney.Amount / Pence),
	}
}

func FromSquareLineItemToDomain(item square.SquareLineItem, product types.Product) types.OrderProduct {
	return types.OrderProduct{
		SquareVarId:  item.CatalogObjectID,
		Quantity:     item.Quantity,
		ItemNumber:   product.Barcode,
		SKU:          product.SKU,
		Title:        product.Title,
		PricePerUnit: product.Price,
	}
}
