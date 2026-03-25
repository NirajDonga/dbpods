package core

import "context"

type DBProvisioner interface {
	CreateTenantDatabase(ctx context.Context, tenantID string, password string) error
}
