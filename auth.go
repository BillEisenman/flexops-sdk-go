package flexops

import (
	"context"
	"fmt"
)

// AuthService acquires and revokes user tokens. API-key callers do not need it.
type AuthService struct{ c *Client }

// Login signs in and, on success, makes the client send data.token as its bearer token
// (replacing any API key). A response with RequiresTwoFactor set carries no token.
func (s *AuthService) Login(ctx context.Context, email, password string) (ApiResponse[LoginResponse], error) {
	resp, err := decode[ApiResponse[LoginResponse]](s.c.http.post(ctx, "/api/account/login", LoginRequest{Email: email, Password: password}))
	if err == nil && resp.Data.Token != "" {
		s.c.SetAccessToken(resp.Data.Token)
	}
	return resp, err
}

// RefreshToken exchanges the current access token and its refresh token for new ones;
// the Gateway requires both. accessToken defaults to the client's current bearer token;
// with neither, it returns an error before any request. On success the client sends the
// new data.token.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string, accessToken ...string) (ApiResponse[LoginResponse], error) {
	if len(accessToken) > 1 {
		return ApiResponse[LoginResponse]{}, fmt.Errorf("flexops: RefreshToken takes at most one access token")
	}
	token := s.c.http.accessToken
	if len(accessToken) == 1 && accessToken[0] != "" {
		token = accessToken[0]
	}
	if token == "" {
		return ApiResponse[LoginResponse]{}, fmt.Errorf("flexops: RefreshToken needs the current access token (pass it or log in first)")
	}
	body := map[string]string{"token": token, "refreshToken": refreshToken}
	resp, err := decode[ApiResponse[LoginResponse]](s.c.http.post(ctx, "/api/account/refresh-token", body))
	if err == nil && resp.Data.Token != "" {
		s.c.SetAccessToken(resp.Data.Token)
	}
	return resp, err
}
func (s *AuthService) Logout(ctx context.Context) error {
	_, err := s.c.http.post(ctx, "/api/account/logout", nil)
	return err
}
