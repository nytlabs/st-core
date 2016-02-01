package core

import (
	"bufio"
	"errors"
	"os"
)

func StdinInterface() SourceSpec {
	return SourceSpec{
		Name: "stdin",
		Type: STDIN,
		New:  NewStdin,
	}
}

func (stdin Stdin) GetType() SourceType {
	return STDIN
}

type stdinMsg struct {
	msg string
	err error
}

type Stdin struct {
	quit        chan chan error
	fromScanner chan stdinMsg
}

func NewStdin() Source {
	stdin := &Stdin{
		quit:        make(chan chan error),
		fromScanner: make(chan stdinMsg),
	}
	return stdin
}

func (stdin *Stdin) Serve() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		txt := scanner.Text()
		stdin.fromScanner <- stdinMsg{txt, nil}
	}
	err := scanner.Err()
	if err != nil {
		stdin.fromScanner <- stdinMsg{"", err}
	} else {
		stdin.fromScanner <- stdinMsg{"", errors.New("EOF")}
	}
}

func (stdin *Stdin) Stop() {
}

func (stdin Stdin) ReceiveMessage(i chan Interrupt) (string, Interrupt, error) {
	select {
	case msg := <-stdin.fromScanner:
		//log.Println("FROM SCANNER", msg.err == nil, msg.err)
		return msg.msg, nil, msg.err
	case f := <-i:
		return "", f, nil
	}
}

func StdinReceive() Spec {
	return Spec{
		Name:    "stdinReceive",
		Outputs: []Pin{Pin{"msg", STRING}},
		Source:  STDIN,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			stdin := s.(*Stdin)
			msg, f, err := stdin.ReceiveMessage(i)
			if err != nil {
				out[0] = NewError("EOF")
				return nil
			}
			if f != nil {
				return f
			}
			out[0] = msg
			return nil
		},
	}
}
