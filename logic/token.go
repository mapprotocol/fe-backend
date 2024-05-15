package logic

import (
	"github.com/mapprotocol/ceffu-fe-backend/dao"
	"github.com/mapprotocol/ceffu-fe-backend/entity"
	"github.com/mapprotocol/ceffu-fe-backend/resource/log"
	"github.com/mapprotocol/ceffu-fe-backend/resp"
)

func SupportedTokens(chainID uint64, symbol string, page, size int) (ret []*entity.SupportedTokensResponse, count int64, code int) {
	list, count, err := dao.NewSupportedToken(chainID, symbol).Find(nil, dao.Paginate(page, size))
	if err != nil {
		fields := map[string]interface{}{
			"page":  page,
			"size":  size,
			"error": err,
		}
		log.Logger().WithFields(fields).Error("failed to get supported chain list")
		return nil, 0, resp.CodeInternalServerError
	}

	length := len(list)
	if length == 0 {
		return []*entity.SupportedTokensResponse{}, count, resp.CodeSuccess
	}

	ret = make([]*entity.SupportedTokensResponse, 0, length)
	for _, c := range list {
		ret = append(ret, &entity.SupportedTokensResponse{
			ChainID:  c.ChainID,
			Symbol:   c.Symbol,
			Name:     c.Name,
			Decimals: c.Decimal,
		})
	}
	return ret, count, resp.CodeSuccess
}

func DepositAddress(chainID uint64, tokenSymbol string) (ret []*entity.DepositAddressResponse, code int) {

	return ret, resp.CodeSuccess
}
