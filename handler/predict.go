package handler

import (
	"context"
	"fmt"
	"github.com/PogpogPapaya/backend-api.git/pb"
	"github.com/PogpogPapaya/backend-api.git/utils"
	"github.com/gofiber/fiber/v2"
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

	buf, err := utils.ResizeImage(imgFile, 150, 150)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot resize image", err.Error()))
	}

	req := &pb.PredictionRequest{Image: *buf}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := h.papayaServiceClient.Predict(ctx, req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("cannot get prediction", err.Error()))
	}

	predictionRes := PredictionResponseDto{
		Classification: res.GetLabel(),
		Confidence:     fmt.Sprintf("%v", res.GetConfidence()),
	}

	return c.Status(fiber.StatusOK).JSON(predictionRes)

}

type PredictionResponseDto struct {
	Classification string `json:"classification"`
	Confidence     string `json:"confidence"`
}
