package testgrp

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/vikaskumar1187/publisher_saasv2/services/publisher/business/web/v1/auth"
	"github.com/vikaskumar1187/publisher_saasv2/services/publisher/foundation/logger"
	"github.com/vikaskumar1187/publisher_saasv2/services/publisher/foundation/web"
)

// Handlers manages the set of check endpoints.
type Handlers struct {
	build string
	log   *logger.Logger
	auth  *auth.Auth
}

// New constructs a Handlers api for the check group.
func New(build string, log *logger.Logger, auth *auth.Auth) *Handlers {
	return &Handlers{
		build: build,
		log:   log,
		auth:  auth,
	}
}

// Test is our example route.
func (h *Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	statusstr := "OK"
	code := http.StatusOK

	if n := rand.Intn(100); n%2 == 0 {
		statusstr = "NOT OK"
		code = http.StatusInternalServerError
	}

	status := struct {
		Status string
	}{
		Status: statusstr,
	}

	return web.Respond(ctx, w, status, code)
}

// Test is our example route.
func (h *Handlers) TestAuth(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	statusstr := "OK"
	code := http.StatusOK

	status := struct {
		Status string
	}{
		Status: statusstr,
	}

	return web.Respond(ctx, w, status, code)
}
