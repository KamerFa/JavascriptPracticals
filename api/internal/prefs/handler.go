package prefs

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/prefs", h.getPrefs)
	mux.HandleFunc("PUT /api/prefs", h.savePrefs)
	mux.HandleFunc("GET /api/bookmark", h.getBookmark)
	mux.HandleFunc("PUT /api/bookmark", h.setBookmark)
	mux.HandleFunc("DELETE /api/bookmark", h.deleteBookmark)
}

func (h *Handler) getPrefs(w http.ResponseWriter, r *http.Request) {
	sid := getSessionID(r)
	p, err := h.store.GetPrefs(sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, p)
}

func (h *Handler) savePrefs(w http.ResponseWriter, r *http.Request) {
	sid := getSessionID(r)
	var p Prefs
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := h.store.SavePrefs(sid, &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, p)
}

func (h *Handler) getBookmark(w http.ResponseWriter, r *http.Request) {
	sid := getSessionID(r)
	b, err := h.store.GetBookmark(sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if b == nil {
		writeJSON(w, nil)
		return
	}
	writeJSON(w, b)
}

func (h *Handler) setBookmark(w http.ResponseWriter, r *http.Request) {
	sid := getSessionID(r)
	var b Bookmark
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := h.store.SetBookmark(sid, &b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, b)
}

func (h *Handler) deleteBookmark(w http.ResponseWriter, r *http.Request) {
	sid := getSessionID(r)
	if err := h.store.DeleteBookmark(sid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]bool{"ok": true})
}

func getSessionID(r *http.Request) string {
	c, err := r.Cookie("session_id")
	if err != nil || c.Value == "" {
		return ""
	}
	return c.Value
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
