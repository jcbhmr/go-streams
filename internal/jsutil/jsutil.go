//go:build js

package jsutil

import (
	"errors"
	"fmt"
	"iter"
	"sync"
	"syscall/js"
)

var jsPromise = sync.OnceValue(func() js.Value {
	return js.Global().Get("Promise")
})
var jsObject = sync.OnceValue(func() js.Value {
	return js.Global().Get("Object")
})
var jsReflect = sync.OnceValue(func() js.Value {
	return js.Global().Get("Reflect")
})
var jsSymbol = sync.OnceValue(func() js.Value {
	return js.Global().Get("Symbol")
})

func Await2(value js.Value) (result js.Value, err error) {
	valueType := value.Type()
	if !(valueType == js.TypeObject || valueType == js.TypeFunction) {
		return value, nil
	}
	jsThen := value.Get("then")
	if jsThen.Type() != js.TypeFunction {
		return value, nil
	}
	resolveChan := make(chan js.Value, 1)
	rejectChan := make(chan js.Value, 1)
	jsResolve := js.FuncOf(func(this js.Value, args []js.Value) any {
		resolveChan <- args[0]
		return nil
	})
	jsReject := js.FuncOf(func(this js.Value, args []js.Value) any {
		rejectChan <- args[0]
		return nil
	})
	jsReflect().Call("apply", jsThen, value, []any{jsResolve, jsReject})
	select {
	case resultValue := <-resolveChan:
		close(resolveChan)
		close(rejectChan)
		return resultValue, nil
	case rejectValue := <-rejectChan:
		close(resolveChan)
		close(rejectChan)
		return js.Value{}, js.Error{Value: rejectValue}
	}
}

func Await(value js.Value) js.Value {
	value, err := Await2(value)
	if err != nil {
		panic(err)
	}
	return value
}

func AsString(value js.Value) string {
	valueType := value.Type()
	if valueType == js.TypeString {
		return value.String()
	} else {
		var err error = &js.ValueError{Method: "jsutil.AsString", Type: valueType}
		panic(err)
	}
}

type JSValuer interface {
	JSValue() js.Value
}

func BetterValueOf(value any) js.Value {
	switch v := value.(type) {
	case JSValuer:
		return v.JSValue()
	case *js.Error:
		return v.Value
	default:
		return js.ValueOf(v)
	}
}

func Lift(value js.Value) any {
	switch value.Type() {
	case js.TypeUndefined, js.TypeNull:
		return nil
	case js.TypeBoolean:
		return value.Bool()
	case js.TypeNumber:
		return value.Float()
	case js.TypeString:
		return value.String()
	case js.TypeSymbol:
		return value
	case js.TypeObject:
		jsThen := value.Get("then")
		if jsThen.Type() == js.TypeFunction {
			return Lift(Await(value))
		}
		jsProto := jsObject().Call("getPrototypeOf", value)
		if jsProto.Equal(jsObject().Get("prototype")) || jsProto.Equal(js.Null()) {
			goMap := make(map[string]any)
			entries := jsObject().Call("entries", value)
			entriesLen := entries.Length()
			for i := range entriesLen {
				entry := entries.Index(i)
				key := AsString(entry.Index(0))
				val := Lift(entry.Index(1))
				goMap[key] = val
			}
			return goMap
		}
		return value
	case js.TypeFunction:
		return func(args ...any) any {
			for i := range args {
				args[i] = BetterValueOf(args[i])
			}
			return Lift(value.Invoke(args...))
		}
	default:
		panic(fmt.Sprintf("unexpected js.Value type: %s", value.Type()))
	}
}

func AsyncIterableToSeq(asyncIterable js.Value) iter.Seq[js.Value] {
	return func(yield func(js.Value) bool) {
		jsAsyncIteratorMethod := jsReflect().Call("get", asyncIterable, jsSymbol().Get("asyncIterator"))
		jsAsyncIterator := jsReflect().Call("apply", jsAsyncIteratorMethod, asyncIterable, []any{})
		for jsValue := range AsyncIteratorToSeq(jsAsyncIterator) {
			if !yield(jsValue) {
				return
			}
		}
	}
}

func AsyncIteratorToSeq(asyncIterator js.Value) iter.Seq[js.Value] {
	return func(yield func(js.Value) bool) {
		defer func() {
			r := recover()
			if r != nil {
				jsThrow := asyncIterator.Get("throw")
				if jsThrow.Type() == js.TypeFunction {
					args := []any{}
					if err, ok := r.(error); ok {
						var jsError *js.Error
						if errors.As(err, &jsError) {
							args = append(args, jsError.Value)
						}
					}
					Await(jsReflect().Call("apply", jsThrow, asyncIterator, args))
				}
			}
			jsReturn := asyncIterator.Get("return")
			if jsReturn.Type() == js.TypeFunction {
				Await(jsReflect().Call("apply", jsReturn, asyncIterator, []any{}))
			}
			if r != nil {
				panic(r)
			}
		}()
		for {
			jsIteratorResult := Await(asyncIterator.Call("next"))
			done := jsIteratorResult.Get("done").Bool()
			jsValue := jsIteratorResult.Get("value")
			if done {
				return
			}
			if !yield(jsValue) {
				return
			}
		}
	}
}
