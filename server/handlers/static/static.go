package static

import (
	"io/fs"
	"jinovatka/assert"
	"jinovatka/utils"
	"log/slog"
	"net/http"
)

func NewStaticHandler(log *slog.Logger, fs fs.FS) *StaticHandler {
	assert.Must(log != nil, "NewStaticEmbedHandler: log can't be nil")
	fileServer := http.FileServerFS(fs)
	return &StaticHandler{
		FileServer: fileServer,
		Log:        log,
	}
}

type StaticHandler struct {
	FileServer http.Handler
	Log        *slog.Logger
}

func (handler *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.FileServer.ServeHTTP(w, r)
	slog.Info("StaticHandler responded", utils.LogRequestInfo(r))
}
