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
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kava-labs/rosetta-kava/kava"

	"github.com/coinbase/rosetta-sdk-go/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const (
	defaultSuggestedFeeMultiplier = float64(1)
	defaultGasAdjustment          = float64(0.5)
)

// ConstructionPreprocess implements the /construction/preprocess endpoint.
func (s *ConstructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	request *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	if len(request.Operations) == 0 {
		return nil, ErrNoOperations
	}

	// TODO: improve operation parsing -- very basic for first pass
	//
	// currently, only supports a single transfer with one currency
	// should support multiple transfers, multiple currencies, and
	// staking operations
	//
	// in addition, parsing logic needs to be refactored with improved
	// testing around invalid cases, and related operations
	msgs, rerr := parseOperationMsgs(request.Operations)
	if rerr != nil {
		return nil, rerr
	}

	txBody := tx.TxBody{}
	anys, err := convertMsgsToAnys(msgs)
	if err != nil {
		return nil, wrapErr(ErrKava, err)
	}
	txBody.Messages = anys
	txBody.Memo = getMemoFromMetadata(request.Metadata)

	encodedTxBody, err := s.encodingConfig.Marshaler.MarshalJSON(&txBody)
	if err != nil {
		return nil, wrapErr(ErrKava, err)
	}

	options := map[string]interface{}{
		"tx_body":                  string(encodedTxBody),
		"gas_adjustment":           getGasAdjustmentFromMetadata(request.Metadata),
		"suggested_fee_multiplier": suggestedMultiplerOrDefault(request.SuggestedFeeMultiplier),
	}

	// TODO: can improve to include other fee options such as payer
	encodedMaxFee, rerr := getMaxFeeAndEncodeOption(request.MaxFee)
	if rerr != nil {
		return nil, rerr
	}
	if encodedMaxFee != nil {
		options["max_fee"] = *encodedMaxFee
	}

	requiredPublicKeys := []*types.AccountIdentifier{}

	for _, msg := range msgs {
		signers := msg.GetSigners()
		seenSigners := make(map[string]bool)

		// TODO: add test cases for multiple signers
		for _, signer := range signers {
			addr := signer.String()

			if !seenSigners[addr] {
				requiredPublicKeys = append(requiredPublicKeys, &types.AccountIdentifier{
					Address: addr,
				})

				seenSigners[addr] = true
			}
		}
	}

	return &types.ConstructionPreprocessResponse{
		Options:            options,
		RequiredPublicKeys: requiredPublicKeys,
	}, nil
}

func parseOperationMsgs(ops []*types.Operation) ([]sdk.Msg, *types.Error) {
	if len(ops) != 2 {
		return nil, wrapErr(ErrUnclearIntent, errors.New("invalid number of operations, expected 2"))
	}

	sendMsg := banktypes.MsgSend{}

	for _, op := range ops {
		if op.Type != kava.TransferOpType {
			return nil, wrapErr(ErrUnclearIntent, fmt.Errorf("invalid opeartion type, only '%s' allowed", kava.TransferOpType))
		}

		value, err := types.AmountValue(op.Amount)
		if err != nil {
			return nil, ErrInvalidCurrencyAmount
		}

		if value.Sign() == 0 {
			return nil, ErrInvalidCurrencyAmount
		}

		if value.Sign() > 0 {
			to, err := getAddressFromAccount(op.Account)
			if err != nil {
				return nil, err
			}

			sendMsg.ToAddress = to.String()

			coin, err := amountToCoin(op.Amount)
			if err != nil {
				return nil, ErrInvalidCurrencyAmount
			}
			sendMsg.Amount = sdk.NewCoins(coin)
		}

		if value.Sign() < 0 {
			from, err := getAddressFromAccount(op.Account)
			if err != nil {
				return nil, err
			}

			sendMsg.FromAddress = from.String()
		}
	}

	return []sdk.Msg{&sendMsg}, nil
}

func suggestedMultiplerOrDefault(multiplier *float64) float64 {
	if multiplier == nil {
		return defaultSuggestedFeeMultiplier
	}

	return *multiplier
}

func getMemoFromMetadata(metadata map[string]interface{}) string {
	if rawMemo, exists := metadata["memo"]; exists {
		if memo, ok := rawMemo.(string); ok {
			return memo
		}
	}

	return ""
}

func getGasAdjustmentFromMetadata(metadata map[string]interface{}) float64 {
	if rawAdjustment, exists := metadata["gas_adjustment"]; exists {
		if adjustment, ok := rawAdjustment.(float64); ok {
			return adjustment
		}
	}

	return defaultGasAdjustment
}

func getMaxFeeAndEncodeOption(amounts []*types.Amount) (*string, *types.Error) {
	if len(amounts) == 0 {
		return nil, nil
	}

	var maxFee sdk.Coins
	for _, feeAmount := range amounts {
		coin, err := amountToCoin(feeAmount)
		if err != nil {
			return nil, err
		}
		maxFee = maxFee.Add(coin)
	}

	b, err := json.Marshal(maxFee)
	if err != nil {
		return nil, wrapErr(ErrKava, err)
	}

	encodedMaxFee := string(b)
	return &encodedMaxFee, nil
}

func amountToCoin(amount *types.Amount) (sdk.Coin, *types.Error) {
	value, ok := sdk.NewIntFromString(amount.Value)
	if !ok {
		return sdk.Coin{}, ErrInvalidCurrencyAmount
	}

	denom, ok := kava.Denoms[amount.Currency.Symbol]
	if !ok {
		return sdk.Coin{}, ErrUnsupportedCurrency
	}

	currency, ok := kava.Currencies[denom]
	if !ok {
		return sdk.Coin{}, ErrUnsupportedCurrency
	}

	if currency.Symbol != amount.Currency.Symbol ||
		currency.Decimals != amount.Currency.Decimals {
		return sdk.Coin{}, ErrUnsupportedCurrency
	}

	return sdk.NewCoin(denom, value), nil
}

func getAddressFromAccount(account *types.AccountIdentifier) (sdk.AccAddress, *types.Error) {
	if account == nil || account.Address == "" {
		return nil, ErrInvalidAddress
	}

	addr, err := sdk.AccAddressFromBech32(account.Address)
	if err != nil {
		return nil, ErrInvalidAddress
	}

	return addr, nil
}

func convertMsgsToAnys(msgs []sdk.Msg) ([]*codectypes.Any, error) {
	var anys []*codectypes.Any

	for _, msg := range msgs {
		any, err := codectypes.NewAnyWithValue(msg)
		if err != nil {
			return anys, nil
		}
		anys = append(anys, any)
	}

	return anys, nil
}
