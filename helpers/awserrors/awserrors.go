package awserrors

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/go-playground/errors"
)

// AWSErrors helps classify io related errors
func AWSErrors(w *errors.Wrapped, err error) (cont bool) {
	switch e := err.(type) {
	case awserr.BatchError:
		w.WithTypes("Transient", "Batch").WithTags(errors.T("aws_error_code", e.Code()))
		return

	case awserr.BatchedErrors:
		w.WithTypes("Transient", "Batch")
		return

	case awserr.RequestFailure:
		w.WithTypes("Transient", "Request").WithTags(
			errors.T("status_code", e.StatusCode()),
			errors.T("request_id", e.RequestID()),
		)
		return

	case awserr.Error:
		w.WithTypes("General", "Error").WithTags(errors.T("aws_error_code", e.Code()))
		return
	}
	return true
}
