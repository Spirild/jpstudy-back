package jpliteservice

import (
	"math/rand"
	"strings"
	"time"

	"translasan-lite/common"
	"translasan-lite/core"
	dbdata "translasan-lite/db/generated"
	"translasan-lite/utils"

	thirdservice "translasan-lite/ThirdService"
	dbservice "translasan-lite/db"
	pbdata "translasan-lite/proto/generated"

	"go.uber.org/zap"
)

type JpLiteService struct {
	core.BaseComponent
	// 但是和一开始注册不同，这里是用了才注册。方便追踪运行链路（共享一条链路）
}

var jpLiteServiceInstant *JpLiteService

func GetJpLiteServiceInstant() *JpLiteService {
	if jpLiteServiceInstant == nil {
		jpLiteServiceInstant = &JpLiteService{}
		jpLiteServiceInstant.init()
	}
	return jpLiteServiceInstant
}

func (js *JpLiteService) init() {
	n := core.GetDefaultNode()
	(&js.BaseComponent).Init(n, &core.ServiceConfig{})
}

func (js *JpLiteService) GetJpLiteTable(req *pbdata.JpWordReq) (*pbdata.JpWordRsp, error) {
	// 获取轻词表
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

	res := ds.SelectJpTableLite(int(req.Level), req.Word, req.Common.User)
	rsp.Common.Total = int32(len(res))
	res = utils.Paginate(res, int(req.Common.Page), int(req.Common.Size))
	for _, r := range res {
		rsp.JpWordList = append(rsp.JpWordList, &pbdata.JpWord{
			WordId:      r.Id,
			Word:        r.Vocabulary,
			Spell:       r.Character,
			Translation: r.Translator,
			Example:     r.SampleSentence,
		})
	}
	rsp.ErrorCode = int32(pbdata.ErrorCode_SUCCESS)
	rsp.ErrorMessage = "SUCCESS"

	return rsp, nil
}

func (js *JpLiteService) RememberJpWord(req *pbdata.RememberJpWordReq) (*pbdata.CommonRsp, error) {
	// 这次记住了该单词
	rsp := &pbdata.CommonRsp{}
	ds, err := js.getDatabaseServiceClient()
	if err != nil {
		js.Log.Error("error", zap.Error(err))
		rsp.ErrorCode = int32(pbdata.ErrorCode_DATABASE)
		rsp.ErrorMessage = "database problem"
		return rsp, nil
	}
	err = ds.JpWordLevelUpdate(int(req.WordId), 1)
	if err != nil {
		js.Log.Error("error", zap.Error(err))
		rsp.ErrorCode = int32(pbdata.ErrorCode_DATABASE)
		rsp.ErrorMessage = "database problem"
		return rsp, nil
	}
	rsp.ErrorCode = int32(pbdata.ErrorCode_SUCCESS)
	rsp.ErrorMessage = "SUCCESS"

	return rsp, nil
}

func (js *JpLiteService) ForgetJpWord(req *pbdata.ForgetJpWordReq) (*pbdata.CommonRsp, error) {
	// 这次忘记了该单词
	rsp := &pbdata.CommonRsp{}
	ds, err := js.getDatabaseServiceClient()
	if err != nil {
		js.Log.Error("error", zap.Error(err))
		rsp.ErrorCode = int32(pbdata.ErrorCode_DATABASE)
		rsp.ErrorMessage = "database problem"
		return rsp, nil
	}
	err = ds.JpWordLevelUpdate(int(req.WordId), -1)
	if err != nil {
		js.Log.Error("error", zap.Error(err))
		rsp.ErrorCode = int32(pbdata.ErrorCode_DATABASE)
		rsp.ErrorMessage = "database problem"
		return rsp, nil
	}
	rsp.ErrorCode = int32(pbdata.ErrorCode_SUCCESS)
	rsp.ErrorMessage = "SUCCESS"

	return rsp, nil
}

func (js *JpLiteService) SaveJpWord(req *pbdata.SaveJpWordReq) (*pbdata.CommonRsp, error) {
	// 增改一体
	rsp := &pbdata.CommonRsp{}
	ds, err := js.getDatabaseServiceClient()
	if err != nil {
		js.Log.Error("error", zap.Error(err))
		rsp.ErrorCode = int32(pbdata.ErrorCode_DATABASE)
		rsp.ErrorMessage = "database problem"
		return rsp, nil
	}
	word := &dbdata.CharacterVocabularyTable{
		Vocabulary:     req.JpWord.Word,
		Character:      req.JpWord.Spell,
		Translator:     req.JpWord.Translation,
		SampleSentence: req.JpWord.Example,
		MemoryLevel:    1,
		UpdateTime:     utils.GetCurrentTimeStr(),
		UserQq:         req.User,
	}
	if req.JpWord.WordId == 0 {
		// 说明是新建的
		err = ds.SelfInsert(word)
	} else {
		word.Id = req.JpWord.WordId
		err = ds.SelfUpdate(word)
	}
	if err != nil {
		js.Log.Error("error", zap.Error(err))
		rsp.ErrorCode = int32(pbdata.ErrorCode_DATABASE)
		rsp.ErrorMessage = "database problem"
		return rsp, nil
	}
	rsp.ErrorCode = int32(pbdata.ErrorCode_SUCCESS)
	rsp.ErrorMessage = "SUCCESS"
	return rsp, nil
}

