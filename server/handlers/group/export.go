package group

import (
	"bytes"
	"jinovatka/assert"
	"jinovatka/entities"
	"jinovatka/server/handlers/httperror"
	"jinovatka/services"
	"jinovatka/utils"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

type ExportGroupHandler struct {
	Log             *slog.Logger
	SeedService     *services.SeedService
	ExporterService *services.ExporterService
	ErrorHandler    *httperror.ErrorHandler
}

func NewExportGroupHandler(
	log *slog.Logger,
	seedService *services.SeedService,
	exporterService *services.ExporterService,
	errorHandler *httperror.ErrorHandler,
) *ExportGroupHandler {
	assert.Must(log != nil, "NewExportGroupHandler: log can't be nil")
	assert.Must(seedService != nil, "NewExportGroupHandler: seedService can't be nil")
	assert.Must(exporterService != nil, "NewExportGroupHandler: exporterService can't be nil")
	assert.Must(errorHandler != nil, "NewExportGroupHandler: errorHandler can't be nil")
	return &ExportGroupHandler{
		Log:             log,
		SeedService:     seedService,
		ExporterService: exporterService,
		ErrorHandler:    errorHandler,
	}
}

func (handler *ExportGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	groupId := r.PathValue("id")

	group, err := handler.SeedService.GetGroup(groupId)
	// TODO: Create common error for services to comunicate that record does not exist so we don't have to break the layer model all the time
	if err == gorm.ErrRecordNotFound {
		handler.Log.Warn("ExportGroupHandler.ServeHTTP group not found", "error", err.Error(), utils.LogRequestInfo(r))
		handler.ErrorHandler.PageNotFound(w, r)
		return
	}
	if err != nil {
		handler.Log.Error("ExportGroupHandler.ServeHTTP failed to fetch SeedsGroup data", "error", err.Error(), utils.LogRequestInfo(r))
		handler.ErrorHandler.InternalServerError(w, r)
		return
	}

	format := r.PathValue("format")
	switch format {
	case "excel":
		handler.RespondExcel(w, r, group)
	case "csv":
		handler.RespondCsv(w, r, group)
	default:
		handler.Log.Warn("ExportGroupHandler.ServeHTTP user requested unknown format", utils.LogRequestInfo(r))
		handler.ErrorHandler.PageNotFound(w, r)
		return
	}
}

func (handler *ExportGroupHandler) RespondExcel(w http.ResponseWriter, r *http.Request, group *entities.SeedsGroup) {
	buffer := new(bytes.Buffer)
	err := handler.ExporterService.GroupToExcel(group, buffer)
	if err != nil {
		handler.Log.Error("ExportGroupHandler.ServeHTTP got error from exporter service", "error", err.Error(), utils.LogRequestInfo(r))
		handler.ErrorHandler.InternalServerError(w, r)
		return
	}

	header := w.Header()
	const XlsxMimetype = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	header.Set(utils.ContentType, XlsxMimetype)
	filename := "export-" + group.ShadowID + ".xlsx"
	header.Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.WriteHeader(http.StatusOK)
	_, _ = buffer.WriteTo(w)
	handler.Log.Info("ExportGroupHandler.ServeHTTP sucessfully responded", utils.LogRequestInfo(r))
}

func (handler *ExportGroupHandler) RespondCsv(w http.ResponseWriter, r *http.Request, group *entities.SeedsGroup) {
	buffer := new(bytes.Buffer)
	err := handler.ExporterService.GroupToCsv(group, buffer)
	if err != nil {
		handler.Log.Error("ExportGroupHandler.ServeHTTP got error from exporter service", "error", err.Error(), utils.LogRequestInfo(r))
		handler.ErrorHandler.InternalServerError(w, r)
		return
	}

	header := w.Header()
	const CsvMimetype = "text/csv"
	header.Set(utils.ContentType, CsvMimetype)
	filename := "export-" + group.ShadowID + ".csv"
	header.Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.WriteHeader(http.StatusOK)
	_, _ = buffer.WriteTo(w)
	handler.Log.Info("ExportGroupHandler.ServeHTTP sucessfully responded", utils.LogRequestInfo(r))
}
