package flexops

import (
	"context"
	"net/url"
)

type WalletService struct{ c *Client }

func (s *WalletService) GetBalance(ctx context.Context) (ApiResponse[WalletBalance], error) { return decode[ApiResponse[WalletBalance]](s.c.http.get(ctx, s.c.wsPath("wallet/balance"), nil)) }
func (s *WalletService) AddFunds(ctx context.Context, amount float64) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("wallet/add-funds"), map[string]float64{"amount": amount})) }
func (s *WalletService) ListTransactions(ctx context.Context, params url.Values) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, s.c.wsPath("wallet/transactions"), params)) }
func (s *WalletService) ConfigureAutoReload(ctx context.Context, config any) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.put(ctx, s.c.wsPath("wallet/auto-reload"), config)) }
