package server

import (
	"context"
	"net/http"
)

func useAppContext(next handlerWithAppContext, appCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, appCtx)
	}
}
