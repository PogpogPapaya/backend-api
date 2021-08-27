package main

import (
	_ "github.com/PogpogPapaya/backend-api.git/docs"
	"github.com/PogpogPapaya/backend-api.git/handler"
	"github.com/PogpogPapaya/backend-api.git/pb"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"os"
	"time"
)

// @title Papaya Ripeness Prediction API
// @version 1.0
// @description This is a sample of papaya ripeness prediction api for CSC340
//
// @license.name MIT
//
// @host https://papaya.cscms.me
// @BasePath /

func main() {
	logger := hclog.Default()

	predictionApiHost := os.Getenv("PREDICTION_API_HOST")
	if len(predictionApiHost) == 0 {
		logger.Error("PREDICTION_API_HOST is required")
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 10 << 20, // 10MB
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := err.Error()
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				msg = e.Message
			}
			logger.Error(msg)

			return c.Status(code).JSON(fiber.Map{"message": msg})
		},
	})

	// Create gRPC client
	conn, err := grpc.Dial(predictionApiHost, grpc.WithInsecure())
	if err != nil {
		logger.Error("unable to dia localhost:8000 for gRPC", err)
		os.Exit(1)
	}
	papayaServiceClient := pb.NewPapayaServiceClient(conn)

	// Create handler
	handlers := handler.NewHandler(papayaServiceClient)

	app.Use(limiter.New(limiter.Config{
		Expiration: time.Second * 5,
		Max:        20,
	}))
	app.Use(cors.New(cors.Config{AllowMethods: "GET POST", AllowOrigins: "*"}))

	app.Get("/swagger/*", swagger.Handler) // default

	app.Get("/api/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":   "pong",
			"timestamp": time.Now(),
		})
	})

	app.Post("/api/papaya/predict", handlers.RipenessPredictHandler)

	if err := app.Listen(":5000"); err != nil {
		logger.Error("Unable to start server on port 5000")
		os.Exit(1)
	}
}

type ErrorMessage struct {
	Message string `json:"message"`
}
