package goplogjson

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/axpira/gop/log"
)

var (
	// LevelFieldName defines the key name for Level field
	LevelFieldName = "level"
	// MessageFieldName defines the key name for Message field
	MessageFieldName = "msg"
	// ErrorFieldName defines the key name for Error field
	ErrorFieldName = "err"
	// TimestampFieldName defines the key name for Timestamp field
	TimestampFieldName = "time"

	TimestampEnabled = true
	TimestampFunc    = time.Now
	TimeFormat       = time.RFC3339
	TimestampFormat  = time.RFC3339
	// DurationFieldUnit defines the unit for time.Duration type fields added
	// using the Dur method.
	DurationFieldUnit = time.Millisecond

	LevelNameFunc = func(l log.Level) string {
		switch l {
		case log.NoLevel:
			return "all"
		case log.TraceLevel:
			return "trace"
		case log.DebugLevel:
			return "debug"
		case log.InfoLevel:
			return "info"
		case log.WarnLevel:
			return "warn"
		case log.ErrorLevel:
			return "error"
		case log.FatalLevel:
			return "fatal"
		case log.PanicLevel:
			return "panic"
		case log.DisabledLevel:
			return "disabled"
		}
		return "unknown"
	}
)

type field struct {
	buf []byte
}

func appendKey(buf []byte, key string) []byte {
	return append(appendString(append(buf, ','), key), ':')
}

func (f *field) Ctx(ctx context.Context) log.FieldBuilder {
	return f
}

func (f *field) Str(key string, value string) log.FieldBuilder {
	f.buf = appendString(appendKey(f.buf, key), value)
	return f
}

func (f *field) Bool(key string, value bool) log.FieldBuilder {
	f.buf = strconv.AppendBool(appendKey(f.buf, key), value)
	return f
}

func (f *field) Bytes(key string, value []byte) log.FieldBuilder {
	f.buf = appendKey(f.buf, key)
	f.buf = append(f.buf, '"')
	for i := 0; i < len(value); i++ {
		if !noEscapeTable[value[i]] {
			f.buf = appendBytesComplex(f.buf, value, i)
			f.buf = append(f.buf, '"')
			return f
		}
	}
	f.buf = append(f.buf, value...)
	f.buf = append(f.buf, '"')
	return f
}

func (f *field) Int64(key string, value int64) log.FieldBuilder {
	f.buf = strconv.AppendInt(appendKey(f.buf, key), value, 10)
	return f
}

func (f *field) Uint64(key string, value uint64) log.FieldBuilder {
	f.buf = strconv.AppendUint(appendKey(f.buf, key), value, 10)
	return f
}

func (f *field) Float32(key string, value float32) log.FieldBuilder {
	f.buf = strconv.AppendFloat(appendKey(f.buf, key), float64(value), 'f', -1, 32)
	return f
}

func (f *field) Float64(key string, value float64) log.FieldBuilder {
	f.buf = strconv.AppendFloat(appendKey(f.buf, key), value, 'f', -1, 64)
	return f
}

func (f *field) Msg(msg string) log.FieldBuilder {
	return f.Str(MessageFieldName, msg)
}

func (f *field) Marshal(key string, value interface{}) log.FieldBuilder {
	if obj, ok := value.(log.LogMarshaler); ok {
		return f.marshalLog(key, obj)
	}
	jsonByteArr, err := json.Marshal(value)
	if err != nil {
		return f.Error(key, err)
	}
	f.buf = append(appendKey(f.buf, key), jsonByteArr...)
	return f
}

func (f *field) marshalLog(key string, value log.LogMarshaler) log.FieldBuilder {
	builder := newField()
	value.MarshalLog(builder)
	return f.Dict(key, builder)
}

func (f *field) Any(key string, value interface{}) log.FieldBuilder {
	switch v := value.(type) {
	case string:
		return f.Str(key, v)
	case bool:
		return f.Bool(key, v)
	case int:
		return f.Int(key, v)
	case int8:
		return f.Int8(key, v)
	case int16:
		return f.Int16(key, v)
	case int32:
		return f.Int32(key, v)
	case int64:
		return f.Int64(key, v)
	case uint:
		return f.Uint(key, v)
	case uint8:
		return f.Uint8(key, v)
	case uint16:
		return f.Uint16(key, v)
	case uint32:
		return f.Uint32(key, v)
	case uint64:
		return f.Uint64(key, v)
	case float32:
		return f.Float32(key, v)
	case float64:
		return f.Float64(key, v)
	case time.Time:
		return f.Time(key, v)
	case time.Duration:
		return f.Dur(key, v)
	case complex64:
		return f.Complex64(key, v)
	case complex128:
		return f.Complex128(key, v)
	case []byte:
		return f.Bytes(key, v)
	case error:
		return f.Error(key, v)
	case fmt.Stringer:
		return f.Stringer(key, v)
	default:
		return f.Marshal(key, v)
	}
}

