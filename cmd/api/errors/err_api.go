package errors

import (
	"errors"
	"net/http"
	"strconv"
)

type ErrAPI interface {
	GetHTTPStatusCode() int
	GetHTTPStatus() string
	GetMessage() string
	ToString() string
}

type errAPIImpl struct {
	Message        string `json:"message"`
	HTTPStatus     string `json:"status"`
	HTTPStatusCode int    `json:"status_code"`
}

func NewErrAPI(httpStatusCode int, httpStatus string, err error) ErrAPI {
	return &errAPIImpl{
		HTTPStatusCode: httpStatusCode,
		HTTPStatus:     httpStatus,
		Message:        err.Error(),
	}
}

func (e errAPIImpl) GetHTTPStatusCode() int {
	return e.HTTPStatusCode
}

func (e errAPIImpl) GetHTTPStatus() string {
	return e.HTTPStatus
}

func (e errAPIImpl) GetMessage() string {
	return e.Message
}

func (e errAPIImpl) ToString() string {
	return strconv.Itoa(e.GetHTTPStatusCode()) + " - " + e.GetMessage()
}

func NewErrAPIInternalServer(err error) ErrAPI {
	code := http.StatusInternalServerError
	return NewErrAPI(code, http.StatusText(code), err)
}

func NewErrAPIInternalServerFromErrAPI(errApi ErrAPI) ErrAPI {
	return NewErrAPIInternalServer(errors.New(strconv.Itoa(errApi.GetHTTPStatusCode()) + " - " + errApi.GetMessage()))
}

func NewBuildAPIRequestAPIError(err error) ErrAPI {
	return NewErrAPIInternalServer(errors.New("Error while trying to build API request: " + err.Error()))
}

func NewErrAPITimeout(err error) ErrAPI {
	const httpTimeoutStatusCode = 499
	return NewErrAPI(httpTimeoutStatusCode, "Timeout", err)
}

func NewErrAPINotFound(err error) ErrAPI {
	code := http.StatusNotFound
	return NewErrAPI(code, http.StatusText(code), err)
}

func NewErrAPIBadRequest(err error) ErrAPI {
	code := http.StatusBadRequest
	return NewErrAPI(code, http.StatusText(code), err)
}

func NewErrAPIUnauthorized(err error) ErrAPI {
	code := http.StatusUnauthorized
	return NewErrAPI(code, http.StatusText(code), err)
}

func NewErrAPIFieldCantBeUpdated(model string) ErrAPI {
	return NewErrAPIBadRequest(errors.New("'" + model + "' can't be updated"))
}

func NewErrAPIMustNotBeBlank(field string) ErrAPI {
	return NewErrAPIBadRequest(errors.New("field '" + field + "' must not be blank"))
}

func NewErrAPIMustNotBeInTheBody(field string) ErrAPI {
	return NewErrAPIBadRequest(errors.New("field '" + field + "' must not be in the body"))
}

func NewErrAPIFieldValueAlreadyExist(field, value string) ErrAPI {
	return NewErrAPIBadRequest(errors.New("field '" + field + "' with value '" + value + "' already exist"))
}
