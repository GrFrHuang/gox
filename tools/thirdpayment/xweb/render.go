package xweb

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/cxuhua/xweb/logging"
	"github.com/cxuhua/xweb/martini"
	"github.com/oxtoacart/bpool"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	MultipartFormData = "multipart/form-data"
	ContentURLEncoded = "application/x-www-form-urlencoded"
	ContentExpires    = "Expires"
	ContentType       = "Content-Type"
	ContentLength     = "Content-Length"
	ProtobufType      = "application/protobuf"
	ContentBinary     = "application/octet-stream"
	ContentText       = "text/plain"
	ContentJSON       = "application/json"
	ContentJPEG       = "image/jpeg"
	ContentPNG        = "image/png"
	ContentHTML       = "text/html"
	ContentXHTML      = "application/xhtml+xml"
	ContentXML        = "text/xml"
	defaultCharset    = "UTF-8"
)

// Provides a temporary buffer to execute templates into and catch errors.
var bufpool *bpool.BufferPool

// Included helper functions for use when rendering html
var helperFuncs = template.FuncMap{
	"yield": func() (string, error) {
		return "", fmt.Errorf("yield called with no layout defined")
	},
	"current": func() (string, error) {
		return "", nil
	},
}

// Render is a service that can be injected into a Martini handler. Render provides functions for easily writing JSON and
// HTML templates out to a http Response.
type Render interface {
	// JSON writes the given status and JSON serialized version of the given value to the http.ResponseWriter.
	JSON(status int, v interface{})
	// HTML renders a html template specified by the name and writes the result and given status to the http.ResponseWriter.
	HTML(status int, name string, v interface{}, htmlOpt ...HTMLOptions)
	// HTML renders a html template specified by the template content and writes the result and given status to the http.ResponseWriter.
	TEMP(status int, template string, v interface{})
	// XML writes the given status and XML serialized version of the given value to the http.ResponseWriter.
	XML(status int, v interface{})
	// Data writes the raw byte array to the http.ResponseWriter.
	Data(status int, v []byte)
	// File write
	File(name string, mod time.Time, file IHttpFile)
	// Text writes the given status and plain text to the http.ResponseWriter.
	Text(status int, v string)
	// Error is a convenience function that writes an http status to the http.ResponseWriter.
	Error(status int)
	// Status is an alias for Error (writes an http status to the http.ResponseWriter)
	Status(status int)
	// Redirect is a convienience function that sends an HTTP redirect. If status is omitted, uses 302 (Found)
	Redirect(location string, status ...int)
	// Template returns the internal *template.Template used to render the HTML
	Template() *template.Template
	// Header exposes the header struct from http.ResponseWriter.
	Header() http.Header
	// SetCookie
	SetCookie(cookie *http.Cookie)
}

// Delims represents a set of Left and Right delimiters for HTML template rendering
type Delims struct {
	// Left delimiter, defaults to {{
	Left string
	// Right delimiter, defaults to }}
	Right string
}

// Options is a struct for specifying configuration options for the render.Renderer middleware
type RenderOptions struct {
	// Directory to load templates. Default is "templates"
	Directory string
	// Layout template name. Will not render a layout if "". Defaults to "".
	Layout string
	// Extensions to parse template files from. Defaults to [".tmpl"]
	Extensions []string
	// Funcs is a slice of FuncMaps to apply to the template upon compilation. This is useful for helper functions. Defaults to [].
	Funcs []template.FuncMap
	// Delims sets the action delimiters to the specified strings in the Delims struct.
	Delims Delims
	// Appends the given charset to the Content-Type header. Default is "UTF-8".
	Charset string
	// Outputs human readable JSON
	IndentJSON bool
	// Outputs human readable XML
	IndentXML bool
	// Prefixes the JSON output with the given bytes.
	PrefixJSON []byte
	// Prefixes the XML output with the given bytes.
	PrefixXML []byte
	// Allows changing of output to XHTML instead of HTML. Default is "text/html"
	HTMLContentType string
}

