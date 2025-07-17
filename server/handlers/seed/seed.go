package seed

import (
	"jinovatka/assert"
	"jinovatka/server/components"
	"jinovatka/services"
	"jinovatka/utils"
	"log/slog"
	"net/http"
)

type SeedHandler struct {
	Log         *slog.Logger
	SeedService *services.SeedService
}

func NewSeedHandler(log *slog.Logger, seedService *services.SeedService) *SeedHandler {
	assert.Must(log != nil, "NewSeedHandler: log can't be nil")
	assert.Must(seedService != nil, "NewSeedHandler: seedService can't be nil")
	return &SeedHandler{
		Log:         log,
		SeedService: seedService,
	}
}

func (handler *SeedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestedID := r.PathValue("id")
	seed, err := handler.SeedService.GetSeed(requestedID)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		handler.Log.Error("SeedHandler.ServeHTTP failed to get Seed data from SeedService", "error", err.Error(), utils.LogRequestInfo(r))
		return
	}
	data := components.NewSeedViewData(seed)
	err = handler.View(w, r, data, "Semínko - "+seed.URL)
	if err != nil {
		handler.Log.Error("SeedHandler.ServeHTTP failed to render view", "error", err.Error(), utils.LogRequestInfo(r))
		return
	}
	handler.Log.Info("SeedHandler.ServeHTTP sucessfully responded", utils.LogRequestInfo(r))
}

func (handler *SeedHandler) View(w http.ResponseWriter, r *http.Request, data *components.SeedViewData, title string) error {
	return components.SeedView(data, title).Render(r.Context(), w)
}

func (handler *SeedHandler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /seed/{id}", handler)
}
