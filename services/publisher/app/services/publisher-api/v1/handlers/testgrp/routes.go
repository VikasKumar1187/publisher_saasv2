package testgrp

import (
	"net/http"

	"github.com/vikaskumar1187/publisher_saasv2/services/publisher/business/web/v1/auth"
	"github.com/vikaskumar1187/publisher_saasv2/services/publisher/business/web/v1/mid"
	"github.com/vikaskumar1187/publisher_saasv2/services/publisher/foundation/logger"
	"github.com/vikaskumar1187/publisher_saasv2/services/publisher/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	UsingWeaver bool
	Build       string
	Log         *logger.Logger
	Auth        *auth.Auth
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	testgrp := New(cfg.Build, cfg.Log, cfg.Auth)

	app.Handle(http.MethodGet, version, "/test", testgrp.Test)
	app.Handle(http.MethodGet, version, "/testauth", testgrp.TestAuth, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth))
}
