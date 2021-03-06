package goselenium

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

// NewSeleniumWebDriver creates a new instance of a Selenium web driver with a
// service URL (usually http://domain:port/wd/hub) and a Capabilities object.
// This method will return validation errors if the Selenium URL is invalid or
// the required capabilities (BrowserName) are not set.
func NewSeleniumWebDriver(serviceURL string, capabilities Capabilities) (WebDriver, error) {
	if serviceURL == "" {
		return nil, errors.New("Provided Selenium URL is invalid")
	}

	urlValid := strings.HasPrefix(serviceURL, "http://") || strings.HasPrefix(serviceURL, "https://")
	if !urlValid {
		return nil, errors.New("Provided Selenium URL is invalid.")
	}

	browser := capabilities.Browser()
	hasBrowserCapability := browser.BrowserName() != ""
	if !hasBrowserCapability {
		return nil, errors.New("An invalid capabilities object was provided.")
	}

	if strings.HasSuffix(serviceURL, "/") {
		serviceURL = strings.TrimSuffix(serviceURL, "/")
	}

	driver := &seleniumWebDriver{
		seleniumURL:  serviceURL,
		capabilities: &capabilities,
		apiService:   &seleniumAPIService{},
	}

	return driver, nil
}

// SessionScriptTimeout creates an appropriate Timeout implementation for the
// script timeout.
func SessionScriptTimeout(to int) Timeout {
	return timeout{
		timeoutType: "script",
		timeout:     to,
	}
}

// SessionPageLoadTimeout creates an appropriate Timeout implementation for the
// page load timeout.
func SessionPageLoadTimeout(to int) Timeout {
	return timeout{
		timeoutType: "page load",
		timeout:     to,
	}
}

// SessionImplicitWaitTimeout creates an appropriate timeout implementation for the
// session implicit wait timeout.
func SessionImplicitWaitTimeout(to int) Timeout {
	return timeout{
		timeoutType: "implicit",
		timeout:     to,
	}
}

// ByIndex accepts an integer that represents what the index of an element is
// and returns the appropriate By implementation.
func ByIndex(index uint) By {
	return by{
		t:     "index",
		value: index,
	}
}

// ByCSSSelector accepts a CSS selector (i.e. ul#id > a) for use in the
// FindElement(s) functions.
func ByCSSSelector(selector string) By {
	return by{
		t:     "css selector",
		value: selector,
	}
}

// ByID is used to find a element by its ID
func ByID(id string) By {
	return by{
		t:     "id",
		value: id,
	}
}

// ByLinkText is used to find an anchor element by its innerText.
func ByLinkText(text string) By {
	return by{
		t:     "link text",
		value: text,
	}
}

// ByPartialLinkText works the same way as ByLinkText but performs a search
// where the link text contains the string passed in instead of a full match.
func ByPartialLinkText(text string) By {
	return by{
		t:     "partial link text",
		value: text,
	}
}

// ByXPath utilises the xpath to find elements (see http://www.guru99.com/xpath-selenium.html).
func ByXPath(path string) By {
	return by{
		t:     "xpath",
		value: path,
	}
}

type seleniumWebDriver struct {
	seleniumURL  string
	sessionID    string
	capabilities *Capabilities
	apiService   apiServicer
}

func (s *seleniumWebDriver) DriverURL() string {
	return s.seleniumURL
}

func (s *seleniumWebDriver) stateRequest(req *request) (*stateResponse, error) {
	response := &stateResponse{
		Status: -1,
	}
	var err error

	resp, err := s.apiService.performRequest(req.url, req.method, req.body)
	if err != nil {
		return nil, newCommunicationError(err, req.callingMethod, req.url, resp)
	}

	err = json.Unmarshal(resp, response)
	if err != nil {
		return nil, newUnmarshallingError(err, req.callingMethod, string(resp))
	}

	response.State = s.convertStatusToStat(response.Status)
	return response, nil
}

