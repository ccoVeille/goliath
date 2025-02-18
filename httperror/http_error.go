package httperror

import (
	"fmt"
	"net/http"

	"github.com/Talento90/goliath/appcontext"
	"github.com/Talento90/goliath/apperror"
)

// Problem Details for HTTP APIs - https://datatracker.ietf.org/doc/html/rfc7807
//
// "type" (string) - A URI reference [RFC3986] that identifies the
// problem type.  This specification encourages that, when
// dereferenced, it provide human-readable documentation for the
// problem type (e.g., using HTML [W3C.REC-html5-20141028]).  When
// this member is not present, its value is assumed to be
// "about:blank".
// "title" (string) - A short, human-readable summary of the problem
// type.  It SHOULD NOT change from occurrence to occurrence of the
// problem, except for purposes of localization (e.g., using
// proactive content negotiation; see [RFC7231], Section 3.4).
// "status" (number) - The HTTP status code ([RFC7231], Section 6)
// generated by the origin server for this occurrence of the problem.
// "detail" (string) - A human-readable explanation specific to this
// occurrence of the problem.
// "instance" (string) - A URI reference that identifies the specific
// occurrence of the problem.  It may or may not yield further
// information if dereferenced.
type ProblemDetails struct {
	Type     string                    `json:"type,omitempty"`
	Title    string                    `json:"title,omitempty"`
	Detail   string                    `json:"detail,omitempty"`
	Status   int                       `json:"status,omitempty"`
	Instance string                    `json:"instance,omitempty"`
	TraceID  string                    `json:"traceId,omitempty"`
	Errors   apperror.ValidationErrors `json:"errors,omitempty"`
}

func (pd ProblemDetails) Error() string {
	return fmt.Sprintf("%s: %s", pd.Type, pd.Title)
}

const UnknownErrorType = "internal_error"

func New(ctx appcontext.AppContext, err error, instance string) ProblemDetails {
	appError, ok := err.(apperror.AppError)

	if !ok {
		return ProblemDetails{
			Type:     UnknownErrorType,
			Title:    "An error ocurred, please contact support.",
			Status:   500,
			TraceID:  ctx.TraceID(),
			Instance: instance,
		}
	}

	return ProblemDetails{
		Type:     appError.Code(),
		Title:    appError.Error(),
		Detail:   appError.Detail(),
		Status:   mapAppErrorToHttpStatusCode(appError),
		Instance: instance,
		TraceID:  ctx.TraceID(),
		Errors:   appError.ValidationErrors(),
	}
}

func mapAppErrorToHttpStatusCode(appError apperror.AppError) int {
	switch appError.Type() {
	case apperror.Validation:
		return http.StatusBadRequest
	case apperror.NotFound:
		return http.StatusNotFound
	case apperror.Permission:
		return http.StatusForbidden
	case apperror.Unauthorized:
		return http.StatusUnauthorized
	case apperror.Conflict:
		return http.StatusConflict
	case apperror.Timeout:
		return http.StatusRequestTimeout
	case apperror.Cancelled:
		return http.StatusAccepted
	default:
		return http.StatusInternalServerError
	}
}
