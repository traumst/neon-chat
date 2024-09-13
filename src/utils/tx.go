package utils

import "net/http"

func FlagTxChages(r *http.Request, hasChanges bool) {
	*r.Context().Value(TxChangesKey).(*bool) = hasChanges
}
