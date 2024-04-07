package model

import "fmt"

type ErrorMsg struct {
	Error string `json:"error"`
}

func NewMsg(msg string) ErrorMsg {
	return ErrorMsg{Error: msg}
}

func NewMsgF(f string, args ...any) ErrorMsg {
	return NewMsg(fmt.Sprintf(f, args...))
}
