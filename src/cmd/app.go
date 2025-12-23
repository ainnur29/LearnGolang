package main

import (
	"flag"
	"learngolang/src/config"
	restHandler "learngolang/src/handler/rest"

	// schedHandler "learngolang/src/handler/scheduler"
	"learngolang/src/preference"
	"learngolang/src/repository"
	"learngolang/src/service"

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
	sql0 = config.InitDB(log, conf.Postgres)

	// Redis Initialization
	redis0 = config.InitRedis(log, conf.Redis, preference.REDIS_APPS)
	redis1 = config.InitRedis(log, conf.Redis, preference.REDIS_AUTH)
	redis2 = config.InitRedis(log, conf.Redis, preference.REDIS_LIMITER)

	// Query Loader Initialization
	queryLoader, err := config.InitQueryLoader(log, conf.Queries)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to load queries")
	}

	// Initialize dependencies
	repository := repository.InitRepository(sql0, redis0, queryLoader, conf.Redis.CacheTTL)
	service := service.InitService(repository)

	// Initialize validator
	config.InitValidator(log)

	// Auth Initialization
	auth := config.InitAuth(log, conf.Auth, redis1)

	// Middleware Initialization
	middleware := config.InitMiddleware(log, auth, redis2)

	// HTTP Gin Initialization
	httpGin := config.InitHttpGin(log, middleware)

	// REST Handler Initialization
	restHandler.InitRestHandler(httpGin, auth, middleware, service)

	// //Scheduler Initialization
	// scheduler = config.InitScheduler(log, conf.Scheduler)
	// schedHandler.InitSchedulerHandler(log, scheduler, service, conf.Scheduler.SchedulerJobs)

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
