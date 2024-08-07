package utils

import "testing"

func TestTrimHexPrefix(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t-1",
			args: args{
				s: "0x1234567890abcdef",
			},
			want: "1234567890abcdef",
		},
		{
			name: "t-2",
			args: args{
				s: "0X1234567890abcdef",
			},
			want: "1234567890abcdef",
		},
		{
			name: "t-3",
			args: args{
				s: "1234567890abcdef",
			},
			want: "1234567890abcdef",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimHexPrefix(tt.args.s); got != tt.want {
				t.Errorf("TrimHexPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
