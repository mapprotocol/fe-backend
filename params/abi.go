package params

var FeRouterABI = `[
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "bytes32",
          "name": "orderId",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "token",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "receiver",
          "type": "address"
        }
      ],
      "name": "Deliver",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "bytes32",
          "name": "orderId",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "bytes32",
          "name": "bridgeId",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "token",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "DeliverAndSwap",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "bytes32",
          "name": "orderId",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "uint64",
          "name": "bridgeId",
          "type": "uint64"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "srcChain",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "bytes",
          "name": "srcToken",
          "type": "bytes"
        },
        {
          "indexed": false,
          "internalType": "string",
          "name": "inAmount",
          "type": "string"
        },
        {
          "indexed": false,
          "internalType": "bytes",
          "name": "sender",
          "type": "bytes"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "chainPoolToken",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "chainPoolTokenAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "dstChain",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "bytes",
          "name": "dstToken",
          "type": "bytes"
        },
        {
          "indexed": false,
          "internalType": "bytes",
          "name": "receiver",
          "type": "bytes"
        },
        {
          "indexed": false,
          "internalType": "uint64",
          "name": "slippage",
          "type": "uint64"
        }
      ],
      "name": "OnReceived",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "orderId",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "token",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "receiver",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "fee",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "feeReceiver",
          "type": "address"
        }
      ],
      "name": "deliver",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "orderId",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "initiator",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "token",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "swapData",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "bridgeData",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "feeData",
          "type": "bytes"
        },
        {
          "internalType": "uint256",
          "name": "fee",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "feeReceiver",
          "type": "address"
        }
      ],
      "name": "deliverAndSwap",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        },
        {
          "components": [
            {
              "internalType": "bytes32",
              "name": "orderId",
              "type": "bytes32"
            },
            {
              "internalType": "uint256",
              "name": "srcChain",
              "type": "uint256"
            },
            {
              "internalType": "bytes",
              "name": "srcToken",
              "type": "bytes"
            },
            {
              "internalType": "bytes",
              "name": "sender",
              "type": "bytes"
            },
            {
              "internalType": "string",
              "name": "inAmount",
              "type": "string"
            },
            {
              "internalType": "address",
              "name": "chainPoolToken",
              "type": "address"
            },
            {
              "internalType": "uint256",
              "name": "dstChain",
              "type": "uint256"
            },
            {
              "internalType": "bytes",
              "name": "dstToken",
              "type": "bytes"
            },
            {
              "internalType": "bytes",
              "name": "receiver",
              "type": "bytes"
            },
            {
              "internalType": "uint64",
              "name": "slippage",
              "type": "uint64"
            }
          ],
          "internalType": "struct IRouter.ReceiverParam",
          "name": "_param",
          "type": "tuple"
        }
      ],
      "name": "onReceived",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`

