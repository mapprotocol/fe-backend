package task

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/fe-backend/params"
	"math/big"
	"strings"
)

const (
	EventNameOnReceived     = "OnReceived"
	EventNameDeliverAndSwap = "DeliverAndSwap"
)

const (
	MethodNameSwapAndBridge = "swapAndBridge"
	MethodNameOnReceived    = "onReceived"
)

var (
	feRouterABI, _     = abi.JSON(strings.NewReader(params.FeRouterABI))
	ButterRouterABI, _ = abi.JSON(strings.NewReader(params.ButterRouterV3))
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
	OrderId              [32]byte
	BridgeId             uint64
	SrcChain             *big.Int
	SrcToken             []byte
	InAmount             string
	Sender               []byte
	ChainPoolToken       common.Address
	ChainPoolTokenAmount *big.Int
	DstChain             *big.Int
	DstToken             []byte
	Receiver             []byte
	Slippage             uint64
}

// DeliverAndSwapEventParams represents a DeliverAndSwap event raised by the fe router contract.
// Solidity: event DeliverAndSwap(bytes32 orderId, bytes32 bridgeId, address token, uint256 amount)
type DeliverAndSwapEventParams struct {
	OrderId  [32]byte
	BridgeId [32]byte
	Token    common.Address
	Amount   *big.Int
}

func DecodeData(data string) (*SwapAndBridgeFunctionParams, error) {
	bs := common.Hex2Bytes(strings.TrimPrefix(data, "0x"))

	swapAndBridge := ButterRouterABI.Methods[MethodNameSwapAndBridge]
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

func UnpackDeliverAndSwap(data []byte) (*DeliverAndSwapEventParams, error) {
	ret := &DeliverAndSwapEventParams{}
	if err := UnpackLog(feRouterABI, EventNameDeliverAndSwap, ret, data); err != nil {
		return nil, err
	}
	return ret, nil
}
