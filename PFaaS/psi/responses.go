package psi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func AlwaysPanic(message any) HandlerFunc {
	return func(r *http.Request) StatusCode {
		panic(message)
	}
}

func Always[T any, PT interface {
	StatusCode
	*T
}]() HandlerFunc {
	var zeroT T
	return func(r *http.Request) StatusCode {
		return PT(&zeroT)
	}
}

type HandlerFunc func(r *http.Request) StatusCode

func Responder(f func(r *http.Request) StatusCode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(f(r), w)
	}
}

func WriteResponse(resp StatusCode, w http.ResponseWriter) {
	status := resp.StatusCode()
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	payload := make(map[string]any)

	valueInfo := reflect.ValueOf(resp)
	for valueInfo.Kind() == reflect.Ptr {
		valueInfo = valueInfo.Elem()
	}
	typeInfo := valueInfo.Type()
	switch data := resp.(type) {
	case error:
		payload["error"] = typeInfo.Name()
		payload["message"] = data.Error()
	default:
		numFields := typeInfo.NumField()
		for i := 0; i < numFields; i++ {
			field := typeInfo.Field(i)

			name, fromTag := field.Tag.Lookup("json")
			if field.Anonymous {
				continue
			}
			if !fromTag {
				name = field.Name
			} else if name == "-" {
				continue
			}

			payload[name] = valueInfo.Field(i).Interface()
		}
	}

	// here to ensure status was not overwritten
	payload["status"] = status
	if err := enc.Encode(payload); err != nil {
		panic(EncodingError{err})
	}
}

type EncodingError struct {
	Inner error
}

func (ee *EncodingError) Error() string {
	return ee.Inner.Error()
}

type StatusCode interface {
	StatusCode() int
}

type RestResponse struct {
	Status int `json:"status"`
}

func (rr *RestResponse) StatusCode() int {
	return rr.Status
}

type NoData struct{}

func (nc *NoData) StatusCode() int {
	return http.StatusNoContent
}

type OkData struct{}

func (os *OkData) StatusCode() int {
	return http.StatusOK
}

type BadRequestError struct {
	Inner error `json:"-"`
}

type NotFoundError struct{}

func (e *NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

func (e *NotFoundError) Error() string {
	return "the requested resource was not found"
}

type MethodNotAllowedError struct{}

func (e *MethodNotAllowedError) StatusCode() int {
	return http.StatusMethodNotAllowed
}

func (e *MethodNotAllowedError) Error() string {
	return "the requested method is not allowed for this resource"
}

type UnauthorizedError struct{}

func (e *UnauthorizedError) StatusCode() int {
	return http.StatusUnauthorized
}

func (e *UnauthorizedError) Error() string {
	return "authorization was not provided"
}

type ForbiddenError struct {
	Username string
}

func (e *ForbiddenError) StatusCode() int {
	return http.StatusForbidden
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("provided auth (%s) is not authorized to perform this action", e.Username)
}

type NotAcceptableError struct{}

func (e *NotAcceptableError) StatusCode() int {
	return http.StatusNotAcceptable
}

func (e *NotAcceptableError) Error() string {
	return "this server cannot provide a response that satisfies this request"
}
