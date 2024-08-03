package templates

import (
    "bytes"
    "html/template"
    "home-app/app/models"
    "home-app/app/dao"
    "home-app/app/util"
    "github.com/solisjose46/pretty-print/debug"
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
    debug.PrintInfo(GetLogin, "Getting login template")

    htmlLogin := util.GetTmplPath(TmplLogin)
    htmlServerResponse := util.GetTmplPath(TmplServerResponse)
    
    tmpl, err := template.ParseFiles(htmlLogin, htmlServerResponse)
    if err != nil {
        debug.PrintError(GetLogin, err)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplLogin, nil)
    if err != nil {
        debug.PrintError(GetLogin, err)
        return "", err
    }

    debug.PrintSucc(GetLogin, "returning home login template")
    return buf.String(), nil
}

func getLoginServerResponse(sr models.ServerResponse) (string, error) {
    debug.PrintInfo(getLoginServerResponse, "Getting login template")

    htmlLogin := util.GetTmplPath(TmplLogin)
    htmlServerResponse := util.GetTmplPath(TmplServerResponse)
    
    tmpl, err := template.ParseFiles(htmlLogin, htmlServerResponse)
    if err != nil {
        debug.PrintError(getLoginServerResponse, err)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplLogin, models.Login{
        ServerResponse: sr,
    })
    if err != nil {
        debug.PrintError(getLoginServerResponse, err)
        return "", err
    }

    debug.PrintSucc(getLoginServerResponse, "returning home login template")
    return buf.String(), nil
}

func PostLogin(username, password string) (string, error) {
    debug.PrintInfo(PostLogin, "Getting post login", "validate input")

	if username == "" || password == "" {
		return getLoginServerResponse(
            models.ServerResponse{
                Message: util.InvalidInput,
                ReturnEndpoint: util.LoginEndpoint,
            },
        )
	}

    debug.PrintInfo(PostLogin, "auth attempt")

    valid, err := dao.ValidateUser(username, password)

    if err != nil {
        debug.PrintError(PostLogin, err)
        return "", err
    }

    if !valid {
        debug.PrintInfo(PostLogin, "user not auth")
        return getLoginServerResponse(
            models.ServerResponse{
                Message: util.InvalidInput,
                ReturnEndpoint: util.LoginEndpoint,
            },
        )
    }

    debug.PrintSucc(PostLogin, "authenticated!")
    return "", nil
}

func GetHome() (string, error) {
    debug.PrintInfo(GetHome, "Getting home template")

    htmlHome := util.GetTmplPath(TmplHome)
    
    tmpl, err := template.ParseFiles(htmlHome)
    if err != nil {
        debug.PrintError(GetHome, err)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplHome, nil)
    if err != nil {
        debug.PrintError(GetHome, err)
        return "", err
    }

    debug.PrintSucc(GetHome, "returning home template")
    return buf.String(), nil
}

func GetFinance(userId string) (string, error) {
    debug.PrintInfo(GetFinance, "Getting finance template")

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
        debug.PrintError(GetFinance, err)
        return "", err
    }

    if err != nil {
        debug.PrintError(GetFinance, err)
        return "", err
    }

    financeTrack, err := BuildFinanceTrack()

    if err != nil {
        debug.PrintError(GetFinance, err)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinance, 
        models.Finance{
            FinanceTrack: financeTrack,
        },
    )

    if err != nil {
        debug.PrintError(GetFinance, err)
        return "", err
    }

    debug.PrintSucc(GetFinance, "returning finance template")
    return buf.String(), nil
}

func GetFinanceTrack(userId string) (string, error) {
    debug.PrintInfo(GetFinanceTrack, "Getting finance track server response template")

    htmlFinanceTrack := util.GetTmplPath(TmplFinanceTrack)
    
    tmpl, err := template.ParseFiles(htmlFinanceTrack)
    if err != nil {
        debug.PrintError(GetFinanceTrack, err)
        return "", err
    }

    financeTrack, err := BuildFinanceTrack()

    if err != nil {
        debug.PrintError(GetFinanceTrack, err)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceTrack, financeTrack)
    if err != nil {
        debug.PrintError(GetFinanceTrack, err)
        return "", err
    }

    debug.PrintSucc(GetFinanceTrack, "returning finance track server response template")
    return buf.String(), nil
}

