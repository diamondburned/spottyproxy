package api

import (
	"context"
	"fmt"

	"github.com/diamondburned/hrt"
	"github.com/go-chi/chi/v5"
	"github.com/mileusna/useragent"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Secret   string `json:"secret"`
}

type LoginResponse struct {
	Token      string `json:"token"`
	DeviceName string `json:"deviceName"`
}

type loginHandler struct {
	*chi.Mux
	sessions *SessionStore
	secret   string
}

func newLoginHandler(sessions *SessionStore, secret string) *loginHandler {
	h := &loginHandler{
		Mux:      chi.NewMux(),
		sessions: sessions,
		secret:   secret,
	}

	h.Use(hrt.Use(hrt.DefaultOpts))
	h.Post("/", hrt.Wrap(h.login))

	return h
}

func (h *loginHandler) login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	request := hrt.RequestFromContext(ctx)

	ua := useragent.Parse(request.UserAgent())
	deviceName := extractDeviceName(ua)

	session, err := h.sessions.LoginAndAcquire(req.Username, req.Password, deviceName)
	if err != nil {
		return nil, err
	}
	session.Release()

	return &LoginResponse{
		Token:      session.token,
		DeviceName: deviceName,
	}, nil
}

func extractDeviceName(ua useragent.UserAgent) string {
	if ua.Device != "" {
		return ua.Device
	}
	if ua.Name != "" && ua.OS != "" {
		return fmt.Sprintf("%s on %s", ua.Name, ua.OS)
	}
	if ua.Name != "" {
		return ua.Name
	}
	if ua.OS != "" {
		return ua.OS
	}
	return ua.String
}
