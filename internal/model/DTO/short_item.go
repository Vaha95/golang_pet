package DTO

// ShortItem holds the resolved data for a single short URL.
type ShortItem struct {
	URL       string  `db:"url"`
	Short     string  `db:"short"`
	DeletedAt *string `db:"deleted_at"`
}
