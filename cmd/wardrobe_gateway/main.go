package main

import (
	"fmt"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/wardrobe-gateway/internal/jwt"
	"github.com/wardrobe-gateway/internal/middleware"
	"github.com/wardrobe-gateway/internal/service/handler"
	"github.com/wardrobe-gateway/internal/service/user_service"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	dataServiceUrl = "http://data-service:8080"
	userServiceUrl = "http://user-service:8083"
)

func main() {
	jwtService := jwt.NewService("asd")
	userService := user_service.NewUserService(jwtService)
	handlerService := handler.NewHandler(userService, jwtService)
	mux := chi.NewRouter()
	mux.Group(func(mux chi.Router) {
		mux.Use(
			middleware.WithGoogleAuth(),
		)
		mux.MethodFunc("POST", "/save-data", serviceHandler(dataServiceUrl))
		mux.MethodFunc("GET", "/get-item", serviceHandler(dataServiceUrl))
		mux.MethodFunc("PUT", "/update-item", serviceHandler(dataServiceUrl))
		mux.MethodFunc("DELETE", "/delete-item", serviceHandler(dataServiceUrl))
		mux.MethodFunc("GET", "/get-user-items", serviceHandler(dataServiceUrl))

		mux.MethodFunc("POST", "/create-capsule", serviceHandler(dataServiceUrl))
		mux.MethodFunc("GET", "/get-capsule", serviceHandler(dataServiceUrl))
		mux.MethodFunc("GET", "/get-user-capsules", serviceHandler(dataServiceUrl))
		mux.MethodFunc("POST", "/update-capsule", serviceHandler(dataServiceUrl))
		mux.MethodFunc("DELETE", "/delete-capsule", serviceHandler(dataServiceUrl))
		mux.MethodFunc("POST", "/add-to-capsule", serviceHandler(dataServiceUrl))
		mux.MethodFunc("DELETE", "/delete-from-capsule", serviceHandler(dataServiceUrl))
		mux.MethodFunc("GET", "/get-all-capsule", serviceHandler(dataServiceUrl))

		mux.MethodFunc("POST", "/create-look", serviceHandler(dataServiceUrl))
		mux.MethodFunc("POST", "/add-to-look", serviceHandler(dataServiceUrl))
		mux.MethodFunc("DELETE", "/delete-look", serviceHandler(dataServiceUrl))
		mux.MethodFunc("GET", "/get-look-data", serviceHandler(dataServiceUrl))

	})
	mux.MethodFunc("POST", "/cut-item", serviceHandler(dataServiceUrl))
	mux.MethodFunc("PUT", "/update-user", serviceHandler(userServiceUrl))
	mux.MethodFunc("POST", "/save-user", serviceHandler(userServiceUrl))
	mux.MethodFunc("GET", "/get-user", serviceHandler(userServiceUrl))
	mux.MethodFunc("GET", "/get-users", serviceHandler(userServiceUrl))
	mux.MethodFunc("POST", "/save-stylist", serviceHandler(userServiceUrl))
	mux.MethodFunc("GET", "/get-stylist", serviceHandler(userServiceUrl))
	//mux.MethodFunc("GET", "/save-data", serviceHandler(userServiceUrl))
	mux.MethodFunc("GET", "/get-stylists", serviceHandler(userServiceUrl))
	mux.MethodFunc("GET", "/get-user-stylists", serviceHandler(userServiceUrl))
	mux.MethodFunc("GET", "/get-stylist-users", serviceHandler(userServiceUrl))
	mux.MethodFunc("POST", "/add-stylist", serviceHandler(userServiceUrl))
	mux.MethodFunc("DELETE", "/remove-stylist", serviceHandler(userServiceUrl))
	mux.MethodFunc("POST", "/create-room", serviceHandler(userServiceUrl))
	mux.MethodFunc("POST", "/save-message", serviceHandler(userServiceUrl))
	mux.MethodFunc("GET", "/get-user-rooms", serviceHandler(userServiceUrl))
	mux.MethodFunc("GET", "/get-stylist-rooms", serviceHandler(userServiceUrl))
	mux.MethodFunc("GET", "/get-chat-messages", serviceHandler(userServiceUrl))

	mux.HandleFunc("/login", handlerService.LoginUser)
	fmt.Println("starting server on :8082")
	httpServer := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", 8082),
		Handler: mux,
	}

	fmt.Printf("listening to http://0.0.0.0:%d/ for debug http", 8082)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("failed to listen on port 8082: %v", err)
	}
}

func serviceHandler(targetURL string) http.HandlerFunc {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = target.Host

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	}
}
