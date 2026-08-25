package contract

import "context"

type IPermissionAccess interface {
	HasPermission(ctx context.Context, userID int, resources []string) (bool, error)
	GetUserPermissions(ctx context.Context, userID int) ([]string, error)
	InvalidateUserPermissions(ctx context.Context, userID int) error
}
