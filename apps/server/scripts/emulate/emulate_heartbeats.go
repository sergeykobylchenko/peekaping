package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"peekaping/src/config"
	"peekaping/src/modules/events"
	"peekaping/src/modules/healthcheck"
	"peekaping/src/modules/heartbeat"
	"peekaping/src/modules/monitor"
	"peekaping/src/modules/monitor_notification"
	"peekaping/src/modules/stats"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func provideConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	return &cfg, nil
}

func provideLogger(cfg *config.Config) (*zap.SugaredLogger, error) {
	zapConfig := zap.NewProductionConfig()
	var logger *zap.Logger
	var err error

	if cfg.Mode == "prod" {
		logger, err = zapConfig.Build()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

func provideMongoDB(cfg *config.Config) (*mongo.Client, error) {
	DBUri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
	)
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(DBUri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}

func main() {
	container := dig.New()

	container.Provide(provideConfig)
	container.Provide(provideLogger)
	container.Provide(provideMongoDB)
	events.RegisterDependencies(container)
	monitor_notification.RegisterDependencies(container)
	stats.RegisterDependencies(container)
	healthcheck.RegisterDependencies(container)
	monitor.RegisterDependencies(container)
	heartbeat.RegisterDependencies(container)

	err := container.Invoke(func(monitorService monitor.Service, heartbeatService heartbeat.Service) {
		ctx := context.Background()

		monitorID := "6840f5d0b9bff35e4ced00a3"

		start := time.Now().Add(-7 * 24 * time.Hour)
		end := time.Now()
		interval := 20 * time.Second
		n := int(end.Sub(start) / interval)
		status := heartbeat.MonitorStatusUp
		downCount := int((time.Minute * 10) / interval)
		for i := 0; i < n; i++ {
			t := start.Add(time.Duration(i) * interval)

			if downCount > 0 {
				downCount--
			} else {
				status = heartbeat.MonitorStatusUp
				if rand.Intn(10000) < 1 {
					status = heartbeat.MonitorStatusDown
					downCount = int((time.Minute * 10) / interval)
				}
			}

			ping := 100 + rand.Intn(100)
			if status == heartbeat.MonitorStatusDown {
				ping = 0
			}

			dto := &heartbeat.CreateUpdateDto{
				MonitorID: monitorID,
				Status:    status,
				Msg:       "emulated",
				Ping:      ping,
				Duration:  0,
				DownCount: 0,
				Retries:   0,
				Important: false,
				Time:      t,
				EndTime:   t.Add(time.Duration(ping) * time.Millisecond),
				Notified:  false,
			}
			_, err := heartbeatService.Create(ctx, dto)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create heartbeat at %s: %v\n", t.Format(time.RFC3339), err)
			} else if i%1000 == 0 {
				fmt.Printf("Created %d/%d heartbeats...\n", i, n)
			}
		}
		fmt.Println("Done!")
	})
	if err != nil {
		log.Fatal(err)
	}
}
