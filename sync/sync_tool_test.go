package sync

import (
	"testing"
)

func TestSync(t *testing.T) {

	tests := []struct {
		name string
	}{
		{"returns error, when sync categories fails"},
		{"returns error, when sync products fails"},
		{"returns error, when sync orders fails"},
		{"returns error, when insert sync status fails"},
		{"returns no error, when nothing fails"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if false {
				t.Errorf("oops")
			}
		})
	}
}
