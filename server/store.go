package server

type Value struct{
	Data string
}

var Store = map[string]Value{}