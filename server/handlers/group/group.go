package group

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

// Hanlder for seed groups. Used to create/show list of seeds to make tracking of progress of individual seeds easier.
type GroupHandler struct {
	Log            *slog.Logger
	SeedService    *services.SeedService
	CaptureService *services.CaptureService
	ErrorHandler   *httperror.ErrorHandler

	// Subhandlers
	SaveGroupHandler   *SaveGroupHandler
	ExportGroupHandler *ExportGroupHandler
}

func NewGroupHandler(
	log *slog.Logger,
	seedService *services.SeedService,
	exporterService *services.ExporterService,
	captureService *services.CaptureService,
	errorHandler *httperror.ErrorHandler,
) *GroupHandler {
	assert.Must(log != nil, "NewGroupHandler: log can't be nil")
	assert.Must(seedService != nil, "NewGroupHandler: seedService can't be nil")
	assert.Must(exporterService != nil, "NewGroupHandler: exporterService can't be nil")
	assert.Must(captureService != nil, "NewGroupHandler: captureService can't be nil")
	assert.Must(errorHandler != nil, "NewGroupHandler: errorHandler can't be nil")
	return &GroupHandler{
		Log:                log,
		SeedService:        seedService,
		ErrorHandler:       errorHandler,
		SaveGroupHandler:   NewSaveGroupHandler(log, seedService, captureService, errorHandler),
		ExportGroupHandler: NewExportGroupHandler(log, seedService, exporterService, errorHandler),
	}
}

func (handler *GroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestedID := r.PathValue("id")
	group, err := handler.SeedService.GetGroup(requestedID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		handler.Log.Warn("GroupHandler.ServeHTTP group not found", "error", err.Error(), utils.LogRequestInfo(r))
		handler.ErrorHandler.PageNotFound(w, r) // Less scary and more informative than 500
		return
	}
	if err != nil {
		handler.Log.Error("GroupHandler.ServeHTTP failed to fetch SeedsGroup data", "error", err.Error(), utils.LogRequestInfo(r))
		handler.ErrorHandler.InternalServerError(w, r)
		return
	}
	data := components.NewGroupViewData(group)
	err = handler.View(w, r, data)
	if err != nil {
		handler.Log.Error("GroupHandler.ServeHTTP failed to render view", "error", err.Error(), utils.LogRequestInfo(r))
		return
	}
	handler.Log.Info("GroupHandler.ServeHTTP sucessfully responded", utils.LogRequestInfo(r))
}

func (handler *GroupHandler) View(w http.ResponseWriter, r *http.Request, data *components.GroupViewData) error {
	return components.GroupView(data).Render(r.Context(), w)
}

func (handler *GroupHandler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /seeds/{id}", handler)
	mux.Handle("POST /seeds/save/", handler.SaveGroupHandler)
	mux.Handle("GET /seeds/export/{id}", handler.ExportGroupHandler)
}
