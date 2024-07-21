package templates

import (
    "bytes"
    "html/template"
    "home-app/app/models"
    "home-app/app/dao"
    "home-app/app/util"
)

const (
    TmplLogin               = "login"
    TmplHome                = "home"
    TmplFinance             = "finance"
    TmplFinanceTrack        = "finance-track"
    TmplFinanceTrackConfirm = "finance-track-confirm"
    TmplFinanceFeed         = "finance-feed"
    TmplFinanceFeedEdit     = "finance-feed-edit"
    TmplFinanceFeedConfirm  = "finance-feed-confirm"
    TmplServerResponse      = "server-response"

)

func GetLogin() (string, error) {
    util.PrintMessage("Getting login template")

    htmlLogin := util.GetTmplPath(TmplLogin)
    htmlServerResponse := util.GetTmplPath(TmplServerResponse)
    
    tmpl, err := template.ParseFiles(htmlLogin, htmlServerResponse)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplLogin, nil)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    util.PrintSuccess("returning home login template")
    return buf.String(), nil
}

func getLoginServerResponse(sr models.ServerResponse) (string, error) {
    util.PrintMessage("Getting login template")

    htmlLogin := util.GetTmplPath(TmplLogin)
    htmlServerResponse := util.GetTmplPath(TmplServerResponse)
    
    tmpl, err := template.ParseFiles(htmlLogin, htmlServerResponse)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplLogin, models.Login{
        ServerResponse: sr,
    })
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    util.PrintSuccess("returning home login template")
    return buf.String(), nil
}

func PostLogin(username, password string) (string, error) {
    util.PrintMessage("Getting post login\n", "validate input")

	if username == "" || password == "" {
		return getLoginServerResponse(
            models.ServerResponse{
                Message: util.InvalidInput,
                ReturnEndpoint: util.LoginEndpoint,
            },
        )
	}

    util.PrintMessage("auth attempt")

    valid, err := dao.ValidateUser(username, password)

    if err != nil {
        util.PrintError(err)
        return "", err
    }

    if !valid {
        return getLoginServerResponse(
            models.ServerResponse{
                Message: util.InvalidInput,
                ReturnEndpoint: util.LoginEndpoint,
            },
        )
    }

    util.PrintSuccess("authenticated!")
    return "", nil
}

func GetHome() (string, error) {
    util.PrintMessage("Getting home template")

    htmlHome := util.GetTmplPath(TmplHome)
    
    tmpl, err := template.ParseFiles(htmlHome)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplHome, nil)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    util.PrintSuccess("returning home template")
    return buf.String(), nil
}


func GetFinance(userId string) (string, error) {
    util.PrintMessage("Getting finance template")

    htmlFinance := util.GetTmplPath(TmplFinance)
    htmlFinanceTrack := util.GetTmplPath(TmplFinanceTrack)
    
    tmpl, err := template.ParseFiles(htmlFinance, htmlFinanceTrack)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeTrack, err := BuildFinanceTrack()

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinance, 
        models.Finance{
            FinanceTrack: financeTrack,
        },
    )
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    util.PrintSuccess("returning finance template")
    return buf.String(), nil
}

func GetFinanceTrack(userId string) (string, error) {
    util.PrintMessage("Getting finance track server response template")

    htmlFinanceTrack := util.GetTmplPath(TmplFinanceTrack)
    
    tmpl, err := template.ParseFiles(htmlFinanceTrack)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeTrack, err := BuildFinanceTrack()

    if err != nil {
        util.PrintError(err)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceTrack, financeTrack)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    util.PrintSuccess("returning finance track server response template")
    return buf.String(), nil
}

func getFinanceTrackServerResponse(sr models.ServerResponse) (string, error) {
    util.PrintMessage("Getting finance track server response template")

    htmlFinanceTrack := util.GetTmplPath(TmplFinanceTrack)
    htmlServerResponse := util.GetTmplPath(TmplServerResponse)
    
    tmpl, err := template.ParseFiles(htmlFinanceTrack, htmlServerResponse)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeTrack, err := BuildFinanceTrack()

    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeTrack.ServerResponse = sr

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceTrack, financeTrack)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    util.PrintSuccess("returning finance track server response template")
    return buf.String(), nil
}

func getFinanceTrackConfirm(ftc models.FinanceTrackConfirm) (string, error) {
    util.PrintMessage("Getting finance track confirm template")

    htmlFinanceTrack := util.GetTmplPath(TmplFinanceTrack)
    htmlFinanceTrackConfirm := util.GetTmplPath(TmplFinanceTrackConfirm)
    
    tmpl, err := template.ParseFiles(htmlFinanceTrack, htmlFinanceTrackConfirm)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeTrack, err := BuildFinanceTrack()

    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeTrack.FinanceTrackConfirm = ftc

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceTrack, financeTrack)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    util.PrintSuccess("returning finance track confirm template")
    return buf.String(), nil
}

func PostFinanceTrack(expense models.Expense) (string, error) {
    util.PrintMessage("Posting expense")

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        return getFinanceTrackServerResponse(
            models.ServerResponse{
                Message: util.InvalidExpenseInput,
                ReturnEndpoint: util.FinanceTrackEndpoint,
            },
        )
    }

    util.PrintMessage("returning finance track confirm")

    return getFinanceTrackConfirm(
        models.FinanceTrackConfirm{
            Expense: expense,
        },
    )
}

