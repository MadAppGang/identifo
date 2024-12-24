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
	maxBodySize int,
	logParams model.LoggerParams,
	httpDetailing model.HTTPDetailing,
	excludeAuth bool,
	exclude ...string,
) negroni.Handler {
	logger := HTTPLogger(component, format, maxBodySize, logParams, httpDetailing, excludeAuth, exclude...)

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
	maxBodySize int,
	logParams model.LoggerParams,
	httpDetailing model.HTTPDetailing,
	excludeAuth bool,
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
				"headers", redactHeaders(r.Header, excludeAuth),
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

		if maxBodySize <= 0 {
			maxBodySize = httpdump.DefaultBodySize
		}

		opts = append(opts, httpdump.WithLimitedBody(maxBodySize))

		hd := httpdump.NewMiddlewareWrapper(dumpReq, dumpResp, opts...)
		return hd
	}

	hl := httplog.LoggerWithFormatterAndName(component, httplog.DefaultLogFormatterWithRequestHeadersAndBody)
	return hl
}

func redactHeaders(headers http.Header, excludeAuth bool) http.Header {
	if !excludeAuth {
		return headers
	}

	result := make(http.Header, len(headers))

	for k, vv := range headers {
		if strings.EqualFold(k, "Authorization") {
			cc := make([]string, len(vv))
			for i, v := range vv {
				cc[i] = redactAuthValue(v)
			}
			result[k] = cc
		} else {
			result[k] = vv
		}

	}

	return result
}

func redactAuthValue(v string) string {
	expectedPrefix := "bearer"

	actualPrefix := ""
	if len(v) >= len(expectedPrefix) {
		actualPrefix = v[:len(expectedPrefix)]
	}

	if strings.EqualFold(actualPrefix, expectedPrefix) {
		if len(v) <= len(expectedPrefix)+1 {
			return actualPrefix + " <empty>"
		}

		return actualPrefix + " <redacted>"
	}

	return "<redacted>"
}
