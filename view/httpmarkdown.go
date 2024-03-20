package view

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	pbdata "translasan-lite/proto/generated"

	"github.com/gin-gonic/gin"
)

var (
	markdown = path.Join(".", "src", "txt", "markdown.md")
)

func (hs *HttpServer) GetMarkdownContent(c *gin.Context) {
	f, _ := os.Open(markdown)
	defer f.Close()
	reader := io.Reader(f)
	contents, _ := io.ReadAll(reader)
	rsp := &pbdata.GetMarkdownRsp{
		ErrorCode:       int32(pbdata.ErrorCode_SUCCESS),
		ErrorMessage:    "SUCCESS",
		MarkdownContent: string(contents),
	}
	t, _ := json.Marshal(rsp)
	c.Data(int(pbdata.ErrorCode_SUCCESS), "application/json", t)
}

func (hs *HttpServer) SaveMarkdownContent(c *gin.Context) {
	var data []byte
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorCode": http.StatusBadRequest,
			"errorMsg":  err.Error(),
		})
		return
	}

	f, _ := os.OpenFile(markdown, os.O_WRONLY|os.O_TRUNC, 0666)
	_, err = f.Write(data)
	if err != nil {
		hs.Log.Error(err.Error())
	}
	defer f.Close()

	rsp := &pbdata.CommonRsp{
		ErrorCode:    int32(pbdata.ErrorCode_SUCCESS),
		ErrorMessage: "SUCCESS",
	}
	t, _ := json.Marshal(rsp)
	c.Data(int(pbdata.ErrorCode_SUCCESS), "application/json", t)
}

func (hs *HttpServer) AskBotDemo(c *gin.Context) {
	ts, err := hs.getThirdServiceClient()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": pbdata.ErrorCode_THIRD,
			"error_msg":  err.Error(),
		})
		return
	}
	botReq := &pbdata.AskBotReq{}
	err = hs.ReadProtoReq(c.Request, botReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": http.StatusBadRequest,
			"error_msg":  err.Error(),
		})
		return
	}

	answer, err := ts.SparkDemo(botReq.Content)
	// answer := "你是一个一个……"
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": pbdata.ErrorCode_THIRD,
			"error_msg":  err.Error(),
		})
		return
	}

	rsp := &pbdata.AskBotRsp{
		ErrorCode:    int32(pbdata.ErrorCode_SUCCESS),
		ErrorMessage: "SUCCESS",
		Content:      answer,
	}
	t, _ := json.Marshal(rsp)
	c.Data(int(pbdata.ErrorCode_SUCCESS), "application/json", t)
}
