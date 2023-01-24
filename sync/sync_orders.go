package sync

import (
	"log"
	"time"

	"github.com/kampanosg/go-lsi/clients/square"
	"github.com/kampanosg/go-lsi/types"
)

func (s *SyncTool) SyncOrders(start time.Time, end time.Time) {
    existingOrders, _ := s.Db.GetOrders()
    newOrders, _ := s.SquareClient.SearchOrders(start, end)
    existingOrdersMap := buildSquareIdToOrderMap(existingOrders)

    ordersToPost := make([]square.SquareOrder, 0)
    for _, newOrder := range newOrders {
        _, ok := existingOrdersMap[newOrder.ID]
        if !ok {
            ordersToPost = append(ordersToPost, newOrder)
        }
    }

    log.Printf("orders to post: %v\n", ordersToPost)
}

func buildSquareIdToOrderMap(orders []types.Order) map[string]types.Order {
    m := make(map[string]types.Order, 0)
    for _, order := range orders {
        m[order.SquareId] = order
    }
    return m
}
