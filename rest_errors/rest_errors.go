package rest_errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RestErr interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
}

type restErr struct {
	ErrMessage string         `json:"message"`
	ErrStatus  int            `json:"status"`
	ErrError   string         `json:"error"`
	ErrCauses  *[]interface{} `json:"causes,omitempty"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: %v",
		e.ErrMessage, e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) Status() int {
	return e.ErrStatus
}

func (e restErr) Causes() []interface{} {
	return *e.ErrCauses
}

func NewRestError(message string, status int, err string, causes []interface{}) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
		ErrCauses:  &causes,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr restErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

func NewBadRequestError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "bad_request",
	}
}

func NewNotFoundError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not_found",
	}
}

func NewUnauthorizedError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   "unauthorized",
	}
}

func NewInternalServerError(message string, err error) RestErr {
	result := restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "internal_server_error",
	}
	if err != nil {
		*result.ErrCauses = append(*result.ErrCauses, err.Error())
	}
	return result
}
func NewMissingPrimayKey(message string) RestErr {
	result := restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "Missing field Id",
	}
	return result
}

type InternalErrors int

const (
	NoRecordFound                   = 452
	InvalidAuthToken InternalErrors = iota + 100
	InvalidEmail
	EmailAndMobileIsEmpty
	InvalidPassword
	InvalidPhone
	InternalServerErrContactAdmin
	OldPasswordIsNotMatching
	FailedWhileUpdatepassword
	InvalidUserIdOrPassword
	UnableToGenerateToken
	InvalidRequestHeader
	DuplicateEntry
	LongData
	IncorrectDateTime
	InvalidForeignKey
	InvalidOperationHours
	InvalidPriceFormat
	MembershipTypeNotFound
	InvalidSlotID
	SlotNotAvailable
	InvalidNamespace
	InvalidReferenceID
	BookingTooAdvance
	InvalidProductIndexFieldMapping //Invalid product index field mapping (bleve while creating schema)
	InvalidJsonBody                 //invalid json body
	IndexSchemaWasNotCreated        //schema was not created
	IndexWasNotCreated              //index was not created
	IndexDocumentIdMissing          //index document id must provided
)
