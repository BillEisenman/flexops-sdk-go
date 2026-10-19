package flexops

import (
	"context"
	"net/http"
)

type WalletService struct{ c *Client }

// GetBalance returns the configured workspace's wallet balance. The Gateway resolves the
// workspace from the X-Workspace-Id header and returns the balance without an envelope.
func (s *WalletService) GetBalance(ctx context.Context) (WalletBalance, error) {
	headers := http.Header{}
	headers.Set("X-Workspace-Id", s.c.WorkspaceID)
	return decode[WalletBalance](s.c.http.do(ctx, http.MethodGet, "/api/wallet/my-balance", nil, nil, headers))
}
