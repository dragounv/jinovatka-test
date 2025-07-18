package index

import (
	"jinovatka/assert"
	"jinovatka/server/components"
	"jinovatka/server/handlers/httperror"
	"jinovatka/utils"
	"log/slog"
	"net/http"
)

// Main handler for routes "/" and "/index.html"
type IndexHandler struct {
	Log          *slog.Logger
	ErrorHandler *httperror.ErrorHandler
}

func NewIndexHandler(log *slog.Logger, errorHandler *httperror.ErrorHandler) *IndexHandler {
	assert.Must(log != nil, "NewIndexHandler: log can't be nil")
	assert.Must(errorHandler != nil, "NewIndexHandler: errorHandler can't be nil")
	return &IndexHandler{
		Log:          log,
		ErrorHandler: errorHandler,
	}
}

func (handler *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Thanks to the http.ServeMux implementation, all NotFound pages
	// are gonna be served by this handler. Also we need to do our own http method check

	// If we aren't looking for index page, return 404 error
	path := r.URL.Path
	if path != "/" && path != "/index.html" {
		handler.ErrorHandler.PageNotFound(w, r)
		return
	}
	// If this isn't GET request, return 405 error
	if r.Method != http.MethodGet {
		handler.ErrorHandler.MethodNotAllowed(w, r)
		return
	}

	// All should be well. We can serve the index page.
	err := handler.View(w, r)
	if err != nil {
		handler.Log.Error("IndexHandler.ServeHTTP failed to render view", "error", err.Error(), utils.LogRequestInfo(r))
		return
	}
	handler.Log.Info("IndexHandler sucessfully responded", utils.LogRequestInfo(r))
}

func (handler *IndexHandler) View(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set(utils.ContentType, utils.TextHTML)
	return components.IndexView().Render(r.Context(), w)
}

func (handler *IndexHandler) Routes(mux *http.ServeMux) {
	mux.Handle("/", handler)
}
