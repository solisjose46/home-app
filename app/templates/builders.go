package templates

import (
    "time"
    "home-app/app/models"
    "home-app/app/dao"
    "github.com/solisjose46/pretty-print/debug"
)

func BuildFinanceTrack() (models.FinanceTrack, error) {
    debug.PrintInfo(BuildFinanceTrack, "building finance track")
    var financeTrack models.FinanceTrack

    financeTrack.Month = time.Now().Format("July 2024")

    categories, err := dao.GetCategoriesForCurrentMonth()

    if err != nil {
        debug.PrintError(BuildFinanceTrack, err)
        return financeTrack, err
    }

    financeTrack.Categories = categories

    debug.PrintSucc(BuildFinanceTrack, "returning finance track")

    return financeTrack, nil
}

func BuildFinanceFeed(userId string) (models.FinanceFeed, error) {
	debug.PrintInfo(BuildFinanceFeed, "building finance feed")

	var financeFeed models.FinanceFeed

	expenses, err := dao.GetExpensesForCurrentMonth(userId)

	if err != nil {
		debug.PrintError(BuildFinanceFeed, err)
		return financeFeed, err
	}

    financeFeed.Expenses = expenses

	debug.PrintSucc(BuildFinanceFeed, "returning finance feed")

	return financeFeed, nil
}