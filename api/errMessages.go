package api

import "errors"

var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrAccountHasTransfers = errors.New("cannot delete account because it has related transfers")
	ErrInternal            = errors.New("internal server error")
	ErrInvalidID           = errors.New("invalid id")
	ErrInvalidReq          = errors.New("invalid request parameters")
)