var ButterRouterV3 = `[
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_bridgeAddress",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_owner",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_wToken",
          "type": "address"
        }
      ],
      "stateMutability": "payable",
      "type": "constructor"
    },
    {
      "inputs": [],
      "name": "BRIDGE_ONLY",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "CALL_BACK_FAIL",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "DATA_EMPTY",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "EMPTY",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "FEE_MISMATCH",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "NATIVE_VALUE_OVERSPEND",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "NOT_CONTRACT",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "NO_APPROVE",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "RECEIVE_LOW",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "SELF_ONLY",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "SWAP_FAIL",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "SWAP_SAME_TOKEN",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "UNSUPPORT_DEX_TYPE",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "ZERO_ADDRESS",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "ZERO_IN",
      "type": "error"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "executor",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "bool",
          "name": "flag",
          "type": "bool"
        }
      ],
      "name": "Approve",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "token",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "receiver",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "integrator",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "routerAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "integratorAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "nativeAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "integratorNative",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "bytes32",
          "name": "transferId",
          "type": "bytes32"
        }
      ],
      "name": "CollectFee",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "previousOwner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "OwnershipTransferStarted",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "previousOwner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "OwnershipTransferred",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "bytes32",
          "name": "orderId",
          "type": "bytes32"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "receiver",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "target",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "originToken",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "swapToken",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "originAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "swapAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "callAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "fromChain",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "toChain",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "bytes",
          "name": "from",
          "type": "bytes"
        }
      ],
      "name": "RemoteSwapAndCall",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "_bridgeAddress",
          "type": "address"
        }
      ],
      "name": "SetBridgeAddress",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "receiver",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "rate",
          "type": "uint256"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "fixedf",
          "type": "uint256"
        }
      ],
      "name": "SetFee",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "_feeManager",
          "type": "address"
        }
      ],
      "name": "SetFeeManager",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "_gasForReFund",
          "type": "uint256"
        }
      ],
      "name": "SetGasForReFund",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "_maxFeeRate",
          "type": "uint256"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "_maxNativeFee",
          "type": "uint256"
        }
      ],
      "name": "SetReferrerMaxFee",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "referrer",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "initiator",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "bytes32",
          "name": "transferId",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "bytes32",
          "name": "orderId",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "originToken",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "bridgeToken",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "originAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "bridgeAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "toChain",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "bytes",
          "name": "to",
          "type": "bytes"
        }
      ],
      "name": "SwapAndBridge",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "referrer",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "initiator",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "bytes32",
          "name": "transferId",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "originToken",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "swapToken",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "originAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "swapAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "receiver",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "target",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "callAmount",
          "type": "uint256"
        }
      ],
      "name": "SwapAndCall",
      "type": "event"
    },
    {
      "inputs": [],
      "name": "acceptOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "approved",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "bridgeAddress",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "components": [
            {
              "internalType": "address",
              "name": "target",
              "type": "address"
            },
            {
              "internalType": "address",
              "name": "approveTo",
              "type": "address"
            },
            {
              "internalType": "uint256",
              "name": "offset",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "extraNativeAmount",
              "type": "uint256"
            },
            {
              "internalType": "address",
              "name": "receiver",
              "type": "address"
            },
            {
              "internalType": "bytes",
              "name": "data",
              "type": "bytes"
            }
          ],
          "internalType": "struct SwapCall.CallbackParam",
          "name": "_callbackParam",
          "type": "tuple"
        },
        {
          "internalType": "address",
          "name": "_callToken",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        }
      ],
      "name": "doRemoteCall",
      "outputs": [
        {
          "internalType": "address",
          "name": "target",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "callAmount",
          "type": "uint256"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "components": [
            {
              "internalType": "address",
              "name": "dstToken",
              "type": "address"
            },
            {
              "internalType": "address",
              "name": "receiver",
              "type": "address"
            },
            {
              "internalType": "address",
              "name": "leftReceiver",
              "type": "address"
            },
            {
              "internalType": "uint256",
              "name": "minAmount",
              "type": "uint256"
            },
            {
              "components": [
                {
                  "internalType": "enum SwapCall.DexType",
                  "name": "dexType",
                  "type": "uint8"
                },
                {
                  "internalType": "address",
                  "name": "callTo",
                  "type": "address"
                },
                {
                  "internalType": "address",
                  "name": "approveTo",
                  "type": "address"
                },
                {
                  "internalType": "uint256",
                  "name": "fromAmount",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes",
                  "name": "callData",
                  "type": "bytes"
                }
              ],
              "internalType": "struct SwapCall.SwapData[]",
              "name": "swaps",
              "type": "tuple[]"
            }
          ],
          "internalType": "struct SwapCall.SwapParam",
          "name": "swapParam",
          "type": "tuple"
        },
        {
          "internalType": "address",
          "name": "_srcToken",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        }
      ],
      "name": "doRemoteSwap",
      "outputs": [
        {
          "internalType": "address",
          "name": "dstToken",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "dstAmount",
          "type": "uint256"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "feeManager",
      "outputs": [
        {
          "internalType": "contract IFeeManager",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "feeReceiver",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "gasForReFund",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_token",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amountAfterFee",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "_feeData",
          "type": "bytes"
        }
      ],
      "name": "getAmountBeforeFee",
      "outputs": [
        {
          "internalType": "address",
          "name": "feeToken",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "beforeAmount",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "nativeFeeAmount",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_inputToken",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_inputAmount",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "_feeData",
          "type": "bytes"
        }
      ],
      "name": "getFee",
      "outputs": [
        {
          "internalType": "address",
          "name": "feeToken",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "tokenFee",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "nativeFee",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "afterFeeAmount",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_inputToken",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_inputAmount",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "_feeData",
          "type": "bytes"
        }
      ],
      "name": "getFeeDetail",
      "outputs": [
        {
          "components": [
            {
              "internalType": "address",
              "name": "feeToken",
              "type": "address"
            },
            {
              "internalType": "address",
              "name": "routerReceiver",
              "type": "address"
            },
            {
              "internalType": "address",
              "name": "integrator",
              "type": "address"
            },
            {
              "internalType": "uint256",
              "name": "routerNativeFee",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "integratorNativeFee",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "routerTokenFee",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "integratorTokenFee",
              "type": "uint256"
            }
          ],
          "internalType": "struct IFeeManager.FeeDetail",
          "name": "feeDetail",
          "type": "tuple"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_token",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amountAfterFee",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "_feeData",
          "type": "bytes"
        }
      ],
      "name": "getInputBeforeFee",
      "outputs": [
        {
          "internalType": "address",
          "name": "_feeToken",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_input",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_fee",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "maxFeeRate",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "maxNativeFee",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "_orderId",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "_srcToken",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_fromChain",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "_from",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "_swapAndCall",
          "type": "bytes"
        }
      ],
      "name": "onReceived",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "owner",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "pendingOwner",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "renounceOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "routerFeeRate",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "routerFixedFee",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address[]",
          "name": "_executors",
          "type": "address[]"
        },
        {
          "internalType": "bool",
          "name": "_flag",
          "type": "bool"
        }
      ],
      "name": "setAuthorization",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_bridgeAddress",
          "type": "address"
        }
      ],
      "name": "setBridgeAddress",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_feeReceiver",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_feeRate",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_fixedFee",
          "type": "uint256"
        }
      ],
      "name": "setFee",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_feeManager",
          "type": "address"
        }
      ],
      "name": "setFeeManager",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "_gasForReFund",
          "type": "uint256"
        }
      ],
      "name": "setGasForReFund",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "_maxFeeRate",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_maxNativeFee",
          "type": "uint256"
        }
      ],
      "name": "setReferrerMaxFee",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "_transferId",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "_initiator",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_srcToken",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "_swapData",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "_bridgeData",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "_permitData",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "_feeData",
          "type": "bytes"
        }
      ],
      "name": "swapAndBridge",
      "outputs": [
        {
          "internalType": "bytes32",
          "name": "orderId",
          "type": "bytes32"
        }
      ],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "_transferId",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "_initiator",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_srcToken",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "_amount",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "_swapData",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "_callbackData",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "_permitData",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "_feeData",
          "type": "bytes"
        }
      ],
      "name": "swapAndCall",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "transferOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "stateMutability": "payable",
      "type": "receive"
    }
  ]`
