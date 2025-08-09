package i18n

import "gluttony/x/httpx"

func Middleware(manager *I18n) httpx.Middleware {
	return func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(c *httpx.Context) error {
			// TODO: Add user setting table containing language (default should be en)
			nextCtx := WithI18nBundle(c.Context(), manager.Bundles["en"])
			nextRequest := c.Request.WithContext(nextCtx)
			c.Request = nextRequest

			return next(c)
		}
	}
}
