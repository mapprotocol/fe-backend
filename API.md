## supported chains

| chain   | chain id           | 
|---------|--------------------|
| TON     | 1360104473493505   |
| Bitcoin | 313230561203979757 |

## route

### request path

/api/v1/route

### request method

GET

### request params

| parameter       | type   | required | default | description                                                        |
|-----------------|--------|----------|---------|--------------------------------------------------------------------|
| fromChainId     | string | Yes      |         | action: 1， fromChainId: 1360104473493505                           |
| toChainId       | string | Yes      |         | action: 2， toChainId: 1360104473493505                             |
| amount          | string | Yes      |         |                                                                    |
| tokenInAddress  | string | Yes      |         |                                                                    |
| tokenOutAddress | string | Yes      |         |                                                                    |
| feeCollector    | string | No       |         |                                                                    |
| feeRatio        | string | No       |         | 200 means 2%                                                       |
| type            | string | Yes      |         |                                                                    |
| slippage        | string | Yes      |         | slippage of swap, a integer in rang [300, 5000], e.g, 300 means 3% |
| action          | number | Yes      |         | 1: to evm, 2: from evm                                             |

### response params

| parameter | type     | description   |
|-----------|----------|---------------|
| code      | number   | response code |
| msg       | string   | response msg  |
| data      | []object | response data |

#### data structure

| parameter | type | description |
|-----------|------|-------------|

### Example

**request**:

```shell
curl --location http://127.0.0.1:8123/api/v1/route?fromChainId=313230561203979757&amount=100&toChainId=56&tokenInAddress=0x0000000000000000000000000000000000000000&tokenOutAddress=0x0000000000000000000000000000000000000000&receiver=0x766f3377497C66c31a5692A435cF3E72Dcc2d4Fc&slippage=300&from=bc1qr2rkrw6a2s79gdc8dyhx3q7k3u6an2wug7wmqk&type=exactIn&action=1&feeCollector=bc1qr2rkrw6a2s79gdc8dyhx3q7k3u6an2wug7wmqk&feeRatio=70&timestamp=1724136612637
```

**response**

