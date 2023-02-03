package sync

import (
	"time"

	"github.com/kampanosg/go-lsi/clients/square"
	"github.com/kampanosg/go-lsi/types"
)

const (
	Pence = 100.0
)

func (s *SyncTool) SyncOrders(start time.Time, end time.Time) error {
	s.logger.Infow("will start syncing orders")

	newOrders, err := s.SquareClient.SearchOrders(start, end)
	if err != nil {
		s.logger.Errorw("unable to retrieve new orders", reasonKey, msgSqErr, errKey, err.Error())
		return err
	}

	s.logger.Infow("found orders", "new", len(newOrders))

	if len(newOrders) > 0 {
		ordersToUpsert := make([]types.Order, 0)

		for _, newOrder := range newOrders {
			if _, err := s.Db.GetOrderBySquareId(newOrder.ID); err == nil {
				continue
			}

			orderProducts := make([]types.OrderProduct, len(newOrder.LineItems))
			for index, item := range newOrder.LineItems {
				var product types.Product
				product, err = s.Db.GetProductByVarId(item.CatalogObjectID)
				if err != nil {
					s.logger.Errorw("unable to retrieve product from db", "variation", item.CatalogObjectID, errKey, err.Error())
					s.logger.Infow("will attempt to retrieve product by name", "name", item.Name)
					product, err = s.Db.GetProductByTitle(item.Name)
					if err != nil {
						s.logger.Errorw("unable to retrieve product from db", "name", item.Name, errKey, err.Error())
						return err
					}
				}
				orderProduct := fromSquareLineItemToDomain(item, product)
				orderProduct.SquareOrderId = newOrder.ID
				orderProducts[index] = orderProduct
			}

			order := fromSquareOrderToDomain(newOrder)
			order.Products = orderProducts
			ordersToUpsert = append(ordersToUpsert, order)
		}

		if _, err := s.LinnworksClient.CreateOrders(ordersToUpsert); err != nil {
			s.logger.Errorw("unable to create orders", reasonKey, msgLwErr, errKey, err.Error())
			return err
		}

		if err := s.Db.InsertOrders(ordersToUpsert); err != nil {
			s.logger.Errorw("unable to insert orders", reasonKey, msgDbErr, errKey, err.Error())
			return err
		}
	}

	return nil
}

func fromSquareOrderToDomain(order square.SquareOrder) types.Order {
	return types.Order{
		SquareID:   order.ID,
		LocationID: order.LocationID,
		CreatedAt:  order.CreatedAt,
		State:      order.State,
		Version:    order.Version,
		TotalMoney: order.TotalMoney.Amount / Pence,
	}
}

func fromSquareLineItemToDomain(item square.SquareLineItem, product types.Product) types.OrderProduct {
	return types.OrderProduct{
		SquareVarId:  item.CatalogObjectID,
		Quantity:     item.Quantity,
		ItemNumber:   product.Barcode,
		SKU:          product.SKU,
		Title:        product.Title,
		PricePerUnit: product.Price,
	}
}
