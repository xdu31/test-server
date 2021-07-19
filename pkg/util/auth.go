package util

import (
	"context"
)

const (
	AccountIDKey = "account-id"
)

// GetAccountID extracts and returns account ID from the context ctx.
func GetAccountID(ctx context.Context) (uint32, error) {

	aid, ok := ctx.Value("account-id").(int)
	if !ok || aid == 0 {
		return 0, UnauthenticatedErr("Unauthorized")
	}

	return uint32(aid), nil
}
