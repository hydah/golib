package httpsvr

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ClientIP returns more real IP address.
func (c *Context) ClientIP() string {
	clientIP := c.Req.Header.Get("X-Forwarded-For")
	if len(clientIP) == 0 {
		clientIP = c.Req.Header.Get("X-Real-IP")
	}
	if len(clientIP) == 0 {
		ipS := strings.Split(c.Req.RemoteAddr, ":")
		if len(ipS) > 0 {
			if ipS[0] != "[" {
				return ipS[0]
			}
		}
		return "127.0.0.1"

	} else {
		idx := strings.Index(clientIP, ",")
		if idx > 0 {
			clientIP = clientIP[:idx]
		}
	}
	return clientIP
}

// SetCookie sets given cookie value to response header.
// ctx.SetCookie(name, value [, MaxAge, Path, Domain, Secure, HttpOnly])
func (c *Context) SetCookie(name, value string, others ...interface{}) {
	cookie := &http.Cookie{}
	cookie.Name = name
	cookie.Value = value

	if len(others) > 0 {
		switch v := others[0].(type) {
		case int:
			cookie.MaxAge = v
		case int64:
			cookie.MaxAge = int(v)
		case int32:
			cookie.MaxAge = int(v)
		}
	}

	// default "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			cookie.Path = v
		}
	} else {
		cookie.Path = "/"
	}

	// default empty
	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			cookie.Domain = v
		}
	}

	// default empty
	if len(others) > 3 {
		switch v := others[3].(type) {
		case bool:
			cookie.Secure = v
		}
	}

	// default false.
	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			cookie.HttpOnly = true
		}
	}

	http.SetCookie(c.Writer, cookie)
}

// GetCookie returns given cookie value from request header.
func (c *Context) GetCookie(name string) string {
	cookie, err := c.Req.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

var cookieSecret string

// SetCookieSecret sets global default secure cookie secret.
func (m *Engine) SetCookieSecret(secret string) {
	cookieSecret = secret
}

// SetSecureCookie sets given cookie value to response header with default secret string.
func (ctx *Context) SetSecureCookie(name, value string, others ...interface{}) {
	ctx.SetBasicSecureCookie(cookieSecret, name, value, others...)
}

// GetSecureCookie returns given cookie value from request header with default secret string.
func (ctx *Context) GetSecureCookie(name string) (string, bool) {
	return ctx.GetBasicSecureCookie(cookieSecret, name)
}

// SetBasicSecureCookie sets given cookie value to response header with secret string.
func (ctx *Context) SetBasicSecureCookie(secret, name, value string, others ...interface{}) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	hm := hmac.New(sha1.New, []byte(secret))
	fmt.Fprintf(hm, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", hm.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")

	ctx.SetCookie(name, cookie, others...)
}

// GetBasicSecureCookie returns given cookie value from request header with secret string.
func (ctx *Context) GetBasicSecureCookie(secret, name string) (string, bool) {
	val := ctx.GetCookie(name)
	if val == "" {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3)
	if len(parts) != 3 {
		return "", false
	}
	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]

	hm := hmac.New(sha1.New, []byte(secret))
	fmt.Fprintf(hm, "%s%s", vs, timestamp)
	if fmt.Sprintf("%02x", hm.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}

func (ctx *Context) GetParamByName(name string) string {
	return ctx.Params.ByName(name)
}

func (ctx *Context) MustQueryInt(key string, d int) int {
	val := ctx.Req.URL.Query().Get(key)
	if val == "" {
		return d
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(err.Error())
	}
	return i
}

func (ctx *Context) MustQueryFloat64(key string, d float64) float64 {
	val := ctx.Req.URL.Query().Get(key)
	if val == "" {
		return d
	}
	f, err := strconv.ParseFloat(ctx.Req.URL.Query().Get(key), 64)
	if err != nil {
		panic(err)
	}
	return f
}

func (ctx *Context) MustQueryString(key, d string) string {
	val := ctx.Req.URL.Query().Get(key)
	if val == "" {
		return d
	}
	return val
}

func (ctx *Context) MustQueryStrings(key string, d []string) []string {
	val := ctx.Req.URL.Query()[key]
	if len(val) == 0 {
		return d
	}
	return val
}

func (ctx *Context) MustQueryTime(key string, layout string, d time.Time) time.Time {
	val := ctx.Req.URL.Query().Get(key)
	if val == "" {
		return d
	}
	t, err := time.Parse(layout, ctx.Req.URL.Query().Get(key))
	if err != nil {
		panic(err)
	}
	return t
}

func (ctx *Context) MustPostInt(key string, d int) int {
	val := ctx.Req.PostFormValue(key)
	if val == "" {
		return d
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return i
}

func (ctx *Context) MustPostFloat64(key string, d float64) float64 {
	val := ctx.Req.PostFormValue(key)
	if val == "" {
		return d
	}
	f, err := strconv.ParseFloat(ctx.Req.URL.Query().Get(key), 64)
	if err != nil {
		panic(err)
	}
	return f
}

func (ctx *Context) MustPostString(key, d string) string {
	val := ctx.Req.PostFormValue(key)
	if val == "" {
		return d
	}
	return val
}

func (ctx *Context) MustPostStrings(key string, d []string) []string {
	if ctx.Req.PostForm == nil {
		ctx.Req.ParseForm()
	}
	val := ctx.Req.PostForm[key]
	if len(val) == 0 {
		return d
	}
	return val
}

func (ctx *Context) MustPostTime(key string, layout string, d time.Time) time.Time {
	val := ctx.Req.PostFormValue(key)
	if val == "" {
		return d
	}
	t, err := time.Parse(layout, ctx.Req.URL.Query().Get(key))
	if err != nil {
		panic(err)
	}
	return t
}
