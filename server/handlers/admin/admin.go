package admin

import (
	"jinovatka/assert"
	"jinovatka/server/components"
	"jinovatka/services"
	"jinovatka/utils"
	"log/slog"
	"net/http"
)

type AdminHandler struct {
	Log         *slog.Logger
	SeedService *services.SeedService
}

func NewAdminHandler(log *slog.Logger) *AdminHandler {
	assert.Must(log != nil, "NewAdminHanlder: log can't be nil")
	return &AdminHandler{
		Log: log,
	}
}

func (handler *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	records, err := handler.SeedService.FindSeeds(&services.FindSeedsArgs{})
	if err != nil {
		handler.Log.Error("FindSeeds failed", slog.String("error", err.Error()))
	}
	handler.View(w, r, components.NewAdminViewData(records, 0, 0)) // TODO: handle error
	handler.Log.Info("AdminHandler responded", utils.LogRequestInfo(r))
}

func (handler *AdminHandler) View(w http.ResponseWriter, r *http.Request, data *components.AdminViewData) error {
	return components.AdminView(data).Render(r.Context(), w)
}

func (handler *AdminHandler) Routes(mux *http.ServeMux) {
	mux.Handle("/admin/", handler)
}
