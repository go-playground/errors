package awserrors

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/go-playground/errors"
)

const (
	transient = "Transient"
)

func init() {
	errors.RegisterHelper(AWSErrors)
}

// AWSErrors helps classify io related errors
func AWSErrors(c errors.Chain, err error) (cont bool) {
	switch e := err.(type) {
	case awserr.BatchedErrors:
		_ = c.AddTypes(transient, "Batch")
		return

	case awserr.RequestFailure:
		_ = c.AddTypes(transient, "Request").AddTags(
			errors.T("status_code", e.StatusCode()),
			errors.T("request_id", e.RequestID()),
		)
		return

	case awserr.Error:
		_ = c.AddTypes("General", "Error").AddTags(errors.T("aws_error_code", e.Code()))
		return
	}
	return true
}
