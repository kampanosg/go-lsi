package sync

import "time"

func (s *SyncTool) SyncOrders(start time.Time, end time.Time) {
    existingOrders, _ := s.Db.GetOrders()
    newOrders, _ := s.SquareClient.SearchOrders(start, end)

}
