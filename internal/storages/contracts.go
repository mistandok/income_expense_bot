package storages

type CategoryStorage interface {
	DeleteCategoryByNameAndTypeForUser(userName string, name string, moneyMovementType MoneyMovementType) error
	CreateCategoryForUserByType(userName string, name string, moneyMovementType MoneyMovementType) error
	GetCategoriesForUserByTYpe(userName string, monetMovementType MoneyMovementType) ([]Category, error)
}
