package templates

import (
	"fmt"
	"time"
)

const (
	financeConfirmCreate = "Are you sure you want to create this expense?"
	financeConfirmUpdate = "Are you sure you want to update this expense?"
)

type serverResponse struct {
	Message string
	IsSuccess bool
}

type expense struct {
	Name string
	User string
	Amount float64
	Date time.Time
}

type category struct {
	Name string
	Limit float64
	Balance float64
}

type financeTrack struct {
	Month string 
	Categories []category
}

type financeFeed struct {
	Expenses []expense
}