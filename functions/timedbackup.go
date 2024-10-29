package functions

import (
	"goback-client/data"
	"log"
	"time"
)

var StopTimedBackup chan struct{}

// 设置定时备份
func TimedBackup(interval time.Duration) {
	// 如果 stopTimedBackup 已经存在，则发送停止信号
	if StopTimedBackup != nil {
		close(StopTimedBackup)
	}

	// 创建新的停止信号通道
	StopTimedBackup = make(chan struct{})

	go func() {
		log.Println("定时备份已启动")
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Println("定时备份中...")
				for _, f := range data.LocalFileList {
					ReBackup(f, []byte(data.Config.Key))
				}
			case <-StopTimedBackup:
				return
			}
		}
	}()
}
