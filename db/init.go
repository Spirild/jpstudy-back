package dbservice

import (
	"translasan-lite/core"
	dbdata "translasan-lite/db/generated"
)

// module register
func init() {
	core.RegisterCompType("DatabaseService", (*DatabaseService)(nil))
}

type IDatabaseService interface {
	WaitDBConnect() <-chan struct{}

	// Create
	SelfInsert(value interface{}) error // this is common insert

	// Read
	SelectJpTableLite(level int, word string, user string) []*dbdata.CharacterVocabularyTable
	SelectJpTableLiteTotal(user string) []*dbdata.CharacterVocabularyTable
	GetMojiTokenMold() *dbdata.MojiToken
	SelectJpTableDetail(start int, end int) ([]*dbdata.JapaneseVocabularyExternalSourceV1, int)

	// Update
	SelfUpdate(value interface{}) error // this is common update
	JpWordLevelUpdate(wordid int, updateVal int) error

	// Del
	JpWordDelete(wordid int) error
}
