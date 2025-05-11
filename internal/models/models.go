package models

import "time"

// Holds reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

//Model which Maps the Users table
type Users struct {
	ID int
	FirstName string
	LastName string
	Email string
	Password string
	AccessLevl int
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Model which Maps the Roms table
type Rooms struct {
	ID int
	RoomName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Restrictions model
type Restrictions struct {
	ID int
	RestrictionName string
	CreatedAt time.Time
	UpdatedAt time.Time
}
//Reservations model
type Reservations struct {
	ID int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate time.Time
	RoomID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Rooms
}

type RoomRestrictions struct {
	ID int
	StartDate time.Time
	EndDate time.Time
	RoomID int
	ReservationId int
	RestrictionId int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Rooms
	Reservation Reservations
	Restriction Restrictions
}