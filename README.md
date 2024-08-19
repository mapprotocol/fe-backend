# fe-backend

## handle process

### ton to evm

| step | stage | status | desc                                                                                                                   |
|------|-------|--------|------------------------------------------------------------------------------------------------------------------------|
| `1`  | 1     | 2      | filter ton to evm event(34a7e0e8) and create order                                                                     |
| `2`  | 2     | 1      | send tx(DeliverAndSwap) to chain pool on evm and updated tx hash to order                                              |
| `3`  | 2     | 2/3    | tx(DeliverAndSwap) confirmed on evm, if tx status is confirmed set status to 2, if tx status is failed set status to 3 |
| `4`  | 2     | 4      |                                                                                                                        |

## evm to ton

| step | stage | status | desc                                                            |
|------|-------|--------|-----------------------------------------------------------------|
| `1`  | 1     | 2      | filter evm to ton event(OnReceived) and create order            |
| `2`  | 2     | 1      | send tx to chain pool on ton ----- and updated tx hash to order |
| `3`  | 2     | 2      | filter evm to ton event(1a6c0a51) and create order              |
| `4`  | 2     | 4      |                                                                 |

## new process

### ton to evm

| step | stage | status | desc                                                                                                                   |
|------|-------|--------|------------------------------------------------------------------------------------------------------------------------|
| `1`  | 1     | 2      | filter ton to evm event(34a7e0e8) and create order                                                                     |
| `2`  | 2     | 1      | prepare send tx(DeliverAndSwap)                                                                                        |
| `3`  | 2     | 2      | send tx(DeliverAndSwap) to chain pool on evm and updated tx hash to order                                              |
| `4`  | 2     | 3/4    | tx(DeliverAndSwap) confirmed on evm, if tx status is confirmed set status to 2, if tx status is failed set status to 3 |
| `5`  | 2     | 5      |                                                                                                                        |

## evm to ton

| step | stage | status | desc                                                 |
|------|-------|--------|------------------------------------------------------|
| `1`  | 1     | 2      | filter evm to ton event(OnReceived) and create order |
| `2`  | 2     | 1      | prepare send tx to chain pool on ton                 |
| `3`  | 2     | 2      | send tx to chain pool on ton                         |
| `4`  | 2     | 3/4    | filter evm to ton event(1a6c0a51) and create order   |
| `5`  | 2     | 5      |                                                      |