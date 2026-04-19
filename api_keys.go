package flexops

import "context"

type ApiKeysService struct{ c *Client }

func (s *ApiKeysService) List(ctx context.Context) (ApiResponse[[]ApiKeyInfo], error) { return decode[ApiResponse[[]ApiKeyInfo]](s.c.http.get(ctx, s.c.wsPath("api-keys"), nil)) }
func (s *ApiKeysService) Create(ctx context.Context, req CreateApiKeyRequest) (ApiResponse[CreateApiKeyResponse], error) { return decode[ApiResponse[CreateApiKeyResponse]](s.c.http.post(ctx, s.c.wsPath("api-keys"), req)) }
func (s *ApiKeysService) Revoke(ctx context.Context, keyID string) error { _, err := s.c.http.del(ctx, s.c.wsPath("api-keys/"+keyID)); return err }
func (s *ApiKeysService) Rotate(ctx context.Context, keyID string) (ApiResponse[CreateApiKeyResponse], error) { return decode[ApiResponse[CreateApiKeyResponse]](s.c.http.post(ctx, s.c.wsPath("api-keys/"+keyID+"/rotate"), nil)) }
