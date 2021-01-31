package HttpRequest

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

type Response struct {
	time int64
	url  string
	R    *http.Response
	body []byte
	Req  *http.Request
}

func (r *Response) Response() *http.Response {
	if r != nil {
		return r.R
	}
	return nil
}

func (r *Response) StatusCode() int {
	if r.R == nil {
		return 0
	}
	return r.R.StatusCode
}

func (r *Response) Time() int64 {
	if r != nil {
		return r.time
	}
	return r.time
}

func (r *Response) ContentLength() int64 {
	b, err := r.Body()
	if err != nil {
		return 0
	}
	return int64(len(b))
}

func (r *Response) Url() string {
	if r != nil {
		return r.url
	}
	return ""
}

func (r *Response) Headers() http.Header {
	if r != nil {
		return r.R.Header
	}
	return nil
}

func (r *Response) Cookies() []*http.Cookie {
	if r != nil {
		return r.R.Cookies()
	}
	return []*http.Cookie{}
}

func (r *Response) Body() ([]byte, error) {
	if r == nil {
		return []byte{}, errors.New("HttpRequest.Response is nil.")
	}

	defer r.R.Body.Close()

	if len(r.body) > 0 {
		return r.body, nil
	}

	if r.R == nil || r.R.Body == nil {
		return nil, errors.New("response or body is nil")
	}

	b, err := ioutil.ReadAll(r.R.Body)
	if err != nil {
		return nil, err
	}
	r.body = b

	return b, nil
}

func (r *Response) Content() (string, error) {
	b, err := r.Body()
	if err != nil {
		return "", nil
	}
	return string(b), nil
}

func (r *Response) Text() string {
	b, err := r.Body()
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (r *Response) Json(v interface{}) error {
	return r.Unmarshal(v)
}

func (r *Response) Unmarshal(v interface{}) error {
	b, err := r.Body()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	return nil
}

func (r *Response) Close() error {
	if r != nil {
		return r.R.Body.Close()
	}
	return nil
}

func (r *Response) Export() (string, error) {
	b, err := r.Body()
	if err != nil {
		return "", err
	}

	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return "", errors.New("illegal json: " + err.Error())
	}

	return Export(i), nil
}

// 转换成gjson.Result,方便解析JSON响应
func (r *Response) Result() gjson.Result {
	text, err := r.Content()
	if err != nil {
		text = "{}"
	}
	return gjson.Parse(text)
}
