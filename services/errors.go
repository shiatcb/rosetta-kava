// Copyright 2021 Kava Labs, Inc.
// Copyright 2020 Coinbase, Inc.
//
// Derived from github.com/coinbase/rosetta-ethereum@f81889b
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"github.com/coinbase/rosetta-sdk-go/types"
)

var (
	// Errors contains all errors that could be returned
	// by this Rosetta implementation.
	Errors = []*types.Error{
		ErrUnimplemented,
		ErrUnavailableOffline,
		ErrKava,
	}

	// ErrUnimplemented is returned when an endpoint
	// is called that is not implemented.
	ErrUnimplemented = &types.Error{
		Code:    0,
		Message: "Endpoint not implemented",
	}

	// ErrUnavailableOffline is returned when an endpoint
	// is called that is not available offline.
	ErrUnavailableOffline = &types.Error{
		Code:    1,
		Message: "Endpoint unavailable offline",
	}

	// ErrKava is returned when kava
	// errors on a request.
	ErrKava = &types.Error{
		Code:    2,
		Message: "Kava error",
	}

	// ErrNoOperations is returned when no operations are provided
	ErrNoOperations = &types.Error{
		Code:    3,
		Message: "No operations provided",
	}

	// ErrInvalidCurrencyAmount is returned when a currency value could not be parsed
	ErrInvalidCurrencyAmount = &types.Error{
		Code:    4,
		Message: "Invalid currency",
	}

	// ErrUnsupportedCurrency is returned when a currency symbol is invalid
	// or the decimals do not match
	ErrUnsupportedCurrency = &types.Error{
		Code:    5,
		Message: "Unsupported concurrency",
	}

	// ErrUnclearIntent is returned when operations
	// provided in /construction/preprocess or /construction/payloads
	// are not valid.
	ErrUnclearIntent = &types.Error{
		Code:    6,
		Message: "Unable to parse intent",
	}

	// ErrInvalidAddress is returned when an account identifier has an invalid address
	ErrInvalidAddress = &types.Error{
		Code:    7,
		Message: "Invalid address",
	}

	// ErrInvalidOptions is returned by the metadata endpoint with invalid options
	ErrInvalidOptions = &types.Error{
		Code:    8,
		Message: "Invalid options",
	}
)

// wrapErr adds details to the types.Error provided. We use a function
// to do this so that we don't accidentially overrwrite the standard
// errors.
func wrapErr(rErr *types.Error, err error) *types.Error {
	newErr := &types.Error{
		Code:      rErr.Code,
		Message:   rErr.Message,
		Retriable: rErr.Retriable,
	}
	if err != nil {
		newErr.Details = map[string]interface{}{
			"context": err.Error(),
		}
	}

	return newErr
}
