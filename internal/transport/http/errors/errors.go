package errors

import "errors"

var ErrorInvalidFormat = errors.New("invalid format : expect form-data")
var ErrorFileNotFound = errors.New("not found file in request")
