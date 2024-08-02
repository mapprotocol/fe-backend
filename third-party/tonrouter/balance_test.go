package tonrouter

import (
	"math/big"
	"reflect"
	"testing"
)

func TestBalance(t *testing.T) {
	tests := []struct {
		name    string
		want    *big.Float
		wantErr bool
	}{
		{
			name:    "t-1",
			want:    big.NewFloat(3.628543),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Balance()
			if (err != nil) != tt.wantErr {
				t.Errorf("Balance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Balance() got = %v, want %v", got, tt.want)
			}
		})
	}
}
