package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// RipenessPredictHandler predict ripeness
// @Summary get papaya ripeness prediction
// @Tags papaya
// @Accept multipart/form-data
// @Produce  json
// @Success 200 {object} handler.PredictionResponseDto "the prediction confidence and classification"
// @Failure 404 {object} main.ErrorMessage "Multipart-form data error"
// @Failure 413 {object} main.ErrorMessage "Payload too large (10MB limit)"
// @Failure 500 {object} main.ErrorMessage "Internal Server Error"
// @Router /api/papaya/predict [post]
func (h *Handler) RipenessPredictHandler(c *fiber.Ctx) error {
	// Get file from FormFile
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprint("cannot get file header from form-data", err.Error()))
	}

	// Open File
	imgFile, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot open file header", err.Error()))
	}

	// Create client and create form-data
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

	// Create request and send the request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/predict", h.predictionApiHost), bytes.NewReader(body.Bytes()))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot create new request", err.Error()))
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot get prediction", err.Error()))
	}

	// Convert prediction res body to PredictionResponseDto
	var predictionRes PredictionResponseDto
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&predictionRes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot parse prediction response to struct", err.Error()))
	}
	if err := res.Body.Close(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot close response body", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(predictionRes)

}

type PredictionResponseDto struct {
	Classification string `json:"classification"`
	Confidence     string `json:"confidence"`
}
