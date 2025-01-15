package handlers

import (
	"go-clickhouse-example/services"
)

type ItemHandler struct {
	DBService   *services.DBService
	NATSService *services.NATSService
}

func NewItemHandler(dbService *services.DBService, natsService *services.NATSService) *ItemHandler {
	return &ItemHandler{DBService: dbService, NATSService: natsService}
}
