package controllers

import (
	"encoding/json"
	"godiscourse/internal/engine"
	"godiscourse/internal/middleware"
	"godiscourse/internal/model"
	"godiscourse/internal/session"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type commentImpl struct {
	poster engine.Poster
}

type commentRequest struct {
	TopicID string `json:"topic_id"`
	Body    string `json:"body"`
}

func registerComment(p engine.Poster, router *httptreemux.TreeMux) {
	impl := &commentImpl{poster: p}

	router.POST("/comments", impl.create)
	router.POST("/comments/:id", impl.update)
	router.POST("/comments/:id/delete", impl.destory)
	router.GET("/topics/:id/comments", impl.comments)
}

func (impl *commentImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body commentRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}

	if comment, err := impl.poster.CreateComment(r.Context(), &model.CommentInfo{
		UserID:  middleware.CurrentUser(r).UserID,
		TopicID: body.TopicID,
		Body:    body.Body,
	}); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderComment(w, r, comment)
	}
}

func (impl *commentImpl) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body commentRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}

	if comment, err := impl.poster.CreateComment(r.Context(), &model.CommentInfo{
		UserID:  middleware.CurrentUser(r).UserID,
		TopicID: params["id"],
		Body:    body.Body,
	}); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderComment(w, r, comment)
	}
}

func (impl *commentImpl) destory(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if err := impl.poster.DeleteComment(r.Context(), params["id"], middleware.CurrentUser(r).UserID); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}

func (impl *commentImpl) comments(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))

	if topic, err := impl.poster.GetTopicByID(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if topic == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if comments, err := impl.poster.GetCommentsByTopicID(r.Context(), topic.TopicID, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderComments(w, r, comments)
	}
}
