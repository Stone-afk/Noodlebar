// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package interactive

import (
	"context"
	"sync"

	"github.com/ecodeclub/mq-api"
	"github.com/ecodeclub/webook/internal/interactive/internal/event"
	"github.com/ecodeclub/webook/internal/interactive/internal/repository"
	"github.com/ecodeclub/webook/internal/interactive/internal/repository/dao"
	"github.com/ecodeclub/webook/internal/interactive/internal/service"
	"github.com/ecodeclub/webook/internal/interactive/internal/web"
	"github.com/ego-component/egorm"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitModule(db *gorm.DB, q mq.MQ) (*Module, error) {
	interactiveDAO := InitTablesOnce(db)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO)
	serviceService := service.NewService(interactiveRepository)
	consumer := initConsumer(serviceService, q)
	handler := web.NewHandler(serviceService)
	module := &Module{
		Svc: serviceService,
		c:   consumer,
		Hdl: handler,
	}
	return module, nil
}

// wire.go:

var HandlerSet = wire.NewSet(
	InitTablesOnce, repository.NewCachedInteractiveRepository, service.NewService, web.NewHandler,
)

var once = &sync.Once{}

func InitTablesOnce(db *egorm.Component) dao.InteractiveDAO {
	once.Do(func() {
		_ = dao.InitTables(db)
	})
	return dao.NewInteractiveDAO(db)
}

func initConsumer(svc service.Service, q mq.MQ) *event.Consumer {
	consumer, err := event.NewSyncConsumer(svc, q)
	if err != nil {
		panic(err)
	}
	consumer.Start(context.Background())
	return consumer
}
