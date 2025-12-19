package main

import (
	"flag"
	"learngolang/src/config"
	"learngolang/src/repository"
	"learngolang/src/service"
	"learngolang/src/util"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	minJitter int
	maxJitter int
	sql0      *sqlx.DB
	redis0    *redis.Client
	redis1    *redis.Client
	redis2    *redis.Client
	scheduler *config.Scheduler
	app       config.App
)

func init() {
	flag.IntVar(&minJitter, "minSleep", DefaultMinJitter, "min. sleep duration during app initialization")
	flag.IntVar(&maxJitter, "maxSleep", DefaultMaxJitter, "max. sleep duration during app initialization")
	flag.Parse()

	sleepWithJitter(minJitter, maxJitter)

	// Config Initialization
	conf, err := InitConfig()
	if err != nil {
		panic(err)
	}

	// Logger Initialization
	log := config.InitLogger(conf.Logger)

	// SQL Initialization
	sqlClient0 = config.InitDB(log, conf.Postgres)

	// Query Loader Initialization
	queryLoader, err := config.InitQueryLoader(log, conf.Queries)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to load queries")
	}

	// Initialize dependencies
	repository := repository.InitRepository(sqlClient0, redisClient0, queryLoader, conf.Redis.CacheTTL)
	service := service.InitService(repository)

	// Initialize validator
	util.Validator()

	// Auth Initialization
	auth := config.InitAuth(log, conf.Auth, redisClient1)

	// Middleware Initialization
	middleware := config.InitMiddleware(log, auth, redisClient2)

	// HTTP Gin Initialization
	httpGin := config.InitHttpGin(log, middleware)

	// REST Handler Initialization
	resthandler.InitRestHandler(httpGin, auth, middleware, service)

	// HTTP Server Initialization
	httpServer := config.InitHttpServer(log, conf.Server, httpGin)

	// App Initialization
	app = config.InitGrace(log, httpServer)
}

func main() {
	defer func() {
		if redis0 != nil {
			redis0.Close()
		}

		if redis1 != nil {
			redis1.Close()
		}

		if redis2 != nil {
			redis2.Close()
		}

		if sql0 != nil {
			sql0.Close()
		}

		if scheduler != nil {
			scheduler.Stop()
		}
	}()

	app.Serve()
}
