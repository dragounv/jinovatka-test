package generator

import (
	"jinovatka/assert"
	"jinovatka/server/components"
	"jinovatka/utils"
	"log/slog"
	"net/http"
)

// The handler for citation generator page
type GeneratorHandler struct {
	Log *slog.Logger
}

func NewGeneratorHandler(log *slog.Logger) *GeneratorHandler {
	assert.Must(log != nil, "NewGeneratorHandler: log can't be nil")
	return &GeneratorHandler{
		Log: log,
	}
}

func (handler *GeneratorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := handler.View(w, r)
	if err != nil {
		handler.Log.Error("GeneratorHandler.ServeHTTP failed to render view", "error", err.Error(), utils.LogRequestInfo(r))
		return
	}
	handler.Log.Info("GeneratorHandler.ServeHTTP sucessfully responded", utils.LogRequestInfo(r))
}

func (handler *GeneratorHandler) View(w http.ResponseWriter, r *http.Request) error {
	return components.GeneratorView().Render(r.Context(), w)
}

func (handler *GeneratorHandler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /generator/", handler)
}
