package tonrouter

import (
	"github.com/mapprotocol/fe-backend/utils"
	"testing"
)

func TestRouteAndSwap(t *testing.T) {
	type args struct {
		request *BridgeSwapRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *RouteData
		wantErr bool
	}{
		//fromChainId=1360104473493505&
		//toChainId=1360104473493505&
		//amount=0.1&
		//slippage=500&
		//tokenInAddress=0x0000000000000000000000000000000000000000&
		//tokenOutAddress=EQBlqsm144Dq6SjbPI4jjZvA1hqTIP3CvHovbIfW_t-SCALE&
		//receiver=UQCcgIOWXxRpCWmQ8n2QLC2crtysNNjIAXzfhuqmVJEBH7Dl
		{
			name: "t-1",
			args: args{
				request: &BridgeSwapRequest{
					Sender:   "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					Receiver: "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					Hash:     "", // invalid hash
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "t-2",
			args: args{
				request: &BridgeSwapRequest{
					Sender:   "0xf855a761f9182c4b22A04753681A1F6324Ed3449", // "invalid sender 0xf855a761f9182c4b22A04753681A1F6324Ed3449"
					Receiver: "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					Hash:     "eyJkaWZmIjoiMCIsImhhc2giOiIiLCJicmlkZ2VGZWUiOnsiYW1vdW50IjoiMCIsInN5bWJvbCI6IlVTRFQifSwiZ2FzRmVlIjp7ImFtb3VudCI6IjEwMDAwMDAiLCJzeW1ib2wiOiJUb24ifSwibWluQW1vdW50T3V0Ijp7ImFtb3VudCI6IjY2OTI5NjciLCJzeW1ib2wiOiJVU0RUIn0sInNyY0NoYWluIjp7ImNoYWluSWQiOiIxMzYwMTA0NDczNDkzNTA1IiwidG9rZW5BbW91bnRJbiI6IjEiLCJ0b2tlbkFtb3VudE91dCI6IjYuNzYwNTczIiwicm91dGUiOlt7ImRleE5hbWUiOiJEZUR1c3QiLCJwYXRoIjpbeyJmZWUiOiIwIiwiaWQiOiJFUUEtWF95bzNmenpiRGJKXzBiekZXS3F0UnVaRklSYTFzSnN2ZVpKMVlwVmlPM3IiLCJ0b2tlbkluIjp7InR5cGUiOiJuYXRpdmUiLCJhZGRyZXNzIjoiMHgwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwIiwibmFtZSI6IlRvbmNvaW4iLCJzeW1ib2wiOiJUT04iLCJpbWFnZSI6Imh0dHBzOi8vYXNzZXRzLmRlZHVzdC5pby9pbWFnZXMvdG9uLndlYnAiLCJkZWNpbWFscyI6OSwiYWxpYXNlZCI6dHJ1ZSwicHJpY2UiOiI3LjY3Iiwic291cmNlIjpudWxsfSwidG9rZW5PdXQiOnsidHlwZSI6ImpldHRvbiIsImFkZHJlc3MiOiJFUUN4RTZtVXRRSktGbkdmYVJPVEtPdDFsWmJEaWlYMWtDaXhSdjdOdzJJZF9zRHMiLCJuYW1lIjoiVGV0aGVyIFVTRCIsInN5bWJvbCI6IlVTRFQiLCJpbWFnZSI6Imh0dHBzOi8vYXNzZXRzLmRlZHVzdC5pby9pbWFnZXMvdXNkdC53ZWJwIiwiZGVjaW1hbHMiOjYsImFsaWFzZWQiOnRydWUsInByaWNlIjoiMC45OTkxIiwic291cmNlIjp7ImNoYWluIjoiZWlwMTU1OjEiLCJhZGRyZXNzIjoiIiwiYnJpZGdlIjoiIiwic3ltYm9sIjoiVVNEVCIsIm5hbWUiOiJUZXRoZXIgVVNEIn19fV19XX0sInRpbWVzdGFtcCI6MTcyMjI0MzY0NDA0NCwidHJhZGVUeXBlIjowLCJicmlkZ2UiOnsidG9DaGFpbklkIjoiMSIsImRlc3RUb2tlbk91dEFkZHJlc3MiOiIweDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLCJ0b3RhbFNsaXBwYWdlIjoiMzAwIn19",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "t-3",
			args: args{
				request: &BridgeSwapRequest{
					Sender:   "UQAiqP-qy6O3Fbe8aGNTMz_aubZDMBqTqE0Nwf7qfyj1nB9c",
					Receiver: "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					Hash:     "eyJkaWZmIjoiMCIsImhhc2giOiIiLCJicmlkZ2VGZWUiOnsiYW1vdW50IjoiMCIsInN5bWJvbCI6IlVTRFQifSwiZ2FzRmVlIjp7ImFtb3VudCI6IjEwMDAwMDAiLCJzeW1ib2wiOiJUb24ifSwibWluQW1vdW50T3V0Ijp7ImFtb3VudCI6IjY2OTI5NjciLCJzeW1ib2wiOiJVU0RUIn0sInNyY0NoYWluIjp7ImNoYWluSWQiOiIxMzYwMTA0NDczNDkzNTA1IiwidG9rZW5BbW91bnRJbiI6IjEiLCJ0b2tlbkFtb3VudE91dCI6IjYuNzYwNTczIiwicm91dGUiOlt7ImRleE5hbWUiOiJEZUR1c3QiLCJwYXRoIjpbeyJmZWUiOiIwIiwiaWQiOiJFUUEtWF95bzNmenpiRGJKXzBiekZXS3F0UnVaRklSYTFzSnN2ZVpKMVlwVmlPM3IiLCJ0b2tlbkluIjp7InR5cGUiOiJuYXRpdmUiLCJhZGRyZXNzIjoiMHgwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwIiwibmFtZSI6IlRvbmNvaW4iLCJzeW1ib2wiOiJUT04iLCJpbWFnZSI6Imh0dHBzOi8vYXNzZXRzLmRlZHVzdC5pby9pbWFnZXMvdG9uLndlYnAiLCJkZWNpbWFscyI6OSwiYWxpYXNlZCI6dHJ1ZSwicHJpY2UiOiI3LjY3Iiwic291cmNlIjpudWxsfSwidG9rZW5PdXQiOnsidHlwZSI6ImpldHRvbiIsImFkZHJlc3MiOiJFUUN4RTZtVXRRSktGbkdmYVJPVEtPdDFsWmJEaWlYMWtDaXhSdjdOdzJJZF9zRHMiLCJuYW1lIjoiVGV0aGVyIFVTRCIsInN5bWJvbCI6IlVTRFQiLCJpbWFnZSI6Imh0dHBzOi8vYXNzZXRzLmRlZHVzdC5pby9pbWFnZXMvdXNkdC53ZWJwIiwiZGVjaW1hbHMiOjYsImFsaWFzZWQiOnRydWUsInByaWNlIjoiMC45OTkxIiwic291cmNlIjp7ImNoYWluIjoiZWlwMTU1OjEiLCJhZGRyZXNzIjoiIiwiYnJpZGdlIjoiIiwic3ltYm9sIjoiVVNEVCIsIm5hbWUiOiJUZXRoZXIgVVNEIn19fV19XX0sInRpbWVzdGFtcCI6MTcyMjI0MzY0NDA0NCwidHJhZGVUeXBlIjowLCJicmlkZ2UiOnsidG9DaGFpbklkIjoiMSIsImRlc3RUb2tlbk91dEFkZHJlc3MiOiIweDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLCJ0b3RhbFNsaXBwYWdlIjoiMzAwIn19",
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BridgeSwap(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouteAndSwap() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("RouteAndSwap() got = %v, want %v", utils.JSON(got), tt.want)
			//}
			t.Logf("RouteAndSwap() got = %v", utils.JSON(got))
		})
	}
}