// HTMLOptions is a struct for overriding some rendering Options for specific HTML call
type HTMLOptions struct {
	// Layout template name. Overrides Options.Layout.
	Layout string
}

// Renderer is a Middleware that maps a render.Render service into the Martini handler chain. An single variadic render.Options
// struct can be optionally provided to configure HTML rendering. The default directory for templates is "templates" and the default
// file extension is ".tmpl".
//
// If MARTINI_ENV is set to "" or "development" then templates will be recompiled on every request. For more performance, set the
// MARTINI_ENV environment variable to "production"
func Renderer(options ...RenderOptions) martini.Handler {
	opt := prepareOptions(options)
	cs := prepareCharset(opt.Charset)
	t := compile(opt)
	bufpool = bpool.NewBufferPool(64)
	return func(res http.ResponseWriter, req *http.Request, c martini.Context, log *logging.Logger) {
		var tc *template.Template
		if martini.Env == martini.Dev {
			// recompile for easy development
			tc = compile(opt)
		} else {
			// use a clone of the initial template
			tc, _ = t.Clone()
		}
		getValueFunc := func(name string) (interface{}, error) {
			if name == "" {
				return nil, errors.New("name args must set")
			}
			typ, ok := c.GetType(name)
			if !ok {
				return nil, errors.New(name + " type not map join context")
			}
			vv := c.Get(typ)
			if !vv.IsValid() {
				return nil, errors.New(name + " type value not valid")
			}
			return vv.Interface(), nil
		}
		tc.Funcs(template.FuncMap{
			"import": func(name string, kv ...string) (template.HTML, error) {
				vv, err := getValueFunc(name)
				if err != nil {
					return "", err
				}
				buf := bufpool.Get()
				err = tc.ExecuteTemplate(buf, name, vv)
				return template.HTML(buf.String()), err
			},
			"value": func(name string) interface{} {
				vv, err := getValueFunc(name)
				if err != nil {
					return err
				}
				return vv
			},
		})
		c.MapTo(&renderer{res, req, tc, opt, cs}, (*Render)(nil))
	}
}

func prepareCharset(charset string) string {
	if len(charset) != 0 {
		return "; charset=" + charset
	}
	return "; charset=" + defaultCharset
}

func prepareOptions(options []RenderOptions) RenderOptions {
	var opt RenderOptions
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.Directory) == 0 {
		opt.Directory = "templates"
	}
	if len(opt.Extensions) == 0 {
		opt.Extensions = []string{".tmpl"}
	}
	if len(opt.HTMLContentType) == 0 {
		opt.HTMLContentType = ContentHTML
	}

	return opt
}

func compile(options RenderOptions) *template.Template {
	dir := options.Directory
	t := template.New(dir)
	t.Delims(options.Delims.Left, options.Delims.Right)
	// parse an initial template in case we don't have any
	template.Must(t.Parse("Martini"))
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		r, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		ext := getExt(r)
		for _, extension := range options.Extensions {
			if ext == extension {

				buf, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}

				name := (r[0 : len(r)-len(ext)])
				tmpl := t.New(filepath.ToSlash(name))
				//for skip error
				tmpl.Funcs(template.FuncMap{
					"import": func() interface{} { return nil },
					"value":  func() interface{} { return nil },
				})
				// add our funcmaps
				for _, funcs := range options.Funcs {
					tmpl.Funcs(funcs)
				}
				// Bomb out if parse fails. We don't want any silent server starts.
				template.Must(tmpl.Funcs(helperFuncs).Parse(string(buf)))
				break
			}
		}
		return nil
	})
	return t
}

func getExt(s string) string {
	if strings.Index(s, ".") == -1 {
		return ""
	}
	return "." + strings.Join(strings.Split(s, ".")[1:], ".")
}

type renderer struct {
	http.ResponseWriter
	req             *http.Request
	t               *template.Template
	opt             RenderOptions
	compiledCharset string
}

func (r *renderer) SetCookie(cookie *http.Cookie) {
	http.SetCookie(r.ResponseWriter, cookie)
}

