package jsonmap

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"io"
	"net/http"
	"reflect"
)

type Options struct {
	ErrorHandler martini.Handler
}

func Bind(valueStruct interface{}, ifacePtr ...interface{}) martini.Handler {
	return func(c martini.Context, req *http.Request, opts *Options) {
		e := &Error{}
		c.Map(e)
		defer func() {
			if !e.Empty() && opts.ErrorHandler != nil {
				c.Invoke(opts.ErrorHandler)
			}
		}()

		sink, isPtr := makeSink(valueStruct)

		if req.Body != nil {
			defer req.Body.Close()
			err := json.NewDecoder(req.Body).Decode(sink.Interface())
			if err != nil && err != io.EOF {
				e.Add(err)
				return
			}
		}

		validate(c, sink)

		if !isPtr {
			sink = sink.Elem()
		}
		mapTo(c, sink, ifacePtr)
	}
}

func makeSink(valueStruct interface{}) (reflect.Value, bool) {
	var isPtr bool
	rt := reflect.TypeOf(valueStruct)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		isPtr = true
	}

	return reflect.New(rt), isPtr
}

func validate(c martini.Context, v reflect.Value) {
	method := v.MethodByName("Validate")
	if method.IsValid() {
		c.Invoke(method.Interface())
	}
}

func mapTo(c martini.Context, v reflect.Value, ifacePtr []interface{}) {
	obj := v.Interface()
	c.Map(obj)
	if len(ifacePtr) > 0 {
		c.MapTo(obj, ifacePtr[0])
	}
}
