package flexops

import "context"

type AuthService struct{ c *Client }

func (s *AuthService) Login(ctx context.Context, email, password string) (ApiResponse[LoginResponse], error) {
	return decode[ApiResponse[LoginResponse]](s.c.http.post(ctx, "/api/Account/login", LoginRequest{Email: email, Password: password}))
}
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, "/api/Account/register", req))
}
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (ApiResponse[LoginResponse], error) {
	return decode[ApiResponse[LoginResponse]](s.c.http.post(ctx, "/api/Account/refresh-token", map[string]string{"refreshToken": refreshToken}))
}
func (s *AuthService) Logout(ctx context.Context) error {
	_, err := s.c.http.post(ctx, "/api/Account/logout", nil)
	return err
}
func (s *AuthService) GetProfile(ctx context.Context) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.get(ctx, "/api/Account/profile", nil))
}
func (s *AuthService) UpdateProfile(ctx context.Context, data any) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.put(ctx, "/api/Account/profile", data))
}
func (s *AuthService) ChangePassword(ctx context.Context, currentPassword, newPassword string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, "/api/Account/change-password", map[string]string{"currentPassword": currentPassword, "newPassword": newPassword}))
}
func (s *AuthService) ForgotPassword(ctx context.Context, email string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, "/api/Account/forgot-password", map[string]string{"email": email}))
}
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, "/api/Account/reset-password", map[string]string{"token": token, "newPassword": newPassword}))
}
func (s *AuthService) VerifyEmail(ctx context.Context, token string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, "/api/Account/verify-email", map[string]string{"token": token}))
}
