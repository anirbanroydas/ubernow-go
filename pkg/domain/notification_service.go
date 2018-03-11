package domain

// notificationService is serive which exposee the interface having a method send, which
// mean any type of notification servie can implement this interface, like email,
// sms, web notificaiton, etc..
type NotificationService interface {
	Send(*CabBookingResponse) error
}
