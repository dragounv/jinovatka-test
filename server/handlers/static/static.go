package static

import (
	"io/fs"
	"jinovatka/assert"
	"jinovatka/utils"
	"log/slog"
	"net/http"
)

// Handler for the "/static" route. Servers files from the fs.FS injected into the constructor function.
type StaticHandler struct {
	FileServer http.Handler
	Log        *slog.Logger
}

func NewStaticHandler(log *slog.Logger, fs fs.FS) *StaticHandler {
	assert.Must(log != nil, "NewStaticEmbedHandler: log can't be nil")
	assert.Must(fs != nil, "NewStaticEmbedHandler: fs can't be nil")
	fileServer := http.FileServerFS(fs)
	return &StaticHandler{
		FileServer: fileServer,
		Log:        log,
	}
}

func (handler *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.FileServer.ServeHTTP(w, r)
	slog.Info("StaticHandler responded", utils.LogRequestInfo(r))
}

func (handler *StaticHandler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /static/", handler)
}
