package main

// 积木是搭好了。。但是感觉嗯，很像被夺舍了。我必须想想自己的设计。
// 重点是为自己所用，如何高效快捷，六经注我

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime/pprof"

	"translasan-lite/core"

	"go.uber.org/zap"
)

var (
	configPath = path.Join(".", "config", "config.json")
)

type ModuleConfigData map[string]interface{}

type Config struct {
	Node    *core.NodeConfig
	Modules []ModuleConfigData
}

func ReadConfig(path string) *Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	config := &Config{}
	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		fmt.Println("failed to unmarshal config", err)
		return nil
	}
	return config
}

func initComponents(node *core.Node, configList []ModuleConfigData) {
	if !core.CheckComponent() {
		panic("no component registered, go generate not called? try build.bat instead of go build.")
	}
	for _, configData := range configList {
		mc := core.NewModuleConfig(configData)

		compName, _ := mc.GetString("Name")

		ignore, _ := mc.GetBool("Ignore")
		if ignore {
			node.GetLogger().Sugar().Info("skip module ", compName)
			continue
		}

		compType, ok := core.ComponentType(compName)
		if !ok {
			node.GetLogger().Warn("failed to get component type with name ", zap.String("comp", compName))
			continue
		}

		compValue := reflect.New(compType.Elem()) // compType is *xxx, compType.Elem() is xxx
		comp, ok := compValue.Interface().(core.IComponent)
		if !ok {
			node.GetLogger().Sugar().Warn("failed to cast to IComponent, name ", compName, ",", compType.Name(), ",", compValue)
			continue
		}

		node.GetLogger().Info("module Init", zap.String("module", compName))
		comp.Init(node, mc)
		node.AddComponent(comp)

		if svc, ok := comp.(core.IService); ok {
			node.AddService(svc.ServiceID(), svc)
		}
	}
}

func RunNode(path string) {
	config := ReadConfig(path)
	if config == nil {
		panic(fmt.Sprintf("failed to read config file, path is %s.", path))
	}

	node := core.NewNode(config.Node)
	initComponents(node, config.Modules)
	node.RunForever()
}

func main() {

	fmt.Println("main()", os.Args)
	var arg1 string
	if len(os.Args) > 1 {
		arg1 = os.Args[1]
	} else {
		arg1 = configPath
	}

	if arg1 == "cpu" {

		file, err := os.Create("./cpu_optimize.pprof")
		// file, err := os.OpenFile("./cpu.pprof", os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("create cpu pprof failed, err: %v\n", err)
			return
		}
		fmt.Println("start cpu pprof")
		pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
		fmt.Println("start trace")
		// trace.Start(tf)
		// defer trace.Stop()
		RunNode(configPath)
	} else {
		RunNode(arg1)
	}
}
