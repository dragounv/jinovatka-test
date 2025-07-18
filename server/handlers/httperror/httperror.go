package httperror

import (
	"jinovatka/assert"
	"jinovatka/server/components"
	"jinovatka/utils"
	"log/slog"
	"net/http"
	"strconv"
)

// This handler is little bit different. It has no ServeHTTP or Routes method.
// It should be used from other handlers to render error pages.
type ErrorHandler struct {
	Log *slog.Logger
}

func NewErrorHandler(log *slog.Logger) *ErrorHandler {
	assert.Must(log != nil, "NewErrorHandler: log can't be nil")
	return &ErrorHandler{
		Log: log,
	}
}

// Serve 404 page.
func (handler *ErrorHandler) PageNotFound(w http.ResponseWriter, r *http.Request) {
	title := "404 - Stránka nenalezena"
	code := http.StatusNotFound
	description := "Stránka nenalezena"
	message := "Vámi hledaná adresa: '" + r.Host + r.URL.Path + "' neexistuje. " +
		"Zkontrolujte zda je adresa správně zadaná a zkuste ji vyhledat znovu. " +
		"Případně se vraťe na předchozí stranu."
	handler.ServeError(w, r, title, code, description, message)
}

// Serve 405 page.
func (handler *ErrorHandler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	title := "405 - Metoda nepodporována"
	code := http.StatusMethodNotAllowed
	description := "Metoda nepodporována"
	message := "HTTP metoda: '" + r.Method + "' není podporována pro tento endpoint." +
		"Pokud vydíte tuto zprávu po odeslání formuláře, tak se může jednat o chybu aplikace." +
		"Prosím kontaktujte provozovatele stránek a popište co se stalo."
	handler.ServeError(w, r, title, code, description, message)
}

// Serve 500 page.
func (handler *ErrorHandler) InternalServerError(w http.ResponseWriter, r *http.Request) {
	// This will likely have special handling in the future, so don't use ServeError and handle the request directly
	title := "500 - Chyba"
	code := strconv.Itoa(http.StatusInternalServerError)
	description := "Chyba na straně serveru"
	message := "Omlouváme se, došlo k chybě a nebyli jsme schopni splnit váš požadavek. Zkuste to prosím později."
	data := components.NewErrorViewData(title, code, description, message)
	err := handler.View(w, r, data)
	if err != nil {
		handler.Log.Error("NewErrorHandler.InternalServerError failed to render view", "error", err.Error(), utils.LogRequestInfo(r))
		return
	}
	errorInfo := slog.Group("error_info", slog.String("code", code), slog.String("message", message))
	handler.Log.Info("NewErrorHandler.InternalServerError sucessfully served error page", utils.LogRequestInfo(r), errorInfo)
}

// Generic error handler.
// 'title' is optional and can be left empty.
// 'code' must be the http response codeof the error.
// 'description' should contain human readable description of the error (Like: Page Not Found).
// 'message' short explanation of what happened and instructions on hoe to proceed (if possible) for user.
func (handler *ErrorHandler) ServeError(w http.ResponseWriter, r *http.Request, title string, code int, description, message string) {
	if code < 200 || code > 599 {
		handler.Log.Warn("ErrorHandler.ServeError recieved unusual error code, this may be bug pls fix!", "code", code)
	}
	strcode := strconv.Itoa(code)
	if title == "" {
		title = strcode + " - Chyba"
	}
	if description == "" {
		description = "Něco se pokazilo :("
	}
	data := components.NewErrorViewData(title, strcode, description, message)
	err := handler.View(w, r, data)
	if err != nil {
		handler.Log.Error("NewErrorHandler.ServeError failed to render view", "error", err.Error(), utils.LogRequestInfo(r))
		return
	}
	errorInfo := slog.Group("error_info", slog.String("code", strcode), slog.String("message", message))
	handler.Log.Info("NewErrorHandler.ServeError sucessfully served error page", utils.LogRequestInfo(r), errorInfo)
}

func (handler *ErrorHandler) View(w http.ResponseWriter, r *http.Request, data *components.ErrorViewData) error {
	w.Header().Set(utils.ContentType, utils.TextHTML)
	return components.ErrorView(data).Render(r.Context(), w)
}
