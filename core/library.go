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
		// mechanisms
		Delay(),
		Set(),
		Get(),
		Keys(),
		Log(),
		Sink(),
		Latch(),
		Gate(),
		Identity(),
		Append(),
		Tail(),
		Head(),
		Init(),
		Last(),
		Pusher(),
		Merge(),
		Len(),
		Timestamp(),

		// monads
		Exp(),
		Floor(),
		Ceil(),
		Log10(),
		Ln(),
		Sqrt(),
		Sin(),
		Cos(),
		Tan(),

		// dyads
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

		// random sources
		UniformRandom(),
		NormalRandom(),
		ZipfRandom(),
		PoissonRandom(),
		ExponentialRandom(),
		BernoulliRandom(),

		// membership
		InArray(),
		HasField(),
		InString(),
		HasPrefix(),
		HasSuffix(),

		// key value
		kvGet(),
		kvSet(),
		kvClear(),
		kvDump(),
		kvDelete(),

		// parsers
		ParseJSON(),

		// NSQ interface
		NSQReceive(),

		// primitive value
		ValueGet(),
		ValueSet(),

		// server
		FromRequest(),

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

		// stateful
		First(),

		// network IO
		HTTPRequest(),

		// IO
		Write(),
		Close(),
		Flush(),

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
	}

	library := make(map[string]Spec)

	for _, s := range b {
		library[s.Name] = s
	}

	return library
}
