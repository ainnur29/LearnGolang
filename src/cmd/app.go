package main

import (
	"flag"
	"golang-bulang-bolang/src/config"
	"golang-bulang-bolang/src/preference"
	"golang-bulang-bolang/src/repository"
	"golang-bulang-bolang/src/service"
	"golang-bulang-bolang/src/util"
	
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	minJitter    int
	maxJitter    int
	sqlClient0   *sqlx.DB
	redisClient0 *redis.Client
	redisClient1 *redis.Client
	redisClient2 *redis.Client
	app          config.App
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
		if sqlClient0 != nil {
			sqlClient0.Close()
		}
	}

	app.Serve()
}