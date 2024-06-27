package middleware

import (
	"net/http"
	"time"

	"github.com/MadAppGang/httplog"
	"github.com/hummerd/httpdump"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/urfave/negroni"
)

func HTTPLogger(component string, logParams model.LoggerParams) negroni.Handler {
	if logParams.Type == model.HTTPLogTypeNone ||
		logParams.Type == "" {
		return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			next(w, r)
		})
	}

	// default to json
	if logParams.Format == model.LogFormatJSON ||
		logParams.Format == "" {
		logger := logging.NewLogger(model.LogFormatJSON, logParams.Level)

		dumpReq := func(r *http.Request, body []byte) {
			logger.Debug("HTTP request",
				logging.FieldComponent, component,
				"method", r.Method,
				"url", r.URL,
				"headers", r.Header,
				"body", string(body))
		}

		dumpResp := func(r *http.Response, body []byte, duration time.Duration) {
			logger.Debug("HTTP response",
				logging.FieldComponent, component,
				"method", r.Request.Method,
				"url", r.Request.URL,
				"status", r.StatusCode,
				"headers", r.Header,
				"body", string(body),
				"duration", duration)
		}

		var opts []httpdump.Option

		if logParams.Type == model.HTTPLogTypeShort {
			// exclude body
			opts = append(opts, httpdump.WithRequestFilters(func(r *http.Request) (dump bool, body bool) {
				return true, false
			}))

			opts = append(opts, httpdump.WithResponseFilters(func(r *http.Request, headers http.Header, status int) (dump bool, body bool) {
				return true, false
			}))
		}

		hd := httpdump.NewMiddlewareWrapper(dumpReq, dumpResp, opts...)

		return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			hd(next).ServeHTTP(w, r)
		})
	}

	hl := httplog.LoggerWithFormatterAndName(component, httplog.DefaultLogFormatterWithRequestHeadersAndBody)

	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		hl(next).ServeHTTP(w, r)
	})
}
