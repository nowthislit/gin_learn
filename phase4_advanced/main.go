package main

import (
	"log"

	"gin-learn/phase4/config"
	"gin-learn/phase4/internal/api"
	"gin-learn/phase4/internal/repository"
	"gin-learn/phase4/internal/service"
	"gin-learn/phase4/pkg/logger"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatal("Failed to init config:", err)
	}

	// 初始化日志
	logger.Init(config.C.Log.Level)
	defer logger.Sync()

	logger.Info("Application starting...",
		logger.String("env", config.C.App.Env),
		logger.String("version", config.C.App.Version),
	)

	// 初始化数据库
	db, err := repository.InitDB()
	if err != nil {
		logger.Fatal("Failed to init database", logger.ErrorField(err))
	}

	// 初始化仓库
	repo := repository.NewRepository(db)

	// 初始化服务
	svc := service.NewService(repo)

	// 启动HTTP服务器
	server := api.NewServer(svc)
	if err := server.Run(config.C.Server.Port); err != nil {
		logger.Fatal("Server failed to start", logger.ErrorField(err))
	}
}
