package models

import (
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type Address struct {
	ID      uuid.UUID `db:"id"`
	Country string    `db:"country"`
	City    string    `db:"city"`
	Street  string    `db:"street"`
}

type Client struct {
	ID               uuid.UUID `db:"id"`
	ClientName       string    `db:"client_name"`
	ClientSurname    string    `db:"client_surname"`
	Birthday         time.Time `db:"birthday"`
	Gender           Gender    `db:"gender"`
	RegistrationDate time.Time `db:"registration_date"`
	AddressID        uuid.UUID `db:"address_id"`
}

type ClientFull struct {
	Client
	Address Address
}

type Category struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type Supplier struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	AddressID   uuid.UUID `db:"address_id"`
	PhoneNumber string    `db:"phone_number"`
}

type SupplierFull struct {
	Supplier
	Address Address
}

type Image struct {
	ID    uuid.UUID `db:"id"`
	Image []byte    `db:"image"`
}

type Product struct {
	ID             uuid.UUID  `db:"id"`
	Name           string     `db:"name"`
	CategoryID     uuid.UUID  `db:"category_id"`
	Price          float64    `db:"price"`
	AvailableStock int        `db:"available_stock"`
	LastUpdateDate time.Time  `db:"last_update_date"`
	SupplierID     uuid.UUID  `db:"supplier_id"`
	ImageID        *uuid.UUID `db:"image_id"`
}

type ProductFull struct {
	Product
	CategoryName string
}
