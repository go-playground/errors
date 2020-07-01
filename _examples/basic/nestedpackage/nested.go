package nestedpackagee

import (
	"io"

	"github.com/go-playground/errors/v5"
)

func GetUser(userID string) error {
	err := io.EOF
	return errors.Wrap(err, "failed to do something").AddTypes("Permanent").AddTags(errors.T("user_id", userID))
}
