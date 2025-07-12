package index

import (
	"jinovatka/assert"
	"jinovatka/server/components"
	"jinovatka/services"
	"jinovatka/utils"
	"log/slog"
	"net/http"
)

// Main handler for routes "/" and "/index.html"
type IndexHandler struct {
	Log *slog.Logger

	// Subhandlers
	SaveSeedHandler *SaveSeedHandler
}

func NewIndexHandler(log *slog.Logger, seedService *services.SeedService) *IndexHandler {
	assert.Must(log != nil, "NewIndexHandler: log can't be nil")
	assert.Must(seedService != nil, "NewIndexHandler: seedService can't be nil")
	return &IndexHandler{
		Log:             log,
		SaveSeedHandler: NewSaveSeedHanlder(log, seedService),
	}
}

func (handler *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.View(w, r) // TODO: handle error
	handler.Log.Info("indexHandler responded", utils.LogRequestInfo(r))
}

func (handler *IndexHandler) View(w http.ResponseWriter, r *http.Request) error {
	return components.IndexView().Render(r.Context(), w)
}

func (handler *IndexHandler) Routes(mux *http.ServeMux) {
	mux.Handle("/", handler)
	mux.Handle("POST /save-seed/", handler.SaveSeedHandler)
}
