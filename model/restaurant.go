package model

import "sync"

type Table struct {
	ID         int
	IsReserved bool
}

type Reservation struct {
	ID           int
	TableIDs     []int
	NumCustomers int
}

type Restaurant struct {
	Tables       []*Table
	Reservations []*Reservation
	Mutex        sync.Mutex
}

type InitializeRequest struct {
	NumTables int `json:"numTables"`
}

type ReserveRequest struct {
	NumCustomers int `json:"numCustomers"`
}
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
