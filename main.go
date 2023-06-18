package main

import (
	"errors"
	"fmt"
	"github.com/Thalerngsak/restaurant/model"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"sync"
)

var (
	restaurant = &model.Restaurant{}
)

func main() {

	r := gin.Default()

	r.POST("/initialize", InitializeTables)

	r.POST("/reserve", ReserveTable)

	r.POST("/cancel", CancelReservation)

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Server error:", err)
	}

}

func InitializeTables(c *gin.Context) {

	restaurant = getRestaurant()

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

func ReserveTable(c *gin.Context) {
	restaurant = getRestaurant()

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

	tableIDs, err := findAvailableTables(restaurant, req.NumCustomers)
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

func CancelReservation(c *gin.Context) {
	restaurant = getRestaurant()
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

	reservation, err := findReservationByID(restaurant, req.BookingID)
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

	restaurant.Reservations = removeReservationByID(restaurant.Reservations, req.BookingID)

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"numFreedTables": len(reservation.TableIDs),
			"numRemaining":   len(restaurant.Tables) + len(reservation.TableIDs),
		},
	})
}

func findAvailableTables(restaurant *model.Restaurant, numCustomers int) ([]int, error) {

	var tableIDs []int
	seatsPerTable := 4
	numberOfTable := int(math.Ceil(float64(numCustomers) / float64(seatsPerTable)))

	for _, table := range restaurant.Tables {
		if !table.IsReserved {
			tableIDs = append(tableIDs, table.ID)
			if len(tableIDs) >= numberOfTable {
				break
			}
		}
	}

	if len(tableIDs) < numberOfTable {
		return nil, errors.New("Not enough tables for the reservation")
	}

	return tableIDs, nil
}

func findReservationByID(restaurant *model.Restaurant, bookingID int) (*model.Reservation, error) {
	for _, reservation := range restaurant.Reservations {
		if reservation.ID == bookingID {
			return reservation, nil
		}
	}
	return nil, errors.New("Booking ID not found")
}

func removeReservationByID(reservations []*model.Reservation, bookingID int) []*model.Reservation {

	for i, reservation := range reservations {
		if reservation.ID == bookingID {
			return append(reservations[:i], reservations[i+1:]...)
		}
	}
	return reservations
}

func getRestaurant() *model.Restaurant {

	if restaurant == nil {
		restaurant = &model.Restaurant{
			Tables:       []*model.Table{},
			Reservations: []*model.Reservation{},
			Mutex:        sync.Mutex{},
		}
	}
	return restaurant
}
