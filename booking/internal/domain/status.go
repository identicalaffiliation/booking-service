package domain

const (
	Active    BookingStatus = "active"
	Cancelled BookingStatus = "cancelled"
)

type BookingStatus string

const (
	Admin  UserRole = "admin"
	Client UserRole = "client"
)

type UserRole string
