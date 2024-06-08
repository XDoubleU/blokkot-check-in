package main

import (
	_ "time/tzdata"

	"check-in/api/internal/config"
	"check-in/api/internal/services"

	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/XDoubleU/essentia/pkg/logger"
)

type application struct {
	config   config.Config
	services services.Services
}

//	@title			Check-In API
//	@version		1.0
//	@license.name	GPL-3.0
//	@Accept			json
//	@Produce		json

func main() {
	cfg := config.New()

	db, err := postgres.Connect(
		cfg.DB.Dsn,
		cfg.DB.MaxConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	spandb := postgres.SpanDB{
		DB: db,
	}

	logger.GetLogger().Printf("connected to database")

	app := &application{
		config:   cfg,
		services: services.New(spandb),
	}

	app.config.Print()

	err = http_tools.Serve(app.config.Port, app.routes(), app.config.Env)
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
}
