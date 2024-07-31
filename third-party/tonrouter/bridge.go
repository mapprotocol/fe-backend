package tonrouter

//import (
//	"encoding/json"
//	"fmt"
//	uhttp "github.com/mapprotocol/fe-backend/utils/http"
//	"github.com/mapprotocol/fe-backend/utils/reqerror"
//	"strconv"
//)
//
//type BridgeRequest struct {
//	TokenInAddress  string `json:"tokenInAddress"`
//	TokenOutAddress string `json:"tokenOutAddress"`
//	Sender          string `json:"sender"`
//	Receiver        string `json:"receiver"`
//	Amount          string `json:"amount"`
//	ToChainId       string `json:"toChainId"`
//	Slippage        uint64 `json:"slippage"`
//}
//
//type BridgeResponse struct {
//	Errno   int    `json:"errno"`
//	Message string `json:"message"`
//	Data    struct {
//		TxParams *TxParams `json:"txParams"`
//	} `json:"data"`
//}
//
//func Bridge(req *BridgeRequest) (*TxParams, error) {
//	params := fmt.Sprintf(
//		"tokenInAddress=%s&tokenOutAddress=%s&sender=%s&receiver=%s&amount=%s&toChainId=%s&slippage=%d",
//		req.TokenInAddress, req.TokenOutAddress, req.Sender, req.Receiver, req.Amount, req.ToChainId, req.Slippage)
//	url := fmt.Sprintf("%s%s?%s", endpoint, PathBridge, params)
//	fmt.Println("============================== route url: ", url)
//	ret, err := uhttp.Get(url, nil, nil)
//	if err != nil {
//		return nil, reqerror.NewExternalRequestError(
//			url,
//			reqerror.WithError(err),
//		)
//	}
//	response := BridgeResponse{}
//	if err := json.Unmarshal(ret, &response); err != nil {
//		return nil, err
//	}
//	if response.Errno != SuccessCode {
//		return nil, reqerror.NewExternalRequestError(
//			url,
//			reqerror.WithCode(strconv.Itoa(response.Errno)),
//			reqerror.WithMessage(response.Message),
//		)
//	}
//	return response.Data.TxParams, nil
//}
