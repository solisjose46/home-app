package templates

import (
    "bytes"
    "errors"
    "html/template"
    "home-app/app/models"
    "home-app/app/util"
    "github.com/solisjose46/pretty-print/debug"
)

func (parser *TemplateParser) GetFinance(financeTrack *models.FinanceTrack, financeFeed *models.FinanceFeed) (*string, error) {
    debug.PrintInfo(parser.GetFinance, "Getting finance template")

    if financeTrack != nil && financeFeed != nil {
        return nil, errors.New("both finance track and finance feed cannot be non-nil")
    }

    if financeTrack == nil && financeFeed == nil {
        return nil, errors.New("both finance track and finance feed cannot be nil")
    }

    htmlFinance := util.GetTmplPath(TmplFinance)
    htmlFinanceTrack := util.GetTmplPath(TmplFinanceTrack)
    htmlFinanceTrackConfirm := util.GetTmplPath(TmplFinanceTrackConfirm)
    htmlServerResponse := util.GetTmplPath(TmplServerResponse)

    htmlFinanceFeed := util.GetTmplPath(TmplFinanceFeed)
    htmlFinanceFeedEdit := util.GetTmplPath(TmplFinanceFeedEdit)
    htmlFinanceFeedConfirm := util.GetTmplPath(TmplFinanceFeedConfirm)
    
    tmpl, err := template.ParseFiles(
        htmlFinance, htmlFinanceTrack,
        htmlFinanceTrackConfirm, htmlServerResponse,
        htmlFinanceFeed, htmlFinanceFeedEdit,
        htmlFinanceFeedConfirm,
    )

    if err != nil {
        debug.PrintError(parser.GetFinance, err)
        return nil, err
    }

    var buf bytes.Buffer
    finance := &models.Finance{
        FinanceTrack: financeTrack,
        FinanceFeed: financeFeed,
    }

    err = tmpl.ExecuteTemplate(&buf, TmplFinance, finance)

    if err != nil {
        debug.PrintError(parser.GetFinance, err)
        return nil, err
    }

    tmplString := buf.String()

    debug.PrintSucc(parser.GetFinance, "returning finance template")
    return &tmplString, nil
}

func (parser *TemplateParser) GetFinanceTrack(financeTrack *models.FinanceTrack) (*string, error) {
    debug.PrintInfo(parser.GetFinanceTrack, "Getting finance track confirm template")

    htmlFinanceTrack := util.GetTmplPath(TmplFinanceTrack)
    htmlFinanceTrackConfirm := util.GetTmplPath(TmplFinanceTrackConfirm)
    htmlServerResponse := util.GetTmplPath(TmplServerResponse)
    
    tmpl, err := template.ParseFiles(
        htmlFinanceTrack, htmlFinanceTrackConfirm,
        htmlServerResponse,
    )
    
    if err != nil {
        debug.PrintError(parser.GetFinanceTrack, err)
        return nil, err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceTrack, financeTrack)
    if err != nil {
        debug.PrintError(parser.GetFinanceTrack, err)
        return nil, err
    }

    tmplString := buf.String()

    debug.PrintSucc(parser.GetFinanceTrack, "returning finance track confirm template")
    return &tmplString, nil
}

func (parser *TemplateParser) PostFinanceTrack(expense *models.Expense) (*string, error) {
    debug.PrintInfo(parser.PostFinanceTrack, "Posting expense")

    financeTrack, err := parser.GetFinanceTrackData()

    if err != nil {
        debug.PrintError(parser.PostFinanceTrack, err)
        return nil, errors.New("failed to get finance track data")
    }

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        financeTrack.ServerResponse = &models.ServerResponse{
            Message: util.InvalidExpenseInput,
            ReturnEndpoint: util.FinanceTrackEndpoint,
        }
    } else {
        financeTrack.FinanceTrackConfirm = &models.FinanceTrackConfirm{
            Expense: expense,
        }
    }

    debug.PrintInfo(parser.PostFinanceTrack, "returning finance track confirm")

    return parser.GetFinanceTrack(financeTrack)
}

func (parser *TemplateParser) PostFinanceTrackConfirm(expense *models.Expense) (*string, error) {
    debug.PrintInfo(parser.PostFinanceTrackConfirm, "Posting expense")

    financeTrack, err := parser.GetFinanceTrackData()

    if err != nil {
        debug.PrintError(parser.PostFinanceTrackConfirm, err)
        return nil, errors.New("failed to get finance track data")
    }

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        financeTrack.ServerResponse = &models.ServerResponse{
            Message: util.InvalidExpenseInput,
            ReturnEndpoint: util.FinanceTrackEndpoint,
        }

        return parser.GetFinanceTrack(financeTrack)
    }

    succ, err := parser.dao.AddExpense(expense)

    if err != nil {
        debug.PrintError(parser.PostFinanceTrackConfirm, err)
        return nil, errors.New("failed to add expense")
    }

    serverResponse := &models.ServerResponse{
        ReturnEndpoint: util.FinanceTrackEndpoint,
    }

    if !succ {
        debug.PrintInfo(parser.PostFinanceTrackConfirm, "fail to add expense")
        serverResponse.Message = util.FailedToAddExpense
    } else {
        debug.PrintInfo(parser.PostFinanceTrackConfirm, "expense added!")
        serverResponse.Message = util.SuccAddExpense
    }

    financeTrack.ServerResponse = serverResponse

    return parser.GetFinanceTrack(financeTrack)
}

