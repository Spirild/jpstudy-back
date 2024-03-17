package view

import (
	// "encoding/json"
	"net/http"
	jpdetailservice "translasan-lite/controller/JpDetailService"
	pbdata "translasan-lite/proto/generated"
	"translasan-lite/utils"

	"github.com/gin-gonic/gin"
)

func (hs *HttpServer) GetJpDetailTable(c *gin.Context) {
	JpDetailController := jpdetailservice.GetJpDetailServiceInstant()
	jpReq := &pbdata.JpWordReq{}
	err := hs.ReadProtoReq(c.Request, jpReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"errorMsg":  err.Error(),
		})
		return
	}

	rsp, err := JpDetailController.GetJpDetailTable(jpReq)
	if err != nil {
		c.JSON(int(pbdata.ErrorCode_UNKNOWN), gin.H{
			"errorCode": int(pbdata.ErrorCode_UNKNOWN),
			"errorMsg":  err.Error(),
		})
		return
	}
	// 这里的下划线找机会确认什么原因
	t, _ := utils.SelfMarshal(*rsp)
	c.Data(int(pbdata.ErrorCode_SUCCESS), "application/json", t)
}
