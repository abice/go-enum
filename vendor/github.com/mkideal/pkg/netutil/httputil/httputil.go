package httputil

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/mkideal/pkg/option"
)

var (
	ErrUnableToMarshalForm = errors.New("unable to marshal form")
)

const (
	charsetUTF8 = "charset=utf-8"

	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

const (
	HeaderAccept                        = "Accept"
	HeaderAcceptEncoding                = "Accept-Encoding"
	HeaderAllow                         = "Allow"
	HeaderAuthorization                 = "Authorization"
	HeaderContentDisposition            = "Content-Disposition"
	HeaderContentEncoding               = "Content-Encoding"
	HeaderContentLength                 = "Content-Length"
	HeaderContentType                   = "Content-Type"
	HeaderCookie                        = "Cookie"
	HeaderSetCookie                     = "Set-Cookie"
	HeaderIfModifiedSince               = "If-Modified-Since"
	HeaderLastModified                  = "Last-Modified"
	HeaderLocation                      = "Location"
	HeaderUpgrade                       = "Upgrade"
	HeaderVary                          = "Vary"
	HeaderWWWAuthenticate               = "WWW-Authenticate"
	HeaderXForwardedProto               = "X-Forwarded-Proto"
	HeaderXHTTPMethodOverride           = "X-HTTP-Method-Override"
	HeaderXForwardedFor                 = "X-Forwarded-For"
	HeaderXRealIP                       = "X-Real-IP"
	HeaderServer                        = "Server"
	HeaderOrigin                        = "Origin"
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderXCSRFToken              = "X-CSRF-Token"
)

func IP(req *http.Request) string {
	ips := getproxy(req)
	if len(ips) > 0 && ips[0] != "" {
		rip := strings.Split(ips[0], ":")
		if len(rip) > 0 {
			return rip[0]
		}
	}
	ip := strings.Split(req.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func getproxy(req *http.Request) []string {
	if ips := req.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

type Result struct {
	Response   *http.Response
	StatusCode int
	Data       []byte
	Error      error
}

func (result Result) Ok() bool { return result.StatusCode == http.StatusOK }
func (result Result) Status() string {
	if result.Response == nil {
		if result.Error != nil {
			return result.Error.Error()
		}
		return ""
	}
	return result.Response.Status
}

func Get(url string) Result {
	resp, err := http.Get(url)
	return readResultFromResponse(resp, err)
}

func PostForm(url string, values url.Values) Result {
	resp, err := http.PostForm(url, values)
	return readResultFromResponse(resp, err)
}

func readResultFromResponse(resp *http.Response, err error) Result {
	result := Result{
		Response: resp,
		Error:    err,
	}
	if err != nil {
		return result
	}
	result.StatusCode = resp.StatusCode
	defer resp.Body.Close()
	result.Data, result.Error = ioutil.ReadAll(resp.Body)
	return result
}

func Response(w http.ResponseWriter, status int, acceptType string, value interface{}, debug ...bool) error {
	switch {
	case strings.HasPrefix(acceptType, MIMEApplicationXML):
		return XMLResponse(w, status, value, debug...)
	case strings.HasPrefix(acceptType, MIMEApplicationForm):
		return FormResponse(w, status, value)
	default:
		return JSONResponse(w, status, value, debug...)
	}
}

func JSONResponse(w http.ResponseWriter, status int, value interface{}, debug ...bool) error {
	return ResponseWithMarshaler(w, status, MIMEApplicationJSONCharsetUTF8, value, json.Marshal, json.MarshalIndent, debug...)
}

func XMLResponse(w http.ResponseWriter, status int, value interface{}, debug ...bool) error {
	return ResponseWithMarshaler(w, status, MIMEApplicationXMLCharsetUTF8, value, xml.Marshal, xml.MarshalIndent, debug...)
}

func FormResponse(w http.ResponseWriter, status int, value interface{}) error {
	return ResponseWithMarshaler(w, status, MIMEApplicationForm, value, MarshalForm, nil)
}

func TextResponse(w http.ResponseWriter, status int, value string) error {
	w.WriteHeader(status)
	_, err := w.Write([]byte(value))
	return err
}

type MarshalFunc func(interface{}) ([]byte, error)
type MarshalIndentFunc func(interface{}, string, string) ([]byte, error)

type FormMarshaler interface {
	MarshalForm() ([]byte, error)
}

func MarshalForm(v interface{}) ([]byte, error) {
	if marshaler, ok := v.(FormMarshaler); ok {
		return marshaler.MarshalForm()
	}
	return nil, ErrUnableToMarshalForm
}

func ResponseWithMarshaler(
	w http.ResponseWriter,
	status int,
	contentType string,
	value interface{},
	marshal MarshalFunc,
	marshalIndent MarshalIndentFunc,
	debug ...bool) error {
	var (
		b   []byte
		err error
	)
	if option.Bool(false, debug...) && marshalIndent != nil {
		b, err = marshalIndent(value, "", "  ")
	} else {
		b, err = marshal(value)
	}
	if err != nil {
		return err
	}
	return BlobResponse(w, status, contentType, b)
}

func BlobResponse(w http.ResponseWriter, status int, contentType string, b []byte) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", contentType)
	_, err := w.Write(b)
	return err
}
