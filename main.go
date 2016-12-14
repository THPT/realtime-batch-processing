package main

import (
	"log"
	"strconv"
	"time"

	"realtime-batch-processing/redis"
)

const (
	videoTrendingKey  = "video_trending"
	videoViewCountKey = "video_view"
)

func main() {
	redis.InitRedis()
	defer redis.CloseRedis()

	// Recounting all first
	log.Println("Recounting trending video")
	recountingTrending()
	log.Println("Done")
	log.Println("Ticker is running...")

	ticker := time.NewTicker(time.Minute * 1).C
	for {
		select {
		case <-ticker:
			log.Println("Tick")
			updateVideoCount()
			log.Println("Tick done")
			log.Println("---")
		}
	}
}

func recountingTrending() {
	// Reset trending key
	if res := redis.Redis.Del(videoTrendingKey); res != nil {
		if err := res.Err(); err != nil {
			panic(err)
		}
	}

	now := time.Now()
	min := now.Minute()
	hour := now.Hour()
	timer := min + hour*60

	count := 0
	for count < 60*6 {
		pMin := timer - count
		if pMin < 0 {
			pMin += 60 * 24
		}
		updateVideoCountAtMin(strconv.Itoa(pMin), 1)
		count++
	}
}

func updateVideoCount() {
	now := time.Now()
	min := now.Minute()
	hour := now.Hour()
	timer := min + hour*60

	lastSixHour := timer - 6*60
	if lastSixHour < 0 {
		lastSixHour += 60 * 24
	}
	updateVideoCountAtMin(strconv.Itoa(timer), 1)
	updateVideoCountAtMin(strconv.Itoa(lastSixHour), -1)
}

func updateVideoCountAtMin(min string, minus int) {
	videoViewKey := videoViewCountKey + "_" + min

	if res := redis.Redis.HGetAll(videoViewKey); res != nil {
		mapCounting, err := res.Result()
		if err != nil {
			log.Println(err)
			return
		}

		for key, val := range mapCounting {
			count, _ := strconv.Atoi(val)
			if res := redis.Redis.ZIncrBy(videoTrendingKey, float64(minus*count), key); res != nil {
				if res.Err() != nil {
					log.Println(err)
				}
			}
		}
	}

}
