package main

import "github.com/golang/protobuf/jsonpb"

var marshaler = &jsonpb.Marshaler{
	EnumsAsInts:  false,
	EmitDefaults: false,
	Indent:       "  ",
	OrigName:     true,
}
