package view

import (
	"encoding/json"
	"net/http"
	jpliteservice "translasan-lite/controller/JpLiteService"
	pbdata "translasan-lite/proto/generated"
	"translasan-lite/utils"

	"github.com/gin-gonic/gin"
)

func (hs *HttpServer) GetJpLiteTable(c *gin.Context) {
	JpLiteController := jpliteservice.GetJpLiteServiceInstant()
	jpLiteReq := &pbdata.JpWordReq{}
	err := hs.ReadProtoReq(c.Request, jpLiteReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"errorMsg":  err.Error(),
		})
		return
	}

	rsp, err := JpLiteController.GetJpLiteTable(jpLiteReq)
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

func (hs *HttpServer) RememberJpWord(c *gin.Context) {
	JpLiteController := jpliteservice.GetJpLiteServiceInstant()
	rememberReq := &pbdata.RememberJpWordReq{}
	err := hs.ReadProtoReq(c.Request, rememberReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"errorMsg":  err.Error(),
		})
		return
	}

	rsp, err := JpLiteController.RememberJpWord(rememberReq)
	if err != nil {
		c.JSON(int(pbdata.ErrorCode_UNKNOWN), gin.H{
			"errorCode": int(pbdata.ErrorCode_UNKNOWN),
			"errorMsg":  err.Error(),
		})
		return
	}
	t, _ := json.Marshal(rsp)
	c.Data(int(pbdata.ErrorCode_SUCCESS), "application/json", t)
}

func (hs *HttpServer) ForgetJpWord(c *gin.Context) {
	JpLiteController := jpliteservice.GetJpLiteServiceInstant()
	forgetReq := &pbdata.ForgetJpWordReq{}
	err := hs.ReadProtoReq(c.Request, forgetReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"errorMsg":  err.Error(),
		})
		return
	}

	rsp, err := JpLiteController.ForgetJpWord(forgetReq)
	if err != nil {
		c.JSON(int(pbdata.ErrorCode_UNKNOWN), gin.H{
			"errorCode": int(pbdata.ErrorCode_UNKNOWN),
			"errorMsg":  err.Error(),
		})
		return
	}
	t, _ := json.Marshal(rsp)
	c.Data(int(pbdata.ErrorCode_SUCCESS), "application/json", t)
}

func (hs *HttpServer) SaveJpWord(c *gin.Context) {
	// 增改一体
	JpLiteController := jpliteservice.GetJpLiteServiceInstant()
	SaveReq := &pbdata.SaveJpWordReq{}
	err := hs.ReadProtoReq(c.Request, SaveReq)
	if err != nil {
		hs.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"errorMsg":  err.Error(),
		})
		return
	}

	rsp, err := JpLiteController.SaveJpWord(SaveReq)
	if err != nil {
		c.JSON(int(pbdata.ErrorCode_UNKNOWN), gin.H{
			"errorCode": int(pbdata.ErrorCode_UNKNOWN),
			"errorMsg":  err.Error(),
		})
		return
	}
	t, _ := json.Marshal(rsp)
	c.Data(int(pbdata.ErrorCode_SUCCESS), "application/json", t)
}

func (hs *HttpServer) DeleteJpWord(c *gin.Context) {
	// 增改一体
	JpLiteController := jpliteservice.GetJpLiteServiceInstant()
	DelReq := &pbdata.DeleteJpWordReq{}
	err := hs.ReadProtoReq(c.Request, DelReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"errorMsg":  err.Error(),
		})
		return
	}

	rsp, err := JpLiteController.DeleteJpWord(DelReq)
	if err != nil {
		c.JSON(int(pbdata.ErrorCode_UNKNOWN), gin.H{
			"errorCode": int(pbdata.ErrorCode_UNKNOWN),
			"errorMsg":  err.Error(),
		})
		return
	}
	t, _ := json.Marshal(rsp)
	c.Data(int(pbdata.ErrorCode_SUCCESS), "application/json", t)
}

func (hs *HttpServer) TranslateJpWord(c *gin.Context) {
	// 增改一体
	JpLiteController := jpliteservice.GetJpLiteServiceInstant()
	translateReq := &pbdata.TranslateJpReq{}
	err := hs.ReadProtoReq(c.Request, translateReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"errorMsg":  err.Error(),
		})
		return
	}

	rsp, err := JpLiteController.TranslateJpWord(translateReq)
	if err != nil {
		c.JSON(int(pbdata.ErrorCode_UNKNOWN), gin.H{
			"errorCode": int(pbdata.ErrorCode_UNKNOWN),
			"errorMsg":  err.Error(),
		})
		return
	}
	t, _ := utils.SelfMarshal(*rsp)
	c.Data(int(pbdata.ErrorCode_SUCCESS), "application/json", t)
}
