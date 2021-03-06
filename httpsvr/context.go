package httpsvr

import (
	"errors"
	"log"
	"math"
	"net/http"
	"reflect"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/hydah/golib/httpsvr/render"
	"github.com/hydah/golib/logger"
)

type Context struct {
	Writer      ResponseWriter
	Req         *http.Request
	Session     Session
	Keys        map[string]interface{}
	Params      httprouter.Params
	Engine      *Engine
	writer      writer
	handlers    []HandlerFunc
	controllers []reflect.Type
	index       int8
	HtmlEngine
}

const (
	abortIndex = math.MaxInt8 / 2
)

// Next should be used only in the middlewares.
// It executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		if len(c.controllers) > 0 && c.controllers[c.index] != nil {
			ctrl, ok := reflect.New(c.controllers[c.index]).Interface().(IController)
			if !ok {
				logger.Error("controller is not IController, %v", c.controllers[c.index].Name())
				return
			}
			ctrl.InitCtx(c)
			ctrl.InitBase(c)
			ctrl.InitApp(c)

			if ctrl.Prepare(c) {
				switch c.Req.Method {
				case "GET":
					ctrl.Get(c)
				case "POST":
					ctrl.Post(c)
				case "PATCH":
					ctrl.Patch(c)
				case "PUT":
					ctrl.Put(c)
				case "OPTIONS":
					ctrl.Options(c)
				case "HEAD":
					ctrl.Head(c)
				default:
					logger.Error("method: %s, controller handler Not Implemented, %v", c.Req.Method, c.controllers[c.index].Name())
				}
			}
			ctrl.Finish(c)

		} else {
			c.handlers[c.index](c)
		}
	}
}

// Sets a new pair key/value just for the specified context.
func (c *Context) Set(key string, item interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = item
}

// Get returns the value for the given key or an error if the key does not exist.
func (c *Context) Get(key string) (interface{}, error) {
	if c.Keys != nil {
		value, ok := c.Keys[key]
		if ok {
			return value, nil
		}
	}
	return nil, errors.New("Key does not exist.")
}

// MustGet returns the value for the given key or panics if the value doesn't exist.
func (c *Context) MustGet(key string) interface{} {
	value, err := c.Get(key)
	if err != nil || value == nil {
		log.Panicf("Key %s doesn't exist", value)
	}
	return value
}

// SetHeader sets a response header.
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// Forces the system to do not continue calling the pending handlers in the chain.
func (c *Context) Abort() {
	c.index = abortIndex
}

// Redirect returns a HTTP redirect to the specific location. default for 302
func (c *Context) Redirect(location string, status ...int) {
	c.SetHeader("Location", location)
	if status != nil {
		http.Redirect(c.Writer, c.Req, location, status[0])
	} else {
		http.Redirect(c.Writer, c.Req, location, 302)
	}
}

// Serializes the given struct as JSON into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/json".
func (c *Context) Json(data interface{}, status ...int) {
	c.executeRender(data, c.Writer, render.JSON{}, status...)
}

// Serializes the given struct as JSONP into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/javascript".
func (c *Context) Jsonp(callback string, data interface{}, status ...int) {
	c.executeRender(data, c.Writer, render.JSONP{Callback: callback}, status...)
}

// Serializes the given struct as XML into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/xml".
func (c *Context) Xml(data interface{}, status ...int) {
	c.executeRender(data, c.Writer, render.XML{}, status...)
}

// Writes the given string into the response body and sets the Content-Type to "text/plain".
func (c *Context) Text(data string, status ...int) {
	c.executeRender(data, c.Writer, render.TEXT{}, status...)
}

// Writes the given bytes into the response body and sets the Content-Type to "text/plain".
func (c *Context) Bytes(data []byte, status ...int) {
	c.executeRender(data, c.Writer, render.BYTES{}, status...)
}

// Writes the given html into the response body and sets the Content-Type to "text/html".
func (c *Context) Html(data string, status ...int) {
	c.executeRender(data, c.Writer, render.HTML{}, status...)
}

func (c *Context) Video(data []byte, status ...int) {
	c.executeRender(data, c.Writer, render.VIDEO{}, status...)
}

func (c *Context) Image(data []byte, status ...int) {
	c.executeRender(data, c.Writer, render.IMAGE{}, status...)
}

func (c *Context) Raw(data []byte, status ...int) {
	c.executeRender(data, c.Writer, render.RAW{}, status...)
}

// IP _
func (c *Context) IP() (ip string) {
	ips := c.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		return ips[0]
	}
	ipS := strings.Split(c.Req.RemoteAddr, ":")
	if len(ipS) > 0 {
		if ipS[0] != "[" {
			return ipS[0]
		}
	}
	return "127.0.0.1"
}

// Proxy _
func (c *Context) Proxy() (proxy []string) {
	if ips := c.Req.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

// UserAgent _
func (c *Context) UserAgent() string {
	return c.Req.Header.Get("User-Agent")
}

func (c *Context) executeRender(data interface{}, w http.ResponseWriter, render render.Render, status ...int) {
	if status != nil {
		c.Writer.WriteHeader(status[0])
	}
	if err := render.Render(data, w); err != nil {
		c.Writer.WriteHeader(500)
		c.Abort()
	}
}

func (c *Engine) createContext(w http.ResponseWriter, req *http.Request, params httprouter.Params, handlers []HandlerFunc, controllers []reflect.Type) *Context {
	ctx := c.pool.Get().(*Context)
	ctx.Writer = &ctx.writer
	ctx.Req = req
	ctx.Keys = nil
	ctx.Params = params
	ctx.handlers = handlers
	ctx.controllers = controllers
	ctx.writer.reset(w)
	ctx.index = -1
	return ctx
}

func (c *Engine) reuseContext(ctx *Context) {
	c.pool.Put(ctx)
}
