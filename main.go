package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	ctx    = context.Background()
	client *redis.Client
)

func init() {
	// read dotenv file
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("No .env file found")
	}
	// setup logging
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to info
	if !ok {
		lvl = "info"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.InfoLevel
	}
	// set global log level
	logrus.SetLevel(ll)
	logrus.Info("Logging initialized with level: ", lvl)
}

func main() {
	redisUrl, success := os.LookupEnv("REDIS_URL")
	if !success {
		logrus.Fatal("REDIS_URL not set")
		return
	}
	baseFilePath, success := os.LookupEnv("BASE_FILE_PATH")
	if !success {
		logrus.Fatal("BASE_FILE_PATH not set")
		return
	}
	client = redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})
	logrus.Info("Started worker...")
	for {
		time.Sleep(10 * time.Second)
		logrus.Info("Waiting for results of source generate-riddle-result...")
		jobResultRaw, err := client.RPop(ctx, "generate-riddle-result").Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			logrus.Errorf("Redis Error: %v", err)
			continue
		}

		var success JobSuccess
		err = json.Unmarshal([]byte(jobResultRaw), &success)

		if err != nil {
			logrus.Errorf("Error unmarshalling job result: %v", err)
			continue
		}

		duration := success.FinishedAt.Sub(success.StartedAt).Truncate(time.Millisecond)

		multipliedDuration := duration * time.Duration(success.ParallelCount)

		logrus.Infof("✅ Success for %s, actual duration: %v, duration for single thread run: %v", success.SuperSolution, duration, multipliedDuration)

		filename := fmt.Sprintf("%s/riddle-%s.json", baseFilePath, uuid.New().String())
		file, err := os.Create(filename)
		if err != nil {
			logrus.Errorf("❌ Failed to create file: %v", err)
			continue
		}

		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		if err := enc.Encode(success.Output); err != nil {
			file.Close()
			logrus.Errorf("❌ Failed to write to file: %v", err)
			continue
		}
		file.Close()

		logrus.Infof("📁 Result saved in %s", filename)

	}

}
