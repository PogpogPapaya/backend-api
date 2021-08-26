package handler

type Handler struct {
	predictionApiHost string
}

func NewHandler(predictionApiHost string) *Handler {
	return &Handler{predictionApiHost: predictionApiHost}
}
