package saveurl

import "errors"

// ErrorParseRequestURI is returned when the input URL is not a valid URI.
var ErrorParseRequestURI = errors.New("parse request URI is failed")

// ErrorSaveToStorage is returned when the URL cannot be written to the storage backend.
var ErrorSaveToStorage = errors.New("save URL to storage is failed")

// ErrorUrlAlreadyExists is returned when the URL was already shortened.
var ErrorUrlAlreadyExists = errors.New("short URL already exists")
