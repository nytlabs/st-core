package core

import "log"

type stcoreError struct {
	S string `json:"core"`
}

func (e *stcoreError) Error() string {
	log.Println(e.S)
	return e.S
}

func NewError(s string) *stcoreError {
	return &stcoreError{
		S: s,
	}
}

// Library is the set of all core block Specs
// TODO: should just "Build" a global variable so we don't have to iterate
// over all funcs every time we need the library
func GetLibrary() map[string]Spec {
	b := []Spec{
		// mechanism
		First(),
		Delay(),
		Log(),
		Sink(),
		Latch(),
		Gate(),
		Identity(),
		Timestamp(),

		// object
		Set(),
		Get(),
		Keys(),
		Merge(),
		HasField(),

		// array
		Append(),
		Tail(),
		Head(),
		Init(),
		Last(),
		Len(),
		InArray(),

		// maths
		Exp(),
		Floor(),
		Ceil(),
		Log10(),
		Ln(),
		Sqrt(),
		Sin(),
		Cos(),
		Tan(),
		Addition(),
		Subtraction(),
		Multiplication(),
		Division(),
		Exponentiation(),
		Modulation(),
		GreaterThan(),
		LessThan(),
		EqualTo(),
		NotEqualTo(),

		// random
		UniformRandom(),
		NormalRandom(),
		ZipfRandom(),
		PoissonRandom(),
		ExponentialRandom(),
		BernoulliRandom(),

		// string
		InString(),
		HasPrefix(),
		HasSuffix(),
		StringConcat(),
		StringSplit(),

		// key value
		kvGet(),
		kvSet(),
		kvClear(),
		kvDump(),
		kvDelete(),

		// parsers
		ParseJSON(),

		// NSQ interface
		NSQConsumerConnect(),
		NSQConsumerReceive(),

		// primitive value
		ValueGet(),
		ValueSet(),

		// list
		listGet(),
		listSet(),
		listShift(),
		listAppend(),
		listPop(),
		listDump(),

		// priority queue
		pqPush(),
		pqPop(),
		pqPeek(),
		pqLen(),
		pqClear(),

		// network IO
		HTTPRequest(),

		//assertions
		IsBoolean(),
		IsNumber(),
		IsString(),
		IsArray(),
		IsObject(),
		IsError(),

		// conversion
		ToString(),
		ToNumber(),

		//logic
		And(),
		Or(),
		Not(),

		//string functions
		StringConcat(),
		StringSplit(),

		// websocket
		wsClientConnect(),
		wsClientReceive(),
		wsClientSend(),

		// stdin
		StdinReceive(),
	}

	library := make(map[string]Spec)

	for _, s := range b {
		library[s.Name] = s
	}

	return library
}
