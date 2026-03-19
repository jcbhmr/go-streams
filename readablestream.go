package streams

import "iter"

type ReadableStream interface {
	Locked() bool
	Cancel(reason any) error
	GetReader(options *ReadableStreamGetReaderOptions) ReadableStreamReader
	PipeThrough(transform ReadableWritablePair, options *StreamPipeOptions) ReadableStream
	PipeTo(destination WritableStream, options *StreamPipeOptions)
	Tee() []ReadableStream
	All(options *ReadableStreamIteratorOptions) iter.Seq[any]
	Values(options *ReadableStreamIteratorOptions) iter.Seq[any]
	InternalReadableStream()
}

type ReadableStreamReader = any

type ReadableStreamReaderMode string

const (
	ReadableStreamReaderModeBYOB ReadableStreamReaderMode = "byob"
)

type ReadableStreamGetReaderOptions struct {
	Mode ReadableStreamReaderMode
}

type ReadableStreamIteratorOptions struct {
	PreventCancel bool
}

type ReadableWritablePair struct {
	Readable ReadableStream
	Writable WritableStream
}

type StreamPipeOptions struct {
	PreventClose  bool
	PreventAbort  bool
	PreventCancel bool
	// Signal        AbortSignal
}

type UnderlyingSource struct {
	Start                 UnderlyingSourceStartCallback
	Pull                  UnderlyingSourcePullCallback
	Cancel                UnderlyingSourceCancelCallback
	Type                  ReadableStreamType
	AutoAllocateChunkSize uint64
}

type ReadableStreamController = any

type UnderlyingSourceStartCallback = func(controller ReadableStreamController) any
type UnderlyingSourcePullCallback = func(controller ReadableStreamController)
type UnderlyingSourceCancelCallback = func(reason any)

type ReadableStreamType string

const (
	ReadableStreamTypeBytes ReadableStreamType = "bytes"
)

type ReadableStreamGenericReader interface {
	Closed() bool
	Cancel(reason any)
	InternalReadableStreamGenericReader()
}

type ReadableStreamDefaultReader interface {
	Read() ReadableStreamReadResult
	ReleaseLock()
	ReadableStreamGenericReader
}

type ReadableStreamReadResult struct {
	Value any
	Done  bool
}

type ReadableStreamBYOBReader interface {
	Read(view []byte, options *ReadableStreamBYOBReaderReadOptions) ReadableStreamReadResult
	ReleaseLock()
	ReadableStreamGenericReader
}

type ReadableStreamBYOBReaderReadOptions struct {
	Min uint64
}

type ReadableStreamDefaultController interface {
	DesiredSize() *float64
	Close()
	Enqueue(chunk any)
	Error(e any)
}

type ReadableStreamByteStreamController interface {

}