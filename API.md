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
curl --location '127.0.0.1:8181/api/v1/route?fromChainId=1&toChainId=22776&amount=1&tokenInAddress=0x0000000000000000000000000000000000000000&tokenOutAddress=0x0000000000000000000000000000000000000000&type=exactIn&slippage=100'
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
        "hash": "0xa7d73ffac909efb9c97b57b71faad84178f8eae439911009e5bb6fd77059c68b",
        "tokenIn": {
          "chainId": "22776",
          "address": "0x0000000000000000000000000000000000000000",
          "name": "MAPO",
          "decimals": 18,
          "symbol": "MAPO",
          "icon": "https://cdn.befiwalletdao.com/image/icon_local_map_checked_3gfyyv.png"
        },
        "tokenOut": {
          "chainId": "1360104473493505",
          "address": "0x0000000000000000000000000000000000000000",
          "name": "Toncoin",
          "decimals": 9,
          "symbol": "TON",
          "icon": "https://assets.dedust.io/images/usdt.webp"
        },
        "amountIn": "300.0",
        "amountOut": "0.426384366",
        "path": [
          {
            "name": "Butter",
            "amountIn": "300.0",
            "amountOut": "2.950721",
            "tokenIn": {
              "chainId": "22776",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "MAPO",
              "decimals": 18,
              "symbol": "MAPO",
              "icon": "https://cdn.befiwalletdao.com/image/icon_local_map_checked_3gfyyv.png"
            },
            "tokenOut": {
              "chainId": "137",
              "address": "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
              "name": "Tether USD",
              "decimals": 6,
              "symbol": "USDT",
              "icon": "https://files.mapprotocol.io/bridge/usdt.png"
            }
          },
          {
            "name": "MAP FE",
            "amountIn": "2.950721",
            "amountOut": "2.950721",
            "tokenIn": {
              "chainId": "137",
              "address": "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
              "name": "Tether USD",
              "decimals": 6,
              "symbol": "USDT",
              "icon": "https://files.mapprotocol.io/bridge/usdt.png"
            },
            "tokenOut": {
              "chainId": "1360104473493505",
              "address": "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs",
              "name": "Tether USD",
              "decimals": 6,
              "symbol": "USDT",
              "icon": "https://assets.dedust.io/images/usdt.webp"
            }
          },
          {
            "name": "DeDust",
            "amountIn": "2.950721",
            "amountOut": "0.426384366",
            "tokenIn": {
              "chainId": "1360104473493505",
              "address": "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs",
              "name": "Tether USD",
              "decimals": 6,
              "symbol": "USDT",
              "icon": "https://assets.dedust.io/images/usdt.webp"
            },
            "tokenOut": {
              "chainId": "1360104473493505",
              "address": "0x0000000000000000000000000000000000000000",
              "name": "Toncoin",
              "decimals": 9,
              "symbol": "TON",
              "icon": "https://assets.dedust.io/images/usdt.webp"
            }
          }
        ],
        "gasFee": {
          "amount": "0.472988216091",
          "symbol": "MAPO",
          "chainId": ""
        },
        "bridgeFee": {
          "amount": "0.6",
          "symbol": "USDT",
          "chainId": ""
        },
        "protocolFee": {
          "amount": "0",
          "symbol": "USDT",
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

| parameter | type   | required | default | description                                                        |
|-----------|--------|----------|---------|--------------------------------------------------------------------|
| srcChain  | string | Yes      |         |                                                                    |
| srcToken  | string | Yes      |         |                                                                    |
| sender    | string | Yes      |         |                                                                    |
| amount    | string | Yes      |         |                                                                    |
| decimal   | number | Yes      |         |                                                                    |
| dstChain  | string | Yes      |         |                                                                    |
| dstToken  | string | Yes      |         |                                                                    |
| receiver  | string | Yes      |         |                                                                    |
| hash      | string | Yes      |         | the route hash returned by /api/v1/route                           |
| slippage  | string | Yes      |         | slippage of swap, a integer in rang [300, 5000], e.g, 300 means 3% |

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
| dstChain     | string | Yes      |         |                                                                  |
| dstToken     | string | Yes      |         |                                                                  |
| receiver     | string | Yes      |         |                                                                  |
| action       | number | Yes      |         | swap direction, 1: to evm, 2: from evm                           |
| hash         | string | Yes      |         | the route hash returned by /api/v1/route                         |                
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

### ton to evm:

1. /api/v1/route
2. /api/v1/swap --> tx data
3. send tx to ton

### evm to ton

1. /api/v1/route
2. /api/v1/swap --> tx data
3. send tx to evm