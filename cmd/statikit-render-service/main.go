package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/sony/gobreaker"
	"github.com/zackattackz/statikit-render-svc/internal/adapters/cache"
	"github.com/zackattackz/statikit-render-svc/internal/middleware"
	"github.com/zackattackz/statikit-render-svc/internal/models"
	"github.com/zackattackz/statikit-render-svc/internal/service"
)

type deps struct {
	logger log.Logger
	cache  service.CacheService
	svc    service.RenderService
}

func run(deps) {

}

func main() {

	var (
		redisAddr = flag.String("redis-addr", "", "Redis cache address")
	)

	flag.Parse()

	redisAddrEnv, exists := os.LookupEnv("RENDERSVC_REDIS_ADDR")
	if *redisAddr == "" {
		if exists {
			*redisAddr = redisAddrEnv
		} else {
			panic("No RENDERSVC_REDIS_ADDR defined")
		}
	}

	logger := log.NewLogfmtLogger(os.Stderr)
	//logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     *redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	}).WithContext(ctx)

	cacheSvc := cache.NewRedisCacheService(redisClient, time.Hour*24)
	cacheEndpoint := service.EndPointFromCacheService(cacheSvc)
	cacheEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(cacheEndpoint)
	cacheSvc = service.CacheServiceFromEndpoint(context.Background(), cacheEndpoint)

	var svc service.RenderService
	svc = service.NewRenderService(cacheSvc)
	svc = middleware.LoggingMW(logger)(svc)

	err := svc.Render("Hello {{.Data.Wut}}!", models.Schema{Data: map[string]any{"Name": "Joe"}}, os.Stdout)
	if err != nil {
		logger.Log(err)
	}
	cancel()
}
