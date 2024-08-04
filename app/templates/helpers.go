package templates

import (
    "time"
    "home-app/app/models"
    "github.com/solisjose46/pretty-print/debug"
    "errors"
)

func (parser *TemplateParser) GetFinanceTrackData() (*models.FinanceTrack, error) {
    debug.PrintInfo(parser.GetFinanceTrackData, "building finance track")

    month := time.Now().Format("July 2024")

    categories, err := parser.dao.GetCategoriesForCurrentMonth()

    if err != nil {
        debug.PrintError(parser.GetFinanceTrackData, err)
        return nil, errors.New("error getting categories for current month")
    }

    debug.PrintSucc(parser.GetFinanceTrackData, "returning finance track")

    return &models.FinanceTrack{
        Month: month,
        Categories: categories,
    }, nil
}

func (parser *TemplateParser) GetFinanceFeedData(userId string) (*models.FinanceFeed, error) {
	debug.PrintInfo(parser.GetFinanceFeedData, "building finance feed")

	expenses, err := parser.dao.GetExpensesForCurrentMonth(userId)

	if err != nil {
		debug.PrintError(parser.GetFinanceFeedData, err)
		return nil, errors.New("error feed data")
	}

	debug.PrintSucc(parser.GetFinanceFeedData, "returning finance feed")

	return &models.FinanceFeed{
        Expenses: expenses,
    }, nil
}