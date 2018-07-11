package httpsvr

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type IController interface {
	InitCtx(ctx *Context)
	InitBase(ctx *Context)
	InitApp(ctx *Context)
	Prepare(ctx *Context) bool

	Get(ctx *Context)
	Post(ctx *Context)
	Delete(ctx *Context)
	Patch(ctx *Context)
	Put(ctx *Context)
	Options(ctx *Context)
	Head(ctx *Context)
	Finish(ctx *Context)
}

type Controller struct {
	Ctx *Context
}

func (c *Controller) InitCtx(ctx *Context) {
	c.Ctx = ctx
}

func (c *Controller) InitBase(ctx *Context) {
}

func (c *Controller) InitApp(ctx *Context) {
}

func (c *Controller) Prepare(ctx *Context) bool {
	return true
}

func (c *Controller) Get(ctx *Context) {
	ctx.Writer.WriteHeader(http.StatusMethodNotAllowed)
	ctx.Abort()
}

func (c *Controller) Post(ctx *Context) {
	ctx.Writer.WriteHeader(http.StatusMethodNotAllowed)
	ctx.Abort()
}

func (c *Controller) Delete(ctx *Context) {
	ctx.Writer.WriteHeader(http.StatusMethodNotAllowed)
	ctx.Abort()
}

func (c *Controller) Patch(ctx *Context) {
	ctx.Writer.WriteHeader(http.StatusMethodNotAllowed)
	ctx.Abort()
}

func (c *Controller) Put(ctx *Context) {
	ctx.Writer.WriteHeader(http.StatusMethodNotAllowed)
	ctx.Abort()
}

func (c *Controller) Options(ctx *Context) {
	ctx.Writer.WriteHeader(http.StatusMethodNotAllowed)
	ctx.Abort()
}

func (c *Controller) Head(ctx *Context) {
	ctx.Writer.WriteHeader(http.StatusMethodNotAllowed)
	ctx.Abort()
}

func (c *Controller) Finish(ctx *Context) {
	c.Ctx = nil
}

var maxMemory = 1 << 26

// Input 从 request 中获取参数.
func (c *Controller) Input() (input url.Values) {
	ct := c.Ctx.Req.Header.Get("Content-Type")
	if strings.Contains(ct, "multipart/form-data") {
		c.Ctx.Req.ParseMultipartForm(int64(maxMemory))
	} else {
		c.Ctx.Req.ParseForm()
	}
	return c.Ctx.Req.Form
}

// GetString 获取 query 参数.
func (c *Controller) GetString(key string) (value string) {
	return c.Input().Get(key)
}

// GetStrings 从 request 中获取输入的参数数组,
// 应用于如 checkbox(input[type=chackbox]), 多选框等情况的表单提交.
// Parameters:
// - key:    要获取的输入值的 key 值.
// Return:
// - values: 要获取的输入值的 key 值对应的值数组.
func (c *Controller) GetStrings(key string) (values []string) {
	r := c.Ctx.Req
	if r.Form == nil {
		return
	}
	vs := r.Form[key]
	if len(vs) > 0 {
		return vs
	}
	return
}

func (c *Controller) GetStringms(key string) (values map[string][]string) {
	r := c.Ctx.Req
	if r.Form == nil {
		return
	}
	values = make(map[string][]string, 0)
	for k, v := range r.Form {
		ks := strings.Split(k, "]")
		if len(ks) == 2 && ks[1] == key {
			kss := strings.Split(ks[0], "[")
			if len(kss) == 2 {
				values[kss[1]] = v
			}
		}
	}
	return
}

func (c *Controller) GetStringm(key string) (values map[string]string) {
	r := c.Ctx.Req
	if r.Form == nil {
		return
	}
	values = make(map[string]string, 0)
	for k, v := range r.Form {
		ks := strings.Split(k, "]")
		if len(ks) == 2 && ks[1] == key {
			kss := strings.Split(ks[0], "[")
			if len(kss) == 2 {
				if len(v) > 0 {
					values[kss[1]] = v[0]
				}
			}
		}
	}
	return
}

// GetInt 从 url 参数中获得 int 变量值.
// Parameters:
// - key:   要获取的输入值的 key 值.
// Return:
// - value: 要获取的输入值的 key 值对应的值.
// - err:
func (c *Controller) GetInt(key string) (value int64, err error) {
	return strconv.ParseInt(c.Input().Get(key), 10, 64)
}

// GetBool 获得表单提交的 bool 类型变量值.
// Parameters:
// - key:   要获取的输入值的 key 值.
// Return:
// - value: 要获取的输入值的 key 值对应的值.
// - err:
func (c *Controller) GetBool(key string) (value bool, err error) {
	return strconv.ParseBool(c.Input().Get(key))
}

// GetFloat 获得 表单提交的 float64 类型变量值.
// Parameters:
// - key:   要获取的输入值的 key 值.
// Return:
// - value: 要获取的输入值的 key 值对应的值.
// - err:
func (c *Controller) GetFloat(key string) (value float64, err error) {
	return strconv.ParseFloat(c.Input().Get(key), 64)
}

// GetFile 获得上传的文件.
// Parameters:
// - key:    要获取的输入值的 key 值.
// Return:
// - file:   上传的文件.
// - header: 上传文件的头部信息.
// - err:
func (c *Controller) GetFile(key string) (file multipart.File, header *multipart.FileHeader, err error) {
	return c.Ctx.Req.FormFile(key)
}

// SaveToFile 将上传的文件转存成文件.
// Parameters:
// - fromfile: 上传的文件.
// - tofile:   转存的文件.
// Return:
// - err:
func (c *Controller) SaveToFile(fromfile, tofile string) (err error) {
	file, _, err := c.Ctx.Req.FormFile(fromfile)
	if err != nil {
		return
	}
	defer file.Close()
	f, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	io.Copy(f, file)
	return
}
