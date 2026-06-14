package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zynu/server/pkg/config"
	"github.com/zynu/server/pkg/logger"
)

type VideoHandler struct {
	cfg *config.Config
	log *logger.Logger
}

func NewVideoHandler(cfg *config.Config, log *logger.Logger) *VideoHandler {
	return &VideoHandler{cfg: cfg, log: log}
}

func (h *VideoHandler) StreamRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.handleStream)
	return mux
}

func (h *VideoHandler) UploadRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/presign", h.handlePresignURL)
	return mux
}

// handleStream serves video segments for HLS streaming.
// Path: /stream/{videoId}/{quality}/{segment}
func (h *VideoHandler) handleStream(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	videoID := parts[0]
	quality := parts[1]
	segment := parts[2]

	allowedQualities := map[string]bool{"360p": true, "480p": true, "720p": true, "1080p": true}
	if !allowedQualities[quality] {
		http.Error(w, "invalid quality", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(h.cfg.StorageDir, videoID, quality, segment)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	if strings.HasSuffix(segment, ".m3u8") {
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	} else if strings.HasSuffix(segment, ".ts") {
		w.Header().Set("Content-Type", "video/mp2t")
	}

	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("X-Video-ID", videoID)
	http.ServeFile(w, r, filePath)

	h.log.Infof("Streamed %s/%s/%s", videoID, quality, segment)
}

// handlePresignURL returns a short-lived pre-signed URL for direct upload.
func (h *VideoHandler) handlePresignURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	videoID := r.URL.Query().Get("video_id")
	if videoID == "" {
		http.Error(w, "missing video_id", http.StatusBadRequest)
		return
	}

	expiresIn := 3600
	uploadURL := fmt.Sprintf("%s/upload/%s?expires=%d&token=presign_%s",
		h.cfg.BaseURL, videoID, expiresIn, videoID)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"upload_url":%q,"expires_in":%s}`, uploadURL, strconv.Itoa(expiresIn))
}
