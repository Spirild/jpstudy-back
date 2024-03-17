package dbservice

import (
	dbdata "translasan-lite/db/generated"

	"gorm.io/gorm"
)

func (ds *DatabaseService) JpWordLevelUpdate(wordid int, updateVal int) error {
	var result *gorm.DB
	if updateVal == 1 {
		result = ds.db.Model(&dbdata.CharacterVocabularyTable{}).Where("id = ?", wordid).Update("memory_level", gorm.Expr("memory_level + ?", updateVal))
	} else if updateVal == -1 {
		result = ds.db.Model(&dbdata.CharacterVocabularyTable{}).Where("id = ? AND memory_level > ?", wordid, 1).Update("memory_level", gorm.Expr("memory_level + ?", updateVal))
	}

	return result.Error
}

func (ds *DatabaseService) JpWordDelete(wordid int) error {
	// 逻辑删除, 暂不彻底删除
	result := ds.db.Model(&dbdata.CharacterVocabularyTable{}).Where("id = ?", wordid).Update("is_del", 1)
	return result.Error
}
