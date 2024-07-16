package filter

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"strconv"

	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
)

const pathGetLogs = "/v1/mos/list"

const successCode = 200

const projectID = 1

type GetLogsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Total int                    `json:"total"`
		List  []*GetLogsResponseItem `json:"list"`
	} `json:"data"`
}

type GetLogsResponseItem struct {
	Id              uint64 `json:"id"`
	ProjectId       int    `json:"project_id"`
	ChainId         int    `json:"chain_id"`
	EventId         int    `json:"event_id"`
	TxHash          string `json:"tx_hash"`
	ContractAddress string `json:"contract_address"`
	Topic           string `json:"topic"`
	BlockNumber     int    `json:"block_number"`
	BlockHash       string `json:"block_hash"`
	LogIndex        int    `json:"log_index"`
	LogData         string `json:"log_data"`
	TxIndex         int    `json:"tx_index"`
	TxTimestamp     int    `json:"tx_timestamp"`
}

func GetLogs(id uint64, chainID, topic string, limit uint8) ([]*GetLogsResponseItem, error) {
	params := fmt.Sprintf("id=%d&project_id=%d&chain_id=%s&topic=%s&limit=%d", id, projectID, chainID, topic, limit)
	//url := fmt.Sprintf("%s%s?%s", viper.GetStringMap("endpoints")["filter"], pathGetLogs, params)
	url := fmt.Sprintf("%s%s?%s", viper.GetStringMap("endpoints")["filter"], pathGetLogs, params)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := GetLogsResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		return nil, err
	}
	if response.Code != successCode {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(strconv.Itoa(response.Code)),
			reqerror.WithMessage(response.Message),
		)
	}
	return response.Data.List, nil
}
