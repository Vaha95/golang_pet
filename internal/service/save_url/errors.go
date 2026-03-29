package saveurl

import "errors"

var ErrorParseRequestURI = errors.New("parse request URI is failed")
var ErrorSaveToStorage = errors.New("save URL to storage is failed")
var ErrorUrlAlreadyExists = errors.New("short URL already exists")