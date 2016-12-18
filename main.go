package main

import (
	"log"
	"realtime-batch-processing/infra"
	"strconv"
	"time"
)

const (
	videoTrendingKey  = "video_trending"
	videoViewCountKey = "video_view"
)

type VideoViewCount struct {
	VideoID   string
	ViewCount int64
	CreatedAt time.Time
}

func main() {
	infra.InitRedis()
	defer infra.CloseRedis()

	infra.Init()
	defer infra.CloseDB()

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
	if res := infra.Redis.Del(videoTrendingKey); res != nil {
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
		updateVideoCountAtMin(strconv.Itoa(pMin), now, 1, false)
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
	last := now.Add(6 * 60 * time.Minute)

	//TODO Better to scan previous time in first restart in order not to miss any event in infra
	updateVideoCountAtMin(strconv.Itoa(timer), now, 1, true)
	updateVideoCountAtMin(strconv.Itoa(lastSixHour), last, -1, false)
}

func updateVideoCountAtMin(min string, rawTime time.Time, minus int, saveToPg bool) {
	videoViewKey := videoViewCountKey + "_" + min

	if res := infra.Redis.HGetAll(videoViewKey); res != nil {
		mapCounting, err := res.Result()
		if err != nil {
			log.Println(err)
			return
		}

		//Update trending video
		for key, val := range mapCounting {
			count, _ := strconv.Atoi(val)
			if res := infra.Redis.ZIncrBy(videoTrendingKey, float64(minus*count), key); res != nil {
				if res.Err() != nil {
					log.Println(err)
				}
			}
		}

		//Save video count
		log.Println(saveToPg)
		if saveToPg {
			log.Println(mapCounting)
			for videoId, count := range mapCounting {
				viewCount, _ := strconv.Atoi(count)
				// Get video if exists
				videoViewCount := VideoViewCount{
					VideoID:   videoId,
					ViewCount: int64(viewCount),
					CreatedAt: rawTime,
				}

				err := infra.MySQL.Save(&videoViewCount).Error
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}

}
