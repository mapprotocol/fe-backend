package http

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func sign(data string, secret string) (string, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return "", err
	}

	hashed := sha512.Sum512([]byte(data))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA512, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func TestSign(t *testing.T) {
	data := "example data"
	secret := "example secret key"
	signed, err := sign(data, secret)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Signature:", signed)
	}
}

func TestCreatePrimeWallet(t *testing.T) {
	type args struct {
		url     string
		headers http.Header
		body    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				url: "/open-api/v1/wallet/create",
				headers: http.Header{
					"open-apikey":  []string{""},
					"signature":    []string{""},
					"Content-Type": []string{"application/json"},
				},
				body: strings.NewReader(`{
    "requestId": 0,
    "timestamp": 0,
    "walletName": "neoiss",
    "walletType": 10
}`),
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Post(tt.args.url, tt.args.headers, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Post() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateSubWallet(t *testing.T) {
	type args struct {
		url     string
		headers http.Header
		body    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				url: "/open-api/v1/subwallet/create",
				headers: http.Header{
					"open-apikey":  []string{""},
					"signature":    []string{""},
					"Content-Type": []string{"application/json"},
				},
				body: nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Post(tt.args.url, tt.args.headers, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSubWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSubWallet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepositFlow(t *testing.T) {
	type args struct {
		url     string
		headers http.Header
		body    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				url: "/open-api/v1/subwallet/deposit/address",
				headers: http.Header{
					"open-apikey":  []string{""},
					"signature":    []string{""},
					"Content-Type": []string{"application/json"},
				},
				body: nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Post(tt.args.url, tt.args.headers, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSubWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSubWallet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithdrawalFlowTransfer(t *testing.T) {
	type args struct {
		url     string
		headers http.Header
		body    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				url: "/open-api/v1/subwallet/transfer",
				headers: http.Header{
					"open-apikey":  []string{""},
					"signature":    []string{""},
					"Content-Type": []string{"application/json"},
				},
				body: nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Post(tt.args.url, tt.args.headers, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSubWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSubWallet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithdrawalFlowWithdrawal(t *testing.T) {
	type args struct {
		url     string
		headers http.Header
		body    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				url: "/open-api/v1/subwallet/withdrawal",
				headers: http.Header{
					"open-apikey":  []string{""},
					"signature":    []string{""},
					"Content-Type": []string{"application/json"},
				},
				body: nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Post(tt.args.url, tt.args.headers, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSubWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSubWallet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssetList(t *testing.T) {
	type args struct {
		url     string
		headers http.Header
		body    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				url: "/open-api/v1/subwallet/asset/list",
				headers: http.Header{
					"open-apikey":  []string{""},
					"signature":    []string{""},
					"Content-Type": []string{"application/json"},
				},
				body: nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Post(tt.args.url, tt.args.headers, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSubWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSubWallet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExchange(t *testing.T) {
	type args struct {
		url     string
		headers http.Header
		body    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				url: "/open-api/v1/wallet/transfer/exchange",
				headers: http.Header{
					"open-apikey":  []string{""},
					"signature":    []string{""},
					"Content-Type": []string{"application/json"},
				},
				body: nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Post(tt.args.url, tt.args.headers, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSubWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSubWallet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName(t *testing.T) {
	fmt.Println("============================== ", time.Now().UnixMilli())
}
