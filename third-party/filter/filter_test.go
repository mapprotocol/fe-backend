package filter

import (
	"github.com/mapprotocol/fe-backend/utils"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	endpoint := os.Getenv("FILTER_ROUTER_ENDPOINT")
	if utils.IsEmpty(endpoint) {
		log.Fatal("FILTER_ROUTER_ENDPOINT environment variable is not set")
	}
	Domain = endpoint

	m.Run()
}

func TestGetLogs(t *testing.T) {
	type args struct {
		id      uint64
		chainID string
		topic   string
		limit   uint8
	}
	tests := []struct {
		name    string
		args    args
		want    []*GetLogsResponseItem
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				id:      1,
				chainID: "1",
				topic:   "0xca1cf8cebf88499429cca8f87cbca15ab8dafd06702259a5344ddce89ef3f3a5",
				limit:   1,
			},
			want: []*GetLogsResponseItem{
				{
					Id:              473165,
					ProjectId:       1,
					ChainId:         1,
					EventId:         1,
					TxHash:          "0x8d2ff9340ce64b02941869ac0ffc67800aa89f391d94ef7ff2d7e51f94789781",
					ContractAddress: "0xfeB2b97e4Efce787c08086dC16Ab69E063911380",
					Topic:           "0xca1cf8cebf88499429cca8f87cbca15ab8dafd06702259a5344ddce89ef3f3a5,0x0000000000000000000000000000000000000000000000000000000000000001,0x00000000000000000000000000000000000000000000000000000000000058f8",
					BlockNumber:     19854610,
					BlockHash:       "",
					LogIndex:        358,
					LogData:         "34b94be87a7efac67876c25f36a8053286d23fc2168da068d62c147d5185375000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000002f5c589198f9c00000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000014c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000000000000000000000000000000000000000000000014b557571e2f3328ff5dd513e864b95c96e5fdadd90000000000000000000000000000000000000000000000000000000000000000000000000000000000000014b557571e2f3328ff5dd513e864b95c96e5fdadd90000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
					TxIndex:         111,
					TxTimestamp:     1715525015,
				},
			},

			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLogs(tt.args.id, tt.args.chainID, tt.args.topic, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLogs() got = %#v, want %v", got, tt.want)
			}
		})
	}
}
