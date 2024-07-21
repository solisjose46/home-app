package templates

import (
    "fmt"
    "bytes"
    "html/template"
    "time"
    "home-app/app/models"
    "home-app/app/dao"
    "home-app/app/util"
)

func BuildFinanceTrack() (models.FinanceTrack, error) {
    util.PrintMessage("building finance track")
    var financeTrack models.FinanceTrack

    financeTrack.month := time.Now().Format("July 2024")

    categories, err := dao.GetCategoriesForCurrentMonth()

    if err != nil {
        util.PrintError(err)
        return nil, err
    }

    financeTrack.Categories = categories

    return financeTrack

    util.PrintSuccess("returning finance track")
}

BuildFinanceFeed(userId string) (models.FinanceFeed, error) {
	util.PrintMessage("building finance feed")

	var financeFeed models.FinanceFeed

	expenses, err := dao.GetExpensesForCurrentMonth(userId)

	if err != nil {
		util.PrintError(err)
		return nil, err
	}

    financeFeed.Expenses = expenses

	uitil.PrintSuccess("returning finance feed")

	return financeFeed
}