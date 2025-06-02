package main

import (
	"github.com/ademaxweb/mfa-go-core/pkg/cfg"
	"github.com/ademaxweb/mfa-go-core/pkg/handler"
	"github.com/ademaxweb/mfa-go-core/pkg/srv"
	"log"
	"os"
	"time"
	"users/pkg/api"
	"users/pkg/db"
)

func main() {
	h := handler.New()
	h.HealthRoute()

	dbConnStr := cfg.GetStringEnv("DB_CONNECTION", "")

	d, err := db.NewPgsqlDB(dbConnStr)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	api.New(os.Stdout, d).RegisterRoutes(h)

	service := srv.New(srv.Config{
		Port:    cfg.GetIntEnv("SERVICE_PORT", 80),
		Timeout: time.Duration(cfg.GetIntEnv("SERVICE_TIMEOUT", 5)) * time.Second,
		Handler: h,
		Writer:  os.Stdout,
	})

	err = service.Start()
	if err != nil {
		log.Fatal(err)
	}
}
