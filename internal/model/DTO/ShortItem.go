package DTO

type ShortItem struct {
	URL string `db:"url"`
	Short string `db:"short"`
	DeletedAt *string `db:"deleted_at"`
}
