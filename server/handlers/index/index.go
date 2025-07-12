package index

import (
	"jinovatka/assert"
	"jinovatka/server/components"
	"jinovatka/utils"
	"log/slog"
	"net/http"
)

func NewIndexHandler(log *slog.Logger) *IndexHandler {
	assert.Must(log != nil, "NewIndexHandler: log can't be nil")
	return &IndexHandler{
		Log: log,
	}
}

type IndexHandler struct {
	Log *slog.Logger
}

func (handler *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.View(w, r) // TODO: handle error
	handler.Log.Info("indexHandler responded", utils.LogRequestInfo(r))
}

func (handler *IndexHandler) View(w http.ResponseWriter, r *http.Request) error {
	return components.IndexView().Render(r.Context(), w)
}
