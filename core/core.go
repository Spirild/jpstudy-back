package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"sync"
)

var (
	componentTypes = make(map[string]reflect.Type)
)

func RegisterCompType(name string, c IComponent) {
	t := reflect.TypeOf(c)
	componentTypes[name] = t
}

func ComponentType(name string) (reflect.Type, bool) {
	val, ok := componentTypes[name]
	return val, ok
}

func CheckComponent() bool {
	return len(componentTypes) > 0
}

type BaseComponent struct {
	// 此类为注册服务的通用基础类
	Config *ServiceConfig
	Log    *zap.Logger
	node   *Node
	Locker sync.Mutex
}

func (bc *BaseComponent) Name() string {
	// 该部分决定的是服务的类。但我不一定非得用这种实现方式
	// 即，自行定义一个注册表即可
	// 只是获取名字。。似乎可以直接跳过
	name, ok := bc.Config.GetString("Name")
	if !ok {
		return reflect.TypeOf(bc).String()
	}
	return name
}

func (bc *BaseComponent) Init(n *Node, cfg *ServiceConfig) {
	// 初始化，将公共类传入
	bc.Config = cfg
	bc.Log = n.GetLogger()
	bc.node = n
}

func (bc *BaseComponent) Update(ctx context.Context) error {
	return nil
}

func (bc *BaseComponent) FindService(serviceID int) (IService, bool) {
	// 同样的，这里相当于自行实现了一个反射
	if bc.node.services == nil {
		return nil, false
	}
	svc, ok := bc.node.services[serviceID]
	return svc, ok
}

func (bc *BaseComponent) Blocking(ctx context.Context) {
	// for
	{
		select {
		case <-ctx.Done():
			// bc.Log.Info("Blocking done!", zap.String("Module", bc.Name()))
			return
		}
	}
}
func (bs *BaseComponent) NodeID() int {
	return bs.node.Id
}
func (bs *BaseComponent) ServerID() int {
	return bs.node.SID
}

func (bc *BaseComponent) RunHttpServer(ctx context.Context, server IHttpServer) error {
	ch := make(chan error, 2)
	var wg sync.WaitGroup

	defer func() {
		wg.Wait()
		close(ch)
	}()

	go func(ch chan error) {
		wg.Add(1)
		defer func() {
			if err := recover(); err != nil {
				err1 := fmt.Errorf("%v", err)
				bc.Log.Error("RunHttpServer panic", zap.Error(err1))
			}
			wg.Done()
		}()
		err := server.Run()

		// notify outside that the http server encountered error
		if err != nil {
			select {
			case ch <- err:
			default:
			}
		}
	}(ch)

	select {
	case err := <-ch:
		// bc.Log.Info("bc notify err")
		server.Shutdown()
		return err
	case <-ctx.Done():
		// bc.Log.Info("bc done")
		server.Shutdown()
		return nil
	}
}

func (bc *BaseComponent) RunSocketServer(ctx context.Context, server ISocketServer) error {
	ch := make(chan error, 2)
	var wg sync.WaitGroup

	defer func() {
		wg.Wait()
		close(ch)
	}()

	go func(ch chan error) {
		wg.Add(1)
		defer func() {
			if err := recover(); err != nil {
				err1 := fmt.Errorf("%v", err)
				bc.Log.Error("RunSocketServer panic", zap.Error(err1))
			}
			wg.Done()
		}()
		err := server.Run(ctx)

		// notify outside that the http server encountered error
		if err != nil {
			select {
			case ch <- err:
			default:
			}
		}
	}(ch)

	select {
	case err := <-ch:
		// bc.Log.Info("bc notify err")
		server.Shutdown()
		return err
	case <-ctx.Done():
		// bc.Log.Info("bc done")
		server.Shutdown()
		return nil
	}
}

func (bc *BaseComponent) SendData() {

}

func initLogger(config *LogConfig) (*zap.Logger, error) {
	if config == nil {
		return nil, errors.New("nil log config")
	}
	// level
	var level zapcore.Level
	if config.Level == "debug" {
		level = zap.DebugLevel
	} else if config.Level == "info" {
		level = zap.InfoLevel
	} else if config.Level == "warn" {
		level = zap.WarnLevel
	} else if config.Level == "error" {
		level = zap.ErrorLevel
	} else {
		level = zap.InfoLevel
	}

	// opt, encoder and core
	var opts []zap.Option
	var encoder zapcore.Encoder
	var cores []zapcore.Core
	var core zapcore.Core
	if config.Env == "dev" {
		opts = []zap.Option{
			zap.Development(),
			zap.AddCaller(),
			zap.AddStacktrace(zap.WarnLevel),
		}
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else if config.Env == "prod" {
		opts = []zap.Option{
			zap.AddCaller(),
			zap.AddStacktrace(zap.WarnLevel),
		}
		encoderConfig := zap.NewProductionEncoderConfig()
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// file syncer
	lumberjack_info := NewLogWritter(config.Filename, config.FileSplitTime)

	writer := zapcore.AddSync(lumberjack_info)
	if config.Stdout == true {
		writer = zap.CombineWriteSyncers(writer, os.Stderr)
	}
	core = zapcore.NewCore(encoder, writer, level)
	cores = append(cores, core)
	if config.ErrFile {
		lumberjack_error := NewLogWritter(config.ErrorFileName, config.FileSplitTime)
		writer = zapcore.AddSync(lumberjack_error)
		core = zapcore.NewCore(encoder, writer, zap.ErrorLevel)
		cores = append(cores, core)
	}
	logger := zap.New(zapcore.NewTee(cores...), opts...)

	// TODO: daily log file split
	return logger, nil

}

func endLogger(logger *zap.Logger) {
	if logger != nil {
		logger.Sync()
	}
}
