package ws

import (
	"context"
	"encoding/json"
	"go-ride/shared/messaging"
	pd "go-ride/shared/proto/driver"
	"go-ride/shared/types"
	"log"
	"net/http"
	"time"
)

type DriverWSHandler struct {
	connManager  *messaging.ConnectionManager
	driverClient pd.DriverServiceClient
}

func NewDriverWSHandler(cm *messaging.ConnectionManager, dc pd.DriverServiceClient) *DriverWSHandler {
	return &DriverWSHandler{
		connManager:  cm,
		driverClient: dc,
	}
}

func (h *DriverWSHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := h.connManager.Upgrade(w, r)
	if err != nil {
		log.Printf("[WS] Falha ao fazer upgrade da conexão: %v", err)
		return
	}
	defer conn.Close()

	h.connManager.Add(userID, conn)
	log.Printf("[WS] driver %s connected", userID)

	defer func() {
		h.connManager.Remove(userID)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := h.driverClient.UpdateStatus(ctx, &pd.UpdateStatusRequest{
			DriverID: userID,
			Status:   pd.DriverStatusType_OFFLINE,
		})
		if err != nil {
			log.Printf("[WS] error while putting the driver with ID %s OFFLINE: %v", userID, err)
		} else {
			log.Printf("[WS] driver with ID %s disconnected (OFFLINE)", userID)
		}
	}()

	_, err = h.driverClient.UpdateStatus(r.Context(), &pd.UpdateStatusRequest{
		DriverID: userID,
		Status:   pd.DriverStatusType_ONLINE,
	})
	if err != nil {
		log.Printf("[WS] error while putting the driver with ID %s as ONLINE: %v", userID, err)
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var payload types.Coordinate
		if err := json.Unmarshal(message, &payload); err != nil {
			continue
		}

		_, err = h.driverClient.UpdateStatus(r.Context(), &pd.UpdateStatusRequest{
			DriverID: userID,
			Status:   pd.DriverStatusType_ONLINE,
			ActualLocation: &pd.Coordinate{
				Latitude:  payload.Latitude,
				Longitude: payload.Longitude,
			},
		})
		if err != nil {
			log.Printf("[WS] failed to update driver location: %v", err)
		}
	}
}
