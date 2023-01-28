package sync

import (
	"time"

	"github.com/kampanosg/go-lsi/transformers"
	"github.com/kampanosg/go-lsi/types"
)

func (s *SyncTool) SyncOrders(start time.Time, end time.Time) error {
	existingOrders, err := s.Db.GetOrders()
	if err != nil {
		return err
	}

	newOrders, err := s.SquareClient.SearchOrders(start, end)
	if err != nil {
		return err
	}
	existingOrdersMap := buildSquareIdToOrderMap(existingOrders)

	ordersToUpsert := make([]types.Order, 0)
	for _, newOrder := range newOrders {
		_, ok := existingOrdersMap[newOrder.ID]
		if !ok {

			orderProducts := make([]types.OrderProduct, len(newOrder.LineItems))
			for index, item := range newOrder.LineItems {
				product, err := s.Db.GetProductByVarId(item.CatalogObjectID)
				if err != nil {
					continue
				}

				orderProduct := transformers.FromSquareLineItemToDomain(item, product)
				orderProduct.SquareOrderId = newOrder.ID
				orderProducts[index] = orderProduct
			}

			order := transformers.FromSquareOrderToDomain(newOrder)
			order.Products = orderProducts
			ordersToUpsert = append(ordersToUpsert, order)
		}
	}

	if _, err := s.LinnworksClient.CreateOrders(ordersToUpsert); err != nil {
		return err
	}

	if err := s.Db.InsertOrders(ordersToUpsert); err != nil {
		return err
	}

	return nil
}

func buildSquareIdToOrderMap(orders []types.Order) map[string]types.Order {
	m := make(map[string]types.Order, 0)
	for _, order := range orders {
		m[order.SquareId] = order
	}
	return m
}
