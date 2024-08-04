package handle

import (
	"github.com/solisjose46/pretty-print/debug"
	"home-app/app/models"
	"home-app/app/util"
	"net/http"
	"strconv"
)

func (handler *Handler) FinanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	debug.PrintInfo(handler.FinanceHandler, GET, util.FinanceEndpoint)

	session, err := handler.cookieStore.Get(r, SessionName)

	if err != nil {
		debug.PrintError(handler.FinanceHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		debug.PrintInfo(handler.FinanceHandler, "redirecting to login")
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	financeTrack, err := handler.parser.GetFinanceTrackData()

	if err != nil {
		debug.PrintError(handler.FinanceHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	response, err := handler.parser.GetFinance(financeTrack, nil)

	if err != nil {
		debug.PrintError(handler.FinanceHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(*response))

	debug.PrintSucc(handler.FinanceHandler, GET, util.FinanceTrackEndpoint)
}

func (handler *Handler) FinanceTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	session, err := handler.cookieStore.Get(r, SessionName)

	if err != nil {
		debug.PrintError(handler.FinanceTrackHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId := session.Values[SessionUserId].(string)

	if r.Method == http.MethodGet {
		debug.PrintInfo(handler.FinanceTrackHandler, GET, util.FinanceTrackEndpoint)

		financeTrack, err := handler.parser.GetFinanceTrackData()

		if err != nil {
			debug.PrintError(handler.FinanceTrackHandler, err)
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		response, err := handler.parser.GetFinanceTrack(financeTrack)

		if err != nil {
			debug.PrintError(handler.FinanceTrackHandler, err)
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		w.Write([]byte(*response))
		debug.PrintSucc(handler.FinanceTrackHandler, GET, util.FinanceTrackEndpoint)
		return
	}

	debug.PrintInfo(handler.FinanceTrackHandler, POST, util.FinanceTrackEndpoint)

	amount, err := strconv.ParseFloat(r.FormValue(FormAmount), 64)

	if err != nil {
		debug.PrintError(handler.FinanceTrackHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := &models.Expense{
		UserId:   userId,
		Name:     r.FormValue(FormName),
		Amount:   amount,
		Category: r.FormValue(FormCategory),
	}

	response, err := handler.parser.PostFinanceTrack(expense)

	if err != nil {
		debug.PrintError(handler.FinanceTrackHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	
	w.Write([]byte(*response))
	debug.PrintSucc(handler.FinanceTrackHandler, POST, util.FinanceTrackEndpoint)
}

func (handler *Handler) FinanceTrackConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	debug.PrintInfo(handler.FinanceTrackConfirmHandler, POST, util.FinanceTrackConfirmEndpoint)

	session, err := handler.cookieStore.Get(r, SessionName)

	if err != nil {
		debug.PrintError(handler.FinanceTrackConfirmHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue(FormAmount), 64)
	
	if err != nil {
		debug.PrintError(handler.FinanceTrackConfirmHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := &models.Expense{
		UserId:   session.Values[SessionUserId].(string),
		Name:     r.FormValue(FormName),
		Amount:   amount,
		Category: r.FormValue(FormCategory),
	}

	response, err := handler.parser.PostFinanceTrackConfirm(expense)

	if err != nil {
		debug.PrintError(handler.FinanceTrackConfirmHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(*response))
	debug.PrintSucc(handler.FinanceTrackConfirmHandler, POST, util.FinanceTrackConfirmEndpoint)
}

func (handler *Handler) FinanceFeedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	session, err := handler.cookieStore.Get(r, SessionName)

	if err != nil {
		debug.PrintError(handler.FinanceFeedHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	userId, ok:= session.Values[SessionUserId].(string)
	if !ok {
		debug.PrintError(handler.FinanceFeedHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		debug.PrintInfo(handler.FinanceFeedHandler, GET, util.FinanceFeedEndpoint)

		financeFeed, err := handler.parser.GetFinanceFeedData(userId)
		if err != nil {
			debug.PrintError(handler.FinanceFeedHandler, err)
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		response, err := handler.parser.GetFinanceFeed(financeFeed)
		if err != nil {
			debug.PrintError(handler.FinanceFeedHandler, err)
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		debug.PrintSucc(handler.FinanceFeedHandler, GET, util.FinanceFeedEndpoint)
		w.Write([]byte(*response))
		return
	}

	debug.PrintInfo(handler.FinanceFeedHandler, POST, util.FinanceFeedEndpoint)

	amount, err := strconv.ParseFloat(r.FormValue(FormAmount), 64)
	if err != nil {
		debug.PrintError(handler.FinanceFeedHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := &models.Expense{
		ExpenseId: r.FormValue(FormExpenseId),
		UserId:    session.Values[SessionUserId].(string),
		Name:      r.FormValue(FormName),
		Amount:    amount,
		Category:  r.FormValue(FormCategory),
	}

	response, err := handler.parser.PostFinanceFeed(expense)

	if err != nil {
		debug.PrintError(handler.FinanceFeedHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(*response))
	debug.PrintSucc(handler.FinanceFeedHandler, POST, util.FinanceFeedEndpoint)
}

func (handler *Handler) FinanceFeedEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	debug.PrintInfo(handler.FinanceFeedEditHandler, POST, util.FinanceFeedEditEndpoint)

	session, err := handler.cookieStore.Get(r, SessionName)

	if err != nil {
		debug.PrintError(handler.FinanceFeedEditHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue(FormAmount), 64)
	if err != nil {
		debug.PrintError(handler.FinanceFeedEditHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := &models.Expense{
		ExpenseId: r.FormValue(FormExpenseId),
		UserId:    session.Values[SessionUserId].(string),
		Name:      r.FormValue(FormName),
		Amount:    amount,
		Category:  r.FormValue(FormCategory),
	}

	response, err := handler.parser.PostFinanceFeedEdit(expense)
	if err != nil {
		debug.PrintError(handler.FinanceFeedEditHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(*response))
	debug.PrintSucc(handler.FinanceFeedEditHandler, POST, util.FinanceFeedEditEndpoint)
}

func (handler *Handler) FinanceFeedConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	debug.PrintInfo(handler.FinanceFeedConfirmHandler, POST, util.FinanceFeedConfirmEndpoint)

	session, err := handler.cookieStore.Get(r, SessionName)

	if err != nil {
		debug.PrintError(handler.FinanceFeedConfirmHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue(FormAmount), 64)
	if err != nil {
		debug.PrintError(handler.FinanceFeedConfirmHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	expense := &models.Expense{
		ExpenseId: r.FormValue(FormExpenseId),
		UserId:    session.Values[SessionUserId].(string),
		Name:      r.FormValue(FormName),
		Amount:    amount,
		Category:  r.FormValue(FormCategory),
	}

	response, err := handler.parser.PostFinanceFeedConfirm(expense)

	if err != nil {
		debug.PrintError(handler.FinanceFeedConfirmHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(*response))
	debug.PrintSucc(handler.FinanceFeedConfirmHandler, POST, util.FinanceFeedConfirmEndpoint)
}