package compile

import "errors"

var ErrNoRegistry = errors.New("registry not initialized")
var ErrNoModelObjects = errors.New("no model objects provided")
var ErrNoModelObject = errors.New("no model object provided")
