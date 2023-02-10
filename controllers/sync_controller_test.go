package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	mock "github.com/kampanosg/go-lsi/clients/db/mock_db"
	"github.com/kampanosg/go-lsi/sync"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

var (
	errSyncNotFound = errors.New("unable to retrieve sync status")
	errSyncFailed   = errors.New("syncing failed")
	from            = time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	to              = time.Date(2009, 11, 18, 20, 34, 58, 651387237, time.UTC)
)

func TestSyncStatusRequest(t *testing.T) {
	tests := []struct {
		name   string
		method string
		dbRes  types.SyncStatus
		dbErr  error
		status int
	}{
		{"fail when receive POST req", http.MethodPost, types.SyncStatus{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive CONNECT req", http.MethodConnect, types.SyncStatus{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive DELETE req", http.MethodDelete, types.SyncStatus{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive HEAD req", http.MethodHead, types.SyncStatus{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive OPTIONS req", http.MethodOptions, types.SyncStatus{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive PATCH req", http.MethodPatch, types.SyncStatus{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive PUT req", http.MethodPut, types.SyncStatus{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive TRACE req", http.MethodTrace, types.SyncStatus{}, nil, http.StatusMethodNotAllowed},
		{"fail when db returns error", http.MethodGet, types.SyncStatus{}, errSyncNotFound, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			logger, _ := zap.NewDevelopment()
			db := new(mock.MockDb)
			st := new(sync.MockSyncTool)

			db.On("GetLastSyncStatus").Return(tt.dbRes, tt.dbErr)

			ctrl := NewSyncController(st, db, logger.Sugar())

			router := mux.NewRouter()
			router.HandleFunc("/", ctrl.HandleSyncStatusRequest)

			req, err := http.NewRequest(tt.method, "/", nil)
			if err != nil {
				t.Errorf("threw error when calling endpoint, got %s", err.Error())
			}

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.status {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.status)
			}
		})
	}
}

func TestSyncRequest(t *testing.T) {

	tests := []struct {
		name    string
		method  string
		req     types.SyncRequest
		toolErr error
		status  int
	}{
		{"fail when receive GET req", http.MethodGet, types.SyncRequest{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive CONNECT req", http.MethodConnect, types.SyncRequest{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive DELETE req", http.MethodDelete, types.SyncRequest{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive HEAD req", http.MethodHead, types.SyncRequest{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive OPTIONS req", http.MethodOptions, types.SyncRequest{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive PATCH req", http.MethodPatch, types.SyncRequest{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive PUT req", http.MethodPut, types.SyncRequest{}, nil, http.StatusMethodNotAllowed},
		{"fail when receive TRACE req", http.MethodTrace, types.SyncRequest{}, nil, http.StatusMethodNotAllowed},
		{"fail when sync tool returns error", http.MethodPost, types.SyncRequest{From: from, To: to}, errSyncFailed, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			logger, _ := zap.NewDevelopment()
			db := new(mock.MockDb)
			st := new(sync.MockSyncTool)

			st.On("Sync", tt.req.From, tt.req.To).Return(tt.toolErr)

			ctrl := NewSyncController(st, db, logger.Sugar())

			router := mux.NewRouter()
			router.HandleFunc("/", ctrl.HandleSyncRequest)

			b, err := json.Marshal(tt.req)
			if err != nil {
				t.Errorf("threw unexpected error, got %s", err.Error())
			}

			req, err := http.NewRequest(tt.method, "/", bytes.NewBuffer(b))
			if err != nil {
				t.Errorf("threw error when calling endpoint, got %s", err.Error())
			}

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.status {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.status)
			}
		})
	}
}
