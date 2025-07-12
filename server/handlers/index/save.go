package index

import (
	"jinovatka/assert"
	"jinovatka/services"
	"jinovatka/utils"
	"log/slog"
	"net/http"
)

func NewSaveSeedHanlder(log *slog.Logger, seedService *services.SeedService) *SaveSeedHandler {
	assert.Must(log != nil, "NewSaveSeedHanlder: log can't be nil")
	assert.Must(seedService != nil, "NewSaveSeedHanlder: seedService can't be nil")
	return &SaveSeedHandler{
		Log:         log,
		SeedService: seedService,
	}
}

type SaveSeedHandler struct {
	Log           *slog.Logger
	SeedService   *services.SeedService
	SucessHandler http.Handler
}

func (handler *SaveSeedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const urlKey = "url-list"
	// TODO: Check that server has correct setting for request size.
	seedURL := r.FormValue(urlKey)
	_, err := handler.SeedService.Save(seedURL, true)
	// TODO: Return different error pages/messages when different errors are recived. This should help user understant what they did wrong.
	if err != nil {
		handler.Log.Error("SaveSeed failed", slog.String("error", err.Error()), utils.LogRequestInfo(r))
		http.Error(w, "error when saving seed", http.StatusInternalServerError)
		return
	}
	handler.Log.Info("SaveHandler responded", utils.LogRequestInfo(r))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
