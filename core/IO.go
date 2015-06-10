package core

import (
	"encoding/json"
	"io"
	"net/http"
)

func Write() Spec {
	return Spec{
		Name: "write",
		Inputs: []Pin{
			Pin{"writer", WRITER},
			Pin{"msg", ANY},
		},
		Outputs: []Pin{
			Pin{"writer", WRITER},
		},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {

			writer, ok := in[0].(io.Writer)
			if !ok {
				out[0] = NewError("writer must implement io.Writer")
				return nil
			}

			data, err := json.Marshal(in[1])
			if err != nil {
				out[0] = NewError("could not marshal msg")
				return nil
			}

			_, err = writer.Write(data)
			if err != nil {
				out[0] = NewError("could not write data to writer")
				return nil
			}

			out[0] = writer

			return nil

		},
	}
}

func Close() Spec {
	return Spec{
		Name: "close",
		Inputs: []Pin{
			Pin{"writer", WRITER},
		},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {

			writer, ok := in[0].(io.Closer)
			if !ok {
				out[0] = NewError("writer must implement io.Writer")
				return nil
			}

			err := writer.Close()
			if err != nil {
				out[0] = NewError("could not close writer")
			}

			return nil

		},
	}
}

func Flush() Spec {
	return Spec{
		Name: "flush",
		Inputs: []Pin{
			Pin{"writer", WRITER},
		},
		Outputs: []Pin{
			Pin{"writer", WRITER},
		},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {

			writer, ok := in[0].(http.Flusher)
			if !ok {
				out[0] = NewError("writer must implement http.Flusher")
				return nil
			}

			writer.Flush()

			out[0] = writer

			return nil

		},
	}
}
