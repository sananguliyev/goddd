package presenter

import (
	"context"
	"flag"
	"fmt"
	"github.com/SananGuliyev/goddd/application/handler"
	"github.com/SananGuliyev/goddd/application/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type httpServer struct {
	router *mux.Router
}

func (s *httpServer) Run() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "gracefully waiting for live connections")
	flag.Parse()

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "80"
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handlers.CombinedLoggingHandler(os.Stdout, s.router),
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait until the timeout deadline.
	go func() {
		if err := srv.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()
	<-ctx.Done()

	log.Println("Shutting down")
	os.Exit(0)
}

func getRouter(recipeHandler handler.RecipeHandler) *mux.Router {
	security := middleware.NewSecurityMiddleware()

	router := mux.NewRouter()

	recipes := router.PathPrefix("/recipes").Subrouter()
	protectedRoutes := recipes.PathPrefix("").Subrouter()
	publicRoutes := recipes.PathPrefix("").Subrouter()

	// Protected routes
	protectedRoutes.Use(security.Auth)
	protectedRoutes.HandleFunc("", recipeHandler.Create).Methods(http.MethodPost)
	protectedRoutes.HandleFunc("/{id}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		recipeHandler.Update(writer, request, vars["id"])
	}).Methods(http.MethodPut, http.MethodPatch)
	protectedRoutes.HandleFunc("/{id}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		recipeHandler.Delete(writer, request, vars["id"])
	}).Methods(http.MethodDelete)

	// Public routes
	publicRoutes.HandleFunc("", recipeHandler.List).Methods(http.MethodGet)
	publicRoutes.HandleFunc("/{id}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		recipeHandler.Get(writer, request, vars["id"])
	}).Methods(http.MethodGet)
	publicRoutes.HandleFunc("/{id}/rating", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		recipeHandler.Rate(writer, request, vars["id"])
	}).Methods(http.MethodPost)

	// Health check
	router.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "OK")
	}).Methods(http.MethodGet)

	return router
}

func NewHttpServer(recipeHandler handler.RecipeHandler) Server {
	router := getRouter(recipeHandler)

	return &httpServer{router: router}
}