func (js *JpLiteService) DeleteJpWord(req *pbdata.DeleteJpWordReq) (*pbdata.CommonRsp, error) {
	// 删除单词
	rsp := &pbdata.CommonRsp{}
	ds, err := js.getDatabaseServiceClient()
	if err != nil {
		js.Log.Error("error", zap.Error(err))
		rsp.ErrorCode = int32(pbdata.ErrorCode_DATABASE)
		rsp.ErrorMessage = "database problem"
		return rsp, nil
	}
	err = ds.JpWordDelete(int(req.WordId))
	if err != nil {
		js.Log.Error("error", zap.Error(err))
		rsp.ErrorCode = int32(pbdata.ErrorCode_DATABASE)
		rsp.ErrorMessage = "database problem"
		return rsp, nil
	}
	rsp.ErrorCode = int32(pbdata.ErrorCode_SUCCESS)
	rsp.ErrorMessage = "SUCCESS"

	return rsp, nil
}

func (js *JpLiteService) TranslateJpWord(req *pbdata.TranslateJpReq) (*pbdata.TranslateJpRsp, error) {
	// 翻译单词
	rsp := &pbdata.TranslateJpRsp{}
	ts, err := js.getThirdServiceClient()
	if err != nil {
		js.Log.Error("error", zap.Error(err))
		rsp.ErrorCode = int32(pbdata.ErrorCode_THIRD)
		rsp.ErrorMessage = "third service problem"
		return rsp, nil
	}
	transRes, err := ts.MojiTranlate(req.Word)
	if err != nil {
		js.Log.Error("error", zap.Error(err))
		rsp.ErrorCode = int32(pbdata.ErrorCode_THIRD)
		rsp.ErrorMessage = "third server problem"
		return rsp, nil
	}

	ds, err := js.getDatabaseServiceClient()
	if err != nil {
		js.Log.Error("error", zap.Error(err))
		rsp.ErrorCode = int32(pbdata.ErrorCode_DATABASE)
		rsp.ErrorMessage = "database problem"
		return rsp, nil
	}
	for _, res := range transRes {
		titleSplit := strings.Split(res.Title, "|")
		word := titleSplit[0]
		spell := ""
		if len(titleSplit) > 1 {
			spell = titleSplit[1]
		}
		dbWord := &dbdata.CharacterVocabularyTable{
			Vocabulary:     word,
			Character:      spell,
			SampleSentence: res.Excerpt,
			UpdateTime:     utils.GetCurrentTimeStr(),
			MemoryLevel:    1,
			UserQq:         req.User,
		}
		check := ds.SelectJpTableLite(0, word, req.User)
		if len(check) == 0 {
			ds.SelfInsert(dbWord)
		} else {
			dbWord.IsDel = 0
			ds.SelfUpdate(dbWord)
		}

		rsp.WordList = append(rsp.WordList, &pbdata.JpWord{
			Word:        word,
			Spell:       spell,
			Translation: res.Excerpt,
		})
	}

	rsp.ErrorCode = int32(pbdata.ErrorCode_SUCCESS)
	rsp.ErrorMessage = "SUCCESS"
	return rsp, nil
}

func (js *JpLiteService) GetJpCardList(req *pbdata.CommonReq) (*pbdata.JpWordRsp, error) {
	pageRsp := &pbdata.PageRsp{
		Size: req.Size,
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

	res := ds.SelectJpTableLiteTotal(req.User)
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	// 生成一个[min, max]范围内的随机整数
	start := rnd.Intn(len(res) - int(req.Size))
	for _, r := range res[start : start+int(req.Size)] {
		rsp.JpWordList = append(rsp.JpWordList, &pbdata.JpWord{
			WordId:      r.Id,
			Word:        r.Vocabulary,
			Spell:       r.Character,
			Translation: r.Translator,
			Example:     r.SampleSentence,
		})
	}
	rsp.ErrorCode = int32(pbdata.ErrorCode_SUCCESS)
	rsp.ErrorMessage = "SUCCESS"

	return rsp, nil
}

func (js *JpLiteService) getDatabaseServiceClient() (dbservice.IDatabaseService, error) {
	svc, ok := js.FindService(common.ServiceIdDatabase)
	if !ok {
		return nil, errNoDatabaseService
	}
	ds, ok := svc.(dbservice.IDatabaseService)
	if !ok {
		return nil, errInvalidDatabaseService
	}
	return ds, nil
}

func (js *JpLiteService) getThirdServiceClient() (thirdservice.IThirdService, error) {
	svc, ok := js.FindService(common.ServiceIdThird)
	if !ok {
		return nil, errNoThirdService
	}
	ts, ok := svc.(thirdservice.IThirdService)
	if !ok {
		return nil, errInvalidThirdService
	}
	return ts, nil
}
