package routemobile

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/madappgang/identifo/v2/model"
)

const baseURLuae = "https://sms.rmlconnect.net%s"

// SMSService sends SMS via RouteMobile service.
type SMSService struct {
	username   string
	password   string
	source     string
	baseURL    string
	httpClient *http.Client
}

// NewSMSService creates, inits and returns RouteMobile-backed SMS service.
func NewSMSService(settings model.RouteMobileServiceSettings) (*SMSService, error) {
	s := &SMSService{
		username: settings.Username,
		password: settings.Password,
		source:   settings.Source,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
	switch {
	case settings.Region == model.RouteMobileRegionUAE:
		s.baseURL = baseURLuae
	default:
		return nil, fmt.Errorf("Unknown RouteMobile region %s", settings.Region)
	}
	return s, nil
}

// SendSMS sends SMS messages using RouteMobile service.
func (ss *SMSService) SendSMS(recipient, message string) error {
	queryParams := fmt.Sprintf("username=%s&password=%s&type=0&dlr=0&destination=%s&source=%s&message=%s", ss.username, ss.password, strings.TrimPrefix(recipient, "+"), url.QueryEscape(ss.source), url.QueryEscape(message))
	url := fmt.Sprintf(ss.baseURL, fmt.Sprintf("/bulksms/bulksms?%s", queryParams))

	resp, err := ss.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respString := string(respBytes)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("%s. %d", respString, resp.StatusCode)
	}
	if !strings.HasPrefix(respString, "1701") {
		return fmt.Errorf("Error from RouteMobile API: '%s'. Please refer to the RouteMobile documentation", respString)
	}
	return nil
}
