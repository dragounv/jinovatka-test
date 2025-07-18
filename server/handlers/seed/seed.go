package seed

import (
	"errors"
	"jinovatka/assert"
	"jinovatka/server/components"
	"jinovatka/server/handlers/httperror"
	"jinovatka/services"
	"jinovatka/utils"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

type SeedHandler struct {
	Log          *slog.Logger
	SeedService  *services.SeedService
	ErrorHandler *httperror.ErrorHandler
}

func NewSeedHandler(log *slog.Logger, seedService *services.SeedService, errorHandler *httperror.ErrorHandler) *SeedHandler {
	assert.Must(log != nil, "NewSeedHandler: log can't be nil")
	assert.Must(seedService != nil, "NewSeedHandler: seedService can't be nil")
	assert.Must(errorHandler != nil, "NewSeedHandler: errorHandler can't be nil")
	return &SeedHandler{
		Log:          log,
		SeedService:  seedService,
		ErrorHandler: errorHandler,
	}
}

func (handler *SeedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestedID := r.PathValue("id")
	seed, err := handler.SeedService.GetSeed(requestedID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		handler.Log.Warn("SeedHandler.ServeHTTP seed not found", "error", err.Error(), utils.LogRequestInfo(r))
		handler.ErrorHandler.PageNotFound(w, r) // Less scary and more informative than 500
		return
	}
	if err != nil {
		handler.Log.Error("SeedHandler.ServeHTTP failed to get Seed data from SeedService", "error", err.Error(), utils.LogRequestInfo(r))
		handler.ErrorHandler.InternalServerError(w, r)
		return
	}
	data := components.NewSeedViewData(seed, "Semínko - "+seed.URL)
	err = handler.View(w, r, data)
	if err != nil {
		handler.Log.Error("SeedHandler.ServeHTTP failed to render view", "error", err.Error(), utils.LogRequestInfo(r))
		return
	}
	handler.Log.Info("SeedHandler.ServeHTTP sucessfully responded", utils.LogRequestInfo(r))
}

func (handler *SeedHandler) View(w http.ResponseWriter, r *http.Request, data *components.SeedViewData) error {
	return components.SeedView(data).Render(r.Context(), w)
}

func (handler *SeedHandler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /seed/{id}", handler)
}
