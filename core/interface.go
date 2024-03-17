package core

import (
	"context"
)

type IComponent interface {
	Name() string
	// ID() int
	Init(n *Node, cfg *ServiceConfig)
	Run(ctx context.Context) error
	Update(ctx context.Context) error

	FindService(serviceID int) (IService, bool)
}

type IHttpServer interface {
	Run() error
	Shutdown()
}

type ISocketServer interface {
	Run(ctx context.Context) error
	Shutdown()
}

type IService interface {
	ServiceID() int
}
