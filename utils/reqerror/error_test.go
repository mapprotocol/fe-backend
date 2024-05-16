package reqerror

import (
	"reflect"
	"testing"
)

func TestNewExternalRequestError(t *testing.T) {
	type args struct {
		path string
		opts []ErrorOption
	}
	tests := []struct {
		name string
		args args
		want *ExternalRequestError
	}{
		{
			name: "t-1",
			args: args{
				path: "/api/v1/users",
				opts: []ErrorOption{},
			},
			want: &ExternalRequestError{
				Path: "/api/v1/users",
			},
		},
		{
			name: "t-2",
			args: args{
				path: "/api/v1/users",
				opts: []ErrorOption{
					WithMethod("GET"),
				},
			},
			want: &ExternalRequestError{
				Path:   "/api/v1/users",
				Method: "GET",
			},
		},
		{
			name: "t-3",
			args: args{
				path: "/api/v1/users",
				opts: []ErrorOption{
					WithMethod("GET"),
					WithCode("4001"),
					WithMessage("failed to get users"),
				},
			},
			want: &ExternalRequestError{
				Path:    "/api/v1/users",
				Method:  "GET",
				Code:    "4001",
				Message: "failed to get users",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExternalRequestError(tt.args.path, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExternalRequestError() = %v, want %v", got, tt.want)
			}
		})
	}
}