func (s *seleniumWebDriver) convertStatusToStat(status int) string {
	statMapping := map[int]string{
		-1: "UNKNOWNE STATE",
		0:  "SUCCESS",
		6:  "NO_SUCH_SESSION",
		7:  "NO_SUCH_ELEMENT",
		8:  "NO_SUCH_FRAME",
		9:  "UNKNOWN_COMMAND",
		10: "STALE_ELEMENT_REFERENCE",
		11: "ELEMENT_NOT_VISIBLE",
		12: "INVALID_ELEMENT_STATE",
		13: "UNHANDLED_ERROR",
		15: "ELEMENT_NOT_SELECTABLE",
		17: "JAVASCRIPT_ERROR",
		19: "XPATH_LOOKUP_ERROR",
		21: "TIMEOUT",
		23: "NO_SUCH_WINDOW",
		24: "INVALID_COOKIE_DOMAIN",
		25: "UNABLE_TO_SET_COOKIE",
		26: "UNEXPECTED_ALERT_PRESENT",
		27: "NO_ALERT_PRESENT",
		28: "ASYNC_SCRIPT_TIMEOUT",
		29: "INVALID_ELEMENT_COORDINATES",
		30: "IME_NOT_AVAILABLE",
		31: "IME_ENGINE_ACTIVATION_FAILED",
		32: "INVALID_SELECTOR_ERROR",
		33: "SESSION_NOT_CREATED",
		34: "MOVE_TARGET_OUT_OF_BOUNDS",
		51: "INVALID_XPATH_SELECTOR",
		52: "INVALID_XPATH_SELECTOR_RETURN_TYPER",
		60: "ELEMENT_NOT_INTERACTABLE",
		61: "INVALID_ARGUMENT",
		62: "NO_SUCH_COOKIE",
		63: "UNABLE_TO_CAPTURE_SCREEN",
		64: "ELEMENT_CLICK_INTERCEPTED",
	}
	return statMapping[status]

}

func (s *seleniumWebDriver) valueRequest(req *request) (*valueResponse, error) {
	var response valueResponse
	var err error

	resp, err := s.apiService.performRequest(req.url, req.method, req.body)
	if err != nil {
		return nil, newCommunicationError(err, req.callingMethod, req.url, resp)
	}

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, newUnmarshallingError(err, req.callingMethod, string(resp))
	}

	return &response, nil
}

func (s *seleniumWebDriver) elementRequest(req *elRequest) ([]byte, error) {
	b := map[string]interface{}{
		"using": req.by.Type(),
		"value": req.by.Value(),
	}
	bJSON, err := json.Marshal(b)
	if err != nil {
		return nil, newMarshallingError(err, req.callingMethod, bJSON)
	}

	body := bytes.NewReader(bJSON)
	resp, err := s.apiService.performRequest(req.url, req.method, body)
	if err != nil {
		return nil, newCommunicationError(err, req.callingMethod, req.url, resp)
	}

	return resp, nil
}

func (s *seleniumWebDriver) scriptRequest(script string, url string, method string) (*ExecuteScriptResponse, error) {
	r := map[string]interface{}{
		"script": script,
		"args":   []string{""},
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, newMarshallingError(err, method, r)
	}
	body := bytes.NewReader(b)
	resp, err := s.valueRequest(&request{
		url:           url,
		method:        "POST",
		body:          body,
		callingMethod: method,
	})
	if err != nil {
		return nil, err
	}

	return &ExecuteScriptResponse{State: resp.State, Response: resp.Value}, nil
}

type timeout struct {
	timeoutType string
	timeout     int
}

func (t timeout) Type() string {
	return t.timeoutType
}

func (t timeout) Timeout() int {
	return t.timeout
}

type request struct {
	url           string
	method        string
	body          io.Reader
	callingMethod string
}

type elRequest struct {
	url           string
	by            By
	method        string
	callingMethod string
}

type stateResponse struct {
	State  string `json:"state"`
	Status int    `json:"status"`
}

type valueResponse struct {
	State string `json:"state"`
	Value string `json:"value"`
}

type by struct {
	t     string
	value interface{}
}

func (b by) Type() string {
	return b.t
}

func (b by) Value() interface{} {
	return b.value
}
