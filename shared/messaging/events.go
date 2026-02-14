package messaging

import (
	pbt "go-ride/shared/proto/trip"
)

const (
	FindAvailableDriversQueue = "find_available_drivers"
	NotifyNoDriversFoundQueue = "notify_no_drivers_found"
	NotifyDriverAssignQueue   = "notify_driver_assign_queue"
	DeadLetterQueue           = "dead_letter_queue"
)

type TripEventData struct {
	Trip *pbt.Trip `json:"trip"`
}
