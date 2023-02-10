package sync

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockSyncTool struct {
	mock.Mock
}

func (m *MockSyncTool) Sync(from time.Time, to time.Time) error {
	args := m.Called(from, to)
	return args.Error(0)
}

func (m *MockSyncTool) SyncCategories() error {
	args := m.Called()
	return args.Error(0)

}

func (m *MockSyncTool) SyncProducts() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSyncTool) SyncOrders(start time.Time, end time.Time) error {
	args := m.Called(start, end)
	return args.Error(0)
}
