package dbdata

//手动配置一波表名。。暂时找不到别的比较好的方法

func (*CharacterVocabularyTable) TableName() string {
	return "character_vocabulary_table"
}

func (*JapaneseVocabularyExternalSourceV1) TableName() string {
	return "japanese_vocabulary_external_source_v1"
}

func (*MojiToken) TableName() string {
	return "moji_token"
}
