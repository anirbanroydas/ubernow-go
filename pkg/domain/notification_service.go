package domain

// NotificationService is serive which is an interface which exposee the a method Send, which
// takes a pointer to a CabBookingResponse as input and returns an error.
// It means any type of notification servie can implement this interface, like email,
// sms, web notificaiton, etc..
type NotificationService interface {
	Send(*CabBookingResponse) error
}
