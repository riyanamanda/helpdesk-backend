package ihs

import "errors"

var (
	ErrPatientNotFound    = errors.New("patient not found")
	ErrPatientNotEligible = errors.New("patient is not eligible for resubmission")
)
