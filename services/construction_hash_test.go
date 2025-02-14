// Copyright 2021 Kava Labs, Inc.
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

// TODO: setup test data for hash
//func TestConstructionHash(t *testing.T) {
//	servicer, _ := setupConstructionAPIServicer()
//
//	type errArgs struct {
//		expectErr     bool
//		expectErrCode int32
//	}
//
//	testCases := []struct {
//		testFixtureFile string
//		expectedTxHash  string
//		errs            errArgs
//	}{
//		{
//			testFixtureFile: "msg-send.json",
//			expectedTxHash:  "4E218DC828F45B7112F7CF6B328563045B5307B07D8602549389553F3B27D997",
//			errs: errArgs{
//				expectErr: false,
//			},
//		},
//		{
//			testFixtureFile: "msg-create-cdp.json",
//			expectedTxHash:  "02C44611CD6898E89839F34830A089AD67A1FDA59D809EABA24B5A4B236849BB",
//			errs: errArgs{
//				expectErr: false,
//			},
//		},
//		{
//			testFixtureFile: "msg-hard-deposit.json",
//			expectedTxHash:  "E47E8BB9FA3C90B925D46C75DA03BB316ABB9D04CB647854AC215CB7C743368C",
//			errs: errArgs{
//				expectErr: false,
//			},
//		},
//		{
//			testFixtureFile: "multiple-msgs.json",
//			expectedTxHash:  "4F5EB96A9F29554F2BF0E01059268B1919D5702C29440B017E5C656547725F4C",
//			errs: errArgs{
//				expectErr: false,
//			},
//		},
//		{
//			testFixtureFile: "long-memo.json",
//			expectedTxHash:  "C25EBDC1FB86BEE1F21FB1F0A97925A64ECF838B424D4E57758751806A100FBF",
//			erjs: errArgs{
//				expectErr: false,
//			},
//		},
//		{
//			testFixtureFile: "unsigned-msg-send.json",
//			expectedTxHash:  "",
//			errs: errArgs{
//				expectErr:     true,
//				expectErrCode: ErrInvalidTx.Code,
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		// Load signed transaction from file
//		relPath, err := filepath.Rel(
//			"services",
//			fmt.Sprintf("kava/test-fixtures/txs/%s", tc.testFixtureFile),
//		)
//		require.NoError(t, err)
//		bz, err := ioutil.ReadFile(relPath)
//		require.NoError(t, err)
//
//		cdc := app.MakeEncodingConfig().Amino
//		var stdtx legacytx.StdTx
//		err = cdc.UnmarshalJSON(bz, &stdtx)
//		require.NoError(t, err)
//
//		payload, err := cdc.MarshalBinaryLengthPrefixed(stdtx)
//		require.NoError(t, err)
//		signedTx := hex.EncodeToString(payload)
//
//		networkIdentifier := &types.NetworkIdentifier{
//			Blockchain: "Kava",
//			Network:    "testing",
//		}
//
//		request := &types.ConstructionHashRequest{
//			NetworkIdentifier: networkIdentifier,
//			SignedTransaction: signedTx,
//		}
//
//		// Check that response contains expected tx hash
//		ctx := context.Background()
//		response, rosettaErr := servicer.ConstructionHash(ctx, request)
//		if tc.errs.expectErr {
//			require.NotNil(t, rosettaErr)
//			require.Equal(t, tc.errs.expectErrCode, rosettaErr.Code)
//		} else {
//			require.Nil(t, rosettaErr)
//			require.Equal(t, tc.expectedTxHash, response.TransactionIdentifier.Hash)
//		}
//	}
//}
