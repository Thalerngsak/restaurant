package handler

import (
	"github.com/Thalerngsak/restaurant/model"
	"github.com/Thalerngsak/restaurant/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type restaurantHandler struct {
	service service.RestaurantService
}

func NewRestaurantHandler(service service.RestaurantService) *restaurantHandler {
	return &restaurantHandler{
		service: service,
	}
}

func (h *restaurantHandler) InitializeTables(c *gin.Context) {

	restaurant := h.service.GetRestaurant()

	restaurant.Mutex.Lock()

	defer restaurant.Mutex.Unlock()

	if len(restaurant.Tables) > 0 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Tables are already initialized",
		})
		return
	}

	var req model.InitializeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	if req.NumTables <= 0 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Number of tables should be greater than 0",
		})
		return
	}

	for i := 0; i < req.NumTables; i++ {
		table := &model.Table{
			ID:         i + 1,
			IsReserved: false,
		}
		restaurant.Tables = append(restaurant.Tables, table)
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Tables initialized successfully",
	})
}

func (h *restaurantHandler) ReserveTable(c *gin.Context) {
	restaurant := h.service.GetRestaurant()

	restaurant.Mutex.Lock()
	defer restaurant.Mutex.Unlock()

	if len(restaurant.Tables) == 0 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Tables are not initialized",
		})
		return
	}
	var req = model.ReserveRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	if req.NumCustomers <= 0 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Number of customers should be greater than 0",
		})
		return
	}

	tableIDs, err := h.service.FindAvailableTables(restaurant, req.NumCustomers)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	reservation := &model.Reservation{
		ID:           len(restaurant.Reservations) + 1,
		TableIDs:     tableIDs,
		NumCustomers: req.NumCustomers,
	}

	for _, tableID := range tableIDs {
		restaurant.Tables[tableID-1].IsReserved = true
	}

	restaurant.Reservations = append(restaurant.Reservations, reservation)

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"bookingID":       reservation.ID,
			"numBookedTables": len(tableIDs),
			"numRemaining":    len(restaurant.Tables) - len(tableIDs),
		},
	})
}

func (h *restaurantHandler) CancelReservation(c *gin.Context) {
	restaurant := h.service.GetRestaurant()
	restaurant.Mutex.Lock()
	defer restaurant.Mutex.Unlock()

	if len(restaurant.Tables) == 0 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Tables are not initialized",
		})
		return
	}

	var req = model.CancelRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	reservation, err := h.service.FindReservationByID(restaurant, req.BookingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	for _, tableID := range reservation.TableIDs {
		restaurant.Tables[tableID-1].IsReserved = false
	}

	restaurant.Reservations = h.service.RemoveReservationByID(restaurant.Reservations, req.BookingID)

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"numFreedTables": len(reservation.TableIDs),
			"numRemaining":   len(restaurant.Tables) + len(reservation.TableIDs),
		},
	})
}
