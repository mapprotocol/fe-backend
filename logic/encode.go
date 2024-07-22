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

const CallbackParam = "callbackParam"

var (
	feRouterABI      = abi.ABI{}
	callbackParamABI = abi.ABI{}
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

// OnReceivedFunctionParams
// Solidity: function onReceived(uint256 _amount, bytes32 _orderId, address _token, address _from, bytes _to) returns()
type OnReceivedFunctionParams struct {
	Amount  *big.Int
	OrderId [32]byte
	Token   common.Address
	From    common.Address
	To      []byte
}

func PackInput(abi abi.ABI, abiMethod string, params ...interface{}) ([]byte, error) {
	input, err := abi.Pack(abiMethod, params...)
	if err != nil {
		return nil, err
	}
	return input, nil
}

// PackOnReceived pack onReceived function params
// amount: token amount
// orderId: order id
// token: token address
// from: sender address
// to: receiver address on bitcoin
// Solidity: function onReceived(uint256 _amount, bytes32 _orderId, address _token, address _from, bytes _to) returns()
func PackOnReceived(amount *big.Int, orderId [32]byte, token common.Address, from common.Address, to []byte) ([]byte, error) {
	return PackInput(feRouterABI, MethodNameOnReceived, amount, orderId, token, from, to)
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
