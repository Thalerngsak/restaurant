package service

import (
	"errors"
	"fmt"
	"github.com/Thalerngsak/restaurant/model"
	"github.com/Thalerngsak/restaurant/repository"
	"math"
	"os"
	"strconv"
	"sync"
)

var (
	restaurant = &model.Restaurant{}
)

type restaurantService struct {
	repo repository.RestaurantRepository
}

func NewRestaurantService(repo repository.RestaurantRepository) RestaurantService {
	return &restaurantService{repo: repo}
}

func (s *restaurantService) GetRestaurant() *model.Restaurant {

	if restaurant == nil {
		restaurant = &model.Restaurant{
			Tables:       []*model.Table{},
			Reservations: []*model.Reservation{},
			Mutex:        sync.Mutex{},
		}
	}
	return restaurant
}

func (s *restaurantService) FindAvailableTables(restaurant *model.Restaurant, numCustomers int) ([]int, error) {

	var tableIDs []int
	seatsPerTableStr := os.Getenv("SEAT_PER_TABLE")
	seatsPerTable, err := strconv.Atoi(seatsPerTableStr)
	if err != nil {
		fmt.Println("Error converting seatsPerTable to int:", err)
		return nil, errors.New("Not found seat per table setting")
	}

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

func (s *restaurantService) FindReservationByID(restaurant *model.Restaurant, bookingID int) (*model.Reservation, error) {
	for _, reservation := range restaurant.Reservations {
		if reservation.ID == bookingID {
			return reservation, nil
		}
	}
	return nil, errors.New("Booking ID not found")
}

func (s *restaurantService) RemoveReservationByID(reservations []*model.Reservation, bookingID int) []*model.Reservation {

	for i, reservation := range reservations {
		if reservation.ID == bookingID {
			return append(reservations[:i], reservations[i+1:]...)
		}
	}
	return reservations
}
