package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/wallets/pkg/core"
	"github.com/formancehq/wallets/pkg/wallet/numscript"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHoldsConfirm(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	hold := core.NewDebitHold(walletID, "bank", "USD")

	req := newRequest(t, http.MethodPost, "/wallets/"+walletID+"/holds/"+hold.ID+"/confirm", nil)
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*sdk.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetHoldAccount(hold.ID), account)
			balances := map[string]int32{
				"USD": 100,
			}
			return &sdk.AccountWithVolumesAndBalances{
				Address:  testEnv.Chart().GetHoldAccount(hold.ID),
				Metadata: hold.LedgerMetadata(testEnv.Chart()),
				Balances: &balances,
			}, nil
		}),
		WithRunScript(func(ctx context.Context, name string, script sdk.Script) (*sdk.ScriptResult, error) {
			require.EqualValues(t, sdk.Script{
				Plain: numscript.BuildConfirmHoldScript(false, "USD"),
				Vars: map[string]interface{}{
					"hold": testEnv.Chart().GetHoldAccount(hold.ID),
					"amount": map[string]any{
						"amount": uint64(100),
						"asset":  "USD",
					},
				},
				Metadata: core.WalletTransactionBaseMetadata(),
			}, script)
			return &sdk.ScriptResult{}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
}

func TestHoldsPartialConfirm(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	hold := core.NewDebitHold(walletID, "bank", "USD")

	req := newRequest(t, http.MethodPost, "/wallets/"+walletID+"/holds/"+hold.ID+"/confirm", ConfirmHoldRequest{
		Amount: 50,
	})
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*sdk.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetHoldAccount(hold.ID), account)
			balances := map[string]int32{
				"USD": 100,
			}
			volumes := map[string]map[string]int32{
				"USD": {
					"input": 100,
				},
			}
			return &sdk.AccountWithVolumesAndBalances{
				Address:  testEnv.Chart().GetHoldAccount(hold.ID),
				Metadata: hold.LedgerMetadata(testEnv.Chart()),
				Balances: &balances,
				Volumes:  &volumes,
			}, nil
		}),
		WithRunScript(func(ctx context.Context, name string, script sdk.Script) (*sdk.ScriptResult, error) {
			require.EqualValues(t, sdk.Script{
				Plain: numscript.BuildConfirmHoldScript(false, "USD"),
				Vars: map[string]interface{}{
					"hold": testEnv.Chart().GetHoldAccount(hold.ID),
					"amount": map[string]any{
						"amount": uint64(50),
						"asset":  "USD",
					},
				},
				Metadata: core.WalletTransactionBaseMetadata(),
			}, script)
			return &sdk.ScriptResult{}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
}

func TestHoldsPartialConfirmWithFinal(t *testing.T) {
	t.Parallel()

	walletID := uuid.NewString()
	hold := core.NewDebitHold(walletID, "bank", "USD")

	req := newRequest(t, http.MethodPost, "/wallets/"+walletID+"/holds/"+hold.ID+"/confirm", ConfirmHoldRequest{
		Amount: 50,
		Final:  true,
	})
	rec := httptest.NewRecorder()

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*sdk.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetHoldAccount(hold.ID), account)
			balances := map[string]int32{
				"USD": 100,
			}
			volumes := map[string]map[string]int32{
				"USD": {
					"input": 100,
				},
			}
			return &sdk.AccountWithVolumesAndBalances{
				Address:  testEnv.Chart().GetHoldAccount(hold.ID),
				Metadata: hold.LedgerMetadata(testEnv.Chart()),
				Balances: &balances,
				Volumes:  &volumes,
			}, nil
		}),
		WithRunScript(func(ctx context.Context, name string, script sdk.Script) (*sdk.ScriptResult, error) {
			require.EqualValues(t, sdk.Script{
				Plain: numscript.BuildConfirmHoldScript(true, "USD"),
				Vars: map[string]interface{}{
					"hold": testEnv.Chart().GetHoldAccount(hold.ID),
					"amount": map[string]any{
						"amount": uint64(50),
						"asset":  "USD",
					},
				},
				Metadata: core.WalletTransactionBaseMetadata(),
			}, script)
			return &sdk.ScriptResult{}, nil
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Result().StatusCode)
}
