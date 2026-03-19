//go:build !js

package streams

type readableStream struct {
	controller  ReadableStreamController
	detached    bool
	disturbed   bool
	reader      ReadableStreamReader
	state       readableStreamState
	storedError error
}

type readableStreamState string

const (
	readableStreamStateReadable readableStreamState = "readable"
	readableStreamStateClosed   readableStreamState = "closed"
	readableStreamStateErrored  readableStreamState = "errored"
)

func newReadableStream(underylingSource any, strategy *QueuingStrategy) ReadableStream {
	// The new ReadableStream(underlyingSource, strategy) constructor steps are:

	//  1. If underlyingSource is missing, set it to null.
	//  2. Let underlyingSourceDict be underlyingSource, converted to an IDL value of type UnderlyingSource.
	//
	// Note: We cannot declare the underlyingSource argument as having the
	// UnderlyingSource type directly, because doing so would lose the reference
	// to the original object. We need to retain the object so we can invoke the
	// various methods on it.

	// 3. Perform ! InitializeReadableStream(this).
	this := &readableStream{}
	initializeReadableStream(this)

	// 4. If underlyingSourceDict["type"] is "bytes":
	if typer, ok := underylingSource.(interface{ Type() ReadableStreamType }); ok && typer.Type() == ReadableStreamTypeBytes {
		// 1. If strategy["size"] exists, throw a RangeError exception.
		if strategy.Size != nil {
			
	}

}

func initializeReadableStream(stream *readableStream) {
	// InitializeReadableStream(stream) performs the following steps:

	// 1. Set stream.[[state]] to "readable".
	stream.state = readableStreamStateReadable

	// 2. Set stream.[[reader]] and stream.[[storedError]] to undefined.
	stream.reader = nil
	stream.storedError = nil

	// 3. Set stream.[[disturbed]] to false.
	stream.disturbed = false
}
