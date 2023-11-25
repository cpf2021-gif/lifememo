package main

import (
	"context"
	"flag"
	"fmt"
	"lifememo/application/like/cron/internal/config"
	"lifememo/application/like/cron/internal/logic"
	"lifememo/application/like/cron/internal/svc"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/cron.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()

	l := logic.NewUpdateLogic(ctx, svcCtx)

	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(30).Second().Do(l.Update)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Starting cron server...")
	s.StartBlocking()
}
