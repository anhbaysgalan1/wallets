package api

import (
	"errors"
	"net/http"

	"github.com/formancehq/wallets/pkg/wallet"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (m *MainHandler) DebitWalletHandler(w http.ResponseWriter, r *http.Request) {
	data := &wallet.DebitRequest{}
	if err := render.Bind(r, data); err != nil {
		badRequest(w, ErrorCodeValidation, err)
		return
	}

	hold, err := m.funding.Debit(r.Context(), wallet.Debit{
		WalletID:     chi.URLParam(r, "walletID"),
		DebitRequest: *data,
	})
	if err != nil {
		switch {
		case errors.Is(err, wallet.ErrInsufficientFundError):
			badRequest(w, ErrorCodeInsufficientFund, wallet.ErrInsufficientFundError)
		default:
			internalError(w, r, err)
		}
		return
	}

	if hold == nil {
		noContent(w)
		return
	}

	created(w, hold)
}
