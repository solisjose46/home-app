package handle

import (
	"github.com/solisjose46/pretty-print/debug"
	"github.com/gorilla/sessions"
	"home-app/app/dao"
	"home-app/app/templates"
	"home-app/app/util"
	"net/http"
)

type Handler struct {
	dao *dao.Dao
	cookieStore *sessions.CookieStore
	parser *templates.TemplateParser
}

func NewHandler(dao *dao.Dao, cookieStore *sessions.CookieStore) *Handler {
	return &Handler{
		dao: dao,
		cookieStore: cookieStore,
		parser: templates.NewTemplateParser(dao),
	}
}

func (handler *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	session, err := handler.cookieStore.Get(r, SessionName)

	if err != nil {
		debug.PrintError(handler.LoginHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] != nil {
		http.Redirect(w, r, util.HomeEndpoint, http.StatusSeeOther)
		return
	}

    if r.Method == http.MethodGet {
		debug.PrintInfo(handler.LoginHandler, GET, util.LoginEndpoint)

        response, err := handler.parser.GetLogin(nil)
        if err != nil {
			debug.PrintError(handler.LoginHandler, err)
            http.Error(w, InternalServerError, http.StatusInternalServerError)
            return
        }
        w.Write([]byte(*response))
		debug.PrintSucc(handler.LoginHandler, GET, util.LoginEndpoint)
        return
    }

	debug.PrintInfo(handler.LoginHandler, POST, util.LoginEndpoint)

	username := r.FormValue(FormUsername)
	password := r.FormValue(FormPassword)

	response, err := handler.parser.PostLogin(username, password)

	if err != nil {
		debug.PrintError(handler.LoginHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if response != nil {
		debug.PrintSucc(handler.LoginHandler, POST, util.LoginEndpoint, "Server response")
		w.Write([]byte(*response))
		return
	}

	userId, err := handler.dao.GetUserId(username)
	if err != nil {
		debug.PrintError(handler.LoginHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	session.Values[SessionUsername] = username
	session.Values[SessionUserId] = userId
	err = session.Save(r, w)

	if err != nil {
		debug.PrintError(handler.LoginHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	debug.PrintSucc(handler.LoginHandler, POST, util.LoginEndpoint)
	http.Redirect(w, r, util.HomeEndpoint, http.StatusSeeOther)
}

func (handler *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	debug.PrintInfo(handler.LogoutHandler, POST, util.LogoutEndpoint)

	session, err := handler.cookieStore.Get(r, SessionName)

	if err != nil {
		debug.PrintError(handler.LogoutHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = Expire
	session.Save(r, w)

	debug.PrintSucc(handler.LogoutHandler, POST, util.LogoutEndpoint)
	http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
}

func (handler *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	debug.PrintInfo(handler.HomeHandler, GET, util.HomeEndpoint)

	session, err := handler.cookieStore.Get(r, SessionName)

	if err != nil {
		debug.PrintError(handler.HomeHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if session.Values[SessionUsername] == nil {
		debug.PrintInfo(handler.HomeHandler, "redirecting to login")
		http.Redirect(w, r, util.LoginEndpoint, http.StatusSeeOther)
		return
	}

	response, err := handler.parser.GetHome()
	
	if err != nil {
		debug.PrintError(handler.HomeHandler, err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	debug.PrintSucc(handler.HomeHandler, GET, util.HomeEndpoint)
}