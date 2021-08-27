package handler

import "github.com/PogpogPapaya/backend-api.git/pb"

type Handler struct {
	papayaServiceClient pb.PapayaServiceClient
}

func NewHandler(papayaServiceClient pb.PapayaServiceClient) *Handler {
	return &Handler{
		papayaServiceClient: papayaServiceClient,
	}
}