func PostFinanceTrackConfirm(expense models.Expense) (string, error) {
    util.PrintMessage("Posting expense")

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        return getFinanceTrackServerResponse(
            models.ServerResponse{
                Message: util.InvalidExpenseInput,
                ReturnEndpoint: util.FinanceTrackEndpoint,
            },
        )
    }

    succ, err := dao.AddExpense(expense)
    if err != nil {
        util.PrintError(err)
        return "", nil
    }

    sr := models.ServerResponse{
        ReturnEndpoint: util.FinanceTrackEndpoint,
    }

    if !succ {
        util.PrintMessage("fail to add expense")
        sr.Message = util.FailedToAddExpense
    } else {
        util.PrintSuccess("expense added!")
        sr.Message = util.SuccAddExpense
    }

    return getFinanceTrackServerResponse(sr)
}

func GetFinanceFeed(userId string) (string, error) {
    util.PrintMessage("getting finance feed template")

    htmlFinance := util.GetTmplPath(TmplFinance)
    htmlFInanceFeed := util.GetTmplPath(TmplFinanceFeed)
    
    tmpl, err := template.ParseFiles(htmlFinance, htmlFInanceFeed)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeFeed, err := BuildFinanceFeed(userId)

    if err != nil {
        util.PrintError(err)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinance, 
        models.Finance{
            FinanceFeed: financeFeed,
        },
    )
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    util.PrintMessage("returning finance feed template")
    return buf.String(), nil
}

func getFinanceFeedServerResponse(userId string, serverResponse models.ServerResponse) (string, error) {
    util.PrintMessage("getting finance feed server response")

    htmlFinanceFeed := util.GetTmplPath(TmplFinanceFeed)
    htmlServerResponse := util.GetTmplPath(TmplServerResponse)

    tmpl, err := template.ParseFiles(htmlFinanceFeed, htmlServerResponse)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeFeed, err := BuildFinanceFeed(userId)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeFeed.ServerResponse = serverResponse

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceFeed, financeFeed)
    if err != nil {
        return "", err
    }

    util.PrintSuccess("returning finance feed server response")
    return buf.String(), nil
}

func getFinanceFeedEdit(expense models.Expense) (string, error) {
    util.PrintMessage("getting finance feed edit")

    htmlFinanceFeed := util.GetTmplPath(TmplFinanceFeed)
    htmlFinanceFeedEdit := util.GetTmplPath(TmplFinanceFeedEdit)

    tmpl, err := template.ParseFiles(htmlFinanceFeed, htmlFinanceFeedEdit)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeFeed, err := BuildFinanceFeed(expense.UserId)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeFeed.FinanceFeedEdit = models.FinanceFeedEdit{
        Expense: expense,
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceFeed, financeFeed)
    if err != nil {
        return "", err
    }

    util.PrintSuccess("returning finance feed edit")
    return buf.String(), nil
}

func getFinanceFeedConfirm(expense models.Expense) (string, error) {
    util.PrintMessage("getting finance feed confirm")

    htmlFinanceFeed := util.GetTmplPath(TmplFinanceFeed)
    htmlFinanceFeedConfirm := util.GetTmplPath(TmplFinanceFeedConfirm)

    tmpl, err := template.ParseFiles(htmlFinanceFeed, htmlFinanceFeedConfirm)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    oldExpense, err := dao.GetExpense(expense.ExpenseId)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeFeed, err := BuildFinanceFeed(expense.UserId)
    if err != nil {
        util.PrintError(err)
        return "", err
    }

    financeFeed.FinanceFeedConfirm = models.FinanceFeedConfirm{
        OldExpense: oldExpense,
        NewExpense: expense,
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceFeed, financeFeed)
    if err != nil {
        return "", err
    }

    util.PrintSuccess("returning finance feed confirm")
    return buf.String(), nil
}

func PostFinanceFeed(expense models.Expense) (string, error) {
    util.PrintMessage("Posting finance feed")

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        return getFinanceFeedServerResponse(
            expense.UserId,
            models.ServerResponse{
                Message: util.InvalidExpenseInput,
                ReturnEndpoint: util.FinanceFeedEndpoint,
            },
        )
    }

    util.PrintSuccess("returning finance feed edit")

    return getFinanceFeedEdit(expense)
}

func PostFinanceFeedEdit(expense models.Expense) (string, error) {
    util.PrintMessage("Posting finance feed edit")

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        return getFinanceFeedServerResponse(
            expense.UserId,
            models.ServerResponse{
                Message: util.InvalidExpenseInput,
                ReturnEndpoint: util.FinanceFeedEndpoint,
            },
        )
    }

    util.PrintSuccess("returning finance feed confirm")

    return getFinanceFeedEdit(expense)
}

func PostFinanceFeedConfirm(expense models.Expense) (string, error) {
    util.PrintMessage("Posting finance feed confirm")

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        return getFinanceFeedServerResponse(
            expense.UserId,
            models.ServerResponse{
                Message: util.InvalidExpenseInput,
                ReturnEndpoint: util.FinanceFeedEndpoint,
            },
        )
    }

    success, err := dao.UpdateExpense(expense)

    if err != nil {
        util.PrintError(err)
        return "", nil
    }

    if !success {
        return getFinanceFeedServerResponse(
            expense.UserId,
            models.ServerResponse{
                Message: util.FailedToUpdateExpense,
                ReturnEndpoint: util.FinanceFeedEndpoint,
            },
        )
    }

    util.PrintSuccess("returning finance feed server response")

    return getFinanceFeedServerResponse(
        expense.UserId,
        models.ServerResponse{
            Message: util.SuccUpdateExpense,
            ReturnEndpoint: util.FinanceFeedEndpoint,
        },
    )
}