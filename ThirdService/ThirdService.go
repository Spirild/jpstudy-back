package thirdservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"translasan-lite/common"
	"translasan-lite/core"
	dbservice "translasan-lite/db"
	pbdata "translasan-lite/proto/generated"
)

type ThirdService struct {
	core.BaseComponent
	url     string
	headers map[string]string
	client  *http.Client

	sparkUrl  string
	appid     string
	apiSecret string
	apiKey    string
}

func (ts *ThirdService) ServiceID() int {
	return common.ServiceIdThird
}

func (ts *ThirdService) Init(n *core.Node, cfg *core.ServiceConfig) {
	(&ts.BaseComponent).Init(n, cfg)
	ts.url = "https://api.mojidict.com/parse/functions/union-api"
	ts.headers = map[string]string{
		"accept":             "*/*",
		"accept-encoding":    "deflate, br",
		"accept-language":    "zh-CN,zh;q=0.9",
		"content-length":     "407",
		"content-type":       "text/plain",
		"origin":             "https://www.mojidict.com",
		"referer":            "https://www.mojidict.com/",
		"sec-ch-ua":          "\"Google Chrome\";v=\"113\",\"Chromium\";v=\"113\",\"Not-A.Brand\";v=\"24\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko), Chrome/111.0.0.0 Safari/537.36",
	}
	ts.client = &http.Client{}

	ts.sparkUrl, _ = cfg.GetString("hostUrl")
	ts.appid, _ = cfg.GetString("appid")
	ts.apiSecret, _ = cfg.GetString("apiSecret")
	ts.apiKey, _ = cfg.GetString("apiKey")
}

func (ts *ThirdService) Run(ctx context.Context) error {

	<-ctx.Done()
	ts.Log.Info("ThirdService stops running")

	return nil
}

func (ts *ThirdService) MojiTranlate(searchContent string) ([]*pbdata.MojiResponseWord, error) {
	var res []*pbdata.MojiResponseWord

	db, err := ts.getDatabaseServiceClient()
	if err != nil {
		ts.Log.Error(err.Error())
		return nil, err
	}
	requestBodyMold := db.GetMojiTokenMold()

	requestBody := fmt.Sprintf(requestBodyMold.Content, searchContent, searchContent)

	req, err := http.NewRequest("POST", ts.url, bytes.NewBufferString(requestBody))
	if err != nil {
		ts.Log.Error(err.Error())
		return nil, err
	}

	for k, v := range ts.headers {
		req.Header.Set(k, v)
	}

	rsp, err := ts.client.Do(req)
	if err != nil {
		ts.Log.Error(err.Error())
		return nil, err
	}
	defer rsp.Body.Close()

	content, err := io.ReadAll(rsp.Body)
	if err != nil {
		ts.Log.Error(err.Error())
		return nil, err
	}
	res = analyzeMojiRsp(content)

	return res, nil
}

func analyzeMojiRsp(jsondata []byte) []*pbdata.MojiResponseWord {
	var result map[string]interface{}

	json.Unmarshal(jsondata, &result)
	tmp := result["result"].(map[string]interface{})
	tmp = tmp["results"].(map[string]interface{})
	tmp = tmp["search-all"].(map[string]interface{})
	tmp = tmp["result"].(map[string]interface{})
	tmp = tmp["word"].(map[string]interface{})
	checkpoint := tmp["searchResult"].([]interface{})
	var res []*pbdata.MojiResponseWord
	for _, c := range checkpoint {
		res = append(res, &pbdata.MojiResponseWord{
			Title:   c.(map[string]interface{})["title"].(string),
			Excerpt: c.(map[string]interface{})["excerpt"].(string),
		})
	}
	return res
}

func (ts *ThirdService) getDatabaseServiceClient() (dbservice.IDatabaseService, error) {
	svc, ok := ts.FindService(common.ServiceIdDatabase)
	if !ok {
		return nil, common.ErrorInstance.ErrNoDatabaseService
	}
	ds, ok := svc.(dbservice.IDatabaseService)
	if !ok {
		return nil, common.ErrorInstance.ErrInvalidDatabaseService
	}
	return ds, nil
}
