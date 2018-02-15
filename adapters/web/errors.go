package web

import (
	"fmt"
	"net/http"

	"github.com/r3code/clean-go/engine"
)

var (
	// ErrInvalidInputParams show that params not fit the request form
	ErrInvalidInputParams = NewHTTPAPIError(http.StatusBadRequest,
		"Invalid input params", "InvalidInputParams", "")
)

// HTTPAPIError presents App Error at web interface
// enrich existent error with extra data requered for http response
type HTTPAPIError struct {
	Code      int    `json:"code"`
	Status    string `json:"status"`
	Message   string `json:"reason"`
	Reason    string `json:"message"`
	MoreInfo  string `json:"more_info,omitempty"`
	DebugInfo string `json:"debug_info,omitempty"`
}

func (h HTTPAPIError) Error() string {
	return fmt.Sprintf("%#v", h)
}

// NewHTTPAPIError creates new api error object
func NewHTTPAPIError(code int, message string, reason string, moreInfo string) *HTTPAPIError {
	return &HTTPAPIError{
		Code:      code,
		Status:    http.StatusText(code),
		Message:   message,
		Reason:    reason,
		MoreInfo:  moreInfo,
		DebugInfo: "",
	}
}

type apiError map[error]*HTTPAPIError

func (a *apiError) RegisterError(appErr error, code int, reason string, moreInfo string) {
	(*a)[appErr] = NewHTTPAPIError(code, appErr.Error(), reason, moreInfo)
}

var apiErrors = make(apiError)

func registerHTTPAPIErrors() {
	apiErrors.RegisterError(engine.ErrGreetingAlredyExists, http.StatusBadRequest, "GreetingAlredyExists", "/help/api/#1")
	apiErrors.RegisterError(engine.ErrForbiddenWordPresent, http.StatusBadRequest, "ForbiddenWordPresent ", "/help/api/#2")
	apiErrors.RegisterError(engine.ErrNotFound, http.StatusBadRequest, "NotFound", "/help/api/not-found")
	apiErrors.RegisterError(engine.ErrInvalidRequestData, http.StatusBadRequest, "InvalidRequestData ", "/help/api/#3")

}

// GetHTTPAPIError creates from application error an HTTPAPIError that has information for http response to be used in the web-adapter
func GetHTTPAPIError(appError error) (apiErr *HTTPAPIError, errFound bool) {
	apiErr, errFound = apiErrors[appError]
	return apiErr, errFound
}

func init() {
	registerHTTPAPIErrors()
}
