package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/MadAppGang/httplog"
	"github.com/hummerd/httpdump"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/urfave/negroni"
)

func NegroniHTTPLogger(
	component string,
	format string,
	logParams model.LoggerParams,
	httpDetailing model.HTTPDetailing,
	exclude ...string,
) negroni.Handler {
	logger := HTTPLogger(component, format, logParams, httpDetailing, exclude...)

	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		logger(next).ServeHTTP(w, r)
	})
}

func emptyMiddleware(next http.Handler) http.Handler {
	return next
}

func HTTPLogger(
	component string,
	format string,
	logParams model.LoggerParams,
	httpDetailing model.HTTPDetailing,
	exclude ...string,
) func(http.Handler) http.Handler {
	if httpDetailing == model.HTTPLogNone ||
		httpDetailing == "" {
		return emptyMiddleware
	}

	// default to json
	if format == model.LogFormatJSON ||
		format == "" {
		logger := logging.NewLogger(model.LogFormatJSON, logParams.Level)

		dumpReq := func(r *http.Request, body []byte) {
			logger.Debug("HTTP request",
				logging.FieldComponent, component,
				"method", r.Method,
				"url", r.URL.String(),
				"headers", r.Header,
				"body", string(body))
		}

		dumpResp := func(r *http.Response, body []byte, duration time.Duration) {
			logger.Debug("HTTP response",
				logging.FieldComponent, component,
				"method", r.Request.Method,
				"url", r.Request.URL.String(),
				"status", r.StatusCode,
				"headers", r.Header,
				"body", string(body),
				"duration", duration)
		}

		var opts []httpdump.Option

		logBody := func(path string) bool {
			if httpDetailing != model.HTTPLogDump {
				return false
			}

			path = strings.ToLower(path)

			for _, e := range exclude {
				if strings.Contains(path, e) {
					return false
				}
			}

			return true
		}

		// exclude body
		opts = append(opts, httpdump.WithRequestFilters(func(r *http.Request) (dump bool, body bool) {
			return true, logBody(r.URL.Path)
		}))

		opts = append(opts, httpdump.WithResponseFilters(func(r *http.Request, headers http.Header, status int) (dump bool, body bool) {
			return true, logBody(r.URL.Path)
		}))

		hd := httpdump.NewMiddlewareWrapper(dumpReq, dumpResp, opts...)
		return hd
	}

	hl := httplog.LoggerWithFormatterAndName(component, httplog.DefaultLogFormatterWithRequestHeadersAndBody)
	return hl
}
