## route

### request path

/api/v1/route

### request method

GET

### request params

| parameter       | type   | required | default | description                                                      |
|-----------------|--------|----------|---------|------------------------------------------------------------------|
| fromChainId     | string | Yes      |         |                                                                  |
| toChainId       | string | Yes      |         |                                                                  |
| amount          | string | Yes      |         |                                                                  |
| tokenInAddress  | string | Yes      |         |                                                                  |
| tokenOutAddress | string | Yes      |         |                                                                  |
| type            | string | Yes      |         |                                                                  |
| slippage        | string | Yes      |         | slippage of swap, a integer in rang [0, 5000], e.g, 100 means 1% |

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
curl --location '127.0.0.1:8181/api/v1/route?fromChainId=1&toChainId=22776&amount=1&tokenInAddress=0x0000000000000000000000000000000000000000&tokenOutAddress=0x0000000000000000000000000000000000000000&type=exactIn&slippage=100'
```

**response**

```json
{
  "code": 2000,
  "msg": "Success",
  "data": {
    "items": [
      {
        "hash": "",
        "tokenIn": {
          "chainId": "",
          "address": "",
          "name": "Bitcoin",
          "decimals": 8,
          "symbol": "BTC",
          "icon": ""
        },
        "tokenOut": {
          "chainId": "",
          "address": "",
          "name": "Ethereum",
          "decimals": 18,
          "symbol": "ETH",
          "icon": ""
        },
        "amountIn": "1.0",
        "amountOut": "11.3",
        "path": [
          {
            "name": "AbcBridge",
            "type": "cross",
            "tokenIn": {
              "chainId": "",
              "address": "",
              "name": "Bitcoin",
              "decimals": 8,
              "symbol": "BTC",
              "icon": ""
            },
            "tokenOut": {
              "chainId": "",
              "address": "",
              "name": "",
              "decimals": 18,
              "symbol": "WBTC",
              "icon": ""
            },
            "amountIn": "1",
            "amountOut": "1"
          },
          {
            "name": "Butter",
            "type": "swap",
            "tokenIn": {
              "chainId": "",
              "address": "",
              "name": "",
              "decimals": 18,
              "symbol": "WBTC",
              "icon": ""
            },
            "tokenOut": {
              "chainId": "",
              "address": "",
              "name": "Ethereum",
              "decimals": 18,
              "symbol": "ETH",
              "icon": ""
            },
            "amountIn": "1",
            "amountOut": "11.3"
          }
        ],
        "gasFee": {
          "amount": "0.000005",
          "symbol": "BTC"
        },
        "bridgeFee": {
          "amount": "0.000005",
          "symbol": "ETH"
        },
        "protocolFee": {
          "amount": "0",
          "symbol": "USDT"
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

| parameter | type   | required | default | description                                                      |
|-----------|--------|----------|---------|------------------------------------------------------------------|
| srcChain  | string | Yes      |         |                                                                  |
| srcToken  | string | Yes      |         |                                                                  |
| sender    | string | Yes      |         |                                                                  |
| amount    | string | Yes      |         |                                                                  |
| dstToken  | string | Yes      |         |                                                                  |
| receiver  | string | Yes      |         |                                                                  |
| hash      | string | Yes      |         | the route hash returned by /api/v1/swap                          |
| slippage  | string | Yes      |         | slippage of swap, a integer in rang [0, 5000], e.g, 100 means 1% |

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
  "data": {}
}
```

## create order

### request path

/api/v1/order/create

### request method

POST

### request params

| parameter    | type   | required | default | description                                                      |
|--------------|--------|----------|---------|------------------------------------------------------------------|
| srcChain     | string | Yes      |         |                                                                  |
| srcToken     | string | Yes      |         |                                                                  |
| sender       | string | Yes      |         |                                                                  |
| amount       | string | Yes      |         |                                                                  |
| dstToken     | string | Yes      |         |                                                                  |
| receiver     | string | Yes      |         |                                                                  |
| action       | number | Yes      |         | swap direction, 1: to evm, 2: from evm                           |
| **slippage** | string | Yes      |         | slippage of swap, a integer in rang [0, 5000], e.g, 100 means 1% |

### response params

| parameter | type   | description   |
|-----------|--------|---------------|
| code      | number | response code |
| msg       | string | response msg  |
| data      | object | response data |

#### data structure

| parameter | type   | description ---- |
|-----------|--------|------------------|
| orderId   | number |                  |
| relayer   | string | relayer address  |

### Example

**request**:

```shell
curl --location '127.0.0.1:8181/api/v1/order/create' \
--header 'Content-Type: application/json' \
--data '{
    "srcChain": "0x3948cddbbe5889e5de5d8d8f91a5cc6619909af4",
    "srcToken": "0x0000000000000000000000000000000000000000"
    ...
```

**response**

```json
{
  "code": 2000,
  "msg": "Success",
  "data": {
    "orderId": 14723,
    "relayer": "tb1ptuad7rdwycax553fp9fjg75ly2tv2065asl8uwcm3uxuve52tggqwclxst"
  }
}
```

## update order

### request path

/api/v1/order/update

### request method

POST

### request params

| parameter | type   | required | default | description |
|-----------|--------|----------|---------|-------------|
| orderId   | number | Yes      |         |             |
| inTxHash  | string | Yes      |         |             |

### response params

| parameter | type   | description   |
|-----------|--------|---------------|
| code      | number | response code |
| msg       | string | response msg  |
| data      | object | response data |

### Example

**request**:

```shell
curl --location '127.0.0.1:8181/api/v1/order/update' \
--header 'Content-Type: application/json' \
--data '{
    "orderId": 14723,
    "inTxHash": "764111ece8e33adcabf5dce7d1a57886d20ff44b06e29c5298542763d20d22cb"
    ...
```

**response**

```json
{
  "code": 2000,
  "msg": "Success",
  "data": {}
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
curl --location 'http://127.0.0.1:8181/api/v1/order/list?sender=tb1p862kth24h9gvz3vlt0g76uyxwaswqra4fr5njz0tjyq85u6dx6jszjc68l'
```

**response**

```json
{
  "code": 2000,
  "msg": "Success",
  "data": {
    "page": 1,
    "total": 60,
    "items": [
      {
        "orderId": 1
      }
    ]
  }
}
```

## order detail

### request path

/api/v1/order/detail

### request method

GET

### request params

| parameter | type   | required | default | description |
|-----------|--------|----------|---------|-------------|
| orderId   | number | Yes      |         |             |

### response params

| parameter | type   | description   |
|-----------|--------|---------------|
| code      | number | response code |
| msg       | string | response msg  |
| data      | object | response data |

#### data structure

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
curl --location 'http://127.0.0.1:8181/api/v1/order/list?sender=tb1p862kth24h9gvz3vlt0g76uyxwaswqra4fr5njz0tjyq85u6dx6jszjc68l'
```

**response**

```json
{
  "code": 2000,
  "msg": "Success",
  "data": {
    "orderId": 1
  }
}
```

## 接口调用顺序

### btc to evm:

1. /api/v1/route
2. /api/v1/order/create

### evm to btc:

1. /api/v1/route
2. /api/v1/swap 该接口会创建一个 order 并构建 swap 交易的 data
3. /api/v1/update