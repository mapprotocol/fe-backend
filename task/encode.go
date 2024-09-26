package task

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func PackInput(abi abi.ABI, abiMethod string, params ...interface{}) ([]byte, error) {
	input, err := abi.Pack(abiMethod, params...)
	if err != nil {
		return nil, err
	}
	return input, nil
}

//// EncodeButterData
//// Solidity: function onReceived(uint256 _amount, (bytes32,uint256,address,bytes,uint256,uint256,address,bytes) _param) returns()
//func EncodeButterData(initiator, srcToken common.Address, swapData, bridgeData, feeData []byte) ([]byte, error) {
//	return PackInput(feRouterABI, "", initiator, srcToken, swapData, bridgeData, feeData)
//}

func EncodeButterData(initiator, dstToken common.Address, swapData, bridgeData, feeData []byte) ([]byte, error) {
	addressType, _ := abi.NewType("address", "string", []abi.ArgumentMarshaling{})
	bytesType, _ := abi.NewType("bytes", "string", nil)

	args := abi.Arguments{
		{Type: addressType},
		{Type: addressType},
		{Type: bytesType},
		{Type: bytesType},
		{Type: bytesType},
	}
	packed, err := args.Pack(initiator, dstToken, swapData, bridgeData, feeData)
	if err != nil {
		return nil, err
	}
	return packed, nil
}
