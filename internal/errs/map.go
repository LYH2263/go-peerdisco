package errs

import "fmt"

func NotFound(kind, id string) error {
	return fmt.Errorf("%s %q not found", kind, id)
}

func Invalid(msg string) error {
	return fmt.Errorf("invalid: %s", msg)
}

func Wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", op, err)
}