func (f *field) Fields(m map[string]interface{}) log.FieldBuilder {
	for k, v := range m {
		f.Any(k, v)
	}
	return f
}

func (f *field) Complex64(key string, value complex64) log.FieldBuilder {
	return f.Complex128(key, complex128(value))
}

func (f *field) Complex128(key string, value complex128) log.FieldBuilder {
	return f.Str(key, strconv.FormatComplex(value, 'g', 4, 128))
}

func (f *field) Dict(key string, fi log.FieldBuilder) log.FieldBuilder {
	f.buf = appendKey(f.buf, key)
	fi1 := fi.(*field)
	fi1.buf[0] = '{'
	f.buf = append(f.buf, fi1.buf...)
	f.buf = append(f.buf, '}')
	fi1.discard()
	return f
}

func (f *field) Int(key string, value int) log.FieldBuilder {
	return f.Int64(key, int64(value))
}

func (f *field) Int8(key string, value int8) log.FieldBuilder {
	return f.Int64(key, int64(value))
}

func (f *field) Int16(key string, value int16) log.FieldBuilder {
	return f.Int64(key, int64(value))
}

func (f *field) Int32(key string, value int32) log.FieldBuilder {
	return f.Int64(key, int64(value))
}

func (f *field) Uint(key string, value uint) log.FieldBuilder {
	return f.Uint64(key, uint64(value))
}

func (f *field) Uint8(key string, value uint8) log.FieldBuilder {
	return f.Uint64(key, uint64(value))
}

func (f *field) Uint16(key string, value uint16) log.FieldBuilder {
	return f.Uint64(key, uint64(value))
}

func (f *field) Uint32(key string, value uint32) log.FieldBuilder {
	return f.Uint64(key, uint64(value))
}

func (f *field) Timef(key string, value time.Time, format string) log.FieldBuilder {
	if value.IsZero() {
		return f
	}
	f.buf = appendKey(f.buf, key)
	f.buf = append(f.buf, '"')
	f.buf = value.AppendFormat(f.buf, format)
	f.buf = append(f.buf, '"')
	return f
}

func (f *field) Time(key string, value time.Time) log.FieldBuilder {
	return f.Timef(key, value, TimeFormat)
}

func (f *field) Dur(key string, value time.Duration) log.FieldBuilder {
	return f.Int64(key, int64(value/DurationFieldUnit))
}

func (f *field) Stringer(key string, value fmt.Stringer) log.FieldBuilder {
	return f.Str(key, value.String())
}

func (f *field) Error(key string, value error) log.FieldBuilder {
	if value == nil {
		return f
	}
	return f.Str(key, value.Error())
}

func (f *field) Err(value error) log.FieldBuilder {
	if value == nil {
		return f
	}
	return f.Str(ErrorFieldName, value.Error())
}

func (f *field) Msgf(format string, args ...interface{}) log.FieldBuilder {
	return f.Msg(fmt.Sprintf(format, args...))
}

func (f *field) Update(l log.Logger) log.Logger {
	ll := l.(*logger).clone()
	ll.buf = append(l.(*logger).buf, f.buf...)
	putField(f)
	return ll
}

func (f *field) send(out io.Writer) {
	if TimestampEnabled {
		f.Timef(TimestampFieldName, TimestampFunc(), TimestampFormat)
	}
	f.buf[0] = '{'
	out.Write(append(f.buf, "}\n"...))
	putField(f)
}

func (f *field) discard() {
	putField(f)
}

func newField() *field {
	e := fieldPool.Get().(*field)
	e.buf = e.buf[:0]
	return e
}

var fieldPool = &sync.Pool{
	New: func() interface{} {
		return &field{
			buf: make([]byte, 0, 500),
		}
	},
}

func putField(e *field) {
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16 // 64KiB
	if cap(e.buf) > maxSize {
		return
	}
	fieldPool.Put(e)
}
