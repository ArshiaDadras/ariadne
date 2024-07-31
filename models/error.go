package models

import "errors"

var ErrNodeExists = errors.New("node already exists")
var ErrEdgeExists = errors.New("edge already exists")
var ErrNodeNotFound = errors.New("node not found")
var ErrEdgeNotFound = errors.New("edge not found")
