package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"
	_ "swagger/proxy/docs"
)

//	@title						HugoProxyWithYandexGeoApi
//	@version					1.1
//	@description				test API server for hugoProxy
//	@host						localhost:8080
//	@basePath					/api
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization

func main() {

	r := chi.NewRouter()
	proxy := NewReverseProxy("hugo", "1313")
	r.Use(middleware.Logger, middleware.Recoverer)
	r.Use(proxy.ReverseProxy)
	r.Post("/api/register", registerHandler)
	r.Post("/api/login", loginHandler)
	r.Route("/api", func(r chi.Router) {
		r.Use(JwtAuthMiddleware)
		r.Post("/address/search", geoFromAddressHandler)
		r.Post("/address/geocode", addressFromGeoHandler)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json")))
	log.Fatal(http.ListenAndServe(":8080", r))
}
