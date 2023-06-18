package service

import "github.com/Thalerngsak/restaurant/model"

type RestaurantService interface {
	GetRestaurant() *model.Restaurant
	FindAvailableTables(restaurant *model.Restaurant, numCustomers int) ([]int, error)
	FindReservationByID(restaurant *model.Restaurant, bookingID int) (*model.Reservation, error)
	RemoveReservationByID(reservations []*model.Reservation, bookingID int) []*model.Reservation
}
