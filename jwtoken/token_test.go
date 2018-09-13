package jwtoken

import (
	"reflect"
	"testing"

	"github.com/madappgang/identifo/model"
)

func TestNewTokenService(t *testing.T) {
	tests := []struct {
		name string
		want model.TokenService
	}{
		{"create new service", NewTokenService()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTokenService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenService() = %v, want %v", got, tt.want)
			}
		})
	}
}
