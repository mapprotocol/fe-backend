package tonrouter

import (
	"testing"
)

func TestBalance(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "t-1",
			want:    "0.071713",
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
			if got != tt.want {
				t.Errorf("Balance() got = %v, want %v", got, tt.want)
			}
		})
	}
}
