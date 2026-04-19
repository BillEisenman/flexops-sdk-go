package flexops

import "context"

type WorkspacesService struct{ c *Client }

func (s *WorkspacesService) List(ctx context.Context) (ApiResponse[[]Workspace], error) {
	return decode[ApiResponse[[]Workspace]](s.c.http.get(ctx, "/api/workspaces", nil))
}
func (s *WorkspacesService) Get(ctx context.Context, workspaceID string) (ApiResponse[Workspace], error) {
	if workspaceID == "" { workspaceID = s.c.WorkspaceID }
	return decode[ApiResponse[Workspace]](s.c.http.get(ctx, "/api/workspaces/"+workspaceID, nil))
}
func (s *WorkspacesService) Create(ctx context.Context, req CreateWorkspaceRequest) (ApiResponse[Workspace], error) {
	return decode[ApiResponse[Workspace]](s.c.http.post(ctx, "/api/workspaces", req))
}
func (s *WorkspacesService) Update(ctx context.Context, data any) (ApiResponse[Workspace], error) {
	return decode[ApiResponse[Workspace]](s.c.http.put(ctx, "/api/workspaces/"+s.c.WorkspaceID, data))
}
func (s *WorkspacesService) ListMembers(ctx context.Context) (ApiResponse[[]WorkspaceMember], error) {
	return decode[ApiResponse[[]WorkspaceMember]](s.c.http.get(ctx, s.c.wsPath("members"), nil))
}
func (s *WorkspacesService) InviteMember(ctx context.Context, email, role string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("members/invite"), map[string]string{"email": email, "role": role}))
}
func (s *WorkspacesService) RemoveMember(ctx context.Context, userID string) error {
	_, err := s.c.http.del(ctx, s.c.wsPath("members/"+userID))
	return err
}
func (s *WorkspacesService) UpdateMemberRole(ctx context.Context, userID, role string) (ApiResponse[any], error) {
	return decode[ApiResponse[any]](s.c.http.put(ctx, s.c.wsPath("members/"+userID+"/role"), map[string]string{"role": role}))
}
