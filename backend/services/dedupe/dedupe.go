package dedupe

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"titles.run/services/interfaces"

	"github.com/redis/go-redis/v9"
)

type DedupeServiceImpl struct {
	rdb *redis.Client
}

var ctx = context.Background()

func NewDedupeService() interfaces.DedupeService {
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
		if err != nil {
			log.Fatalf("Could not connect to Redis: %v", err)
		} else {
			return &DedupeServiceImpl{rdb}
		}

		log.Printf("Redis not ready (attempt %d/%d): %v\n", attempt, maxRetries, err)
		time.Sleep(retryInterval)
	}

	log.Fatalln("Database not ready after max retries")
	return nil
}

func (d *DedupeServiceImpl) Close() {
	d.rdb.Close()
}

func (d *DedupeServiceImpl) AddActivity(id int64) error {
	err := d.rdb.Set(ctx, strconv.FormatInt(id, 10), "", 10*time.Minute).Err()
	if err != nil {
		log.Printf("Failed to set key: %v\n", err)
		return errors.New("failed to set key")
	}

	return nil
}

func (d *DedupeServiceImpl) DedupeActivity(id int64) (bool, error) {
	exists, err := d.rdb.Exists(ctx, strconv.FormatInt(id, 10)).Result()
	if err != nil {
		log.Printf("Failed to check if key exists: %v\n", err)
		return false, errors.New("failed to check if key exists")
	}

	return exists > 0, nil
}
