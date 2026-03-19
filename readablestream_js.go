package streams

import (
	"iter"
	"sync"
	"syscall/js"

	"go.jcbhmr.com/streams/internal/jsutil"
)

var jsReadableStream = sync.OnceValue(func() js.Value {
	return js.Global().Get("ReadableStream")
})
var jsReadableStreamBYOBReader = sync.OnceValue(func() js.Value {
	return js.Global().Get("ReadableStreamBYOBReader")
})
var jsReadableStreamDefaultReader = sync.OnceValue(func() js.Value {
	return js.Global().Get("ReadableStreamDefaultReader")
})

type readableStream js.Value
type readableStreamBYOBReader js.Value
type readableStreamDefaultReader js.Value

func newReadableStream() ReadableStream {
	return readableStream(jsReadableStream().New())
}

func readableStreamFrom(asyncIterable any) ReadableStream {
	return readableStream(jsReadableStream().Call("from", asyncIterable))
}

func (r readableStream) Locked() bool {
	return js.Value(r).Call("locked").Bool()
}

func (r readableStream) Cancel(reason any) error {
	jsPromise := js.Value(r).Call("cancel", reason)
	_, err := jsutil.Await2(jsPromise)
	return err
}

func (r readableStream) GetReader(options *ReadableStreamGetReaderOptions) ReadableStreamReader {
	args := []any{}
	if options != nil {
		args = append(args, jsutil.BetterValueOf(options))
	}
	jsReader := js.Value(r).Call("getReader", args...)
	if !jsReadableStreamBYOBReader().IsUndefined() && jsReader.InstanceOf(jsReadableStreamBYOBReader()) {
		return readableStreamBYOBReader(jsReader)
	} else {
		return readableStreamDefaultReader(jsReader)
	}
}

func (r readableStream) PipeThrough(transform *ReadableWritablePair, options *StreamPipeOptions) ReadableStream {
	args := []any{jsutil.BetterValueOf(transform)}
	if options != nil {
		args = append(args, jsutil.BetterValueOf(options))
	}
	jsReadable := js.Value(r).Call("pipeThrough", args...)
	return readableStream(jsReadable)
}

func (r readableStream) PipeTo(destination WritableStream, options *StreamPipeOptions) error {
	args := []any{jsutil.BetterValueOf(destination)}
	if options != nil {
		args = append(args, jsutil.BetterValueOf(options))
	}
	jsPromise := js.Value(r).Call("pipeTo", args...)
	_, err := jsutil.Await2(jsPromise)
	return err
}

func (r readableStream) Tee() []ReadableStream {
	jsReadables := js.Value(r).Call("tee")
	readable1 := readableStream(jsReadables.Index(0))
	readable2 := readableStream(jsReadables.Index(1))
	return []ReadableStream{readable1, readable2}
}

func (r readableStream) All(options *ReadableStreamIteratorOptions) iter.Seq[any] {
	return func(yield func(any) bool) {
		args := []any{}
		if options != nil {
			args = append(args, jsutil.BetterValueOf(options))
		}
		jsAsyncIterator := js.Value(r).Call("values", args...)
		for jsValue := range jsutil.AsyncIteratorToSeq(jsAsyncIterator) {
			if !yield(jsutil.Lift(jsValue)) {
				return
			}
		}
	}
}

func (r readableStream) Values(options *ReadableStreamIteratorOptions) iter.Seq[any] {
	return r.All(options)
}


