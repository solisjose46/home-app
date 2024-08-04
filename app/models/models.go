package models

import "time"

type ServerResponse struct {
	Message        string
	ReturnEndpoint string
}

type FinanceFeedEdit struct {
	Expense *Expense
}

type FinanceFeedConfirm struct {
	OldExpense *Expense
	NewExpense *Expense
}

type FinanceTrackConfirm struct {
	Expense *Expense
}
type FinanceFeed struct {
	ServerResponse    *ServerResponse
	FinanceFeedEdit   *FinanceFeedEdit
	FinanceFeedConfirm *FinanceFeedConfirm
	Expenses          *[]Expense
}

type FinanceTrack struct {
	Month           string
	Categories      *[]Category
	ServerResponse  *ServerResponse
	FinanceTrackConfirm *FinanceTrackConfirm
}

type Finance struct {
	FinanceTrack *FinanceTrack
	FinanceFeed *FinanceFeed
}

type Login struct {
	ServerResponse *ServerResponse
}

type Expense struct {
	ExpenseId string
	Name      string
	Amount    float64
	Category  string
	Username  string
	UserId    string
	Datetime  time.Time
	IsOwner   bool
}

type Category struct {
	Name    string
	Balance float64
	Limit   float64
}