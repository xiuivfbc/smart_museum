package my_init

import (
	"github.com/robfig/cron/v3"
	"group_ten_server/dao"
	"log"
)

func StartDailyEntryCountJob() {
	c := cron.New()
	// 每天0点执行
	c.AddFunc("0 0 * * *", func() {
		err := dao.InsertNewDailyEntryCount()
		if err != nil {
			log.Println("每日统计插入失败:", err)
		}
	})
	c.Start()
}