```json
{
  "code": 2000,
  "msg": "Success",
  "data": {
    "total": 4,
    "items": [
      {
        "hash": "0x0000000000000000000000000000000000000000000000000000000000022776",
        "tokenIn": {
          "chainId": "313230561203979757",
          "address": "0x0000000000000000000000000000000000000000",
          "name": "Bitcoin",
          "decimals": 8,
          "symbol": "BTC",
          "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
        },
        "tokenOut": {
          "chainId": "56",
          "address": "0x0000000000000000000000000000000000000000",
          "name": "BNB",
          "decimals": 18,
          "symbol": "BNB",
          "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/binance-smart-chain/0x0000000000000000000000000000000000000000/logo.png"
        },
        "amountIn": "98.6",
        "amountOut": "0.080763104168497019",
        "path": [
          {
            "name": "MAP FE",
            "amountIn": "98.6",
            "amountOut": "98.6",
            "tokenIn": {
              "chainId": "313230561203979757",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
            },
            "tokenOut": {
              "chainId": "313230561203979757",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
            }
          },
          {
            "name": "MAP FE",
            "amountIn": "98.6",
            "amountOut": "98.6",
            "tokenIn": {
              "chainId": "313230561203979757",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
            },
            "tokenOut": {
              "chainId": "137",
              "address": "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
              "name": "Wrapped MATIC",
              "decimals": 18,
              "symbol": "WMATIC",
              "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/polygon/0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270/logo.png"
            }
          },
          {
            "name": "Butter",
            "amountIn": "98.6",
            "amountOut": "0.080763104168497019",
            "tokenIn": {
              "chainId": "137",
              "address": "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
              "name": "Wrapped MATIC",
              "decimals": 18,
              "symbol": "WMATIC",
              "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/polygon/0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270/logo.png"
            },
            "tokenOut": {
              "chainId": "56",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "BNB",
              "decimals": 18,
              "symbol": "BNB",
              "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/binance-smart-chain/0x0000000000000000000000000000000000000000/logo.png"
            }
          }
        ],
        "gasFee": {
          "amount": "0.00007614",
          "symbol": "BTC",
          "chainId": ""
        },
        "bridgeFee": {
          "amount": "0.7",
          "symbol": "BTC",
          "chainId": ""
        },
        "protocolFee": {
          "amount": "0.7",
          "symbol": "BTC",
          "chainId": ""
        }
      },
      {
        "hash": "0x0000000000000000000000000000000000000000000000000000000000022776",
        "tokenIn": {
          "chainId": "313230561203979757",
          "address": "0x0000000000000000000000000000000000000000",
          "name": "Bitcoin",
          "decimals": 8,
          "symbol": "BTC",
          "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
        },
        "tokenOut": {
          "chainId": "56",
          "address": "0x0000000000000000000000000000000000000000",
          "name": "BNB",
          "decimals": 18,
          "symbol": "BNB",
          "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/binance-smart-chain/0x0000000000000000000000000000000000000000/logo.png"
        },
        "amountIn": "98.6",
        "amountOut": "0.080634266305625526",
        "path": [
          {
            "name": "MAP FE",
            "amountIn": "98.6",
            "amountOut": "98.6",
            "tokenIn": {
              "chainId": "313230561203979757",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
            },
            "tokenOut": {
              "chainId": "313230561203979757",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
            }
          },
          {
            "name": "MAP FE",
            "amountIn": "98.6",
            "amountOut": "98.6",
            "tokenIn": {
              "chainId": "313230561203979757",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
            },
            "tokenOut": {
              "chainId": "137",
              "address": "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
              "name": "Wrapped MATIC",
              "decimals": 18,
              "symbol": "WMATIC",
              "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/polygon/0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270/logo.png"
            }
          },
          {
            "name": "Butter",
            "amountIn": "98.6",
            "amountOut": "0.080634266305625526",
            "tokenIn": {
              "chainId": "137",
              "address": "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
              "name": "Wrapped MATIC",
              "decimals": 18,
              "symbol": "WMATIC",
              "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/polygon/0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270/logo.png"
            },
            "tokenOut": {
              "chainId": "56",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "BNB",
              "decimals": 18,
              "symbol": "BNB",
              "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/binance-smart-chain/0x0000000000000000000000000000000000000000/logo.png"
            }
          }
        ],
        "gasFee": {
          "amount": "0.00007614",
          "symbol": "BTC",
          "chainId": ""
        },
        "bridgeFee": {
          "amount": "0.7",
          "symbol": "BTC",
          "chainId": ""
        },
        "protocolFee": {
          "amount": "0.7",
          "symbol": "BTC",
          "chainId": ""
        }
      },
      {
        "hash": "0x0000000000000000000000000000000000000000000000000000000000022776",
        "tokenIn": {
          "chainId": "313230561203979757",
          "address": "0x0000000000000000000000000000000000000000",
          "name": "Bitcoin",
          "decimals": 8,
          "symbol": "BTC",
          "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
        },
        "tokenOut": {
          "chainId": "56",
          "address": "0x0000000000000000000000000000000000000000",
          "name": "BNB",
          "decimals": 18,
          "symbol": "BNB",
          "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/binance-smart-chain/0x0000000000000000000000000000000000000000/logo.png"
        },
        "amountIn": "98.6",
        "amountOut": "0.080113867232187782",
        "path": [
          {
            "name": "MAP FE",
            "amountIn": "98.6",
            "amountOut": "98.6",
            "tokenIn": {
              "chainId": "313230561203979757",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
            },
            "tokenOut": {
              "chainId": "313230561203979757",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
            }
          },
          {
            "name": "MAP FE",
            "amountIn": "98.6",
            "amountOut": "98.6",
            "tokenIn": {
              "chainId": "313230561203979757",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
            },
            "tokenOut": {
              "chainId": "137",
              "address": "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
              "name": "Wrapped MATIC",
              "decimals": 18,
              "symbol": "WMATIC",
              "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/polygon/0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270/logo.png"
            }
          },
          {
            "name": "Butter",
            "amountIn": "98.6",
            "amountOut": "0.080113867232187782",
            "tokenIn": {
              "chainId": "137",
              "address": "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
              "name": "Wrapped MATIC",
              "decimals": 18,
              "symbol": "WMATIC",
              "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/polygon/0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270/logo.png"
            },
            "tokenOut": {
              "chainId": "56",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "BNB",
              "decimals": 18,
              "symbol": "BNB",
              "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/binance-smart-chain/0x0000000000000000000000000000000000000000/logo.png"
            }
          }
        ],
        "gasFee": {
          "amount": "0.00007614",
          "symbol": "BTC",
          "chainId": ""
        },
        "bridgeFee": {
          "amount": "0.7",
          "symbol": "BTC",
          "chainId": ""
        },
        "protocolFee": {
          "amount": "0.7",
          "symbol": "BTC",
          "chainId": ""
        }
      },
      {
        "hash": "0x0000000000000000000000000000000000000000000000000000000000022776",
        "tokenIn": {
          "chainId": "313230561203979757",
          "address": "0x0000000000000000000000000000000000000000",
          "name": "Bitcoin",
          "decimals": 8,
          "symbol": "BTC",
          "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
        },
        "tokenOut": {
          "chainId": "56",
          "address": "0x0000000000000000000000000000000000000000",
          "name": "BNB",
          "decimals": 18,
          "symbol": "BNB",
          "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/binance-smart-chain/0x0000000000000000000000000000000000000000/logo.png"
        },
        "amountIn": "98.6",
        "amountOut": "0.079833380464649455",
        "path": [
          {
            "name": "MAP FE",
            "amountIn": "98.6",
            "amountOut": "98.6",
            "tokenIn": {
              "chainId": "313230561203979757",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
            },
            "tokenOut": {
              "chainId": "313230561203979757",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
            }
          },
          {
            "name": "MAP FE",
            "amountIn": "98.6",
            "amountOut": "98.6",
            "tokenIn": {
              "chainId": "313230561203979757",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg"
            },
            "tokenOut": {
              "chainId": "137",
              "address": "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
              "name": "Wrapped MATIC",
              "decimals": 18,
              "symbol": "WMATIC",
              "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/polygon/0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270/logo.png"
            }
          },
          {
            "name": "Butter",
            "amountIn": "98.6",
            "amountOut": "0.079833380464649455",
            "tokenIn": {
              "chainId": "137",
              "address": "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
              "name": "Wrapped MATIC",
              "decimals": 18,
              "symbol": "WMATIC",
              "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/polygon/0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270/logo.png"
            },
            "tokenOut": {
              "chainId": "56",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "BNB",
              "decimals": 18,
              "symbol": "BNB",
              "icon": "https://s3.amazonaws.com/map-static-file/mapSwap/binance-smart-chain/0x0000000000000000000000000000000000000000/logo.png"
            }
          }
        ],
        "gasFee": {
          "amount": "0.00007614",
          "symbol": "BTC",
          "chainId": ""
        },
        "bridgeFee": {
          "amount": "0.7",
          "symbol": "BTC",
          "chainId": ""
        },
        "protocolFee": {
          "amount": "0.7",
          "symbol": "BTC",
          "chainId": ""
        }
      }
    ]
  }
}
```

