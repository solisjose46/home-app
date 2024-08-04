package templates

import (
    "time"
    "home-app/app/models"
    "home-app/app/dao"
    "github.com/solisjose46/pretty-print/debug"
    "fmt"
)

func BuildFinanceTrack() (*models.FinanceTrack, error) {
    debug.PrintInfo(BuildFinanceTrack, "building finance track")

    month := time.Now().Format("July 2024")

    categories, err := dao.GetCategoriesForCurrentMonth()

    if err != nil {
        debug.PrintError(BuildFinanceTrack, err)
        return nil, err
    }

    debug.PrintSucc(BuildFinanceTrack, "returning finance track")

    return &models.FinanceTrack{
        Month: month,
        Categories: categories,
    }, nil
}

func BuildFinanceFeed(userId string) (*models.FinanceFeed, error) {
	debug.PrintInfo(BuildFinanceFeed, "building finance feed")

	expenses, err := dao.GetExpensesForCurrentMonth(userId)

    expensesCount := len(*expenses)
    debug.PrintInfo(BuildFinanceFeed, fmt.Sprintf("got %d expenses", expensesCount))    

	if err != nil {
		debug.PrintError(BuildFinanceFeed, err)
		return nil, err
	}

	debug.PrintSucc(BuildFinanceFeed, "returning finance feed")

	return &models.FinanceFeed{
        Expenses: expenses,
    }, nil
}