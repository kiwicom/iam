package api

// Error used for errors that will be returned to the client.
type Error struct {
	Message string
	Code    int
}

func (err Error) Error() string {
	return err.Message
}
