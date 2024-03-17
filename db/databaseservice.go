package dbservice

import (
	"context"
	"fmt"

	"translasan-lite/common"
	"translasan-lite/core"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseService struct {
	core.BaseComponent
	db *gorm.DB

	user     string
	password string
	addr     string
	dbname   string

	// 处理有些在注册时就希望从数据库拿到些数据的情况
	connectChan chan struct{}
}

func (ds *DatabaseService) ServiceID() int {
	return common.ServiceIdDatabase
}

func (ds *DatabaseService) Init(n *core.Node, cfg *core.ServiceConfig) {
	(&ds.BaseComponent).Init(n, cfg)
	ds.user, _ = cfg.GetString("user")
	ds.password, _ = cfg.GetString("password")
	ds.addr, _ = cfg.GetString("address")
	ds.dbname, _ = cfg.GetString("dbname")

	ds.connectChan = make(chan struct{})
}

func (ds *DatabaseService) Run(ctx context.Context) error {

	dsn_format := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(dsn_format, "root", "fsFOREVER1022", "39.104.92.255:3306", "japanese_study")
	ds.db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	ds.Log.Info("Database starts running")
	close(ds.connectChan)
	<-ctx.Done()
	ds.Log.Info("Database stops running")

	return nil
}

func (ds *DatabaseService) WaitDBConnect() <-chan struct{} {
	// 启动时等待数据库连接完成的方法
	return ds.connectChan
}

func (ds *DatabaseService) SelfInsert(value interface{}) error {
	// 如果是多条数据会逐一被导入
	res := ds.db.Create(value)
	return res.Error
}

func (ds *DatabaseService) SelfUpdate(value interface{}) error {
	// 如果是多条数据会逐一被导入
	res := ds.db.Save(value)
	return res.Error
}
