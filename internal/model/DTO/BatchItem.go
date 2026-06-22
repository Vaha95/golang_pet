package DTO

// BatchItem represents a single URL entry in a batch save request.
type BatchItem struct {
	ExtId string
	URL   string
	Short string
}
