package application

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/helf4ch/gocrudl/internal/config"
	"github.com/helf4ch/gocrudl/internal/store"
)

type Application struct {
	Log                *slog.Logger
	Config             *config.Config
	Store              *store.Store
	defaultMiddlewares []Middleware
	mux                *http.ServeMux
}

func New(log *slog.Logger, cfg *config.Config, str *store.Store) *Application {
	return &Application{
		Log:    log,
		Config: cfg,
		Store:  str,
		mux:    http.NewServeMux(),
	}
}

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (err *AppError) Error() string {
	return err.Err.Error()
}

type Handler func(Application, http.ResponseWriter, *http.Request) error

type Middleware func(Handler) Handler

func (app *Application) RegisterHandler(pattern string, handler Handler) {
	app.mux.HandleFunc(
		pattern,
		func(w http.ResponseWriter, r *http.Request) {
			handler(*app, w, r)
		},
	)
}

func (app *Application) SetDefaultMiddlewares(middlewares ...Middleware) {
	app.defaultMiddlewares = middlewares
}

func (app *Application) Handle(
	pattern string,
	handler Handler,
	middlewares ...Middleware,
) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	for i := len(app.defaultMiddlewares) - 1; i >= 0; i-- {
		handler = app.defaultMiddlewares[i](handler)
	}

	app.RegisterHandler(pattern, handler)
}

func (app Application) ListenAndServe() error {
	srv := &http.Server{
		Addr:    app.Config.Addr,
		Handler: app.mux,
	}
	app.Log.Info("server started")
	return srv.ListenAndServe()
}

type AppResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	Body  any    `json:"body,omitempty"`
}

func (app Application) Response(
	w http.ResponseWriter,
	ok bool,
	code int,
	error string,
	body any,
) error {
	r := AppResponse{
		Ok:    ok,
		Error: error,
		Body:  body,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		return fmt.Errorf("encoding response: %w", err)
	}

	return nil
}

// 200
func (app Application) ResponseOk(w http.ResponseWriter, body any) error {
	err := app.Response(w, true, http.StatusOK, "", body)
	if err != nil {
		return fmt.Errorf("response ok: %w", err)
	}

	return nil
}

// 201
func (app Application) ResponseCreated(w http.ResponseWriter, body any) error {
	err := app.Response(w, true, http.StatusCreated, "", body)
	if err != nil {
		return fmt.Errorf("response ok: %w", err)
	}

	return nil
}