## swap

### request path

/api/v1/swap

### request method

GET

### request params

| parameter    | type   | required | default | description                                                        |
|--------------|--------|----------|---------|--------------------------------------------------------------------|
| srcChain     | string | Yes      |         |                                                                    |
| srcToken     | string | Yes      |         |                                                                    |
| sender       | string | Yes      |         |                                                                    |
| amount       | string | Yes      |         |                                                                    |
| decimal      | number | Yes      |         |                                                                    |
| dstChain     | string | Yes      |         |                                                                    |
| dstToken     | string | Yes      |         |                                                                    |
| receiver     | string | Yes      |         |                                                                    |
| feeCollector | string | No       |         |                                                                    |
| feeRatio     | string | No       |         | 200 means 2%                                                       |
| hash         | string | Yes      |         | the route hash returned by /api/v1/route                           |
| slippage     | string | Yes      |         | slippage of swap, a integer in rang [300, 5000], e.g, 300 means 3% |

### response params

| parameter | type   | description   |
|-----------|--------|---------------|
| code      | number | response code |
| msg       | string | response msg  |
| data      | object | response data |

#### data structure

| parameter | type   | description |
|-----------|--------|-------------|
| to        | string |             |
| data      | string |             |
| value     | string |             |
| chainId   | string |             |

### Example

**request**:

```shell
curl --location '127.0.0.1:8181/api/v1/swap?srcChain=1&srcToken=0x0000000000000000000000000000000000000000&sender=0x0000000000000000000000000000000000000000&amount=10&dstToken=0x0000000000000000000000000000000000000000&receiver=0x0000000000000000000000000000000000000000...'
```

**response**

```json
{
  "code": 2000,
  "msg": "Success",
  "data": {
    "to": "0xEE3020a308B0E9F6765279C595f17a534CCC7019",
    "data": "0x6e1537da0000000000000000000......",
    "value": "0x1043561a8829300000",
    "chainId": "22776"
  }
}
```

