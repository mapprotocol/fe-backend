package logic

import (
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/entity"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/resp"
)

func SupportedChains(page, size int) (ret []*entity.SupportedChainsResponse, count int64, code int) {
	list, count, err := dao.NewSupportedChain().Find(nil, dao.Paginate(page, size))
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
		return []*entity.SupportedChainsResponse{}, count, resp.CodeSuccess
	}

	ret = make([]*entity.SupportedChainsResponse, 0, length)
	for _, c := range list {
		ret = append(ret, &entity.SupportedChainsResponse{
			ChainID:   c.ChainID,
			ChainName: c.ChainName,
			ChainIcon: c.ChainIcon,
		})
	}
	return ret, count, resp.CodeSuccess
}
