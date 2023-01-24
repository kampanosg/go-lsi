package sync

import (
	"log"
	"time"

	"github.com/kampanosg/go-lsi/transformers"
	"github.com/kampanosg/go-lsi/types"
)

func (s *SyncTool) SyncOrders(start time.Time, end time.Time) {
	existingOrders, _ := s.Db.GetOrders()
	newOrders, _ := s.SquareClient.SearchOrders(start, end)
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

    for _, order := range ordersToUpsert {
        log.Printf("%v\n", order)
    }

}

func buildSquareIdToOrderMap(orders []types.Order) map[string]types.Order {
	m := make(map[string]types.Order, 0)
	for _, order := range orders {
		m[order.SquareId] = order
	}
	return m
}
