package core

// 该部分为实现我具体功能的node节点定义。不用完全复刻仅需先用好我想要的功能

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Node struct {
	Id   int
	SID  int
	Name string

	comp     []IComponent
	services map[int]IService
	// svcManager IServiceManager

	log    *zap.Logger
	config *NodeConfig
	// 最后的config大概类似于拓展, 定义具体功能相关。设计模式抽象工厂。
	// 但，真用不上，跳过吧
}

func (n *Node) GetLogger() *zap.Logger {
	return n.log
}

func NewNode(cfg *NodeConfig) *Node {
	n := Node{
		config: cfg,
		Name:   cfg.Name,
		Id:     cfg.NodeId,
		SID:    cfg.ServerId,
	}

	log, err := initLogger(cfg.Log)
	if err != nil {
		panic(fmt.Sprintf("new node failed to init log, err=%s", err.Error()))
	}
	n.log = log.With(
		zap.Int("Id", n.Id),
		zap.Int("Server ID", n.SID),
		zap.String("node", n.Name))
	n.log = n.log.With(zap.Namespace("msg"))
	return &n
}

func (n *Node) RunForever() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()

	var wg sync.WaitGroup
	errCh := make(chan error, 64)

	// start each component
	for _, comp := range n.comp {
		go func(c IComponent) {
			wg.Add(1)
			defer func() {
				if err := recover(); err != nil {
					err1 := fmt.Errorf("%v", err)
					n.log.Error("module panic", zap.Error(err1), zap.String("module", c.Name()))
				}
				wg.Done()
				n.log.Info("module done.", zap.String("module", c.Name()))
			}()
			n.log.Info("module starting.", zap.String("module", c.Name()))
			err := c.Run(ctx)
			if err != nil {
				n.log.Info("module Run() returns", zap.Error(err), zap.String("module", c.Name()))
				// notify module error without blocking on error channel
				select {
				case errCh <- err:
				default:
				}
			}
		}(comp)
	}

	runtime.Gosched()
	if n.config.FrameMS > 0 {
		go func() {
			wg.Add(1)
			defer func() {
				if err := recover(); err != nil {
					err1 := fmt.Errorf("%v", err)
					n.log.Error("RunForever panic", zap.Error(err1))
				}
				wg.Done()
			}()
			if err := n.Run(ctx); err != nil {
				n.log.Error("node.Run() exception", zap.String("error", err.Error()), zap.String("stack", string(debug.Stack())))
				return
			}
			n.log.Info("node run finished.")
		}()
	}

	n.log.Info("node started!")
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGQUIT)

	select {
	case <-ctx.Done():
		n.log.Info("node ctx done")
	case s := <-sig:
		n.log.Info("signal, canceling...", zap.String("sig", s.String()))
		cancel()
	case err := <-errCh:
		n.log.Info("module failing", zap.Error(err))
		cancel()
	}

	wg.Wait()
	n.shutdown()

	close(errCh)
	n.log.Info("node exit!")
}

func (n *Node) Run(ctx context.Context) error {
	if n.config.FrameMS <= 0 {
		return nil
	}
	ticker := time.NewTicker(time.Millisecond * time.Duration(n.config.FrameMS))
	defer func() {
		ticker.Stop()
		if err := recover(); err != nil {
			err1 := fmt.Errorf("%v", err)
			n.log.Error("node.Run() panic", zap.Error(err1), zap.String("stack", string(debug.Stack())))
		}
	}()

	for {
		select {
		case <-ticker.C:
			n.log.Info("node tick.")
			for _, comp := range n.comp {
				if err := comp.Update(ctx); err != nil {
					return err
				}
			}
			runtime.Gosched()
		case <-ctx.Done():
			n.log.Info("node run finished.")
			return nil
		}
	}
}

func (n *Node) shutdown() {
	endLogger(n.log)
}

func (n *Node) AddService(serviceID int, svc IService) {
	if n.services == nil {
		n.services = make(map[int]IService)
	}
	n.log.Info("AddService", zap.Int("sid", serviceID), zap.String("name", svc.(IComponent).Name()))
	n.services[serviceID] = svc
}

func (n *Node) AddComponent(comp IComponent) {
	n.log.Info("AddComponent", zap.String("comp", comp.Name()))
	n.comp = append(n.comp, comp)
}
