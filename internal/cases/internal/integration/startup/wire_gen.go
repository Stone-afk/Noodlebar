// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/ecodeclub/ginx/session"
	"github.com/ecodeclub/webook/internal/ai"
	"github.com/ecodeclub/webook/internal/cases"
	"github.com/ecodeclub/webook/internal/cases/internal/event"
	"github.com/ecodeclub/webook/internal/cases/internal/repository"
	"github.com/ecodeclub/webook/internal/cases/internal/repository/cache"
	"github.com/ecodeclub/webook/internal/cases/internal/repository/dao"
	"github.com/ecodeclub/webook/internal/cases/internal/service"
	"github.com/ecodeclub/webook/internal/cases/internal/web"
	"github.com/ecodeclub/webook/internal/interactive"
	"github.com/ecodeclub/webook/internal/member"
	testioc "github.com/ecodeclub/webook/internal/test/ioc"
)

// Injectors from wire.go:

func InitModule(syncProducer event.SyncEventProducer, knowledgeBaseProducer event.KnowledgeBaseEventProducer, aiModule *ai.Module, memberModule *member.Module, sp session.Provider, intrModule *interactive.Module) (*cases.Module, error) {
	db := testioc.InitDB()
	caseDAO := cases.InitCaseDAO(db)
	ecacheCache := testioc.InitCache()
	caseCache := cache.NewCaseCache(ecacheCache)
	caseRepo := repository.NewCaseRepo(caseDAO, caseCache)
	mq := testioc.InitMQ()
	interactiveEventProducer, err := event.NewInteractiveEventProducer(mq)
	if err != nil {
		return nil, err
	}
	serviceService := service.NewService(caseRepo, interactiveEventProducer, knowledgeBaseProducer, syncProducer)
	adminCaseHandler := web.NewAdminCaseHandler(serviceService)
	examineDAO := dao.NewGORMExamineDAO(db)
	examineRepository := repository.NewCachedExamineRepository(examineDAO)
	llmService := aiModule.Svc
	examineService := service.NewLLMExamineService(caseRepo, examineRepository, llmService)
	service2 := intrModule.Svc
	service3 := memberModule.Svc
	handler := web.NewHandler(serviceService, examineService, service2, service3, sp)
	caseSetDAO := dao.NewCaseSetDAO(db)
	caseSetRepository := repository.NewCaseSetRepo(caseSetDAO)
	caseSetService := service.NewCaseSetService(caseSetRepository, caseRepo, interactiveEventProducer)
	adminCaseSetHandler := web.NewAdminCaseSetHandler(caseSetService)
	repositoryBaseSvc := aiModule.KnowledgeBaseSvc
	knowledgeBaseService := initKnowledgeBaseSvc(repositoryBaseSvc, caseRepo)
	knowledgeBaseHandler := web.NewKnowledgeBaseHandler(knowledgeBaseService)
	module := &cases.Module{
		AdminHandler:         adminCaseHandler,
		ExamineSvc:           examineService,
		Svc:                  serviceService,
		Hdl:                  handler,
		AdminSetHandler:      adminCaseSetHandler,
		KnowledgeBaseHandler: knowledgeBaseHandler,
	}
	return module, nil
}

func InitExamModule(syncProducer event.SyncEventProducer, knowledgeBaseProducer event.KnowledgeBaseEventProducer, intrModule *interactive.Module, memberModule *member.Module, sp session.Provider, aiModule *ai.Module) (*cases.Module, error) {
	db := testioc.InitDB()
	caseDAO := cases.InitCaseDAO(db)
	ecacheCache := testioc.InitCache()
	caseCache := cache.NewCaseCache(ecacheCache)
	caseRepo := repository.NewCaseRepo(caseDAO, caseCache)
	mq := testioc.InitMQ()
	interactiveEventProducer, err := event.NewInteractiveEventProducer(mq)
	if err != nil {
		return nil, err
	}
	serviceService := service.NewService(caseRepo, interactiveEventProducer, knowledgeBaseProducer, syncProducer)
	caseSetDAO := dao.NewCaseSetDAO(db)
	caseSetRepository := repository.NewCaseSetRepo(caseSetDAO)
	caseSetService := service.NewCaseSetService(caseSetRepository, caseRepo, interactiveEventProducer)
	examineDAO := dao.NewGORMExamineDAO(db)
	examineRepository := repository.NewCachedExamineRepository(examineDAO)
	llmService := aiModule.Svc
	examineService := service.NewLLMExamineService(caseRepo, examineRepository, llmService)
	service2 := intrModule.Svc
	service3 := memberModule.Svc
	handler := web.NewHandler(serviceService, examineService, service2, service3, sp)
	adminCaseSetHandler := web.NewAdminCaseSetHandler(caseSetService)
	adminCaseHandler := web.NewAdminCaseHandler(serviceService)
	examineHandler := web.NewExamineHandler(examineService)
	caseSetHandler := web.NewCaseSetHandler(caseSetService, examineService, service2, sp)
	repositoryBaseSvc := aiModule.KnowledgeBaseSvc
	knowledgeBaseService := initKnowledgeBaseSvc(repositoryBaseSvc, caseRepo)
	knowledgeBaseHandler := web.NewKnowledgeBaseHandler(knowledgeBaseService)
	module := &cases.Module{
		Svc:                  serviceService,
		SetSvc:               caseSetService,
		ExamineSvc:           examineService,
		Hdl:                  handler,
		AdminSetHandler:      adminCaseSetHandler,
		AdminHandler:         adminCaseHandler,
		ExamineHdl:           examineHandler,
		CsHdl:                caseSetHandler,
		KnowledgeBaseHandler: knowledgeBaseHandler,
	}
	return module, nil
}

// wire.go:

func initKnowledgeBaseSvc(svc ai.KnowledgeBaseService, caRepo repository.CaseRepo) service.KnowledgeBaseService {
	return service.NewKnowledgeBaseService(caRepo, svc, "knowledge_id")
}
