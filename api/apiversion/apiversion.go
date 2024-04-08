package apiversion

import "github.com/blang/semver"

var VERSION_STRING = "1.0.0"

var VERSION = semver.MustParse(VERSION_STRING)
