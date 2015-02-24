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
		Log(),
		Sink(),
		Latch(),
		Gate(),
		Identity(),
		Append(),
		Tail(),
		Head(),
		Pusher(),

		// monads
		Exp(),
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

		// stateful
		First(),

		// network IO
		GET(),

		//assertions
		IsBoolean(),
		IsNumber(),
		IsString(),
		IsArray(),
		IsObject(),
		IsError(),
		ToString(),
		ToNumber(),

		//logic
		And(),
		Or(),
		Not(),
	}

	library := make(map[string]Spec)

	for _, s := range b {
		library[s.Name] = s
	}

	return library
}
