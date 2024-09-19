# fe-backend

| 交易方向       | protocol fee 扣除阶段 | protocol fee symbol | protocol fee 规则                                      | bridge fee 扣除阶段 | bridge fee symbol | bridge fee 规则                         |
|------------|-------------------|---------------------|------------------------------------------------------|-----------------|-------------------|---------------------------------------|
| ton to evm | 第一阶段              | USDT                | 根据代扣参数从 ton router 获取 protocol fee， 并返回 protocol fee | 第二阶段            | USDT              | 根据规则计算(比例 + base tx fee)并返回 bridge fe |
| btc to evm | 第一阶段              | WBTC                | 根据代扣参数闪兑服务计算并返回 protocol fee                         | 第二阶段            | WBTC              | 根据规则计算(比例 + base tx fee)并返回 bridge fe |
| from evm   | 第一阶段              | src token symbol    | 根据代扣参数闪兑服务计算并返回 protocol fee                         | 第二阶段            | USDT              | 根据规则计算(比例 + base tx fee)并返回 bridge fe |