func getFinanceTrackServerResponse(sr models.ServerResponse) (string, error) {
    debug.PrintInfo(getFinanceTrackServerResponse, "Getting finance track server response template")

    htmlFinanceTrack := util.GetTmplPath(TmplFinanceTrack)
    htmlServerResponse := util.GetTmplPath(TmplServerResponse)
    
    tmpl, err := template.ParseFiles(htmlFinanceTrack, htmlServerResponse)
    if err != nil {
        debug.PrintError(getFinanceTrackServerResponse, err)
        return "", err
    }

    financeTrack, err := BuildFinanceTrack()

    if err != nil {
        debug.PrintError(getFinanceTrackServerResponse, err)
        return "", err
    }

    financeTrack.ServerResponse = sr

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceTrack, financeTrack)
    if err != nil {
        debug.PrintError(getFinanceTrackServerResponse, err)
        return "", err
    }

    debug.PrintSucc(getFinanceTrackServerResponse, "returning finance track server response template")
    return buf.String(), nil
}

func getFinanceTrackConfirm(ftc models.FinanceTrackConfirm) (string, error) {
    debug.PrintInfo(getFinanceTrackConfirm, "Getting finance track confirm template")

    htmlFinanceTrack := util.GetTmplPath(TmplFinanceTrack)
    htmlFinanceTrackConfirm := util.GetTmplPath(TmplFinanceTrackConfirm)
    
    tmpl, err := template.ParseFiles(htmlFinanceTrack, htmlFinanceTrackConfirm)
    if err != nil {
        debug.PrintError(getFinanceTrackConfirm, err)
        return "", err
    }

    financeTrack, err := BuildFinanceTrack()

    if err != nil {
        debug.PrintError(getFinanceTrackConfirm, err)
        return "", err
    }

    financeTrack.FinanceTrackConfirm = ftc

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceTrack, financeTrack)
    if err != nil {
        debug.PrintError(getFinanceTrackConfirm, err)
        return "", err
    }

    debug.PrintSucc(getFinanceTrackConfirm, "returning finance track confirm template")
    return buf.String(), nil
}

func PostFinanceTrack(expense models.Expense) (string, error) {
    debug.PrintInfo(PostFinanceTrack, "Posting expense")

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        return getFinanceTrackServerResponse(
            models.ServerResponse{
                Message: util.InvalidExpenseInput,
                ReturnEndpoint: util.FinanceTrackEndpoint,
            },
        )
    }

    debug.PrintSucc(PostFinanceTrack, "returning finance track confirm")

    return getFinanceTrackConfirm(
        models.FinanceTrackConfirm{
            Expense: expense,
        },
    )
}

func PostFinanceTrackConfirm(expense models.Expense) (string, error) {
    debug.PrintInfo(PostFinanceTrackConfirm, "Posting expense")

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
        debug.PrintError(PostFinanceTrackConfirm, err)
        return "", nil
    }

    sr := models.ServerResponse{
        ReturnEndpoint: util.FinanceTrackEndpoint,
    }

    if !succ {
        debug.PrintInfo(PostFinanceTrackConfirm, "fail to add expense")
        sr.Message = util.FailedToAddExpense
    } else {
        debug.PrintInfo(PostFinanceTrackConfirm, "expense added!")
        sr.Message = util.SuccAddExpense
    }

    return getFinanceTrackServerResponse(sr)
}

func GetFinanceFeed(userId string) (string, error) {
    debug.PrintInfo(GetFinanceFeed, "getting finance feed template")

    htmlFinance := util.GetTmplPath(TmplFinance)
    htmlFInanceFeed := util.GetTmplPath(TmplFinanceFeed)
    
    tmpl, err := template.ParseFiles(htmlFinance, htmlFInanceFeed)
    if err != nil {
        debug.PrintError(GetFinanceFeed, err)
        return "", err
    }

    financeFeed, err := BuildFinanceFeed(userId)

    if err != nil {
        debug.PrintError(GetFinanceFeed, err)
        return "", err
    }

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinance, 
        models.Finance{
            FinanceFeed: financeFeed,
        },
    )
    if err != nil {
        debug.PrintError(GetFinanceFeed, err)
        return "", err
    }

    debug.PrintSucc(GetFinanceFeed, "returning finance feed template")
    return buf.String(), nil
}

func getFinanceFeedServerResponse(userId string, serverResponse models.ServerResponse) (string, error) {
    debug.PrintInfo(getFinanceFeedServerResponse, "getting finance feed server response")

    htmlFinanceFeed := util.GetTmplPath(TmplFinanceFeed)
    htmlServerResponse := util.GetTmplPath(TmplServerResponse)

    tmpl, err := template.ParseFiles(htmlFinanceFeed, htmlServerResponse)
    if err != nil {
        debug.PrintError(getFinanceFeedServerResponse, err)
        return "", err
    }

    financeFeed, err := BuildFinanceFeed(userId)
    if err != nil {
        debug.PrintError(getFinanceFeedServerResponse, err)
        return "", err
    }

    financeFeed.ServerResponse = serverResponse

    var buf bytes.Buffer
    err = tmpl.ExecuteTemplate(&buf, TmplFinanceFeed, financeFeed)
    if err != nil {
        return "", err
    }

    debug.PrintSucc(getFinanceFeedServerResponse, "returning finance feed server response")
    return buf.String(), nil
}

