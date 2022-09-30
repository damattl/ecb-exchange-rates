package models

type apiErrorCode int32

const CURRENCY_NOT_FOUND_ERROR apiErrorCode = 0
const NO_ENTRY_FOUND_ERROR apiErrorCode = 1
const FUTURE_DATE_ERROR apiErrorCode = 2
const FORMAT_NOT_SUPPORTED_ERROR apiErrorCode = 3
const CONVERSION_ERROR apiErrorCode = 4

type APIError struct {
	Message string       `json:"message"`
	Code    apiErrorCode `json:"code"`
}
