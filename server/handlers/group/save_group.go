package group

import (
	"jinovatka/assert"
	"jinovatka/server/handlers/httperror"
	"jinovatka/services"
	"jinovatka/utils"
	"log/slog"
	"net/http"
)

func NewSaveGroupHandler(log *slog.Logger, seedService *services.SeedService, errorHandler *httperror.ErrorHandler) *SaveGroupHandler {
	assert.Must(log != nil, "NewSaveGroupHandler: log can't be nil")
	assert.Must(seedService != nil, "NewSaveGroupHandler: seedService can't be nil")
	assert.Must(errorHandler != nil, "NewSaveGroupHandler: errorHandler can't be nil")
	return &SaveGroupHandler{
		Log:          log,
		SeedService:  seedService,
		ErrorHandler: errorHandler,
	}
}

type SaveGroupHandler struct {
	Log          *slog.Logger
	SeedService  *services.SeedService
	ErrorHandler *httperror.ErrorHandler
}

func (handler *SaveGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const urlKey = "url-list"
	// TODO: Check that server has correct setting for request size.
	seedURL := r.FormValue(urlKey)
	// TODO: Maybe don't return the entire entity?
	group, err := handler.SeedService.Save(seedURL, true)
	// TODO: Return different error pages/messages when different errors are recived. This should help user understant what they did wrong.
	if err != nil {
		handler.Log.Error("SaveGroupHandler.ServeHTTP SeedService returned error when trying to save group", slog.String("error", err.Error()), utils.LogRequestInfo(r))
		handler.ErrorHandler.InternalServerError(w, r)
		return
	}
	http.Redirect(w, r, "/seeds/"+group.ShadowID, http.StatusSeeOther)
	handler.Log.Info("SaveGroupHandler.ServeHTTP sucessfully responded", utils.LogRequestInfo(r))
}
