package models

import "time"


//Model which Maps the Users table
type User struct {
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
type Room struct {
	ID int
	RoomName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Restrictions model
type Restriction struct {
	ID int
	RestrictionName string
	CreatedAt time.Time
	UpdatedAt time.Time
}
//Reservations model
type Reservation struct {
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
	Room Room
}

type RoomRestriction struct {
	ID int
	StartDate time.Time
	EndDate time.Time
	RoomID int
	ReservationID int
	RestrictionID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Room
	Reservation Reservation
	Restriction Restriction
}