func (r *renderer) File(name string, mod time.Time, file IHttpFile) {
	http.ServeContent(r, r.req, name, mod, file)
}

func (r *renderer) JSON(status int, v interface{}) {
	var result []byte
	var err error
	if r.opt.IndentJSON {
		result, err = json.MarshalIndent(v, "", "  ")
	} else {
		result, err = json.Marshal(v)
	}
	if err != nil {
		http.Error(r, err.Error(), 500)
		return
	}
	// json rendered fine, write out the result
	r.Header().Set(ContentType, ContentJSON+r.compiledCharset)
	r.WriteHeader(status)
	if len(r.opt.PrefixJSON) > 0 {
		r.Write(r.opt.PrefixJSON)
	}
	r.Write(result)
}

func (r *renderer) TEMP(status int, template string, data interface{}) {
	buf := &bytes.Buffer{}
	tmp, err := r.Template().Parse(template)
	if err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmp.Execute(buf, data); err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Header().Set(ContentType, r.opt.HTMLContentType+r.compiledCharset)
	r.WriteHeader(status)
	io.Copy(r, buf)
	bufpool.Put(buf)
}

func (r *renderer) HTML(status int, name string, binding interface{}, htmlOpt ...HTMLOptions) {
	opt := r.prepareHTMLOptions(htmlOpt)
	// assign a layout if there is one
	if len(opt.Layout) > 0 {
		r.addYield(name, binding)
		name = opt.Layout
	}
	buf, err := r.execute(name, binding)
	if err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
		return
	}
	// template rendered fine, write out the result
	r.Header().Set(ContentType, r.opt.HTMLContentType+r.compiledCharset)
	r.WriteHeader(status)
	io.Copy(r, buf)
	bufpool.Put(buf)
}

func (r *renderer) XML(status int, v interface{}) {
	var result []byte
	var err error
	if r.opt.IndentXML {
		result, err = xml.MarshalIndent(v, "", "  ")
	} else {
		result, err = xml.Marshal(v)
	}
	if err != nil {
		http.Error(r, err.Error(), 500)
		return
	}
	// XML rendered fine, write out the result
	r.Header().Set(ContentType, ContentXML+r.compiledCharset)
	r.WriteHeader(status)
	if len(r.opt.PrefixXML) > 0 {
		r.Write(r.opt.PrefixXML)
	}
	r.Write(result)
}

func (r *renderer) Data(status int, v []byte) {
	if r.Header().Get(ContentType) == "" {
		r.Header().Set(ContentType, ContentBinary)
	}
	r.WriteHeader(status)
	r.Write(v)
}

func (r *renderer) Text(status int, v string) {
	if r.Header().Get(ContentType) == "" {
		r.Header().Set(ContentType, ContentText+r.compiledCharset)
	}
	r.WriteHeader(status)
	r.Write([]byte(v))
}

// Error writes the given HTTP status to the current ResponseWriter
func (r *renderer) Error(status int) {
	r.WriteHeader(status)
}

func (r *renderer) Status(status int) {
	r.WriteHeader(status)
}

func (r *renderer) Redirect(location string, status ...int) {
	code := http.StatusFound
	if len(status) == 1 {
		code = status[0]
	}
	http.Redirect(r, r.req, location, code)
}

func (r *renderer) Template() *template.Template {
	return r.t
}

func (r *renderer) execute(name string, binding interface{}) (*bytes.Buffer, error) {
	buf := bufpool.Get()
	return buf, r.t.ExecuteTemplate(buf, name, binding)
}

func (r *renderer) addYield(name string, binding interface{}) {
	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf, err := r.execute(name, binding)
			// return safe html here since we are rendering our own template
			return template.HTML(buf.String()), err
		},
		"current": func() (string, error) {
			return name, nil
		},
	}
	r.t.Funcs(funcs)
}

func (r *renderer) prepareHTMLOptions(htmlOpt []HTMLOptions) HTMLOptions {
	if len(htmlOpt) > 0 {
		return htmlOpt[0]
	}

	return HTMLOptions{
		Layout: r.opt.Layout,
	}
}