func (parser *TemplateParser) GetFinanceFeed(financeFeed *models.FinanceFeed) (*string, error) {
    debug.PrintInfo(parser.GetFinanceFeed, "getting finance feed template")

    htmlFinanceFeed := util.GetTmplPath(TmplFinanceFeed)
    htmlFinanceFeedEdit := util.GetTmplPath(TmplFinanceFeedEdit)
    htmlFinanceFeedConfirm := util.GetTmplPath(TmplFinanceFeedConfirm)
    htmlServerResponse := util.GetTmplPath(TmplServerResponse)
    
    tmpl, err := template.ParseFiles(
        htmlFinanceFeed, htmlFinanceFeedEdit,
        htmlFinanceFeedConfirm, htmlServerResponse,
    )
    
    if err != nil {
        debug.PrintError(parser.GetFinanceFeed, err)
        return nil, errors.New("failed to parse finance feed template")
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceFeed, financeFeed)
    
    if err != nil {
        debug.PrintError(parser.GetFinanceFeed, err)
        return nil, err
    }

    tmplString := buf.String()

    debug.PrintSucc(parser.GetFinanceFeed, "returning finance feed template")
    return &tmplString, nil
}

func (parser *TemplateParser) PostFinanceFeed(expense *models.Expense) (*string, error) {
    debug.PrintInfo(parser.PostFinanceFeed, "Posting finance feed")

    financeFeed, err := parser.GetFinanceFeedData(expense.UserId)

    if err != nil {
        debug.PrintError(parser.PostFinanceFeed, err)
        return nil, errors.New("failed to get finance feed data")
    }

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        debug.PrintSucc(parser.PostFinanceFeed, "empty input")
        financeFeed.ServerResponse = &models.ServerResponse{
            Message: util.InvalidExpenseInput,
            ReturnEndpoint: util.FinanceFeedEndpoint,
        }
    } else {
        financeFeed.FinanceFeedEdit = &models.FinanceFeedEdit{
            Expense: expense,
        }
    }

    debug.PrintSucc(parser.PostFinanceFeed, "returning finance feed edit")
    return parser.GetFinanceFeed(financeFeed)
}

func (parser *TemplateParser) PostFinanceFeedEdit(expense *models.Expense) (*string, error) {
    debug.PrintInfo(parser.PostFinanceFeedEdit, "Posting finance feed edit")

    financeFeed, err := parser.GetFinanceFeedData(expense.UserId)

    if err != nil {
        debug.PrintError(parser.PostFinanceFeedEdit, err)
        return nil, errors.New("failed to get finance feed data")
    }

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        debug.PrintInfo(parser.PostFinanceFeedEdit, "empty input")
        financeFeed.ServerResponse = &models.ServerResponse{
            Message: util.InvalidExpenseInput,
            ReturnEndpoint: util.FinanceFeedEndpoint,
        }
        return parser.GetFinanceFeed(financeFeed)
    }

    OldExpense, err := parser.dao.GetExpense(expense.ExpenseId)

    if err != nil {
        debug.PrintError(parser.PostFinanceFeedEdit, err)
        return nil, errors.New("failed to get expense")
    }

    financeFeed.FinanceFeedConfirm = &models.FinanceFeedConfirm{
        OldExpense: OldExpense,
        NewExpense: expense,
    }

    debug.PrintSucc(parser.PostFinanceFeedEdit, "returning finance feed confirm")
    return parser.GetFinanceFeed(financeFeed)
}

func (parser *TemplateParser) PostFinanceFeedConfirm(expense *models.Expense) (*string, error) {
    debug.PrintInfo(parser.PostFinanceFeedConfirm, "Posting finance feed confirm")

    financeFeed, err := parser.GetFinanceFeedData(expense.UserId)
    if err != nil {
        debug.PrintError(parser.PostFinanceFeedConfirm, err)
        return nil, errors.New("failed to get finance feed data")
    }

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        debug.PrintInfo(parser.PostFinanceFeedConfirm, "empty input")
        financeFeed.ServerResponse = &models.ServerResponse{
            Message: util.InvalidExpenseInput,
            ReturnEndpoint: util.FinanceFeedEndpoint,
        }
        return parser.GetFinanceFeed(financeFeed)
    }

    success, err := parser.dao.UpdateExpense(expense)

    if !success {
        debug.PrintInfo(parser.PostFinanceFeedConfirm, "update expense issue")
        financeFeed.ServerResponse = &models.ServerResponse{
            Message: util.FailedToUpdateExpense,
            ReturnEndpoint: util.FinanceFeedEndpoint,
        }
        return parser.GetFinanceFeed(financeFeed)
    }

    if err != nil {
        debug.PrintError(parser.PostFinanceFeedConfirm, err)
        return nil, errors.New("failed to update expense")
    }

    financeFeed.ServerResponse = &models.ServerResponse{
        Message: util.SuccUpdateExpense,
        ReturnEndpoint: util.FinanceFeedEndpoint,
    }

    debug.PrintInfo(parser.PostFinanceFeedConfirm, "returning finance feed confirm server response")
    return parser.GetFinanceFeed(financeFeed)
}