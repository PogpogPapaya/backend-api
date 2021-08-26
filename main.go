package main

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/hashicorp/go-hclog"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

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

	app.Use(limiter.New(limiter.Config{
		Expiration: time.Second * 5,
		Max:        20,
	}))
	app.Use(cors.New(cors.Config{AllowMethods: "GET POST", AllowOrigins: "*", AllowHeaders: "Origin, Content-Type, Accept"}))

	app.Get("/api/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":   "pong",
			"timestamp": time.Now(),
		})
	})

	app.Post("/api/papaya/predict", func(c *fiber.Ctx) error {

		fileHeader, err := c.FormFile("image")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprint("cannot get file header from form-data", err.Error()))
		}

		imgFile, err := fileHeader.Open()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot open file header", err.Error()))
		}

		client := &http.Client{
			Timeout: time.Second * 10,
		}
		body := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(body)
		fw, _ := bodyWriter.CreateFormFile("file", "papaya")

		if _, err := io.Copy(fw, imgFile); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot write to file form", err.Error()))
		}
		if err := bodyWriter.Close(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot close bodyWriter", err.Error()))
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/predict", predictionApiHost), bytes.NewReader(body.Bytes()))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot create new request", err.Error()))
		}
		req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
		res, err := client.Do(req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot get prediction", err.Error()))
		}

		if _, err := io.Copy(c, res.Body); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot copy prediction result to response", err.Error()))
		}
		c.Set("Content-Type", "application/json")
		return nil

	})

	if err := app.Listen(":5000"); err != nil {
		logger.Error("Unable to start server on port 5000")
		os.Exit(1)
	}
}
