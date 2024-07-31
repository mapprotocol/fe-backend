package logic

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/fe-backend/utils/params"
	"math/big"
	"strings"
)

const (
	MethodNameOnReceived = "onReceived"
)

var (
	feRouterABI      = abi.ABI{}
	swapCallbackArgs = abi.Arguments{}
)

func init() {
	feRouter, err := abi.JSON(strings.NewReader(params.FeRouterABI))
	if err != nil {
		panic(err)
	}
	feRouterABI = feRouter

	//callbackParam, err := abi.JSON(strings.NewReader(ABI))
	//if err != nil {
	//	panic(err)
	//}
	//callbackParamABI = callbackParam

	swapCallbackParams, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "target", Type: "address"},
		{Name: "approveTo", Type: "address"},
		{Name: "offset", Type: "uint256"},
		{Name: "extraNativeAmount", Type: "uint256"},
		{Name: "receiver", Type: "address"},
		{Name: "data", Type: "bytes"},
	})
	if err != nil {
		panic(err)
	}

	swapCallbackArgs = abi.Arguments{{Type: swapCallbackParams}}

}

// SwapCallbackParams
/*
Solidity: struct CallbackParam {
	address target;
	address approveTo;
	uint256 offset;
	uint256 extraNativeAmount;
	address receiver;
	bytes data; // encoded onReceived function input params
}
*/
type SwapCallbackParams struct {
	Target            common.Address
	ApproveTo         common.Address
	Offset            *big.Int       // offset bit size, fixed to 36
	ExtraNativeAmount *big.Int       // 0
	Receiver          common.Address // refund address in case of error
	Data              []byte         //  packed onReceived function params
}

type ReceiverParam struct {
	OrderId        [32]byte
	SrcChain       *big.Int
	SrcToken       []byte
	Sender         []byte
	InAmount       string
	ChainPoolToken common.Address
	DstChain       *big.Int
	DstToken       []byte
	Receiver       []byte
	Slippage       uint64
}

func PackInput(abi abi.ABI, abiMethod string, params ...interface{}) ([]byte, error) {
	input, err := abi.Pack(abiMethod, params...)
	if err != nil {
		return nil, err
	}
	return input, nil
}

// PackOnReceived pack onReceived function params
// Solidity: function onReceived(uint256 _amount, (bytes32,uint256,address,bytes,uint256,uint256,address,bytes) _param) returns()
func PackOnReceived(amount *big.Int, params ReceiverParam) ([]byte, error) {
	return PackInput(feRouterABI, MethodNameOnReceived, amount, params)
}

func EncodeSwapCallbackParams(feRouter, receiver common.Address, data []byte) (string, error) {
	offset := 36
	extraNativeAmount := 0

	swapCallbackParams := SwapCallbackParams{
		Target:            feRouter,
		ApproveTo:         feRouter,
		Offset:            big.NewInt(int64(offset)),
		ExtraNativeAmount: big.NewInt(int64(extraNativeAmount)),
		Receiver:          receiver,
		Data:              data,
	}
	packed, err := swapCallbackArgs.Pack(swapCallbackParams)
	if err != nil {
		return "", err
	}

	return "0x" + hex.EncodeToString(packed), err
}
