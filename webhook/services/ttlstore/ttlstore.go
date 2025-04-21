package ttlstore

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"titles.run/webhook/services/interfaces"

	strava "titles.run/strava/models"

	"github.com/redis/go-redis/v9"
)

type TTLStoreServiceImpl struct {
	rdb *redis.Client
}

var ctx = context.Background()

func NewTTLStoreService() interfaces.TTLStoreService {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR environment variable not set")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	maxRetries := 20
	retryInterval := 5 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := rdb.Ping(ctx).Err()
		if err == nil {
			return &TTLStoreServiceImpl{rdb}
		}

		log.Printf("Redis not ready (attempt %d/%d): %v\n", attempt, maxRetries, err)
		time.Sleep(retryInterval)
	}

	log.Fatalln("Redis not ready after max retries")
	return nil
}

func (d *TTLStoreServiceImpl) Close() {
	d.rdb.Close()
}

func (d *TTLStoreServiceImpl) AddActivity(id int64) error {
	err := d.rdb.Set(ctx, "dedupeactivity:"+strconv.FormatInt(id, 10), "", 10*time.Minute).Err()
	if err != nil {
		log.Printf("Failed to set key: %v\n", err)
		return errors.New("failed to set key")
	}

	return nil
}

func (d *TTLStoreServiceImpl) DedupeActivity(id int64) (bool, error) {
	exists, err := d.rdb.Exists(ctx, "dedupeactivity:"+strconv.FormatInt(id, 10)).Result()
	if err != nil {
		log.Printf("Failed to check if key exists: %v\n", err)
		return false, errors.New("failed to check if key exists")
	}

	return exists > 0, nil
}

func (d *TTLStoreServiceImpl) CheckRateLimit(id int64, plan strava.UserPlan) (bool, error) {
	if plan == strava.UserPlanNone {
		return false, nil
	}

	var dailyLimit int
	var monthlyLimit int
	switch plan {
	case strava.UserPlanFree:
		dailyLimit, _ = strconv.Atoi(os.Getenv("LIMIT_FREE_DAILY"))
		monthlyLimit, _ = strconv.Atoi(os.Getenv("LIMIT_FREE_MONTHLY"))
	case strava.UserPlanPro:
		dailyLimit, _ = strconv.Atoi(os.Getenv("LIMIT_PRO_DAILY"))
		monthlyLimit, _ = strconv.Atoi(os.Getenv("LIMIT_PRO_MONTHLY"))
	default:
		return false, errors.New("invalid plan in limit check")
	}

	dailyCount, err := d.rdb.Get(ctx, "ratelimit:"+strconv.FormatInt(id, 10)+":daily").Int()
	if err != nil && err != redis.Nil {
		log.Printf("Failed to get key: %v\n", err)
		return false, errors.New("failed to get key")
	}

	if dailyCount >= dailyLimit {
		return true, nil
	}

	monthlyCount, err := d.rdb.Get(ctx, "ratelimit:"+strconv.FormatInt(id, 10)+":monthly").Int()
	if err != nil && err != redis.Nil {
		log.Printf("Failed to get key: %v\n", err)
		return false, errors.New("failed to get key")
	}

	if monthlyCount >= monthlyLimit {
		return true, nil
	}

	return false, nil
}

func (d *TTLStoreServiceImpl) IncrementRateLimit(id int64) error {
	dailyKey := "ratelimit:" + strconv.FormatInt(id, 10) + ":daily"
	monthlyKey := "ratelimit:" + strconv.FormatInt(id, 10) + ":monthly"

	dailyCount, err := d.rdb.Incr(ctx, dailyKey).Result()
	if err != nil {
		log.Printf("Failed to increment daily limit: %v\n", err)
		return errors.New("failed to increment daily limit")
	}

	monthlyCount, err := d.rdb.Incr(ctx, monthlyKey).Result()
	if err != nil {
		log.Printf("Failed to increment monthly limit: %v\n", err)
		return errors.New("failed to increment monthly limit")
	}

	if dailyCount == 1 {
		err = d.rdb.Expire(ctx, dailyKey, 24*time.Hour).Err()
		if err != nil {
			log.Printf("Failed to set expiration for daily limit: %v\n", err)
			return errors.New("failed to set expiration for daily limit")
		}
	}

	if monthlyCount == 1 {
		err = d.rdb.Expire(ctx, monthlyKey, 30*24*time.Hour).Err()
		if err != nil {
			log.Printf("Failed to set expiration for monthly limit: %v\n", err)
			return errors.New("failed to set expiration for monthly limit")
		}
	}

	return nil
}
