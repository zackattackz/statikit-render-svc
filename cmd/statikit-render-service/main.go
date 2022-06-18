package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/zackattackz/statikit-render-svc/internal/adapters/cache"
	"github.com/zackattackz/statikit-render-svc/internal/middleware"
	"github.com/zackattackz/statikit-render-svc/internal/models"
	"github.com/zackattackz/statikit-render-svc/internal/ports"
	"github.com/zackattackz/statikit-render-svc/internal/service"
)

type deps struct {
	logger log.Logger
	cache  ports.Cache
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

	cache := cache.NewRedisCache(redisClient, time.Hour*24)

	var svc service.RenderService
	svc = service.New(cache)
	svc = middleware.LoggingMW(logger)(svc)

	err := svc.Render("Hello {{.Data.Name}}!", models.Schema{Data: map[string]any{"Name": "Joe"}}, os.Stdout)
	if err != nil {
		logger.Log(err)
	}
	cancel()
}
