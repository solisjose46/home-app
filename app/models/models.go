package models

import "time"

const (
	CategorySdo = "sdo"
	SdoName = "Seattle Dine out"
	CategorySg = "sg"
	SgName = "Seattle Groceries"
	CategoryPdo = "pdo"
	PdoName = "Phx Dine Out"
	CategorySc = "sc"
	ScName = "Seattle Cleaning"
	CategoryPg = "pg"
	PgName = "Phx Groceries"
)

type ServerResponse struct {
	Message        string
	ReturnEndpoint string
}

type Expense struct {
	ExpenseId string
	Name      string
	Amount    float64
	Category  string
	User      string
	UserId    string
	Datetime  time.Time
	IsOwner   bool
}

type Category struct {
	Name    string
	Balance float64
	Limit   float64
}

type Login struct {
	ServerResponse ServerResponse
}

type FinanceTrack struct {
	Month           string
	Categories      []Category
	ServerResponse  ServerResponse
	FinanceTrackConfirm FinanceTrackConfirm
}

type FinanceFeed struct {
	ServerResponse    ServerResponse
	FinanceFeedEdit   FinanceFeedEdit
	FinanceFeedConfirm FinanceFeedConfirm
	Expenses          []Expense
}

type Finance struct {
	FinanceTrack FinanceTrack
	FinanceFeed FinanceFeed
}

type FinanceTrackConfirm struct {
	Expense Expense
}

type FinanceFeedEdit struct {
	Expense Expense
}

type FinanceFeedConfirm struct {
	OldExpense Expense
	NewExpense Expense
}