func getFinanceFeedEdit(expense models.Expense) (string, error) {
    debug.PrintInfo(getFinanceFeedEdit, "getting finance feed edit")

    htmlFinanceFeed := util.GetTmplPath(TmplFinanceFeed)
    htmlFinanceFeedEdit := util.GetTmplPath(TmplFinanceFeedEdit)

    tmpl, err := template.ParseFiles(htmlFinanceFeed, htmlFinanceFeedEdit)
    if err != nil {
        debug.PrintError(getFinanceFeedEdit, err)
        return "", err
    }

    financeFeed, err := BuildFinanceFeed(expense.UserId)
    if err != nil {
        debug.PrintError(getFinanceFeedEdit, err)
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

    debug.PrintSucc(getFinanceFeedEdit, "returning finance feed edit")
    return buf.String(), nil
}

func getFinanceFeedConfirm(expense models.Expense) (string, error) {
    debug.PrintInfo(getFinanceFeedConfirm, "getting finance feed confirm")

    htmlFinanceFeed := util.GetTmplPath(TmplFinanceFeed)
    htmlFinanceFeedConfirm := util.GetTmplPath(TmplFinanceFeedConfirm)

    tmpl, err := template.ParseFiles(htmlFinanceFeed, htmlFinanceFeedConfirm)
    if err != nil {
        debug.PrintError(getFinanceFeedConfirm, err)
        return "", err
    }

    oldExpense, err := dao.GetExpense(expense.ExpenseId)
    if err != nil {
        debug.PrintError(getFinanceFeedConfirm, err)
        return "", err
    }

    financeFeed, err := BuildFinanceFeed(expense.UserId)
    if err != nil {
        debug.PrintError(getFinanceFeedConfirm, err)
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

    debug.PrintSucc(getFinanceFeedConfirm, "returning finance feed confirm")
    return buf.String(), nil
}

func PostFinanceFeed(expense models.Expense) (string, error) {
    debug.PrintInfo(PostFinanceFeed, "Posting finance feed")

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        debug.PrintSucc(PostFinanceFeed, "empty input")
        return getFinanceFeedServerResponse(
            expense.UserId,
            models.ServerResponse{
                Message: util.InvalidExpenseInput,
                ReturnEndpoint: util.FinanceFeedEndpoint,
            },
        )
    }

    debug.PrintSucc(PostFinanceFeed, "returning finance feed edit")
    return getFinanceFeedEdit(expense)
}

func PostFinanceFeedEdit(expense models.Expense) (string, error) {
    debug.PrintInfo(PostFinanceFeedEdit, "Posting finance feed edit")

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        debug.PrintInfo(PostFinanceFeedEdit, "empty input")
        return getFinanceFeedServerResponse(
            expense.UserId,
            models.ServerResponse{
                Message: util.InvalidExpenseInput,
                ReturnEndpoint: util.FinanceFeedEndpoint,
            },
        )
    }

    debug.PrintSucc(PostFinanceFeedEdit, "returning finance feed confirm")
    return getFinanceFeedEdit(expense)
}

func PostFinanceFeedConfirm(expense models.Expense) (string, error) {
    debug.PrintInfo(PostFinanceFeedConfirm, "Posting finance feed confirm")

    if expense.Name == "" || expense.Amount == 0 || expense.Category == "" {
        debug.PrintInfo(PostFinanceFeedConfirm, "empty input")
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
        debug.PrintError(PostFinanceFeedConfirm, err)
        return "", nil
    }

    if !success {
        debug.PrintInfo(PostFinanceFeedConfirm, "update expense issue")
        return getFinanceFeedServerResponse(
            expense.UserId,
            models.ServerResponse{
                Message: util.FailedToUpdateExpense,
                ReturnEndpoint: util.FinanceFeedEndpoint,
            },
        )
    }

    debug.PrintInfo(PostFinanceFeedConfirm, "returning finance feed confirm server response")
    return getFinanceFeedServerResponse(
        expense.UserId,
        models.ServerResponse{
            Message: util.SuccUpdateExpense,
            ReturnEndpoint: util.FinanceFeedEndpoint,
        },
    )
}