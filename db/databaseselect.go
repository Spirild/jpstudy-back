package dbservice

import dbdata "translasan-lite/db/generated"

// 放弃自定义一套灵活通用的查询范式了，但是查询语句还是约束在这里比较好

func (ds *DatabaseService) SelectJpTableLite(level int, word string, user string) []*dbdata.CharacterVocabularyTable {
	// 根据等级和单词查表, word可为空
	var res []*dbdata.CharacterVocabularyTable
	if word == "" {
		ds.db.Where("memory_level = ? AND is_del = ? AND user_qq = ?", level, 0, user).Find(&res)
	} else {
		// 查单词的时候不想考虑等级
		// ds.db.Where("memory_level = ? AND vocabulary = ? AND is_del = ? and user_qq = ?", level, word, 0, user).Find(&res)
		ds.db.Where("vocabulary = ? AND is_del = ? and user_qq = ?", word, 0, user).Find(&res)
	}
	return res
}

func (ds *DatabaseService) SelectJpTableDetail(start int, end int) ([]*dbdata.JapaneseVocabularyExternalSourceV1, int) {
	// 直接通过分页计算开始结束来查（提升性能）
	var res []*dbdata.JapaneseVocabularyExternalSourceV1
	ds.db.Where("id >= ? AND id <= ?", start, end).Find(&res)

	var count int64
	ds.db.Model(&dbdata.JapaneseVocabularyExternalSourceV1{}).Count(&count)

	return res, int(count)
}

func (ds *DatabaseService) GetMojiTokenMold() *dbdata.MojiToken {
	var res *dbdata.MojiToken
	ds.db.First(&res)
	return res
}

func (ds *DatabaseService) SelectJpTableLiteTotal(user string) []*dbdata.CharacterVocabularyTable {
	// 获取轻词本所有内容（未来可能按用户切片）
	var res []*dbdata.CharacterVocabularyTable
	ds.db.Where("is_del = ? and user = ?", 0, user).Find(&res)

	return res

}
