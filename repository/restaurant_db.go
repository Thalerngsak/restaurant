package repository

type restaurantRepository struct {
}

func NewRestaurantDB() RestaurantRepository {
	return &restaurantRepository{}
}
