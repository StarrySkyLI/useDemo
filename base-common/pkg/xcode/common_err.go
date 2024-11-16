package xcode

import (
	"net/http"
)

var (
	UserNotFound         = New(http.StatusUnauthorized, "User not found. ")
	PlatformCodeNotFound = New(http.StatusInternalServerError, "Platform business code not found, please set debug mode")

	DbSelectErr = New(http.StatusInternalServerError, "Not found. ")
	DbUpdateErr = New(http.StatusInternalServerError, "Update fail. ")
	DbCreateErr = New(http.StatusInternalServerError, "Create fail. ")
	DbDeleteErr = New(http.StatusInternalServerError, "Delete fail. ")
)
