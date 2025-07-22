package group

import (
	"bytes"
	"jinovatka/assert"
	"jinovatka/server/handlers/httperror"
	"jinovatka/services"
	"jinovatka/utils"
	"log/slog"
	"net/http"
	"net/url"

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
	// We need Host header to generate URLs
	// TODO: In production this must be taken from config. Don't rely on the host header, anything can be there.
	if r.Host == "" {
		handler.Log.Error("ExportGroupHandler.ServeHTTP missing Host header in request", utils.LogRequestInfo(r))
		handler.ErrorHandler.InternalServerError(w, r)
		return
	}
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

	buffer := new(bytes.Buffer)
	// Remeber to add the http prefix, otherwise the URL library will fail silently!
	urlPrefix, err := url.Parse("http://" + r.Host + "/seed/")
	if err != nil {
		handler.Log.Error("ExportGroupHandler.ServeHTTP colud not parse r.Host to URL", "error", err.Error(), utils.LogRequestInfo(r))
		handler.ErrorHandler.InternalServerError(w, r)
		return
	}
	err = handler.ExporterService.GroupToExcel(group, buffer, urlPrefix)
	if err != nil {
		handler.Log.Error("ExportGroupHandler.ServeHTTP got error from exporter service", "error", err.Error(), utils.LogRequestInfo(r))
		handler.ErrorHandler.InternalServerError(w, r)
		return
	}

	header := w.Header()
	const XlsxMimetype = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	header.Set(utils.ContentType, XlsxMimetype)
	filename := "seminka-" + group.ShadowID + ".xlsx"
	header.Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.WriteHeader(http.StatusOK)
	_, _ = buffer.WriteTo(w)
	handler.Log.Info("ExportGroupHandler.ServeHTTP sucessfully responded", utils.LogRequestInfo(r))
}
