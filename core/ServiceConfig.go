package core

// 该部分是实现注册表功能
// 准确说是比较灵活的对配置json进行拓展, 随后在此进行读取

type ServiceConfig struct {
	m map[string]interface{}
}

type NodeConfig struct {
	NodeId   int
	ServerId int
	Name     string
	Log      *LogConfig
	FrameMS  int
}

func (mc *ServiceConfig) Get(k string) (interface{}, bool) {
	v, ok := mc.m[k]
	return v, ok
}

type LogConfig struct {
	Env           string // dev, prod,
	Filename      string
	ErrorFileName string
	Level         string // debug, info, warn, error
	Stdout        bool
	ErrFile       bool
	FileSplitTime int32 //hour
}

func NewModuleConfig(kv map[string]interface{}) *ServiceConfig {
	cc := ServiceConfig{
		m: kv,
	}
	return &cc
}

func (mc *ServiceConfig) GetString(k string) (string, bool) {
	v, ok := mc.Get(k)
	if !ok {
		return "", false
	}
	return v.(string), true
}

func (mc *ServiceConfig) GetBool(k string) (bool, bool) {
	v, ok := mc.Get(k)
	if !ok {
		return false, false
	}
	return v.(bool), true
}
