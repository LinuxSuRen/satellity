package admin

import (
	"encoding/json"
	"godiscourse/internal/category"
	"godiscourse/internal/session"
	"godiscourse/internal/views"
	"net/http"

	"github.com/dimfeld/httptreemux"
)

type adminCategoryImpl struct {
	category category.CategoryDatastore
}

type categoryRequest struct {
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
	Position    int64  `json:"position"`
}

func RegisterAdminCategory(c category.CategoryDatastore, router *httptreemux.TreeMux) {
	impl := &adminCategoryImpl{category: c}

	router.POST("/admin/categories", impl.create)
	router.GET("/admin/categories", impl.index)
	router.POST("/admin/categories/:id", impl.update)
	router.GET("/admin/categories/:id", impl.show)
}

func (impl *adminCategoryImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	category, err := impl.category.Create(r.Context(), &category.Params{
		Name:        body.Name,
		Alias:       body.Alias,
		Description: body.Description,
		Position:    body.Position,
	})
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderCategory(w, r, category)
}

func (impl *adminCategoryImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	categories, err := impl.category.GetAll(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderCategories(w, r, categories)
}

func (impl *adminCategoryImpl) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}

	category, err := impl.category.Update(r.Context(), params["id"], &category.Params{
		Name:        body.Name,
		Alias:       body.Alias,
		Description: body.Description,
		Position:    body.Position,
	})
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderCategory(w, r, category)
}

func (impl *adminCategoryImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	category, err := impl.category.GetByID(r.Context(), params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderCategory(w, r, category)
}
