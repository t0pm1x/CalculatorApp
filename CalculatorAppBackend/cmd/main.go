package main

import (
	"CalculatorAppBackend/internal/calculationService"
	"CalculatorAppBackend/internal/db"
	"CalculatorAppBackend/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to DB %v", err)
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
	}))

	calcRepo := calculationService.NewCalculationRepository(database)
	calcService := calculationService.NewCalculationService(calcRepo)
	calcHandlers := handlers.NewCalculationHandler(calcService)

	e.Use(middleware.Logger())

	e.GET("/calculations", calcHandlers.GetCalculation)
	e.POST("/calculations", calcHandlers.PostCalculation)
	e.PATCH("/calculations/:id", calcHandlers.PatchCalculation)
	e.DELETE("/calculations/:id", calcHandlers.DeleteCalculation)

	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
