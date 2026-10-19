package flexops

import "context"

type WorkspacesService struct{ c *Client }

// List returns the workspaces the caller belongs to.
func (s *WorkspacesService) List(ctx context.Context) (ApiResponse[[]Workspace], error) {
	return decode[ApiResponse[[]Workspace]](s.c.http.get(ctx, "/api/workspaces", nil))
}

// Get returns a workspace; an empty workspaceID means the configured one.
func (s *WorkspacesService) Get(ctx context.Context, workspaceID string) (ApiResponse[WorkspaceDetail], error) {
	if workspaceID == "" {
		workspaceID = s.c.WorkspaceID
	}
	return decode[ApiResponse[WorkspaceDetail]](s.c.http.get(ctx, "/api/workspaces/"+workspaceID, nil))
}

// ListMembers returns the members, pending invites and seat usage of the configured workspace.
func (s *WorkspacesService) ListMembers(ctx context.Context) (ApiResponse[MemberList], error) {
	return decode[ApiResponse[MemberList]](s.c.http.get(ctx, s.c.wsPath("members"), nil))
}

// InviteMember invites email to the configured workspace with role (sent as its number).
func (s *WorkspacesService) InviteMember(ctx context.Context, email string, role WorkspaceRole) (ApiResponse[any], error) {
	body := struct {
		Email string        `json:"email"`
		Role  WorkspaceRole `json:"role"`
	}{email, role}
	return decode[ApiResponse[any]](s.c.http.post(ctx, s.c.wsPath("members/invite"), body))
}
func (s *WorkspacesService) RemoveMember(ctx context.Context, userID string) error {
	_, err := s.c.http.del(ctx, s.c.wsPath("members/"+userID))
	return err
}

// UpdateMemberRole changes a member's role (sent as the Gateway's newRole number).
func (s *WorkspacesService) UpdateMemberRole(ctx context.Context, userID string, role WorkspaceRole) (ApiResponse[any], error) {
	body := struct {
		NewRole WorkspaceRole `json:"newRole"`
	}{role}
	return decode[ApiResponse[any]](s.c.http.put(ctx, s.c.wsPath("members/"+userID+"/role"), body))
}
