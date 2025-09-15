package my_init

import (
	"context"
	"group_ten_server/config"
	"group_ten_server/dao"
	"group_ten_server/service"
	"log"

	"github.com/robfig/cron/v3"
)

func StartDailyEntryCountJob() {
	c := cron.New()
	// 每天0点执行
	c.AddFunc("0 0 * * *", func() {
		ctx := context.Background()
		v, _ := service.RedisClient.Get(ctx, config.TOTAL_NUM).Int()
		err := dao.InsertNewDailyEntryCount(v)
		if err != nil {
			log.Println("每日统计插入失败:", err)
		} else {
			log.Println("每日统计已更新，今日入馆人数:", v)
		}
		service.RedisClient.Set(ctx, config.TOTAL_NUM, 0, 0)
	})
	c.Start()
}
