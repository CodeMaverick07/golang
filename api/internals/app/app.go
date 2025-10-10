package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codemaverick07/api/internals/api"
	"github.com/codemaverick07/api/internals/store"
	"github.com/codemaverick07/api/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	UserStore := store.NewPostgresUserStore(pgDB)
	TokenStore := store.NewPostgresTokenStore(pgDB)
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	UserHandler := api.NewUserHandler(UserStore, logger)
	TokenHandler := api.NewTokenHandler(TokenStore, UserStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    UserHandler,
		TokenHandler:   TokenHandler,
		DB:             pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "server is running correctly")
}
