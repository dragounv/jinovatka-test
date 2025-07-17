package group

import (
	"jinovatka/assert"
	"jinovatka/server/components"
	"jinovatka/services"
	"jinovatka/utils"
	"log/slog"
	"net/http"
)

// Hanlder for seed groups. Used to create/show list of seeds to make tracking of progress of individual seeds easier.
type GroupHandler struct {
	Log         *slog.Logger
	SeedService *services.SeedService

	// Subhandlers
	SaveGroupHandler *SaveGroupHandler
}

func NewGroupHandler(log *slog.Logger, seedService *services.SeedService) *GroupHandler {
	assert.Must(log != nil, "NewGroupHandler: log can't be nil")
	assert.Must(seedService != nil, "NewGroupHandler: seedService can't be nil")
	return &GroupHandler{
		Log:              log,
		SeedService:      seedService,
		SaveGroupHandler: NewSaveGroupHandler(log, seedService),
	}
}

func (handler *GroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestedID := r.PathValue("id")
	group, err := handler.SeedService.GetGroup(requestedID)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		handler.Log.Error("NewGroupHandler.ServeHTTP failed to fetch SeedsGroup data", "error", err.Error(), utils.LogRequestInfo(r))
		return
	}
	data := components.NewGroupViewData(group)
	err = handler.View(w, r, data)
	if err != nil {
		handler.Log.Error("NewGroupHandler.ServeHTTP failed to render view", "error", err.Error(), utils.LogRequestInfo(r))
		return
	}
	handler.Log.Info("NewGroupHandler.ServeHTTP sucessfully responded", utils.LogRequestInfo(r))
}

func (handler *GroupHandler) View(w http.ResponseWriter, r *http.Request, data *components.GroupViewData) error {
	return components.GroupView(data).Render(r.Context(), w)
}

func (handler *GroupHandler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /seeds/{id}", handler)
	mux.Handle("POST /seeds/save/", handler.SaveGroupHandler)
}
