package DTO

// DeleteBatch is the message sent to the async deletion worker.
type DeleteBatch struct {
	UserId int
	Shorts []string
}
