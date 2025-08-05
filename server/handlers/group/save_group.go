package group

import (
	"errors"
	"jinovatka/assert"
	"jinovatka/server/handlers/httperror"
	"jinovatka/services"
	"jinovatka/utils"
	"log/slog"
	"net/http"
)

func NewSaveGroupHandler(
	log *slog.Logger,
	seedService *services.SeedService,
	captureService *services.CaptureService,
	errorHandler *httperror.ErrorHandler,
) *SaveGroupHandler {
	assert.Must(log != nil, "NewSaveGroupHandler: log can't be nil")
	assert.Must(seedService != nil, "NewSaveGroupHandler: seedService can't be nil")
	assert.Must(captureService != nil, "NewSaveGroupHandler: captureService can't be nil")
	assert.Must(errorHandler != nil, "NewSaveGroupHandler: errorHandler can't be nil")
	return &SaveGroupHandler{
		Log:            log,
		SeedService:    seedService,
		CaptureService: captureService,
		ErrorHandler:   errorHandler,
	}
}

type SaveGroupHandler struct {
	Log            *slog.Logger
	SeedService    *services.SeedService
	CaptureService *services.CaptureService
	ErrorHandler   *httperror.ErrorHandler
}

func (handler *SaveGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const urlKey = "url-list"
	// TODO: Check that server has correct setting for request size.
	seedURL := r.FormValue(urlKey)
	group, err := handler.SeedService.Save(seedURL, true)
	// TODO: Return different error pages/messages when different errors are recived. This should help user understant what they did wrong.
	if errors.Is(err, services.ErrEmptyList) {
		handler.Log.Warn("SaveGroupHandler.ServeHTTP recieved empty seed list", utils.LogRequestInfo(r))
		handler.ErrorHandler.ServeError(w, r, "Prázdný požadavek", 400, "Prázdný požadavek", "Požadavek který jsme obdrželi obsahoval jen prázdné řádky. Prosím vraťe se na hlavní stránku a zadejte platnou URL adresu.")
		return
	}
	if err != nil {
		handler.Log.Error("SaveGroupHandler.ServeHTTP SeedService returned error when trying to save group", slog.String("error", err.Error()), utils.LogRequestInfo(r))
		handler.ErrorHandler.InternalServerError(w, r)
		return
	}

	// Enqueue seeds for capture (this does not change the success of the http request)
	err = handler.CaptureService.CaptureGroup(r.Context(), group)
	if err != nil {
		handler.Log.Error("SaveGroupHandler.ServeHTTP CaptureService returned error when trying to enqueue group", "error", err.Error(), utils.LogRequestInfo(r))
		// Do not return!
	}

	http.Redirect(w, r, "/seeds/"+group.ShadowID, http.StatusSeeOther)
	handler.Log.Info("SaveGroupHandler.ServeHTTP sucessfully responded", utils.LogRequestInfo(r))
}
