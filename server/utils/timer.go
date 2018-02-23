package utils

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func TimeIntervalFileName(interval int) string {
	t := time.Now().Unix()
	if interval < 60 {
		minute, err := strconv.Atoi(time.Unix(t, 0).Format("04"))
		if err != nil {
			log.Println("生成日志文件名失败", err)
			return time.Unix(t, 0).Format("200601021504")
		}
		suffix := fmt.Sprintf("%02d", minute/interval*interval)
		return time.Unix(t, 0).Format("2006010215") + suffix
	} else if interval < 1440 {
		interval = interval / 60
		hour, err := strconv.Atoi(time.Unix(t, 0).Format("15"))
		if err != nil {
			log.Println("生成日志文件名失败", err)
			return time.Unix(t, 0).Format("2006010215") + "00"
		}
		suffix := fmt.Sprintf("%02d", hour/interval*interval)
		return time.Unix(t, 0).Format("20060102") + suffix + "00"
	}
	interval = interval / 1440
	day, err := strconv.Atoi(time.Unix(t, 0).Format("02"))
	if err != nil {
		log.Println("生成日志文件名失败", err)
		return time.Unix(t, 0).Format("20060102") + "0000"
	}
	suffix := fmt.Sprintf("%02d", day/interval*interval)
	return time.Unix(t, 0).Format("200601") + suffix + "0000"
}
