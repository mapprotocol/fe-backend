package task

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/fe-backend/bindings/router"
	"github.com/mapprotocol/fe-backend/params"
	"math/big"
	"strings"
)

const (
	EventNameOnReceived = "OnReceived"
)

const (
	MethodNameSwapAndBridge = "swapAndBridge"
	MethodNameOnReceived    = "onReceived"
)

var (
	feRouterABI, _ = abi.JSON(strings.NewReader(params.FeRouterABI))
)

// SwapAndBridgeFunctionParams
// Solidity: function swapAndBridge(bytes32 _transferId, address _initiator, address _srcToken, uint256 _amount, bytes _swapData, bytes _bridgeData, bytes _permitData, bytes _feeData) payable returns(bytes32 orderId)
type SwapAndBridgeFunctionParams struct {
	TransferId [32]byte
	Initiator  common.Address
	SrcToken   common.Address
	Amount     *big.Int
	SwapData   []byte
	BridgeData []byte
	PermitData []byte
	FeeData    []byte
}

// OnReceivedEventParams represents a OnReceived event raised by the fe router contract.
// Solidity: event OnReceived(bytes32 _orderId, address _token, address _from, bytes to, uint256 _amount, address _caller)
type OnReceivedEventParams struct {
	OrderId [32]byte
	Token   common.Address
	From    common.Address
	To      []byte
	Amount  *big.Int
	Caller  common.Address
}

func DecodeData(data string) (*SwapAndBridgeFunctionParams, error) {
	bs := common.Hex2Bytes(strings.TrimPrefix(data, "0x"))

	getAbi, err := router.RouterMetaData.GetAbi() // todo
	if err != nil {
		return nil, err
	}

	swapAndBridge := getAbi.Methods[MethodNameSwapAndBridge]
	args, err := swapAndBridge.Inputs.Unpack(bs[4:])
	if err != nil {
		return nil, err
	}

	ret := &SwapAndBridgeFunctionParams{}
	if err := swapAndBridge.Inputs.Copy(ret, args); err != nil {
		return nil, err
	}
	return ret, nil
}

func UnpackLog(a abi.ABI, event string, ret interface{}, data []byte) error {
	inputs := a.Events[event].Inputs
	unpack, err := inputs.Unpack(data)
	if err != nil {
		return err
	}

	return inputs.Copy(ret, unpack)
}

func UnpackOnReceived(data []byte) (*OnReceivedEventParams, error) {
	ret := &OnReceivedEventParams{}
	if err := UnpackLog(feRouterABI, EventNameOnReceived, ret, data); err != nil {
		return nil, err
	}
	return ret, nil
}

//struct CallbackParam {
//	address target; // fe-router 地址
//	address approveTo; // fe-router 地址
//	uint256 offset; // 固定 36
//	uint256 extraNativeAmount; // 固定 0
//	address receiver; // 发生错误的退款地址， chain pool 所在链的用户地址(或者我们的一个集中地址)
//	bytes data; // encoded onReceived params
//}
//
//function onReceived(
//	uint256 _amount // chain pool out amount
//	bytes32 _orderId,  // 订单 ID
//	address _token, // chain pool token 地址
//	address _from,  // 用户地址
//	bytes _to,  // btc 地址
//) returns()

type SwapCallCallbackParams struct {
	Target            common.Address
	ApproveTo         common.Address
	Offset            *big.Int
	ExtraNativeAmount *big.Int
	Receiver          common.Address
	Data              []byte //  pack onReceived function params
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

func PackOnReceived(params ...interface{}) ([]byte, error) {
	return PackInput(feRouterABI, MethodNameOnReceived, params...)
}

func PackSwapCallbackParams(feRouter, receiver common.Address, data []byte) (string, error) {
	swapCallbackParams, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "target", Type: "address"},
		{Name: "approveTo", Type: "address"},
		{Name: "offset", Type: "uint256"},
		{Name: "extraNativeAmount", Type: "uint256"},
		{Name: "receiver", Type: "address"},
		{Name: "data", Type: "bytes"},
	})
	if err != nil {
		return "", err
	}

	args := abi.Arguments{{Type: swapCallbackParams}}

	offset := 36
	extraNativeAmount := 0
	packed, err := args.Pack(feRouter, feRouter, offset, extraNativeAmount, receiver, data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(packed), err
}
