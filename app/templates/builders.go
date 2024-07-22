package templates

import (
    "time"
    "home-app/app/models"
    "home-app/app/dao"
    "home-app/app/util"
)

func BuildFinanceTrack() (models.FinanceTrack, error) {
    util.PrintMessage("building finance track")
    var financeTrack models.FinanceTrack

    financeTrack.Month = time.Now().Format("July 2024")

    categories, err := dao.GetCategoriesForCurrentMonth()

    if err != nil {
        util.PrintError(err)
        return financeTrack, err
    }

    financeTrack.Categories = categories

    util.PrintSuccess("returning finance track")

    return financeTrack, nil
}

func BuildFinanceFeed(userId string) (models.FinanceFeed, error) {
	util.PrintMessage("building finance feed")

	var financeFeed models.FinanceFeed

	expenses, err := dao.GetExpensesForCurrentMonth(userId)

	if err != nil {
		util.PrintError(err)
		return financeFeed, err
	}

    financeFeed.Expenses = expenses

	util.PrintSuccess("returning finance feed")

	return financeFeed, nil
}