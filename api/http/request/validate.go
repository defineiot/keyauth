package request

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/json-iterator/go"

	"github.com/defineiot/keyauth/internal/exception"
)

func checkBody(r *http.Request) ([]byte, error) {
	// 检测请求大小
	if r.ContentLength == 0 {
		return nil, exception.NewBadRequest("request body is empty")
	}
	if r.ContentLength > 20971520 {
		return nil, exception.NewBadRequest("the body exceeding the maximum limit, max size 20M")
	}

	// 读取body数据
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return nil, exception.NewBadRequest("read request body error, %s", err)
	}

	return body, nil
}

// CheckObjectBody check json array body
func CheckObjectBody(r *http.Request) (jsoniter.Any, error) {
	body, err := checkBody(r)
	if err != nil {
		return nil, err
	}

	iter := jsoniter.ParseBytes(jsoniter.ConfigDefault, body)
	val := iter.ReadAny()
	if val.ValueType() != jsoniter.ObjectValue {
		return nil, exception.NewBadRequest("body must be an valid json object")
	}

	return val, nil
}

// CheckArrayBody check json object body
func CheckArrayBody(r *http.Request) (*jsoniter.Iterator, error) {
	body, err := checkBody(r)
	if err != nil {
		return nil, err
	}

	data := strings.TrimSpace(string(body))
	dl := len(data)
	if !(data[0] == '[' && data[dl-1] == ']') {
		return nil, exception.NewBadRequest("body must be an valid json array")
	}
	iter := jsoniter.ParseString(jsoniter.ConfigDefault, data)

	return iter, nil
}
