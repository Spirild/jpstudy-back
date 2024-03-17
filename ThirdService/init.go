package thirdservice

import (
	"translasan-lite/core"
	pbdata "translasan-lite/proto/generated"
)

// module register
func init() {
	core.RegisterCompType("ThirdService", (*ThirdService)(nil))
}

type IThirdService interface {
	MojiTranlate(searchContent string) ([]*pbdata.MojiResponseWord, error)
	SparkDemo(question string) (string, error)
}
