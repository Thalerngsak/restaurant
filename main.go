package main

import (
	"github.com/Thalerngsak/restaurant/handler"
	"github.com/Thalerngsak/restaurant/middleware"
	"github.com/Thalerngsak/restaurant/repository"
	"github.com/Thalerngsak/restaurant/service"
	"github.com/Thalerngsak/restaurant/token"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	// Load environment variables from .env file
	if err := godotenv.Load("local.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	tokenMaker := token.NewJWTMaker(os.Getenv("JWT_SECRET_KEY"))

	userStore := repository.NewUserDB()
	userService := service.NewUserService(userStore)
	userHandler := handler.NewUserHandler(userService, tokenMaker)

	restaurantStore := repository.NewRestaurantDB()
	restaurantService := service.NewRestaurantService(restaurantStore)
	restaurantHandler := handler.NewRestaurantHandler(restaurantService)

	r := gin.New()
	r.Use(gin.Logger())
	r.POST("/api/login", userHandler.Login)

	api := r.Group("/api")
	api.Use(middleware.AuthenticationMiddleware(tokenMaker))
	api.POST("/initialize", restaurantHandler.InitializeTables)

	api.POST("/reserve", restaurantHandler.ReserveTable)

	api.POST("/cancel", restaurantHandler.CancelReservation)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}

}
