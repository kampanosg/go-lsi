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

	existingOrders, err := s.Db.GetOrders()
	if err != nil {
		s.logger.Errorw("unable to retrieve existing orders", reasonKey, msgDbErr, errKey, err.Error())
		return err
	}

	newOrders, err := s.SquareClient.SearchOrders(start, end)
	if err != nil {
		s.logger.Errorw("unable to retrieve new orders", reasonKey, msgSqErr, errKey, err.Error())
		return err
	}

	s.logger.Infow("found orders", "existing", len(existingOrders), "new", len(newOrders))

	if len(newOrders) > 0 {
		existingOrdersMap := buildSquareIdToOrderMap(existingOrders)
		ordersToUpsert := make([]types.Order, 0)

		for _, newOrder := range newOrders {
			_, ok := existingOrdersMap[newOrder.ID]
			if !ok {
				orderProducts := make([]types.OrderProduct, len(newOrder.LineItems))
				for index, item := range newOrder.LineItems {
					product, err := s.Db.GetProductByVarId(item.CatalogObjectID)
					if err != nil {
						s.logger.Errorw("unable to retrieve product from db", "variation", item.CatalogObjectID, errKey, err.Error())
						return err
					}
					orderProduct := fromSquareLineItemToDomain(item, product)
					orderProduct.SquareOrderId = newOrder.ID
					orderProducts[index] = orderProduct
				}

				order := fromSquareOrderToDomain(newOrder)
				order.Products = orderProducts
				ordersToUpsert = append(ordersToUpsert, order)
			}
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

func buildSquareIdToOrderMap(orders []types.Order) map[string]types.Order {
	m := make(map[string]types.Order, 0)
	for _, order := range orders {
		m[order.SquareID] = order
	}
	return m
}

func fromSquareOrderToDomain(order square.SquareOrder) types.Order {
	return types.Order{
		SquareID:   order.ID,
		LocationID: order.LocationID,
		CreatedAt:  order.CreatedAt,
		State:      order.State,
		Version:    order.Version,
		TotalMoney: float64(order.TotalMoney.Amount / Pence),
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
