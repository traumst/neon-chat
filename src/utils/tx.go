package utils

import (
	"context"
	"neon-chat/src/consts"
	"net/http"
)

func HasTxChanges(r *http.Request) bool {
	ctx := r.Context()
	changesMade, ok := ctx.Value(consts.TxChangesKey).(*bool)
	if !ok || changesMade == nil {
		return false
	}
	return *changesMade
}

func FlagTxChages(r *http.Request, hasChanges bool) {
	ctx := r.Context()
	changesMade, ok := ctx.Value(consts.TxChangesKey).(*bool)
	if !ok || changesMade == nil {
		changesMade = new(bool)
		*changesMade = hasChanges
		ctx = context.WithValue(ctx, consts.TxChangesKey, changesMade)
		*r = *r.WithContext(ctx)
	} else {
		*changesMade = hasChanges
	}
}
