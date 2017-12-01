package request

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/json-iterator/go"
)

func checkBody(r *http.Request) ([]byte, error) {
	// 检测请求大小
	if r.ContentLength == 0 {
		return nil, errors.New("request body is empty")
	}
	if r.ContentLength > 20971520 {
		return nil, errors.New("the body exceeding the maximum limit, max size 20M")
	}

	// 读取body数据
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read request body error, %s", err)
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
		return nil, fmt.Errorf("body must be an valid json object")
	}

	return val, nil
}

// CheckArrayBody check json object body
func CheckArrayBody(r *http.Request) (jsoniter.Any, error) {
	body, err := checkBody(r)
	if err != nil {
		return nil, err
	}

	iter := jsoniter.ParseBytes(jsoniter.ConfigDefault, body)
	val := iter.ReadAny()
	if val.ValueType() != jsoniter.ArrayValue {
		return nil, fmt.Errorf("body must be an valid json array")
	}

	return val, nil
}
