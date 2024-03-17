package jpdetailservice

// 详情页, 具体功能中午再想想

import (
	"translasan-lite/common"
	"translasan-lite/core"

	// dbdata "translasan-lite/db/generated"
	dbservice "translasan-lite/db"
	pbdata "translasan-lite/proto/generated"

	"go.uber.org/zap"
)

type JpDetailService struct {
	core.BaseComponent
	// 但是和一开始注册不同，这里是用了才注册。方便追踪运行链路（共享一条链路）
}

var jpDetailServiceInstant *JpDetailService

func GetJpDetailServiceInstant() *JpDetailService {
	if jpDetailServiceInstant == nil {
		jpDetailServiceInstant = &JpDetailService{}
		jpDetailServiceInstant.init()
	}
	return jpDetailServiceInstant
}

func (js *JpDetailService) init() {
	n := core.GetDefaultNode()
	(&js.BaseComponent).Init(n, &core.ServiceConfig{})
}

func (js *JpDetailService) GetJpDetailTable(req *pbdata.JpWordReq) (*pbdata.JpWordRsp, error) {
	pageRsp := &pbdata.PageRsp{
		Page: req.Common.Page,
		Size: req.Common.Size,
	}
	rsp := &pbdata.JpWordRsp{
		Common: pageRsp,
	}

	ds, err := js.getDatabaseServiceClient()
	if err != nil {
		js.Log.Error("error", zap.Error(err))
		rsp.ErrorCode = int32(pbdata.ErrorCode_DATABASE)
		rsp.ErrorMessage = "database problem"
		return rsp, nil
	}
	start := (req.Common.Page-1)*req.Common.Size + 1
	end := req.Common.Page * req.Common.Size
	res, total := ds.SelectJpTableDetail(int(start), int(end))
	rsp.Common.Total = int32(total)
	for _, r := range res {
		rsp.JpWordList = append(rsp.JpWordList, &pbdata.JpWord{
			WordId:  r.Id,
			Word:    r.WordName,
			Example: r.WordDesc,
		})
	}
	rsp.ErrorCode = int32(pbdata.ErrorCode_SUCCESS)
	rsp.ErrorMessage = "SUCCESS"

	return rsp, nil
}

func (js *JpDetailService) getDatabaseServiceClient() (dbservice.IDatabaseService, error) {
	svc, ok := js.FindService(common.ServiceIdDatabase)
	if !ok {
		return nil, common.ErrorInstance.ErrNoDatabaseService
	}
	ds, ok := svc.(dbservice.IDatabaseService)
	if !ok {
		return nil, common.ErrorInstance.ErrInvalidDatabaseService
	}
	return ds, nil
}
