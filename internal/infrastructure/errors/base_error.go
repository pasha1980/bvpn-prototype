package errors

import (
	"bvpn-prototype/internal/infrastructure/logger"
	"fmt"
)

type ErrorType string

type Error struct {
	Code int
	Type ErrorType
	Data any
}

func (e Error) Error() string {
	return string(e.Type)
}

var layers = map[string]string{
	"0": "infrastructure",
	"1": "domain",
	"2": "protocol",
}

func (e Error) Layer() string {
	code := fmt.Sprintf("%04d", e.Code)
	layer, ok := layers[code[1:2]]
	if !ok {
		return "undefined"
	}

	return layer
}

var domains = map[string]string{
	"0": "NO",
	"1": "CHAIN",
	"2": "PEER",
	"3": "VPN",
}

func (e Error) Domain() string {
	code := fmt.Sprintf("%04d", e.Code)

	domain, ok := domains[code[2:3]]
	if !ok {
		return "UNDEFINED"
	}

	return domain
}

var levels = map[string]string{
	"0": "WARNING",
	"1": "CRITICAL",
}

func (e Error) Level() string {
	code := fmt.Sprintf("%04d", e.Code)
	level, ok := levels[code[3:4]]
	if !ok {
		return "UNDEFINED"
	}
	return level
}

func (e Error) Log() {
	code := fmt.Sprintf("%04d", e.Code)
	if code[3:4] == "0" {
		return
	}

	logger.LogError(
		fmt.Sprintf(
			"[%s - %s] D:%s %s",
			e.Level(),
			e.Layer(),
			e.Domain(),
			string(e.Type),
		),
	)
}
