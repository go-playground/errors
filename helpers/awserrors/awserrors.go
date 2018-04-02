package awserrors

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/go-playground/errors"
)

// AWSErrors helps classify io related errors
func AWSErrors(c errors.Chain, err error) (cont bool) {
	switch e := err.(type) {
	case awserr.BatchError:
		c.WithTypes("Transient", "Batch").WithTags(errors.T("aws_error_code", e.Code()))
		return

	case awserr.BatchedErrors:
		c.WithTypes("Transient", "Batch")
		return

	case awserr.RequestFailure:
		c.WithTypes("Transient", "Request").WithTags(
			errors.T("status_code", e.StatusCode()),
			errors.T("request_id", e.RequestID()),
		)
		return

	case awserr.Error:
		c.WithTypes("General", "Error").WithTags(errors.T("aws_error_code", e.Code()))
		return
	}
	return true
}
