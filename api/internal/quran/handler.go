package quran

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const upstreamBase = "https://api.alquran.cloud/v1"

type Handler struct {
	client *http.Client
	cache  sync.Map // surahNum -> cachedEntry
}

type cachedEntry struct {
	data      json.RawMessage
	fetchedAt time.Time
}

const cacheTTL = 24 * time.Hour

func NewHandler() *Handler {
	return &Handler{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/surah/{num}", h.getSurah)
}

func (h *Handler) getSurah(w http.ResponseWriter, r *http.Request) {
	numStr := r.PathValue("num")
	num, err := strconv.Atoi(numStr)
	if err != nil || num < 1 || num > 114 {
		http.Error(w, "invalid surah number", http.StatusBadRequest)
		return
	}

	// Check cache
	if entry, ok := h.cache.Load(num); ok {
		ce := entry.(cachedEntry)
		if time.Since(ce.fetchedAt) < cacheTTL {
			w.Header().Set("Content-Type", "application/json")
			w.Write(ce.data)
			return
		}
	}

	// Fetch from upstream
	url := fmt.Sprintf("%s/surah/%d/editions/quran-uthmani,en.transliteration,en.sahih,bs.mlivo,tr.diyanet", upstreamBase, num)
	resp, err := h.client.Get(url)
	if err != nil {
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "read error", http.StatusBadGateway)
		return
	}

	// Cache the response
	h.cache.Store(num, cachedEntry{data: body, fetchedAt: time.Now()})

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
