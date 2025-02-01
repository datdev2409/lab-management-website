package main

import (
	"github.com/datdev2409/lab-admin-go/internal/storage"
	"net/http"
)

type DBConfig struct {
	Addr string
	Port int
}

type Config struct {
	Env  string
	Port string
	DB   *DBConfig
}

type Application struct {
	Config  *Config
	Store   storage.AppStorage
	Handler http.Handler
}

func (app *Application) Init(config *Config, store storage.AppStorage, handler http.Handler) {
	app.Config = config
	app.Store = store
	app.Handler = handler
}

func (app *Application) Start() error {
	err := http.ListenAndServe(app.Config.Port, app.Handler)
	return err
}
