package goplogjson

import (
	"context"
	"fmt"
	"time"

	"github.com/axpira/gop/log"
)

var emptyFieldPtr = &emptyField{}

type emptyField struct {
}

func (f *emptyField) Ctx(context.Context) log.FieldBuilder {
	return f
}

func (f *emptyField) Str(string, string) log.FieldBuilder {
	return f
}

func (f *emptyField) Int64(string, int64) log.FieldBuilder {
	return f
}

func (f *emptyField) Uint64(string, uint64) log.FieldBuilder {
	return f
}

func (f *emptyField) Bool(string, bool) log.FieldBuilder {
	return f
}

func (f *emptyField) Float32(string, float32) log.FieldBuilder {
	return f
}

func (f *emptyField) Float64(string, float64) log.FieldBuilder {
	return f
}

func (f *emptyField) Msg(string) log.FieldBuilder {
	return f
}

func (f *emptyField) Any(string, interface{}) log.FieldBuilder {
	return f
}

func (f *emptyField) Marshal(string, interface{}) log.FieldBuilder {
	return f
}

func (f *emptyField) Fields(map[string]interface{}) log.FieldBuilder {
	return f
}

func (f *emptyField) Dict(string, log.FieldBuilder) log.FieldBuilder {
	return f
}

func (f *emptyField) Int(string, int) log.FieldBuilder {
	return f
}

func (f *emptyField) Int8(string, int8) log.FieldBuilder {
	return f
}

func (f *emptyField) Int16(string, int16) log.FieldBuilder {
	return f
}

func (f *emptyField) Int32(string, int32) log.FieldBuilder {
	return f
}

func (f *emptyField) Uint(string, uint) log.FieldBuilder {
	return f
}

func (f *emptyField) Uint8(string, uint8) log.FieldBuilder {
	return f
}

func (f *emptyField) Uint16(string, uint16) log.FieldBuilder {
	return f
}

func (f *emptyField) Uint32(string, uint32) log.FieldBuilder {
	return f
}

func (f *emptyField) Timef(string, time.Time, string) log.FieldBuilder {
	return f
}

func (f *emptyField) Time(string, time.Time) log.FieldBuilder {
	return f
}

func (f *emptyField) Dur(string, time.Duration) log.FieldBuilder {
	return f
}

func (f *emptyField) Stringer(string, fmt.Stringer) log.FieldBuilder {
	return f
}

func (f *emptyField) Error(string, error) log.FieldBuilder {
	return f
}

func (f *emptyField) Err(error) log.FieldBuilder {
	return f
}

func (f *emptyField) Msgf(string, ...interface{}) log.FieldBuilder {
	return f
}

func (f *emptyField) Complex64(string, complex64) log.FieldBuilder {
	return f
}

func (f *emptyField) Complex128(string, complex128) log.FieldBuilder {
	return f
}

func (f *emptyField) Bytes(string, []byte) log.FieldBuilder {
	return f
}

func (f *emptyField) Update(l log.Logger) log.Logger {
	return l
}

func (f *emptyField) Level(string) log.FieldBuilder {
	return f
}
