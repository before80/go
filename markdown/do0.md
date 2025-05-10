## README {#section-readme}

### Gin Web Framework

![img](https://raw.githubusercontent.com/gin-gonic/logo/master/color.png)



Gin is a web framework written in [Go](https://go.dev/). It features a martini-like API with performance that is up to 40 times faster thanks to [httprouter](https://github.com/julienschmidt/httprouter). If you need performance and good productivity, you will love Gin.

**The key features of Gin are:**

- Zero allocation router
- Fast
- Middleware support
- Crash-free
- JSON validation
- Routes grouping
- Error management
- Rendering built-in
- Extendable

#### Getting started

##### Prerequisites

- **[Go](https://go.dev/)**: any one of the **three latest major** [releases](https://go.dev/doc/devel/release) (we test it with these).

##### Getting Gin

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/gin-gonic/gin"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `gin` package:

```
$ go get -u github.com/gin-gonic/gin
```

##### Running Gin

First you need to import Gin package for using Gin, one simplest example likes the follow `example.go`:

```
package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
```

And use the Go command to run the demo:

```
# run example.go and visit 0.0.0.0:8080/ping on browser
$ go run example.go
```

##### Learn more examples

###### Quick Start

Learn and practice more examples, please read the [Gin Quick Start](https://github.com/gin-gonic/gin/blob/v1.10.0/docs/doc.md) which includes API examples and builds tag.

###### Examples

A number of ready-to-run examples demonstrating various use cases of Gin on the [Gin examples](https://github.com/gin-gonic/examples) repository.

#### Documentation

See [API documentation and descriptions](https://godoc.org/github.com/gin-gonic/gin) for package.

All documentation is available on the Gin website.

- [English](https://gin-gonic.com/docs/)
- [简体中文](https://gin-gonic.com/zh-cn/docs/)
- [繁體中文](https://gin-gonic.com/zh-tw/docs/)
- [日本語](https://gin-gonic.com/ja/docs/)
- [Español](https://gin-gonic.com/es/docs/)
- [한국어](https://gin-gonic.com/ko-kr/docs/)
- [Turkish](https://gin-gonic.com/tr/docs/)
- [Persian](https://gin-gonic.com/fa/docs/)

##### Articles about Gin

A curated list of awesome Gin framework.

- [Tutorial: Developing a RESTful API with Go and Gin](https://go.dev/doc/tutorial/web-service-gin)

#### Benchmarks

Gin uses a custom version of [HttpRouter](https://github.com/julienschmidt/httprouter), [see all benchmarks details](https://github.com/gin-gonic/gin/blob/v1.10.0/BENCHMARKS.md).

| Benchmark name                 | (1)       | (2)             | (3)          | (4)             |
| ------------------------------ | --------- | --------------- | ------------ | --------------- |
| BenchmarkGin_GithubAll         | **43550** | **27364 ns/op** | **0 B/op**   | **0 allocs/op** |
| BenchmarkAce_GithubAll         | 40543     | 29670 ns/op     | 0 B/op       | 0 allocs/op     |
| BenchmarkAero_GithubAll        | 57632     | 20648 ns/op     | 0 B/op       | 0 allocs/op     |
| BenchmarkBear_GithubAll        | 9234      | 216179 ns/op    | 86448 B/op   | 943 allocs/op   |
| BenchmarkBeego_GithubAll       | 7407      | 243496 ns/op    | 71456 B/op   | 609 allocs/op   |
| BenchmarkBone_GithubAll        | 420       | 2922835 ns/op   | 720160 B/op  | 8620 allocs/op  |
| BenchmarkChi_GithubAll         | 7620      | 238331 ns/op    | 87696 B/op   | 609 allocs/op   |
| BenchmarkDenco_GithubAll       | 18355     | 64494 ns/op     | 20224 B/op   | 167 allocs/op   |
| BenchmarkEcho_GithubAll        | 31251     | 38479 ns/op     | 0 B/op       | 0 allocs/op     |
| BenchmarkGocraftWeb_GithubAll  | 4117      | 300062 ns/op    | 131656 B/op  | 1686 allocs/op  |
| BenchmarkGoji_GithubAll        | 3274      | 416158 ns/op    | 56112 B/op   | 334 allocs/op   |
| BenchmarkGojiv2_GithubAll      | 1402      | 870518 ns/op    | 352720 B/op  | 4321 allocs/op  |
| BenchmarkGoJsonRest_GithubAll  | 2976      | 401507 ns/op    | 134371 B/op  | 2737 allocs/op  |
| BenchmarkGoRestful_GithubAll   | 410       | 2913158 ns/op   | 910144 B/op  | 2938 allocs/op  |
| BenchmarkGorillaMux_GithubAll  | 346       | 3384987 ns/op   | 251650 B/op  | 1994 allocs/op  |
| BenchmarkGowwwRouter_GithubAll | 10000     | 143025 ns/op    | 72144 B/op   | 501 allocs/op   |
| BenchmarkHttpRouter_GithubAll  | 55938     | 21360 ns/op     | 0 B/op       | 0 allocs/op     |
| BenchmarkHttpTreeMux_GithubAll | 10000     | 153944 ns/op    | 65856 B/op   | 671 allocs/op   |
| BenchmarkKocha_GithubAll       | 10000     | 106315 ns/op    | 23304 B/op   | 843 allocs/op   |
| BenchmarkLARS_GithubAll        | 47779     | 25084 ns/op     | 0 B/op       | 0 allocs/op     |
| BenchmarkMacaron_GithubAll     | 3266      | 371907 ns/op    | 149409 B/op  | 1624 allocs/op  |
| BenchmarkMartini_GithubAll     | 331       | 3444706 ns/op   | 226551 B/op  | 2325 allocs/op  |
| BenchmarkPat_GithubAll         | 273       | 4381818 ns/op   | 1483152 B/op | 26963 allocs/op |
| BenchmarkPossum_GithubAll      | 10000     | 164367 ns/op    | 84448 B/op   | 609 allocs/op   |
| BenchmarkR2router_GithubAll    | 10000     | 160220 ns/op    | 77328 B/op   | 979 allocs/op   |
| BenchmarkRivet_GithubAll       | 14625     | 82453 ns/op     | 16272 B/op   | 167 allocs/op   |
| BenchmarkTango_GithubAll       | 6255      | 279611 ns/op    | 63826 B/op   | 1618 allocs/op  |
| BenchmarkTigerTonic_GithubAll  | 2008      | 687874 ns/op    | 193856 B/op  | 4474 allocs/op  |
| BenchmarkTraffic_GithubAll     | 355       | 3478508 ns/op   | 820744 B/op  | 14114 allocs/op |
| BenchmarkVulcan_GithubAll      | 6885      | 193333 ns/op    | 19894 B/op   | 609 allocs/op   |

- (1): Total Repetitions achieved in constant time, higher means more confident result
- (2): Single Repetition Duration (ns/op), lower is better
- (3): Heap Memory (B/op), lower is better
- (4): Average Allocations per Repetition (allocs/op), lower is better

#### Middlewares

You can find many useful Gin middlewares at [gin-contrib](https://github.com/gin-contrib).

#### Users

Awesome project lists using [Gin](https://github.com/gin-gonic/gin) web framework.

- [gorush](https://github.com/appleboy/gorush): A push notification server written in Go.
- [fnproject](https://github.com/fnproject/fn): The container native, cloud agnostic serverless platform.
- [photoprism](https://github.com/photoprism/photoprism): Personal photo management powered by Go and Google TensorFlow.
- [lura](https://github.com/luraproject/lura): Ultra performant API Gateway with middlewares.
- [picfit](https://github.com/thoas/picfit): An image resizing server written in Go.
- [dkron](https://github.com/distribworks/dkron): Distributed, fault tolerant job scheduling system.

#### Contributing

Gin is the work of hundreds of contributors. We appreciate your help!

Please see [CONTRIBUTING](https://github.com/gin-gonic/gin/blob/v1.10.0/CONTRIBUTING.md) for details on submitting patches and the contribution workflow.

Collapse ▴

## Documentation {#section-documentation}

### Overview {#pkg-overview}

Package gin implements a HTTP web framework called gin.

See https://gin-gonic.com/ for more information about gin.

### Constants {#pkg-constants}

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L28)

```go
const (
    MIMEJSON              = binding.MIMEJSON
    MIMEHTML              = binding.MIMEHTML
    MIMEXML               = binding.MIMEXML
    MIMEXML2              = binding.MIMEXML2
    MIMEPlain             = binding.MIMEPlain
    MIMEPOSTForm          = binding.MIMEPOSTForm
    MIMEMultipartPOSTForm = binding.MIMEMultipartPOSTForm
    MIMEYAML              = binding.MIMEYAML
    MIMETOML              = binding.MIMETOML
)
```

Content-Type MIME of the most common data formats.

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L76)

```go
const (
    // PlatformGoogleAppEngine when running on Google App Engine. Trust X-Appengine-Remote-Addr
    // for determining the client's IP
    PlatformGoogleAppEngine = "X-Appengine-Remote-Addr"
    // PlatformCloudflare when using Cloudflare's CDN. Trust CF-Connecting-IP for determining
    // the client's IP
    PlatformCloudflare = "CF-Connecting-IP"
    // PlatformFlyIO when running on Fly.io. Trust Fly-Client-IP for determining the client's IP
    PlatformFlyIO = "Fly-Client-IP"
)
```

Trusted platforms

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/mode.go#L18)

```go
const (
    // DebugMode indicates gin mode is debug.
    DebugMode = "debug"
    // ReleaseMode indicates gin mode is release.
    ReleaseMode = "release"
    // TestMode indicates gin mode is test.
    TestMode = "test"
)
```

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/auth.go#L20)

```go
const AuthProxyUserKey = "proxy_user"
```

AuthProxyUserKey is the cookie name for proxy_user credential in basic auth for proxy.

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/auth.go#L17)

```go
const AuthUserKey = "user"
```

AuthUserKey is the cookie name for user credential in basic auth.

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/utils.go#L19)

```go
const BindKey = "_gin-gonic/gin/bindkey"
```

BindKey indicates a default bind key.

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L41)

```go
const BodyBytesKey = "_gin-gonic/gin/bodybyteskey"
```

BodyBytesKey indicates a default body bytes key.

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L44)

```go
const ContextKey = "_gin-gonic/gin/contextkey"
```

ContextKey is the key that a Context returns itself for.

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/mode.go#L16)

```go
const EnvGinMode = "GIN_MODE"
```

EnvGinMode indicates environment name for gin mode.

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/version.go#L8)

```go
const Version = "v1.10.0"
```

Version is the current gin framework's version.

### Variables {#pkg-variables}

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/debug.go#L27)

```go
var DebugPrintFunc func(format string, values ...interface{})
```

DebugPrintFunc indicates debug log output format.

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/debug.go#L24)

```go
var DebugPrintRouteFunc func(httpMethod, absolutePath, handlerName string, nuHandlers int)
```

DebugPrintRouteFunc indicates debug log output format.

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/mode.go#L44)

```go
var DefaultErrorWriter io.Writer = os.Stderr
```

DefaultErrorWriter is the default io.Writer used by Gin to debug errors

[View Source](https://github.com/gin-gonic/gin/blob/v1.10.0/mode.go#L41)

```go
var DefaultWriter io.Writer = os.Stdout
```

DefaultWriter is the default io.Writer used by Gin for debug output and middleware output like Logger() or Recovery(). Note that both Logger and Recovery provides custom ways to configure their output io.Writer. To support coloring in Windows use:

```
import "github.com/mattn/go-colorable"
gin.DefaultWriter = colorable.NewColorableStdout()
```

### Functions {#pkg-functions}

#### func [CreateTestContext](https://github.com/gin-gonic/gin/blob/v1.10.0/test_helpers.go#L10) <- v1.3.0{#CreateTestContext}

```go
func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine)
```

CreateTestContext returns a fresh engine and context for testing purposes

#### func [Dir](https://github.com/gin-gonic/gin/blob/v1.10.0/fs.go#L24) {#Dir}

```go
func Dir(root string, listDirectory bool) http.FileSystem
```

Dir returns a http.FileSystem that can be used by http.FileServer(). It is used internally in router.Static(). if listDirectory == true, then it works the same as http.Dir() otherwise it returns a filesystem that prevents http.FileServer() to list the directory files.

#### func [DisableBindValidation](https://github.com/gin-gonic/gin/blob/v1.10.0/mode.go#L81) {#DisableBindValidation}

```go
func DisableBindValidation()
```

DisableBindValidation closes the default validator.

#### func [DisableConsoleColor](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L164) <- v1.3.0{#DisableConsoleColor}

```go
func DisableConsoleColor()
```

DisableConsoleColor disables color output in the console.

#### func [EnableJsonDecoderDisallowUnknownFields](https://github.com/gin-gonic/gin/blob/v1.10.0/mode.go#L93) <- v1.5.0{#EnableJsonDecoderDisallowUnknownFields}

```go
func EnableJsonDecoderDisallowUnknownFields()
```

EnableJsonDecoderDisallowUnknownFields sets true for binding.EnableDecoderDisallowUnknownFields to call the DisallowUnknownFields method on the JSON Decoder instance.

#### func [EnableJsonDecoderUseNumber](https://github.com/gin-gonic/gin/blob/v1.10.0/mode.go#L87) <- v1.3.0{#EnableJsonDecoderUseNumber}

```go
func EnableJsonDecoderUseNumber()
```

EnableJsonDecoderUseNumber sets true for binding.EnableDecoderUseNumber to call the UseNumber method on the JSON Decoder instance.

#### func [ForceConsoleColor](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L169) <- v1.4.0{#ForceConsoleColor}

```go
func ForceConsoleColor()
```

ForceConsoleColor force color output in the console.

#### func [IsDebugging](https://github.com/gin-gonic/gin/blob/v1.10.0/debug.go#L19) {#IsDebugging}

```go
func IsDebugging() bool
```

IsDebugging returns true if the framework is running in debug mode. Use SetMode(gin.ReleaseMode) to disable debug mode.

#### func [Mode](https://github.com/gin-gonic/gin/blob/v1.10.0/mode.go#L98) {#Mode}

```go
func Mode() string
```

Mode returns current gin mode.

#### func [SetMode](https://github.com/gin-gonic/gin/blob/v1.10.0/mode.go#L57) {#SetMode}

```go
func SetMode(value string)
```

SetMode sets gin mode according to input string.

### Types {#pkg-types}

#### type [Accounts](https://github.com/gin-gonic/gin/blob/v1.10.0/auth.go#L23) {#Accounts}

```go
type Accounts map[string]string
```

Accounts defines a key/value for user/pass list of authorized logins.

#### type [Context](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L55) {#Context}

```go
type Context struct {
    Request *http.Request
    Writer  ResponseWriter

    Params Params

    // Keys is a key/value pair exclusively for the context of each request.
    Keys map[string]any

    // Errors is a list of errors attached to all the handlers/middlewares who used this context.
    Errors errorMsgs

    // Accepted defines a list of manually accepted formats for content negotiation.
    Accepted []string
    // contains filtered or unexported fields

}
```

Context is the most important part of gin. It allows us to pass variables between middleware, manage the flow, validate the JSON of a request and render a JSON response for example.

##### func [CreateTestContextOnly](https://github.com/gin-gonic/gin/blob/v1.10.0/test_helpers.go#L19) <- v1.8.2{#CreateTestContextOnly}

```go
func CreateTestContextOnly(w http.ResponseWriter, r *Engine) (c *Context)
```

CreateTestContextOnly returns a fresh context base on the engine for testing purposes

##### func (*Context) [Abort](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L199) {#Context.Abort}

```go
func (c *Context) Abort()
```

Abort prevents pending handlers from being called. Note that this will not stop the current handler. Let's say you have an authorization middleware that validates that the current request is authorized. If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers for this request are not called.

##### func (*Context) [AbortWithError](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L222) {#Context.AbortWithError}

```go
func (c *Context) AbortWithError(code int, err error) *Error
```

AbortWithError calls `AbortWithStatus()` and `Error()` internally. This method stops the chain, writes the status code and pushes the specified error to `c.Errors`. See Context.Error() for more details.

##### func (*Context) [AbortWithStatus](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L205) {#Context.AbortWithStatus}

```go
func (c *Context) AbortWithStatus(code int)
```

AbortWithStatus calls `Abort()` and writes the headers with the specified status code. For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401).

##### func (*Context) [AbortWithStatusJSON](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L214) <- v1.3.0{#Context.AbortWithStatusJSON}

```go
func (c *Context) AbortWithStatusJSON(code int, jsonObj any)
```

AbortWithStatusJSON calls `Abort()` and then `JSON` internally. This method stops the chain, writes the status code and return a JSON body. It also sets the Content-Type as "application/json".

##### func (*Context) [AddParam](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L413) <- v1.8.0{#Context.AddParam}

```go
func (c *Context) AddParam(key, value string)
```

AddParam adds param to context and replaces path param key with given value for e2e testing purposes Example Route: "/user/:id" AddParam("id", 1) Result: "/user/1"

##### func (*Context) [AsciiJSON](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1010) <- v1.3.0{#Context.AsciiJSON}

```go
func (c *Context) AsciiJSON(code int, obj any)
```

AsciiJSON serializes the given struct as JSON into the response body with unicode to ASCII string. It also sets the Content-Type as "application/json".

##### func (*Context) [Bind](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L640) {#Context.Bind}

```go
func (c *Context) Bind(obj any) error
```

Bind checks the Method and Content-Type to select a binding engine automatically, Depending on the "Content-Type" header different bindings are used, for example:

```
"application/json" --> JSON binding
"application/xml"  --> XML binding
```

It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input. It decodes the json payload into the struct specified as a pointer. It writes a 400 error and sets Content-Type header "text/plain" in the response if input is not valid.

##### func (*Context) [BindHeader](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L671) <- v1.5.0{#Context.BindHeader}

```go
func (c *Context) BindHeader(obj any) error
```

BindHeader is a shortcut for c.MustBindWith(obj, binding.Header).

##### func (*Context) [BindJSON](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L646) {#Context.BindJSON}

```go
func (c *Context) BindJSON(obj any) error
```

BindJSON is a shortcut for c.MustBindWith(obj, binding.JSON).

##### func (*Context) [BindQuery](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L656) <- v1.3.0{#Context.BindQuery}

```go
func (c *Context) BindQuery(obj any) error
```

BindQuery is a shortcut for c.MustBindWith(obj, binding.Query).

##### func (*Context) [BindTOML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L666) <- v1.8.0{#Context.BindTOML}

```go
func (c *Context) BindTOML(obj any) error
```

BindTOML is a shortcut for c.MustBindWith(obj, binding.TOML).

##### func (*Context) [BindUri](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L677) <- v1.4.0{#Context.BindUri}

```go
func (c *Context) BindUri(obj any) error
```

BindUri binds the passed struct pointer using binding.Uri. It will abort the request with HTTP 400 if any error occurs.

##### func (*Context)[BindWith](https://github.com/gin-gonic/gin/blob/v1.10.0/deprecated.go#L17)<- DEPRECATED

```go
func (c *Context) BindWith(obj any, b binding.Binding) error
```

BindWith binds the passed struct pointer using the specified binding engine. See the binding package.

Deprecated: Use MustBindWith or ShouldBindWith.

##### func (*Context) [BindXML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L651) <- v1.4.0{#Context.BindXML}

```go
func (c *Context) BindXML(obj any) error
```

BindXML is a shortcut for c.MustBindWith(obj, binding.BindXML).

##### func (*Context) [BindYAML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L661) <- v1.4.0{#Context.BindYAML}

```go
func (c *Context) BindYAML(obj any) error
```

BindYAML is a shortcut for c.MustBindWith(obj, binding.YAML).

##### func (*Context) [ClientIP](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L802) {#Context.ClientIP}

```go
func (c *Context) ClientIP() string
```

ClientIP implements one best effort algorithm to return the real client IP. It calls c.RemoteIP() under the hood, to check if the remote IP is a trusted proxy or not. If it is it will then try to parse the headers defined in Engine.RemoteIPHeaders (defaulting to [X-Forwarded-For, X-Real-Ip]). If the headers are not syntactically valid OR the remote IP does not correspond to a trusted proxy, the remote IP (coming from Request.RemoteAddr) is returned.

##### func (*Context) [ContentType](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L849) {#Context.ContentType}

```go
func (c *Context) ContentType() string
```

ContentType returns the Content-Type header of the request.

##### func (*Context) [Cookie](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L941) {#Context.Cookie}

```go
func (c *Context) Cookie(name string) (string, error)
```

Cookie returns the named cookie provided in the request or ErrNoCookie if not found. And return the named cookie is unescaped. If multiple cookies match the given name, only one cookie will be returned.

##### func (*Context) [Copy](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L116) {#Context.Copy}

```go
func (c *Context) Copy() *Context
```

Copy returns a copy of the current context that can be safely used outside the request's scope. This has to be used when the context has to be passed to a goroutine.

##### func (*Context) [Data](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1056) {#Context.Data}

```go
func (c *Context) Data(code int, contentType string, data []byte)
```

Data writes some data into the body stream and updates the HTTP code.

##### func (*Context) [DataFromReader](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1064) <- v1.3.0{#Context.DataFromReader}

```go
func (c *Context) DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string)
```

DataFromReader writes the specified reader into the body stream and updates the HTTP code.

##### func (*Context) [Deadline](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1225) {#Context.Deadline}

```go
func (c *Context) Deadline() (deadline time.Time, ok bool)
```

Deadline returns that there is no deadline (ok==false) when c.Request has no Context.

##### func (*Context) [DefaultPostForm](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L510) {#Context.DefaultPostForm}

```go
func (c *Context) DefaultPostForm(key, defaultValue string) string
```

DefaultPostForm returns the specified key from a POST urlencoded form or multipart form when it exists, otherwise it returns the specified defaultValue string. See: PostForm() and GetPostForm() for further information.

##### func (*Context) [DefaultQuery](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L439) {#Context.DefaultQuery}

```go
func (c *Context) DefaultQuery(key, defaultValue string) string
```

DefaultQuery returns the keyed url query value if it exists, otherwise it returns the specified defaultValue string. See: Query() and GetQuery() for further information.

```
GET /?name=Manu&lastname=
c.DefaultQuery("name", "unknown") == "Manu"
c.DefaultQuery("id", "none") == "none"
c.DefaultQuery("lastname", "none") == ""
```

##### func (*Context) [Done](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1233) {#Context.Done}

```go
func (c *Context) Done() <-chan struct{}
```

Done returns nil (chan which will wait forever) when c.Request has no Context.

##### func (*Context) [Err](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1241) {#Context.Err}

```go
func (c *Context) Err() error
```

Err returns nil when c.Request has no Context.

##### func (*Context) [Error](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L236) {#Context.Error}

```go
func (c *Context) Error(err error) *Error
```

Error attaches an error to the current context. The error is pushed to a list of errors. It's a good idea to call Error for each error that occurred during the resolution of a request. A middleware can be used to collect all the errors and push them to a database together, print a log, or append it in the HTTP response. Error will panic if err is nil.

##### func (*Context) [File](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1074) {#Context.File}

```go
func (c *Context) File(filepath string)
```

File writes the specified file into the body stream in an efficient way.

##### func (*Context) [FileAttachment](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1097) <- v1.4.0{#Context.FileAttachment}

```go
func (c *Context) FileAttachment(filepath, filename string)
```

FileAttachment writes the specified file into the body stream in an efficient way On the client side, the file will typically be downloaded with the given filename

##### func (*Context) [FileFromFS](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1079) <- v1.6.0{#Context.FileFromFS}

```go
func (c *Context) FileFromFS(filepath string, fs http.FileSystem)
```

FileFromFS writes the specified file from http.FileSystem into the body stream in an efficient way.

##### func (*Context) [FormFile](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L589) <- v1.3.0{#Context.FormFile}

```go
func (c *Context) FormFile(name string) (*multipart.FileHeader, error)
```

FormFile returns the first file for the provided form key.

##### func (*Context) [FullPath](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L171) <- v1.5.0{#Context.FullPath}

```go
func (c *Context) FullPath() string
```

FullPath returns a matched route full path. For not found routes returns an empty string.

```
router.GET("/user/:id", func(c *gin.Context) {
    c.FullPath() == "/user/:id" // true
})
```

##### func (*Context) [Get](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L272) {#Context.Get}

```go
func (c *Context) Get(key string) (value any, exists bool)
```

Get returns the value for the given key, ie: (value, true). If the value does not exist it returns (nil, false)

##### func (*Context) [GetBool](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L296) <- v1.3.0{#Context.GetBool}

```go
func (c *Context) GetBool(key string) (b bool)
```

GetBool returns the value associated with the key as a boolean.

##### func (*Context) [GetDuration](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L352) <- v1.3.0{#Context.GetDuration}

```go
func (c *Context) GetDuration(key string) (d time.Duration)
```

GetDuration returns the value associated with the key as a duration.

##### func (*Context) [GetFloat64](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L336) <- v1.3.0{#Context.GetFloat64}

```go
func (c *Context) GetFloat64(key string) (f64 float64)
```

GetFloat64 returns the value associated with the key as a float64.

##### func (*Context) [GetHeader](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L901) <- v1.3.0{#Context.GetHeader}

```go
func (c *Context) GetHeader(key string) string
```

GetHeader returns value from request headers.

##### func (*Context) [GetInt](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L304) <- v1.3.0{#Context.GetInt}

```go
func (c *Context) GetInt(key string) (i int)
```

GetInt returns the value associated with the key as an integer.

##### func (*Context) [GetInt64](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L312) <- v1.3.0{#Context.GetInt64}

```go
func (c *Context) GetInt64(key string) (i64 int64)
```

GetInt64 returns the value associated with the key as an integer.

##### func (*Context) [GetPostForm](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L525) {#Context.GetPostForm}

```go
func (c *Context) GetPostForm(key string) (string, bool)
```

GetPostForm is like PostForm(key). It returns the specified key from a POST urlencoded form or multipart form when it exists `(value, true)` (even when the value is an empty string), otherwise it returns ("", false). For example, during a PATCH request to update the user's email:

```
    email=mail@example.com  -->  ("mail@example.com", true) := GetPostForm("email") // set email to "mail@example.com"
	   email=                  -->  ("", true) := GetPostForm("email") // set email to ""
                            -->  ("", false) := GetPostForm("email") // do nothing with email
```

##### func (*Context) [GetPostFormArray](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L554) {#Context.GetPostFormArray}

```go
func (c *Context) GetPostFormArray(key string) (values []string, ok bool)
```

GetPostFormArray returns a slice of strings for a given form key, plus a boolean value whether at least one value exists for the given key.

##### func (*Context) [GetPostFormMap](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L568) <- v1.3.0{#Context.GetPostFormMap}

```go
func (c *Context) GetPostFormMap(key string) (map[string]string, bool)
```

GetPostFormMap returns a map for a given form key, plus a boolean value whether at least one value exists for the given key.

##### func (*Context) [GetQuery](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L455) {#Context.GetQuery}

```go
func (c *Context) GetQuery(key string) (string, bool)
```

GetQuery is like Query(), it returns the keyed url query value if it exists `(value, true)` (even when the value is an empty string), otherwise it returns `("", false)`. It is shortcut for `c.Request.URL.Query().Get(key)`

```
GET /?name=Manu&lastname=
("Manu", true) == c.GetQuery("name")
("", false) == c.GetQuery("id")
("", true) == c.GetQuery("lastname")
```

##### func (*Context) [GetQueryArray](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L481) {#Context.GetQueryArray}

```go
func (c *Context) GetQueryArray(key string) (values []string, ok bool)
```

GetQueryArray returns a slice of strings for a given query key, plus a boolean value whether at least one value exists for the given key.

##### func (*Context) [GetQueryMap](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L495) <- v1.3.0{#Context.GetQueryMap}

```go
func (c *Context) GetQueryMap(key string) (map[string]string, bool)
```

GetQueryMap returns a map for a given query key, plus a boolean value whether at least one value exists for the given key.

##### func (*Context) [GetRawData](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L906) <- v1.3.0{#Context.GetRawData}

```go
func (c *Context) GetRawData() ([]byte, error)
```

GetRawData returns stream data.

##### func (*Context) [GetString](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L288) <- v1.3.0{#Context.GetString}

```go
func (c *Context) GetString(key string) (s string)
```

GetString returns the value associated with the key as a string.

##### func (*Context) [GetStringMap](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L368) <- v1.3.0{#Context.GetStringMap}

```go
func (c *Context) GetStringMap(key string) (sm map[string]any)
```

GetStringMap returns the value associated with the key as a map of interfaces.

##### func (*Context) [GetStringMapString](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L376) <- v1.3.0{#Context.GetStringMapString}

```go
func (c *Context) GetStringMapString(key string) (sms map[string]string)
```

GetStringMapString returns the value associated with the key as a map of strings.

##### func (*Context) [GetStringMapStringSlice](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L384) <- v1.3.0{#Context.GetStringMapStringSlice}

```go
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string)
```

GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.

##### func (*Context) [GetStringSlice](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L360) <- v1.3.0{#Context.GetStringSlice}

```go
func (c *Context) GetStringSlice(key string) (ss []string)
```

GetStringSlice returns the value associated with the key as a slice of strings.

##### func (*Context) [GetTime](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L344) <- v1.3.0{#Context.GetTime}

```go
func (c *Context) GetTime(key string) (t time.Time)
```

GetTime returns the value associated with the key as time.

##### func (*Context) [GetUint](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L320) <- v1.7.0{#Context.GetUint}

```go
func (c *Context) GetUint(key string) (ui uint)
```

GetUint returns the value associated with the key as an unsigned integer.

##### func (*Context) [GetUint64](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L328) <- v1.7.0{#Context.GetUint64}

```go
func (c *Context) GetUint64(key string) (ui64 uint64)
```

GetUint64 returns the value associated with the key as an unsigned integer.

##### func (*Context) [HTML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L970) {#Context.HTML}

```go
func (c *Context) HTML(code int, name string, obj any)
```

HTML renders the HTTP template specified by its file name. It also updates the HTTP code and sets the Content-Type as "text/html". See http://golang.org/doc/articles/wiki/

##### func (*Context) [Handler](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L161) <- v1.3.0{#Context.Handler}

```go
func (c *Context) Handler() HandlerFunc
```

Handler returns the main handler.

##### func (*Context) [HandlerName](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L146) {#Context.HandlerName}

```go
func (c *Context) HandlerName() string
```

HandlerName returns the main handler's name. For example if the handler is "handleGetUsers()", this function will return "main.handleGetUsers".

##### func (*Context) [HandlerNames](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L152) <- v1.4.0{#Context.HandlerNames}

```go
func (c *Context) HandlerNames() []string
```

HandlerNames returns a list of all registered handlers for this context in descending order, following the semantics of HandlerName()

##### func (*Context) [Header](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L892) {#Context.Header}

```go
func (c *Context) Header(key, value string)
```

Header is an intelligent shortcut for c.Writer.Header().Set(key, value). It writes a header in the response. If value == "", this method removes the header `c.Writer.Header().Del(key)`

##### func (*Context) [IndentedJSON](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L979) {#Context.IndentedJSON}

```go
func (c *Context) IndentedJSON(code int, obj any)
```

IndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body. It also sets the Content-Type as "application/json". WARNING: we recommend using this only for development purposes since printing pretty JSON is more CPU and bandwidth consuming. Use Context.JSON() instead.

##### func (*Context) [IsAborted](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L191) {#Context.IsAborted}

```go
func (c *Context) IsAborted() bool
```

IsAborted returns true if the current context was aborted.

##### func (*Context) [IsWebsocket](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L855) <- v1.3.0{#Context.IsWebsocket}

```go
func (c *Context) IsWebsocket() bool
```

IsWebsocket returns true if the request headers indicate that a websocket handshake is being initiated by the client.

##### func (*Context) [JSON](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1004) {#Context.JSON}

```go
func (c *Context) JSON(code int, obj any)
```

JSON serializes the given struct as JSON into the response body. It also sets the Content-Type as "application/json".

##### func (*Context) [JSONP](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L993) <- v1.3.0{#Context.JSONP}

```go
func (c *Context) JSONP(code int, obj any)
```

JSONP serializes the given struct as JSON into the response body. It adds padding to response body to request data from a server residing in a different domain than the client. It also sets the Content-Type as "application/javascript".

##### func (*Context) [MultipartForm](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L604) <- v1.3.0{#Context.MultipartForm}

```go
func (c *Context) MultipartForm() (*multipart.Form, error)
```

MultipartForm is the parsed multipart form, including file uploads.

##### func (*Context) [MustBindWith](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L688) <- v1.3.0{#Context.MustBindWith}

```go
func (c *Context) MustBindWith(obj any, b binding.Binding) error
```

MustBindWith binds the passed struct pointer using the specified binding engine. It will abort the request with HTTP 400 if any error occurs. See the binding package.

##### func (*Context) [MustGet](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L280) {#Context.MustGet}

```go
func (c *Context) MustGet(key string) any
```

MustGet returns the value for the given key if it exists, otherwise it panics.

##### func (*Context) [Negotiate](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1150) {#Context.Negotiate}

```go
func (c *Context) Negotiate(code int, config Negotiate)
```

Negotiate calls different Render according to acceptable Accept format.

##### func (*Context) [NegotiateFormat](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1178) {#Context.NegotiateFormat}

```go
func (c *Context) NegotiateFormat(offered ...string) string
```

NegotiateFormat returns an acceptable Accept format.

##### func (*Context) [Next](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L182) {#Context.Next}

```go
func (c *Context) Next()
```

Next should be used only inside middleware. It executes the pending handlers in the chain inside the calling handler. See example in GitHub.

##### func (*Context) [Param](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L404) {#Context.Param}

```go
func (c *Context) Param(key string) string
```

Param returns the value of the URL param. It is a shortcut for c.Params.ByName(key)

```
router.GET("/user/:id", func(c *gin.Context) {
    // a GET request to /user/john
    id := c.Param("id") // id == "john"
    // a GET request to /user/john/
    id := c.Param("id") // id == "/john/"
})
```

##### func (*Context) [PostForm](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L502) {#Context.PostForm}

```go
func (c *Context) PostForm(key string) (value string)
```

PostForm returns the specified key from a POST urlencoded form or multipart form when it exists, otherwise it returns an empty string `("")`.

##### func (*Context) [PostFormArray](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L534) {#Context.PostFormArray}

```go
func (c *Context) PostFormArray(key string) (values []string)
```

PostFormArray returns a slice of strings for a given form key. The length of the slice depends on the number of params with the given key.

##### func (*Context) [PostFormMap](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L561) <- v1.3.0{#Context.PostFormMap}

```go
func (c *Context) PostFormMap(key string) (dicts map[string]string)
```

PostFormMap returns a map for a given form key.

##### func (*Context) [ProtoBuf](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1037) <- v1.4.0{#Context.ProtoBuf}

```go
func (c *Context) ProtoBuf(code int, obj any)
```

ProtoBuf serializes the given struct as ProtoBuf into the response body.

##### func (*Context) [PureJSON](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1016) <- v1.4.0{#Context.PureJSON}

```go
func (c *Context) PureJSON(code int, obj any)
```

PureJSON serializes the given struct as JSON into the response body. PureJSON, unlike JSON, does not replace special html characters with their unicode entities.

##### func (*Context) [Query](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L426) {#Context.Query}

```go
func (c *Context) Query(key string) (value string)
```

Query returns the keyed url query value if it exists, otherwise it returns an empty string `("")`. It is shortcut for `c.Request.URL.Query().Get(key)`

```
    GET /path?id=1234&name=Manu&value=
	   c.Query("id") == "1234"
	   c.Query("name") == "Manu"
	   c.Query("value") == ""
	   c.Query("wtf") == ""
```

##### func (*Context) [QueryArray](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L464) {#Context.QueryArray}

```go
func (c *Context) QueryArray(key string) (values []string)
```

QueryArray returns a slice of strings for a given query key. The length of the slice depends on the number of params with the given key.

##### func (*Context) [QueryMap](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L488) <- v1.3.0{#Context.QueryMap}

```go
func (c *Context) QueryMap(key string) (dicts map[string]string)
```

QueryMap returns a map for a given query key.

##### func (*Context) [Redirect](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1047) {#Context.Redirect}

```go
func (c *Context) Redirect(code int, location string)
```

Redirect returns an HTTP redirect to the specific location.

##### func (*Context) [RemoteIP](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L840) <- v1.7.0{#Context.RemoteIP}

```go
func (c *Context) RemoteIP() string
```

RemoteIP parses the IP from Request.RemoteAddr, normalizes and returns the IP (without the port).

##### func (*Context) [Render](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L951) {#Context.Render}

```go
func (c *Context) Render(code int, r render.Render)
```

Render writes the response headers and calls render.Render to render data.

##### func (*Context) [SSEvent](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1107) {#Context.SSEvent}

```go
func (c *Context) SSEvent(name string, message any)
```

SSEvent writes a Server-Sent Event into the body stream.

##### func (*Context) [SaveUploadedFile](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L610) <- v1.3.0{#Context.SaveUploadedFile}

```go
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error
```

SaveUploadedFile uploads the form file to specific dst.

##### func (*Context) [SecureJSON](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L986) <- v1.3.0{#Context.SecureJSON}

```go
func (c *Context) SecureJSON(code int, obj any)
```

SecureJSON serializes the given struct as Secure JSON into the response body. Default prepends "while(1)," to response body if the given struct is array values. It also sets the Content-Type as "application/json".

##### func (*Context) [Set](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L260) {#Context.Set}

```go
func (c *Context) Set(key string, value any)
```

Set is used to store a new key/value pair exclusively for this context. It also lazy initializes c.Keys if it was not used previously.

##### func (*Context) [SetAccepted](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1209) {#Context.SetAccepted}

```go
func (c *Context) SetAccepted(formats ...string)
```

SetAccepted sets Accept header data.

##### func (*Context) [SetCookie](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L921) {#Context.SetCookie}

```go
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
```

SetCookie adds a Set-Cookie header to the ResponseWriter's headers. The provided cookie must have a valid Name. Invalid cookies may be silently dropped.

##### func (*Context) [SetSameSite](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L914) <- v1.6.2{#Context.SetSameSite}

```go
func (c *Context) SetSameSite(samesite http.SameSite)
```

SetSameSite with cookie

##### func (*Context) [ShouldBind](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L705) <- v1.3.0{#Context.ShouldBind}

```go
func (c *Context) ShouldBind(obj any) error
```

ShouldBind checks the Method and Content-Type to select a binding engine automatically, Depending on the "Content-Type" header different bindings are used, for example:

```
"application/json" --> JSON binding
"application/xml"  --> XML binding
```

It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input. It decodes the json payload into the struct specified as a pointer. Like c.Bind() but this method does not set the response status code to 400 or abort if input is not valid.

##### func (*Context) [ShouldBindBodyWith](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L760) <- v1.3.0{#Context.ShouldBindBodyWith}

```go
func (c *Context) ShouldBindBodyWith(obj any, bb binding.BindingBody) (err error)
```

ShouldBindBodyWith is similar with ShouldBindWith, but it stores the request body into the context, and reuse when it is called again.

NOTE: This method reads the body before binding. So you should use ShouldBindWith for better performance if you need to call only once.

##### func (*Context) [ShouldBindBodyWithJSON](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L778) <- v1.10.0{#Context.ShouldBindBodyWithJSON}

```go
func (c *Context) ShouldBindBodyWithJSON(obj any) error
```

ShouldBindBodyWithJSON is a shortcut for c.ShouldBindBodyWith(obj, binding.JSON).

##### func (*Context) [ShouldBindBodyWithTOML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L793) <- v1.10.0{#Context.ShouldBindBodyWithTOML}

```go
func (c *Context) ShouldBindBodyWithTOML(obj any) error
```

ShouldBindBodyWithTOML is a shortcut for c.ShouldBindBodyWith(obj, binding.TOML).

##### func (*Context) [ShouldBindBodyWithXML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L783) <- v1.10.0{#Context.ShouldBindBodyWithXML}

```go
func (c *Context) ShouldBindBodyWithXML(obj any) error
```

ShouldBindBodyWithXML is a shortcut for c.ShouldBindBodyWith(obj, binding.XML).

##### func (*Context) [ShouldBindBodyWithYAML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L788) <- v1.10.0{#Context.ShouldBindBodyWithYAML}

```go
func (c *Context) ShouldBindBodyWithYAML(obj any) error
```

ShouldBindBodyWithYAML is a shortcut for c.ShouldBindBodyWith(obj, binding.YAML).

##### func (*Context) [ShouldBindHeader](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L736) <- v1.5.0{#Context.ShouldBindHeader}

```go
func (c *Context) ShouldBindHeader(obj any) error
```

ShouldBindHeader is a shortcut for c.ShouldBindWith(obj, binding.Header).

##### func (*Context) [ShouldBindJSON](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L711) <- v1.3.0{#Context.ShouldBindJSON}

```go
func (c *Context) ShouldBindJSON(obj any) error
```

ShouldBindJSON is a shortcut for c.ShouldBindWith(obj, binding.JSON).

##### func (*Context) [ShouldBindQuery](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L721) <- v1.3.0{#Context.ShouldBindQuery}

```go
func (c *Context) ShouldBindQuery(obj any) error
```

ShouldBindQuery is a shortcut for c.ShouldBindWith(obj, binding.Query).

##### func (*Context) [ShouldBindTOML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L731) <- v1.8.0{#Context.ShouldBindTOML}

```go
func (c *Context) ShouldBindTOML(obj any) error
```

ShouldBindTOML is a shortcut for c.ShouldBindWith(obj, binding.TOML).

##### func (*Context) [ShouldBindUri](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L741) <- v1.4.0{#Context.ShouldBindUri}

```go
func (c *Context) ShouldBindUri(obj any) error
```

ShouldBindUri binds the passed struct pointer using the specified binding engine.

##### func (*Context) [ShouldBindWith](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L751) <- v1.3.0{#Context.ShouldBindWith}

```go
func (c *Context) ShouldBindWith(obj any, b binding.Binding) error
```

ShouldBindWith binds the passed struct pointer using the specified binding engine. See the binding package.

##### func (*Context) [ShouldBindXML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L716) <- v1.4.0{#Context.ShouldBindXML}

```go
func (c *Context) ShouldBindXML(obj any) error
```

ShouldBindXML is a shortcut for c.ShouldBindWith(obj, binding.XML).

##### func (*Context) [ShouldBindYAML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L726) <- v1.4.0{#Context.ShouldBindYAML}

```go
func (c *Context) ShouldBindYAML(obj any) error
```

ShouldBindYAML is a shortcut for c.ShouldBindWith(obj, binding.YAML).

##### func (*Context) [Status](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L885) {#Context.Status}

```go
func (c *Context) Status(code int)
```

Status sets the HTTP response code.

##### func (*Context) [Stream](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1116) {#Context.Stream}

```go
func (c *Context) Stream(step func(w io.Writer) bool) bool
```

Stream sends a streaming response and returns a boolean indicates "Is client disconnected in middle of stream"

##### func (*Context) [String](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1042) {#Context.String}

```go
func (c *Context) String(code int, format string, values ...any)
```

String writes the given string into the response body.

##### func (*Context) [TOML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1032) <- v1.8.0{#Context.TOML}

```go
func (c *Context) TOML(code int, obj any)
```

TOML serializes the given struct as TOML into the response body.

##### func (*Context) [Value](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1251) {#Context.Value}

```go
func (c *Context) Value(key any) any
```

Value returns the value associated with this context for key, or nil if no value is associated with key. Successive calls to Value with the same key returns the same result.

##### func (*Context) [XML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1022) {#Context.XML}

```go
func (c *Context) XML(code int, obj any)
```

XML serializes the given struct as XML into the response body. It also sets the Content-Type as "application/xml".

##### func (*Context) [YAML](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1027) {#Context.YAML}

```go
func (c *Context) YAML(code int, obj any)
```

YAML serializes the given struct as YAML into the response body.

#### type [ContextKeyType](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L46) <- v1.10.0{#ContextKeyType}

```go
type ContextKeyType int
const ContextRequestKey ContextKeyType = 0
```

#### type [Engine](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L89) {#Engine}

```go
type Engine struct {
    RouterGroup

    // RedirectTrailingSlash enables automatic redirection if the current route can't be matched but a
    // handler for the path with (without) the trailing slash exists.
    // For example if /foo/ is requested but a route only exists for /foo, the
    // client is redirected to /foo with http status code 301 for GET requests
    // and 307 for all other request methods.
    RedirectTrailingSlash bool

    // RedirectFixedPath if enabled, the router tries to fix the current request path, if no
    // handle is registered for it.
    // First superfluous path elements like ../ or // are removed.
    // Afterwards the router does a case-insensitive lookup of the cleaned path.
    // If a handle can be found for this route, the router makes a redirection
    // to the corrected path with status code 301 for GET requests and 307 for
    // all other request methods.
    // For example /FOO and /..//Foo could be redirected to /foo.
    // RedirectTrailingSlash is independent of this option.
    RedirectFixedPath bool

    // HandleMethodNotAllowed if enabled, the router checks if another method is allowed for the
    // current route, if the current request can not be routed.
    // If this is the case, the request is answered with 'Method Not Allowed'
    // and HTTP status code 405.
    // If no other Method is allowed, the request is delegated to the NotFound
    // handler.
    HandleMethodNotAllowed bool

    // ForwardedByClientIP if enabled, client IP will be parsed from the request's headers that
    // match those stored at `(*gin.Engine).RemoteIPHeaders`. If no IP was
    // fetched, it falls back to the IP obtained from
    // `(*gin.Context).Request.RemoteAddr`.
    ForwardedByClientIP bool

    // AppEngine was deprecated.
    // Deprecated: USE `TrustedPlatform` WITH VALUE `gin.PlatformGoogleAppEngine` INSTEAD
    // #726 #755 If enabled, it will trust some headers starting with
    // 'X-AppEngine...' for better integration with that PaaS.
    AppEngine bool

    // UseRawPath if enabled, the url.RawPath will be used to find parameters.
    UseRawPath bool

    // UnescapePathValues if true, the path value will be unescaped.
    // If UseRawPath is false (by default), the UnescapePathValues effectively is true,
    // as url.Path gonna be used, which is already unescaped.
    UnescapePathValues bool

    // RemoveExtraSlash a parameter can be parsed from the URL even with extra slashes.
    // See the PR #1817 and issue #1644
    RemoveExtraSlash bool

    // RemoteIPHeaders list of headers used to obtain the client IP when
    // `(*gin.Engine).ForwardedByClientIP` is `true` and
    // `(*gin.Context).Request.RemoteAddr` is matched by at least one of the
    // network origins of list defined by `(*gin.Engine).SetTrustedProxies()`.
    RemoteIPHeaders []string

    // TrustedPlatform if set to a constant of value gin.Platform*, trusts the headers set by
    // that platform, for example to determine the client IP
    TrustedPlatform string

    // MaxMultipartMemory value of 'maxMemory' param that is given to http.Request's ParseMultipartForm
    // method call.
    MaxMultipartMemory int64

    // UseH2C enable h2c support.
    UseH2C bool

    // ContextWithFallback enable fallback Context.Deadline(), Context.Done(), Context.Err() and Context.Value() when Context.Request.Context() is not nil.
    ContextWithFallback bool

    HTMLRender render.HTMLRender
    FuncMap    template.FuncMap
    // contains filtered or unexported fields

}
```

Engine is the framework's instance, it contains the muxer, middleware and configuration settings. Create an instance of Engine, by using New() or Default()

##### func [Default](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L221) {#Default}

```go
func Default(opts ...OptionFunc) *Engine
```

Default returns an Engine instance with the Logger and Recovery middleware already attached.

##### func [New](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L188) {#New}

```go
func New(opts ...OptionFunc) *Engine
```

New returns a new blank Engine instance without any middleware attached. By default, the configuration is: - RedirectTrailingSlash: true - RedirectFixedPath: false - HandleMethodNotAllowed: false - ForwardedByClientIP: true - UseRawPath: false - UnescapePathValues: true

##### func (*Engine) [Delims](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L244) <- v1.3.0{#Engine.Delims}

```go
func (engine *Engine) Delims(left, right string) *Engine
```

Delims sets template left and right delims and returns an Engine instance.

##### func (*Engine) [HandleContext](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L597) <- v1.3.0{#Engine.HandleContext}

```go
func (engine *Engine) HandleContext(c *Context)
```

HandleContext re-enters a context that has been rewritten. This can be done by setting c.Request.URL.Path to your new target. Disclaimer: You can loop yourself to deal with this, use wisely.

##### func (*Engine) [Handler](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L228) <- v1.8.0{#Engine.Handler}

```go
func (engine *Engine) Handler() http.Handler
```

##### func (*Engine) [LoadHTMLFiles](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L273) {#Engine.LoadHTMLFiles}

```go
func (engine *Engine) LoadHTMLFiles(files ...string)
```

LoadHTMLFiles loads a slice of HTML files and associates the result with HTML renderer.

##### func (*Engine) [LoadHTMLGlob](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L257) {#Engine.LoadHTMLGlob}

```go
func (engine *Engine) LoadHTMLGlob(pattern string)
```

LoadHTMLGlob loads HTML files identified by glob pattern and associates the result with HTML renderer.

##### func (*Engine) [NoMethod](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L304) {#Engine.NoMethod}

```go
func (engine *Engine) NoMethod(handlers ...HandlerFunc)
```

NoMethod sets the handlers called when Engine.HandleMethodNotAllowed = true.

##### func (*Engine) [NoRoute](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L298) {#Engine.NoRoute}

```go
func (engine *Engine) NoRoute(handlers ...HandlerFunc)
```

NoRoute adds handlers for NoRoute. It returns a 404 code by default.

##### func (*Engine) [Routes](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L362) {#Engine.Routes}

```go
func (engine *Engine) Routes() (routes RoutesInfo)
```

Routes returns a slice of registered routes, including some useful information, such as: the http method, path and the handler name.

##### func (*Engine) [Run](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L389) {#Engine.Run}

```go
func (engine *Engine) Run(addr ...string) (err error)
```

Run attaches the router to a http.Server and starts listening and serving HTTP requests. It is a shortcut for http.ListenAndServe(addr, router) Note: this method will block the calling goroutine indefinitely unless an error happens.

##### func (*Engine) [RunFd](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L548) <- v1.4.0{#Engine.RunFd}

```go
func (engine *Engine) RunFd(fd int) (err error)
```

RunFd attaches the router to a http.Server and starts listening and serving HTTP requests through the specified file descriptor. Note: this method will block the calling goroutine indefinitely unless an error happens.

##### func (*Engine) [RunListener](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L569) <- v1.5.0{#Engine.RunListener}

```go
func (engine *Engine) RunListener(listener net.Listener) (err error)
```

RunListener attaches the router to a http.Server and starts listening and serving HTTP requests through the specified net.Listener

##### func (*Engine) [RunTLS](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L509) {#Engine.RunTLS}

```go
func (engine *Engine) RunTLS(addr, certFile, keyFile string) (err error)
```

RunTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests. It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router) Note: this method will block the calling goroutine indefinitely unless an error happens.

##### func (*Engine) [RunUnix](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L525) {#Engine.RunUnix}

```go
func (engine *Engine) RunUnix(file string) (err error)
```

RunUnix attaches the router to a http.Server and starts listening and serving HTTP requests through the specified unix socket (i.e. a file). Note: this method will block the calling goroutine indefinitely unless an error happens.

##### func (*Engine) [SecureJsonPrefix](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L250) <- v1.3.0{#Engine.SecureJsonPrefix}

```go
func (engine *Engine) SecureJsonPrefix(prefix string) *Engine
```

SecureJsonPrefix sets the secureJSONPrefix used in Context.SecureJSON.

##### func (*Engine) [ServeHTTP](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L583) {#Engine.ServeHTTP}

```go
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)
```

ServeHTTP conforms to the http.Handler interface.

##### func (*Engine) [SetFuncMap](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L293) <- v1.3.0{#Engine.SetFuncMap}

```go
func (engine *Engine) SetFuncMap(funcMap template.FuncMap)
```

SetFuncMap sets the FuncMap used for template.FuncMap.

##### func (*Engine) [SetHTMLTemplate](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L284) {#Engine.SetHTMLTemplate}

```go
func (engine *Engine) SetHTMLTemplate(templ *template.Template)
```

SetHTMLTemplate associate a template with HTML renderer.

##### func (*Engine) [SetTrustedProxies](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L440) <- v1.7.5{#Engine.SetTrustedProxies}

```go
func (engine *Engine) SetTrustedProxies(trustedProxies []string) error
```

SetTrustedProxies set a list of network origins (IPv4 addresses, IPv4 CIDRs, IPv6 addresses or IPv6 CIDRs) from which to trust request's headers that contain alternative client IP when `(*gin.Engine).ForwardedByClientIP` is `true`. `TrustedProxies` feature is enabled by default, and it also trusts all proxies by default. If you want to disable this feature, use Engine.SetTrustedProxies(nil), then Context.ClientIP() will return the remote address directly.

##### func (*Engine) [Use](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L312) {#Engine.Use}

```go
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes
```

Use attaches a global middleware to the router. i.e. the middleware attached through Use() will be included in the handlers chain for every single request. Even 404, 405, static files... For example, this is the right place for a logger or error management middleware.

##### func (*Engine) [With](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L320) <- v1.10.0{#Engine.With}

```go
func (engine *Engine) With(opts ...OptionFunc) *Engine
```

With returns a new Engine instance with the provided options.

#### type [Error](https://github.com/gin-gonic/gin/blob/v1.10.0/errors.go#L34) {#Error}

```go
type Error struct {
    Err  error
    Type ErrorType
    Meta any

}
```

Error represents a error's specification.

##### func (Error) [Error](https://github.com/gin-gonic/gin/blob/v1.10.0/errors.go#L84) {#Error.Error}

```go
func (msg Error) Error() string
```

Error implements the error interface.

##### func (*Error) [IsType](https://github.com/gin-gonic/gin/blob/v1.10.0/errors.go#L89) {#Error.IsType}

```go
func (msg *Error) IsType(flags ErrorType) bool
```

IsType judges one error.

##### func (*Error) [JSON](https://github.com/gin-gonic/gin/blob/v1.10.0/errors.go#L57) {#Error.JSON}

```go
func (msg *Error) JSON() any
```

JSON creates a properly formatted JSON

##### func (*Error) [MarshalJSON](https://github.com/gin-gonic/gin/blob/v1.10.0/errors.go#L79) {#Error.MarshalJSON}

```go
func (msg *Error) MarshalJSON() ([]byte, error)
```

MarshalJSON implements the json.Marshaller interface.

##### func (*Error) [SetMeta](https://github.com/gin-gonic/gin/blob/v1.10.0/errors.go#L51) {#Error.SetMeta}

```go
func (msg *Error) SetMeta(data any) *Error
```

SetMeta sets the error's meta data.

##### func (*Error) [SetType](https://github.com/gin-gonic/gin/blob/v1.10.0/errors.go#L45) {#Error.SetType}

```go
func (msg *Error) SetType(flags ErrorType) *Error
```

SetType sets the error's type.

##### func (*Error) [Unwrap](https://github.com/gin-gonic/gin/blob/v1.10.0/errors.go#L94) <- v1.7.0{#Error.Unwrap}

```go
func (msg *Error) Unwrap() error
```

Unwrap returns the wrapped error, to allow interoperability with errors.Is(), errors.As() and errors.Unwrap()

#### type [ErrorType](https://github.com/gin-gonic/gin/blob/v1.10.0/errors.go#L16) {#ErrorType}

```go
type ErrorType uint64
```

ErrorType is an unsigned 64-bit error code as defined in the gin spec.

```go
const (
    // ErrorTypeBind is used when Context.Bind() fails.
    ErrorTypeBind ErrorType = 1 << 63
    // ErrorTypeRender is used when Context.Render() fails.
    ErrorTypeRender ErrorType = 1 << 62
    // ErrorTypePrivate indicates a private error.
    ErrorTypePrivate ErrorType = 1 << 0
    // ErrorTypePublic indicates a public error.
    ErrorTypePublic ErrorType = 1 << 1
    // ErrorTypeAny indicates any other error.
    ErrorTypeAny ErrorType = 1<<64 - 1
    // ErrorTypeNu indicates any other error.
    ErrorTypeNu = 2

)
```

#### type [H](https://github.com/gin-gonic/gin/blob/v1.10.0/utils.go#L54) {#H}

```go
type H map[string]any
```

H is a shortcut for map[string]any

##### func (H) [MarshalXML](https://github.com/gin-gonic/gin/blob/v1.10.0/utils.go#L57) {#H.MarshalXML}

```go
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error
```

MarshalXML allows type H to be used with xml.Marshal.

#### type [HandlerFunc](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L48) {#HandlerFunc}

```go
type HandlerFunc func(*Context)
```

HandlerFunc defines the handler used by gin middleware as return value.

##### func [BasicAuth](https://github.com/gin-gonic/gin/blob/v1.10.0/auth.go#L72) {#BasicAuth}

```go
func BasicAuth(accounts Accounts) HandlerFunc
```

BasicAuth returns a Basic HTTP Authorization middleware. It takes as argument a map[string]string where the key is the user name and the value is the password.

##### func [BasicAuthForProxy](https://github.com/gin-gonic/gin/blob/v1.10.0/auth.go#L98) <- v1.10.0{#BasicAuthForProxy}

```go
func BasicAuthForProxy(accounts Accounts, realm string) HandlerFunc
```

BasicAuthForProxy returns a Basic HTTP Proxy-Authorization middleware. If the realm is empty, "Proxy Authorization Required" will be used by default.

##### func [BasicAuthForRealm](https://github.com/gin-gonic/gin/blob/v1.10.0/auth.go#L48) {#BasicAuthForRealm}

```go
func BasicAuthForRealm(accounts Accounts, realm string) HandlerFunc
```

BasicAuthForRealm returns a Basic HTTP Authorization middleware. It takes as arguments a map[string]string where the key is the user name and the value is the password, as well as the name of the Realm. If the realm is empty, "Authorization Required" will be used by default. (see http://tools.ietf.org/html/rfc2617#section-1.2)

##### func [Bind](https://github.com/gin-gonic/gin/blob/v1.10.0/utils.go#L22) {#Bind}

```go
func Bind(val any) HandlerFunc
```

Bind is a helper function for given interface object and returns a Gin middleware.

##### func [CustomRecovery](https://github.com/gin-gonic/gin/blob/v1.10.0/recovery.go#L38) <- v1.7.0{#CustomRecovery}

```go
func CustomRecovery(handle RecoveryFunc) HandlerFunc
```

CustomRecovery returns a middleware that recovers from any panics and calls the provided handle func to handle it.

##### func [CustomRecoveryWithWriter](https://github.com/gin-gonic/gin/blob/v1.10.0/recovery.go#L51) <- v1.7.0{#CustomRecoveryWithWriter}

```go
func CustomRecoveryWithWriter(out io.Writer, handle RecoveryFunc) HandlerFunc
```

CustomRecoveryWithWriter returns a middleware for a given writer that recovers from any panics and calls the provided handle func to handle it.

##### func [ErrorLogger](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L174) {#ErrorLogger}

```go
func ErrorLogger() HandlerFunc
```

ErrorLogger returns a HandlerFunc for any error type.

##### func [ErrorLoggerT](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L179) {#ErrorLoggerT}

```go
func ErrorLoggerT(typ ErrorType) HandlerFunc
```

ErrorLoggerT returns a HandlerFunc for a given error type.

##### func [Logger](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L191) {#Logger}

```go
func Logger() HandlerFunc
```

Logger instances a Logger middleware that will write the logs to gin.DefaultWriter. By default, gin.DefaultWriter = os.Stdout.

##### func [LoggerWithConfig](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L212) <- v1.4.0{#LoggerWithConfig}

```go
func LoggerWithConfig(conf LoggerConfig) HandlerFunc
```

LoggerWithConfig instance a Logger middleware with config.

##### func [LoggerWithFormatter](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L196) <- v1.4.0{#LoggerWithFormatter}

```go
func LoggerWithFormatter(f LogFormatter) HandlerFunc
```

LoggerWithFormatter instance a Logger middleware with the specified log format function.

##### func [LoggerWithWriter](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L204) {#LoggerWithWriter}

```go
func LoggerWithWriter(out io.Writer, notlogged ...string) HandlerFunc
```

LoggerWithWriter instance a Logger middleware with the specified writer buffer. Example: os.Stdout, a file opened in write mode, a socket...

##### func [Recovery](https://github.com/gin-gonic/gin/blob/v1.10.0/recovery.go#L33) {#Recovery}

```go
func Recovery() HandlerFunc
```

Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.

##### func [RecoveryWithWriter](https://github.com/gin-gonic/gin/blob/v1.10.0/recovery.go#L43) {#RecoveryWithWriter}

```go
func RecoveryWithWriter(out io.Writer, recovery ...RecoveryFunc) HandlerFunc
```

RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.

##### func [WrapF](https://github.com/gin-gonic/gin/blob/v1.10.0/utils.go#L40) {#WrapF}

```go
func WrapF(f http.HandlerFunc) HandlerFunc
```

WrapF is a helper function for wrapping http.HandlerFunc and returns a Gin middleware.

##### func [WrapH](https://github.com/gin-gonic/gin/blob/v1.10.0/utils.go#L47) {#WrapH}

```go
func WrapH(h http.Handler) HandlerFunc
```

WrapH is a helper function for wrapping http.Handler and returns a Gin middleware.

#### type [HandlersChain](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L54) {#HandlersChain}

```go
type HandlersChain []HandlerFunc
```

HandlersChain defines a HandlerFunc slice.

##### func (HandlersChain) [Last](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L57) {#HandlersChain.Last}

```go
func (c HandlersChain) Last() HandlerFunc
```

Last returns the last handler in the chain. i.e. the last handler is the main one.

#### type [IRouter](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L27) {#IRouter}

```go
type IRouter interface {
    IRoutes
    Group(string, ...HandlerFunc) *RouterGroup

}
```

IRouter defines all router handle interface includes single and group router.

#### type [IRoutes](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L33) {#IRoutes}

```go
type IRoutes interface {
    Use(...HandlerFunc) IRoutes

    Handle(string, string, ...HandlerFunc) IRoutes
    Any(string, ...HandlerFunc) IRoutes
    GET(string, ...HandlerFunc) IRoutes
    POST(string, ...HandlerFunc) IRoutes
    DELETE(string, ...HandlerFunc) IRoutes
    PATCH(string, ...HandlerFunc) IRoutes
    PUT(string, ...HandlerFunc) IRoutes
    OPTIONS(string, ...HandlerFunc) IRoutes
    HEAD(string, ...HandlerFunc) IRoutes
    Match([]string, string, ...HandlerFunc) IRoutes

    StaticFile(string, string) IRoutes
    StaticFileFS(string, string, http.FileSystem) IRoutes
    Static(string, string) IRoutes
    StaticFS(string, http.FileSystem) IRoutes

}
```

IRoutes defines all router handle interface.

#### type [LogFormatter](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L60) <- v1.4.0{#LogFormatter}

```go
type LogFormatter func(params LogFormatterParams) string
```

LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter

#### type [LogFormatterParams](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L63) <- v1.4.0{#LogFormatterParams}

```go
type LogFormatterParams struct {
    Request *http.Request

    // TimeStamp shows the time after the server returns a response.
    TimeStamp time.Time
    // StatusCode is HTTP response code.
    StatusCode int
    // Latency is how much time the server cost to process a certain request.
    Latency time.Duration
    // ClientIP equals Context's ClientIP method.
    ClientIP string
    // Method is the HTTP method given to the request.
    Method string
    // Path is a path the client requests.
    Path string
    // ErrorMessage is set if error has occurred in processing the request.
    ErrorMessage string

    // BodySize is the size of the Response Body
    BodySize int
    // Keys are the keys set on the request's context.
    Keys map[string]any
    // contains filtered or unexported fields

}
```

LogFormatterParams is the structure any formatter will be handed when time to log comes

##### func (*LogFormatterParams) [IsOutputColor](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L136) <- v1.4.0{#LogFormatterParams.IsOutputColor}

```go
func (p *LogFormatterParams) IsOutputColor() bool
```

IsOutputColor indicates whether can colors be outputted to the log.

##### func (*LogFormatterParams) [MethodColor](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L107) <- v1.4.0{#LogFormatterParams.MethodColor}

```go
func (p *LogFormatterParams) MethodColor() string
```

MethodColor is the ANSI color for appropriately logging http method to a terminal.

##### func (*LogFormatterParams) [ResetColor](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L131) <- v1.4.0{#LogFormatterParams.ResetColor}

```go
func (p *LogFormatterParams) ResetColor() string
```

ResetColor resets all escape attributes.

##### func (*LogFormatterParams) [StatusCodeColor](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L89) <- v1.4.0{#LogFormatterParams.StatusCodeColor}

```go
func (p *LogFormatterParams) StatusCodeColor() string
```

StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.

#### type [LoggerConfig](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L39) <- v1.4.0{#LoggerConfig}

```go
type LoggerConfig struct {
    // Optional. Default value is gin.defaultLogFormatter
    Formatter LogFormatter

    // Output is a writer where logs are written.
    // Optional. Default value is gin.DefaultWriter.
    Output io.Writer

    // SkipPaths is an url path array which logs are not written.
    // Optional.
    SkipPaths []string

    // Skip is a Skipper that indicates which logs should not be written.
    // Optional.
    Skip Skipper

}
```

LoggerConfig defines the config for Logger middleware.

#### type [Negotiate](https://github.com/gin-gonic/gin/blob/v1.10.0/context.go#L1138) {#Negotiate}

```go
type Negotiate struct {
    Offered  []string
    HTMLName string
    HTMLData any
    JSONData any
    XMLData  any
    YAMLData any
    Data     any
    TOMLData any

}
```

Negotiate contains all negotiations data.

#### type [OptionFunc](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L51) <- v1.10.0{#OptionFunc}

```go
type OptionFunc func(*Engine)
```

OptionFunc defines the function to change the default configuration

#### type [Param](https://github.com/gin-gonic/gin/blob/v1.10.0/tree.go#L24) {#Param}

```go
type Param struct {
    Key   string
    Value string

}
```

Param is a single URL parameter, consisting of a key and a value.

#### type [Params](https://github.com/gin-gonic/gin/blob/v1.10.0/tree.go#L32) {#Params}

```go
type Params []Param
```

Params is a Param-slice, as returned by the router. The slice is ordered, the first URL parameter is also the first slice value. It is therefore safe to read values by the index.

##### func (Params) [ByName](https://github.com/gin-gonic/gin/blob/v1.10.0/tree.go#L47) {#Params.ByName}

```go
func (ps Params) ByName(name string) (va string)
```

ByName returns the value of the first Param which key matches the given name. If no matching Param is found, an empty string is returned.

##### func (Params) [Get](https://github.com/gin-gonic/gin/blob/v1.10.0/tree.go#L36) {#Params.Get}

```go
func (ps Params) Get(name string) (string, bool)
```

Get returns the value of the first Param which key matches the given name and a boolean true. If no matching Param is found, an empty string is returned and a boolean false .

#### type [RecoveryFunc](https://github.com/gin-gonic/gin/blob/v1.10.0/recovery.go#L30) <- v1.7.0{#RecoveryFunc}

```go
type RecoveryFunc func(c *Context, err any)
```

RecoveryFunc defines the function passable to CustomRecovery.

#### type [ResponseWriter](https://github.com/gin-gonic/gin/blob/v1.10.0/response_writer.go#L20) {#ResponseWriter}

```go
type ResponseWriter interface {
    http.ResponseWriter
    http.Hijacker
    http.Flusher
    http.CloseNotifier

    // Status returns the HTTP response status code of the current request.
    Status() int

    // Size returns the number of bytes already written into the response http body.
    // See Written()
    Size() int

    // WriteString writes the string into the response body.
    WriteString(string) (int, error)

    // Written returns true if the response body was already written.
    Written() bool

    // WriteHeaderNow forces to write the http header (status code + headers).
    WriteHeaderNow()

    // Pusher get the http.Pusher for server push
    Pusher() http.Pusher

}
```

ResponseWriter ...

#### type [RouteInfo](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L65) {#RouteInfo}

```go
type RouteInfo struct {
    Method      string
    Path        string
    Handler     string
    HandlerFunc HandlerFunc

}
```

RouteInfo represents a request route's specification which contains method and path and its handler.

#### type [RouterGroup](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L55) {#RouterGroup}

```go
type RouterGroup struct {
    Handlers HandlersChain
    // contains filtered or unexported fields

}
```

RouterGroup is used internally to configure router, a RouterGroup is associated with a prefix and an array of handlers (middleware).

##### func (*RouterGroup) [Any](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L147) {#RouterGroup.Any}

```go
func (group *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) IRoutes
```

Any registers a route that matches all the HTTP methods. GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.

##### func (*RouterGroup) [BasePath](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L82) {#RouterGroup.BasePath}

```go
func (group *RouterGroup) BasePath() string
```

BasePath returns the base path of router group. For example, if v := router.Group("/rest/n/v1/api"), v.BasePath() is "/rest/n/v1/api".

##### func (*RouterGroup) [DELETE](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L121) {#RouterGroup.DELETE}

```go
func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes
```

DELETE is a shortcut for router.Handle("DELETE", path, handlers).

##### func (*RouterGroup) [GET](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L116) {#RouterGroup.GET}

```go
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes
```

GET is a shortcut for router.Handle("GET", path, handlers).

##### func (*RouterGroup) [Group](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L72) {#RouterGroup.Group}

```go
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup
```

Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix. For example, all the routes that use a common middleware for authorization could be grouped.

##### func (*RouterGroup) [HEAD](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L141) {#RouterGroup.HEAD}

```go
func (group *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes
```

HEAD is a shortcut for router.Handle("HEAD", path, handlers).

##### func (*RouterGroup) [Handle](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L103) {#RouterGroup.Handle}

```go
func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes
```

Handle registers a new request handle and middleware with the given path and method. The last handler should be the real handler, the other ones should be middleware that can and should be shared among different routes. See the example code in GitHub.

For GET, POST, PUT, PATCH and DELETE requests the respective shortcut functions can be used.

This function is intended for bulk loading and to allow the usage of less frequently used, non-standardized or custom methods (e.g. for internal communication with a proxy).

##### func (*RouterGroup) [Match](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L156) <- v1.9.0{#RouterGroup.Match}

```go
func (group *RouterGroup) Match(methods []string, relativePath string, handlers ...HandlerFunc) IRoutes
```

Match registers a route that matches the specified methods that you declared.

##### func (*RouterGroup) [OPTIONS](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L136) {#RouterGroup.OPTIONS}

```go
func (group *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes
```

OPTIONS is a shortcut for router.Handle("OPTIONS", path, handlers).

##### func (*RouterGroup) [PATCH](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L126) {#RouterGroup.PATCH}

```go
func (group *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes
```

PATCH is a shortcut for router.Handle("PATCH", path, handlers).

##### func (*RouterGroup) [POST](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L111) {#RouterGroup.POST}

```go
func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes
```

POST is a shortcut for router.Handle("POST", path, handlers).

##### func (*RouterGroup) [PUT](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L131) {#RouterGroup.PUT}

```go
func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes
```

PUT is a shortcut for router.Handle("PUT", path, handlers).

##### func (*RouterGroup) [Static](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L197) {#RouterGroup.Static}

```go
func (group *RouterGroup) Static(relativePath, root string) IRoutes
```

Static serves files from the given file system root. Internally a http.FileServer is used, therefore http.NotFound is used instead of the Router's NotFound handler. To use the operating system's file system implementation, use :

```
router.Static("/static", "/var/www")
```

##### func (*RouterGroup) [StaticFS](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L203) {#RouterGroup.StaticFS}

```go
func (group *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) IRoutes
```

StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead. Gin by default uses: gin.Dir()

##### func (*RouterGroup) [StaticFile](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L166) {#RouterGroup.StaticFile}

```go
func (group *RouterGroup) StaticFile(relativePath, filepath string) IRoutes
```

StaticFile registers a single route in order to serve a single file of the local filesystem. router.StaticFile("favicon.ico", "./resources/favicon.ico")

##### func (*RouterGroup) [StaticFileFS](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L175) <- v1.8.0{#RouterGroup.StaticFileFS}

```go
func (group *RouterGroup) StaticFileFS(relativePath, filepath string, fs http.FileSystem) IRoutes
```

StaticFileFS works just like `StaticFile` but a custom `http.FileSystem` can be used instead.. router.StaticFileFS("favicon.ico", "./resources/favicon.ico", Dir{".", false}) Gin by default uses: gin.Dir()

##### func (*RouterGroup) [Use](https://github.com/gin-gonic/gin/blob/v1.10.0/routergroup.go#L65) {#RouterGroup.Use}

```go
func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes
```

Use adds middleware to the group, see example code in GitHub.

#### type [RoutesInfo](https://github.com/gin-gonic/gin/blob/v1.10.0/gin.go#L73) {#RoutesInfo}

```go
type RoutesInfo []RouteInfo
```

RoutesInfo defines a RouteInfo slice.

#### type [Skipper](https://github.com/gin-gonic/gin/blob/v1.10.0/logger.go#L57) <- v1.10.0{#Skipper}

```go
type Skipper func(c *Context) bool
```

Skipper is a function to skip logs based on provided Context