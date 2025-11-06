package redirect

import (
	"errors"
	"jinovatka/assert"
	"jinovatka/server/handlers/httperror"
	"jinovatka/services"
	"jinovatka/utils"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

// This handler redirects shortened archival URLs to wayback
type RedirectHandler struct {
	Log          *slog.Logger
	SeedService  *services.SeedService
	ErrorHandler *httperror.ErrorHandler
}

func NewRedirectHanlder(log *slog.Logger, seedService *services.SeedService, errorHandler *httperror.ErrorHandler) *RedirectHandler {
	assert.Must(log != nil, "NewRedirectHanlder: log can't be nil")
	assert.Must(seedService != nil, "NewRedirectHanlder: seedService can't be nil")
	assert.Must(errorHandler != nil, "NewRedirectHanlder: errorHandler can't be nil")
	return &RedirectHandler{
		Log:          log,
		SeedService:  seedService,
		ErrorHandler: errorHandler,
	}
}

func (handler *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	seedId := r.PathValue("id")
	seed, err := handler.SeedService.GetSeed(seedId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		handler.Log.Warn("RedirectHandler.ServeHTTP seed not found", "error", err.Error(), utils.LogRequestInfo(r))
		handler.ErrorHandler.PageNotFound(w, r) // Less scary and more informative than 500
		return
	}
	if err != nil {
		handler.Log.Error("RedirectHandler.ServeHTTP failed to get Seed data from SeedService", "error", err.Error(), utils.LogRequestInfo(r))
		handler.ErrorHandler.InternalServerError(w, r)
		return
	}

	if seed.ArchivalURL == "" {
		handler.Log.Warn("RedirectHandler.ServeHTTP seed is not harvested or archival URL is missing", "seed", seed.ShadowID, utils.LogRequestInfo(r))
		handler.ErrorHandler.PageNotFound(w, r) // TODO: Make special page explaining what happened
		return
	}

	http.Redirect(w, r, seed.ArchivalURL, http.StatusMovedPermanently)
}

func (handler *RedirectHandler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /archiv/{id}", handler)
}
