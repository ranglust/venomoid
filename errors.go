package venomoid

import "fmt"

var ErrorLookupAndFileMismatchAndAutomaticEnv = fmt.Errorf("key mismatch, either ConfigLookup or ConfigFile or AutomaticEnv must be supplied")
var ErrorMissingConfigFile = fmt.Errorf("config file not found")

type ErrorWrapper struct {
	InternalError error
	Label         string
}

func (e *ErrorWrapper) Error() string {
	return fmt.Sprintf("%s, error: %s", e.Label, e.InternalError.Error())
}
