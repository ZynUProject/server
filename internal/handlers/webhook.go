package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/zynu/server/pkg/config"
	"github.com/zynu/server/pkg/logger"
)

type WebhookHandler struct {
	cfg *config.Config
	log *logger.Logger
}

func NewWebhookHandler(cfg *config.Config, log *logger.Logger) *WebhookHandler {
	return &WebhookHandler{cfg: cfg, log: log}
}

func (h *WebhookHandler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/video", h.handleVideoWebhook)
	return mux
}

type webhookEvent struct {
	Type    string          `json:"type"`
	VideoID string          `json:"video_id"`
	Payload json.RawMessage `json:"payload"`
}

func (h *WebhookHandler) handleVideoWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var event webhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	h.log.Infof("Webhook received: type=%s video_id=%s", event.Type, event.VideoID)

	switch event.Type {
	case "video.processed":
		h.onVideoProcessed(event)
	case "video.failed":
		h.onVideoFailed(event)
	default:
		h.log.Warnf("Unknown webhook type: %s", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) onVideoProcessed(e webhookEvent) {
	h.log.Infof("Video ready: %s", e.VideoID)
}

func (h *WebhookHandler) onVideoFailed(e webhookEvent) {
	h.log.Warnf("Video processing failed: %s", e.VideoID)
}
