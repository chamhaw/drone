package httputils

import (
	"context"
	"github.com/drone/drone/core"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/transport"
	"net/http"
	"time"
)

func PrepareHttpClient(ctx context.Context, user *core.User, renew core.Renewer) (context.Context, *http.Client, error) {
	if user.PrivateToken != "" {
		return ctx, &http.Client{
			Transport: &transport.PrivateToken{
				Token: user.PrivateToken,
			}}, nil
	}
	err := renew.Renew(ctx, user, false)
	if err != nil {
		return ctx, nil, err
	}
	token := &scm.Token{
		Token:   user.Token,
		Refresh: user.Refresh,
	}
	if user.Expiry != 0 {
		token.Expires = time.Unix(user.Expiry, 0)
	}
	ctx = context.WithValue(ctx, scm.TokenKey{}, token)
	return ctx, nil, nil
}
