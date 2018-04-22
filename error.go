package gdax

// Error to handle GDAX errors
type Error struct {
	Message string `json:"message"`
}

// Error function to show the message
func (e Error) Error() string {
	return e.Message
}
