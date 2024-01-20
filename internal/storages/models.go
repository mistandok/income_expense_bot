package storages

type Category struct {
	ID                int64             `json:"id"`
	Name              string            `json:"name"`
	UserID            int64             `json:"user_id"`
	MoneyMovementType MoneyMovementType `json:"money_movement_type"`
}

type MoneyMovementType string

const (
	Income  MoneyMovementType = "income"
	Expense MoneyMovementType = "expense"
)
