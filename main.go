package main

import (
	"backend/adapters"
	"backend/constants"
	"backend/controllers"
	"backend/middlewares"
	"backend/repositories"
	"backend/routes"
	"backend/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	dbUser := os.Getenv(constants.DB_USER)
	dbPass := os.Getenv(constants.DB_PASSWORD)
	dbHost := os.Getenv(constants.DB_HOST)
	dbPort := os.Getenv(constants.DB_PORT)
	dbName := os.Getenv(constants.DB_NAME)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// repositories
	userRepo := repositories.NewUserRepository(db)

	// adapters
	calendarAdapter := adapters.NewGoogleAdapter()

	// services
	authService := services.NewAuthService(userRepo)
	calendarService := services.NewCalendarService(calendarAdapter)

	// controllers
	authController := controllers.NewAuthController(authService)
	CalendarController := controllers.NewCalendarController(calendarService)

	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
		AllowCredentials: true,
	}))

	routes.SetupAuthRoutes(e, authController)

	// routes
	calendarGroup := e.Group("/calendar",
		middlewares.JWTMiddleware(),
		middlewares.TokenRefreshMiddleware(authController.GoogleOAuthConfig, userRepo),
	)
	routes.SetupCalenderRoutes(calendarGroup, CalendarController)

	port := os.Getenv(constants.PORT)
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(e.Start(":" + port))
}