## order list

### request path

/api/v1/order/list

### request method

GET

### request params

| parameter | type   | required | default | description |
|-----------|--------|----------|---------|-------------|
| sender    | string | Yes      |         |             |
| page      | number | Yes      | 1       |             |       
| size      | number | Yes      | 20      |             |       

### response params

| parameter | type   | description   |
|-----------|--------|---------------|
| code      | number | response code |
| msg       | string | response msg  |
| data      | object | response data |

#### data structure

| parameter | type     | description  |
|-----------|----------|--------------|
| page      | number   |              |
| total     | number   | record total |
| items     | []object |              |

#### items structure

| parameter | type   | description                            |
|-----------|--------|----------------------------------------|
| orderId   | number |                                        |
| srcChain  | string |                                        |
| srcToken  | string |                                        |
| sender    | string |                                        |
| inAmount  | string |                                        |
| dstChain  | string |                                        |
| dstToken  | string |                                        |
| receiver  | string |                                        |
| outAmount | string |                                        |
| action    | number | swap direction, 1: to evm, 2: from evm |                                                                                                                                                  |
| stage     | number |                                        |                                                                                                                                                  |
| status    | number |                                        |                                                                                                                                                  |
| createdAt | number |                                        |                                                                                                                                                  |

### Example

**request**:

```shell
curl --location 'http://127.0.0.1:8181/api/v1/order/list?sender=EQAQkdeIVxGcH67fRM6duM2o0KociERYIhMiZW5cZ0kgrM8R'
```

**response**

```json
{
  "code": 2000,
  "msg": "Success",
  "data": {
    "total": 1,
    "items": [
      {
        "orderId": 14,
        "srcChain": "1360104473493505",
        "srcToken": "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs",
        "sender": "EQAQkdeIVxGcH67fRM6duM2o0KociERYIhMiZW5cZ0kgrM8R",
        "inAmount": "5.916",
        "dstChain": "22776",
        "dstToken": "0x0000000000000000000000000000000000000000",
        "receiver": "0x2E9B4be739453cdDbB3641FB61052BA46873D41f",
        "outAmount": "",
        "action": 1,
        "stage": 2,
        "status": 4,
        "createdAt": 1724405186
      }
    ]
  }
}
```

## order detail by order id

### request path

/api/v1/order/detail-by-order-id

### request method

GET

### request params

| parameter | type   | required | default | description |
|-----------|--------|----------|---------|-------------|
| chainId   | string | Yes      |         |             |
| orderId   | number | Yes      |         |             |

### response params

| parameter | type   | description   |
|-----------|--------|---------------|
| code      | number | response code |
| msg       | string | response msg  |
| data      | object | response data |

#### data structure

| parameter | type   | description                                                   |
|-----------|--------|---------------------------------------------------------------|
| srcChain  | string |                                                               |
| srcToken  | string |                                                               |
| sender    | string |                                                               |
| inAmount  | string |                                                               |
| dstChain  | string |                                                               |
| dstToken  | string |                                                               |
| receiver  | string |                                                               |
| outAmount | string |                                                               |
| action    | number | swap direction, 1: to evm, 2: from evm                        |                                                                                                                                                  |
| stage     | number | 1: stage 1, 2: stage 2                                        |                                                                                                                                                  |
| status    | number | 1: tx prepare send, 2: tx sent, 3: tx failed, 4: tx confirmed |                                                                                                                                                  |
| createdAt | number |                                                               |                                                                                                                                                  |

### Example

**request**:

```shell
curl --location 'http://127.0.0.1:8181/api/v1/order/detail-by-order-id?chainId=1360104473493505&orderId=5044031782654955525'
```

**response**

```json
{
  "code": 2000,
  "msg": "Success",
  "data": {
    "srcChain": "1360104473493505",
    "srcToken": "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs",
    "sender": "EQAQkdeIVxGcH67fRM6duM2o0KociERYIhMiZW5cZ0kgrM8R",
    "inAmount": "6.000000",
    "dstChain": "137",
    "dstToken": "0xc2132d05d31c914a87c6611c10748aeb04b58e8f",
    "receiver": "0x2E9B4be739453cdDbB3641FB61052BA46873D41f",
    "outAmount": "",
    "action": 1,
    "stage": 2,
    "status": 4,
    "createdAt": 1726217636
  }
}
```
