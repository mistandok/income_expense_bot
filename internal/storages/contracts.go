package storages

type CategoryStorage interface {
	GetCategoryByName(name string) (Category, error)
	DeleteCategoryByNameForUser(userID string, name string) error
	CreateCategory(name string) error
}
