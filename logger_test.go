package goplogjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/axpira/gop/log"
)

var (
	now, _        = time.Parse(time.RFC3339, "2021-09-26T07:57:36Z")
	dur, _        = time.ParseDuration("1h2m34s")
	complexString = "\"\b\f\n\r\t"
)

func init() {
	TimestampFunc = func() time.Time {
		return now
	}
	TimeFormat = time.RFC3339
}

func TestLoggerFields(t *testing.T) {
	tests := map[string]struct {
		logger func(out io.Writer)
		want   string
	}{
		"must log all fields": {
			logger: func(out io.Writer) {
				l := New(WithOutput(out))
				l.Log(log.InfoLevel, newFields(l, "key_").
					Dict("key_dict", newFields(l, "key_dict_")),
				)
			},
			want: createJson("key_", "info", "dict"),
		},
		"must log complex string": {
			logger: func(out io.Writer) {
				New(WithOutput(out)).
					Log(log.InfoLevel, newField().
						Str("key_complex_"+complexString+"_str", complexString).
						Str("str_no_encoding", `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`).
						Str("str_encoding_first", `"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`).
						Str("str_encoding_middle", `aaaaaaaaaaaaaaaaaaaaaaaaa"aaaaaaaaaaaaaaaaaaaaaaaa`).
						Str("str_encoding_last", `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"`).
						Str("str_multibytes_first", `❤️aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`).
						Str("str_multibytes_middle", `aaaaaaaaaaaaaaaaaaaaaaaaa❤️aaaaaaaaaaaaaaaaaaaaaaaa`).
						Str("str_multibytes_last", `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa❤️`),
					)
			},
			want: `{
				"key_complex_\"\b\f\n\r\t_str":"\"\b\f\n\r\t",
				"level":"info",
				"time":"2021-09-26T07:57:36Z",
				"str_no_encoding": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				"str_encoding_first":    "\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				"str_encoding_middle":   "aaaaaaaaaaaaaaaaaaaaaaaaa\"aaaaaaaaaaaaaaaaaaaaaaaa",
				"str_encoding_last":     "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\"",
				"str_multibytes_first":  "❤️aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				"str_multibytes_middle": "aaaaaaaaaaaaaaaaaaaaaaaaa❤️aaaaaaaaaaaaaaaaaaaaaaaa",
				"str_multibytes_last":   "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa❤️"
			}`,
		},
		"must log complex bytes": {
			logger: func(out io.Writer) {
				New(WithOutput(out)).
					Log(log.InfoLevel, newField().
						Bytes("key_complex_"+complexString+"_bytes", []byte(complexString)).
						Bytes("bytes_no_encoding", []byte(`aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`)).
						Bytes("bytes_encoding_first", []byte(`"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`)).
						Bytes("bytes_encoding_middle", []byte(`aaaaaaaaaaaaaaaaaaaaaaaaa"aaaaaaaaaaaaaaaaaaaaaaaa`)).
						Bytes("bytes_encoding_last", []byte(`aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"`)).
						Bytes("bytes_multibytes_first", []byte(`❤️aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`)).
						Bytes("bytes_multibytes_middle", []byte(`aaaaaaaaaaaaaaaaaaaaaaaaa❤️aaaaaaaaaaaaaaaaaaaaaaaa`)).
						Bytes("bytes_multibytes_last", []byte(`aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa❤️`)),
					)
			},
			want: `{
				"key_complex_\"\b\f\n\r\t_bytes":"\"\b\f\n\r\t",
				"level":"info",
				"time":"2021-09-26T07:57:36Z",
				"bytes_no_encoding": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				"bytes_encoding_first":    "\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				"bytes_encoding_middle":   "aaaaaaaaaaaaaaaaaaaaaaaaa\"aaaaaaaaaaaaaaaaaaaaaaaa",
				"bytes_encoding_last":     "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\"",
				"bytes_multibytes_first":  "❤️aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				"bytes_multibytes_middle": "aaaaaaaaaaaaaaaaaaaaaaaaa❤️aaaaaaaaaaaaaaaaaaaaaaaa",
				"bytes_multibytes_last":   "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa❤️"
			}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			out := new(strings.Builder)
			tc.logger(out)
			want := tc.want
			got := out.String()
			if diff := compareJson(want, got); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

var levelFuncs = map[log.Level]func(log.Logger, log.FieldBuilder, string, string, ...interface{}){
	log.TraceLevel: func(l log.Logger, fb log.FieldBuilder, msg string, format string, args ...interface{}) {
		l.Trc(fb)
		l.Trace(msg)
		l.Tracef(format, args...)
	},
	log.DebugLevel: func(l log.Logger, fb log.FieldBuilder, msg string, format string, args ...interface{}) {
		l.Dbg(fb)
		l.Debug(msg)
		l.Debugf(format, args...)
	},
	log.InfoLevel: func(l log.Logger, fb log.FieldBuilder, msg string, format string, args ...interface{}) {
		l.Inf(fb)
		l.Info(msg)
		l.Infof(format, args...)
	},
	log.WarnLevel: func(l log.Logger, fb log.FieldBuilder, msg string, format string, args ...interface{}) {
		l.Wrn(fb)
		l.Warn(msg)
		l.Warnf(format, args...)
	},
	log.ErrorLevel: func(l log.Logger, fb log.FieldBuilder, msg string, format string, args ...interface{}) {
		l.Err(fb)
		l.Error(msg, nil)
		l.Errorf(format, args...)
	},
	log.FatalLevel: func(l log.Logger, fb log.FieldBuilder, msg string, format string, args ...interface{}) {
		l.Ftl(fb)
		l.Fatal(msg)
		l.Fatalf(format, args...)
	},
	log.PanicLevel: func(l log.Logger, fb log.FieldBuilder, msg string, format string, args ...interface{}) {
		l.Pnc(fb)
		l.Panic(msg)
		l.Panicf(format, args...)
	},
}

func TestMustLogCorrectLevel(t *testing.T) {
	tests := []struct {
		loggerLevel   log.Level
		allowedLevels []log.Level
		denyLevels    []log.Level
	}{
		{
			loggerLevel:   log.NoLevel,
			allowedLevels: []log.Level{log.TraceLevel, log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.FatalLevel, log.PanicLevel},
		},
		{
			loggerLevel:   log.TraceLevel,
			allowedLevels: []log.Level{log.TraceLevel, log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.FatalLevel, log.PanicLevel},
		},
		{
			loggerLevel:   log.DebugLevel,
			allowedLevels: []log.Level{log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.FatalLevel, log.PanicLevel},
			denyLevels:    []log.Level{log.TraceLevel},
		},
		{
			loggerLevel:   log.InfoLevel,
			allowedLevels: []log.Level{log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.FatalLevel, log.PanicLevel},
			denyLevels:    []log.Level{log.TraceLevel, log.DebugLevel},
		},
		{
			loggerLevel:   log.WarnLevel,
			allowedLevels: []log.Level{log.WarnLevel, log.ErrorLevel, log.FatalLevel, log.PanicLevel},
			denyLevels:    []log.Level{log.TraceLevel, log.DebugLevel, log.InfoLevel},
		},
		{
			loggerLevel:   log.ErrorLevel,
			allowedLevels: []log.Level{log.ErrorLevel, log.FatalLevel, log.PanicLevel},
			denyLevels:    []log.Level{log.TraceLevel, log.DebugLevel, log.InfoLevel, log.WarnLevel},
		},
		{
			loggerLevel:   log.FatalLevel,
			allowedLevels: []log.Level{log.FatalLevel, log.PanicLevel},
			denyLevels:    []log.Level{log.TraceLevel, log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel},
		},
		{
			loggerLevel:   log.PanicLevel,
			allowedLevels: []log.Level{log.PanicLevel},
			denyLevels:    []log.Level{log.TraceLevel, log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.FatalLevel},
		},
		{
			loggerLevel: log.DisabledLevel,
			denyLevels:  []log.Level{log.TraceLevel, log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.FatalLevel, log.PanicLevel},
		},
	}

	levelHook = nil
	for _, tc := range tests {
		for _, denyLevel := range tc.denyLevels {
			out := new(strings.Builder)
			l := New(WithOutput(out), WithLevel(tc.loggerLevel))
			fields := l.NewFieldBuilder().Msg("Hello World")
			l.Log(denyLevel, fields)
			if f, ok := levelFuncs[denyLevel]; ok {
				f(l, fields, "Hello World", "Hello %s", "World")
			}
			for i, line := range strings.Split(out.String(), "\n") {
				if line != "" {
					t.Errorf("when log level is %s want no log when send %s and got %q on line %d", LevelNameFunc(tc.loggerLevel), LevelNameFunc(denyLevel), line, i)
				}
			}
		}
		for _, allowedLevel := range tc.allowedLevels {
			out := new(strings.Builder)
			l := New(WithOutput(out), WithLevel(tc.loggerLevel))
			fields := l.NewFieldBuilder().Msg("Hello World")
			l.Log(allowedLevel, fields)
			want := fmt.Sprintf(`{"level":"%s", "msg":"Hello World", "time":"2021-09-26T07:57:36Z"}`, LevelNameFunc(allowedLevel))
			if diff := compareJson(want, out.String()); diff != "" {
				t.Errorf(diff)
			}
			if f, ok := levelFuncs[allowedLevel]; ok {
				f(l, fields, "Hello World", "Hello %s", "World")
				lines := strings.Split(out.String(), "\n")
				if diff := compareJson(want, lines[1]); diff != "" {
					t.Errorf(diff)
				}
				if diff := compareJson(want, lines[2]); diff != "" {
					t.Errorf(diff)
				}
				if diff := compareJson(want, lines[3]); diff != "" {
					t.Errorf(diff)
				}
			}
		}
	}
}

func createMapFields(prefix string) map[string]interface{} {
	return map[string]interface{}{
		prefix + "map_str":        "value str",
		prefix + "map_bool":       false,
		prefix + "map_bytes":      []byte("map bytes"),
		prefix + "map_complex64":  complex(float32(1), float32(2)),
		prefix + "map_complex128": complex(float64(3), float64(2)),
		prefix + "map_dur":        dur,
		prefix + "map_error":      errors.New("map error"),
		prefix + "map_float32":    float32(1.23),
		prefix + "map_float64":    float64(9.87),
		prefix + "map_int":        int(-42),
		prefix + "map_int8":       int8(-8),
		prefix + "map_int16":      int16(-16),
		prefix + "map_int32":      int32(-32),
		prefix + "map_int64":      int64(-64),
		prefix + "map_uint":       uint(42),
		prefix + "map_uint8":      uint8(8),
		prefix + "map_uint16":     uint16(16),
		prefix + "map_uint32":     uint32(32),
		prefix + "map_uint64":     uint64(64),
		prefix + "map_time":       now,
	}
}

func newFields(l log.Logger, prefix string) log.FieldBuilder {
	return l.NewFieldBuilder().
		Bool(prefix+"bool", true).
		Bytes(prefix+"bytes", []byte("bytes")).
		Complex64(prefix+"complex64", complex(float32(1), float32(2))).
		Complex128(prefix+"complex128", complex(float64(3), float64(2))).
		Dur(prefix+"dur", dur).
		Err(errors.New("unknown error")).
		Error(prefix+"error", errors.New("another unknown error")).
		Fields(createMapFields(prefix)).
		Float32(prefix+"float32", 1.23).
		Float64(prefix+"float64", 9.87).
		Int(prefix+"int", -42).
		Int8(prefix+"int8", -8).
		Int16(prefix+"int16", -16).
		Int32(prefix+"int32", -32).
		Int64(prefix+"int64", -64).
		Marshal(prefix+"marshal", "teste").
		Str(prefix+"str", "value str").
		Time(prefix+"time", now).
		Timef(prefix+"timef", now, time.RFC3339).
		Uint(prefix+"uint", 42).
		Uint8(prefix+"uint8", 8).
		Uint16(prefix+"uint16", 16).
		Uint32(prefix+"uint32", 32).
		Uint64(prefix+"uint64", 64)
}

func createJson(prefix string, level string, dictKey string) string {
	dict := ""
	if dictKey != "" {
		dict = fmt.Sprintf(`,"%s":%s`, prefix+dictKey, createJson(prefix+"dict_", level, ""))
	}
	return fmt.Sprintf(`{
	"%[1]sbool":true,
	"%[1]sbytes":"bytes",
	"%[1]scomplex64":"(1+2i)",
	"%[1]scomplex128":"(3+2i)",
	"%[1]sdur":3754000,
	"err":"unknown error",
	"%[1]serror":"another unknown error",
	"%[1]sfloat32":1.23,
	"%[1]sfloat64":9.87,
	"%[1]sint":-42,
	"%[1]sint8":-8,
	"%[1]sint16":-16,
	"%[1]sint32":-32,
	"%[1]sint64":-64,
	"%[1]smarshal":"teste",
	"%[1]sstr":"value str",
	"%[1]stime":"2021-09-26T07:57:36Z",
	"%[1]stimef":"2021-09-26T07:57:36Z",
	"%[1]suint":42,
	"%[1]suint8":8,
	"%[1]suint16":16,
	"%[1]suint32":32,
	"%[1]suint64":64,
	"level":"%[2]s",
	"time":"2021-09-26T07:57:36Z",
	"%[1]smap_bool":false,
	"%[1]smap_bytes":"map bytes",
	"%[1]smap_complex64":"(1+2i)",
	"%[1]smap_complex128":"(3+2i)",
	"%[1]smap_dur":3754000,
	"%[1]smap_error":"map error",
	"%[1]smap_float32":1.23,
	"%[1]smap_float64":9.87,
	"%[1]smap_int":-42,
	"%[1]smap_int8":-8,
	"%[1]smap_int16":-16,
	"%[1]smap_int32":-32,
	"%[1]smap_int64":-64,
	"%[1]smarshal":"teste",
	"%[1]smap_str":"value str",
	"%[1]smap_time":"2021-09-26T07:57:36Z",
	"%[1]smap_uint":42,
	"%[1]smap_uint8":8,
	"%[1]smap_uint16":16,
	"%[1]smap_uint32":32,
	"%[1]smap_uint64":64
	%[3]s
}`, prefix, level, dict)
}

func compareJson(want, got string) string {
	wantMap, err := stringToMap(want)
	if err != nil {
		return fmt.Sprintf("error on convert want to map: %s\n%s", err.Error(), want)
	}
	gotMap, err := stringToMap(got)
	if err != nil {
		return fmt.Sprintf("error on convert got to map: %s\n%s", err.Error(), got)
	}
	return compareMaps(wantMap, gotMap)
}

func compareMaps(want, got map[string]interface{}) string {
	diff := new(strings.Builder)
	fmt.Fprintf(diff, "+++ Want\n")
	fmt.Fprintf(diff, "--- Got\n")
	fmt.Fprintf(diff, "want :%#+v\n", want)
	fmt.Fprintf(diff, "got  :%#+v\n", got)
	fmt.Fprintf(diff, "diff\n")
	equals := true
	for key, wantValue := range want {
		if gotValue, ok := got[key]; ok {
			switch w := want[key].(type) {
			case map[string]interface{}:
				fmt.Fprint(diff, compareMaps(w, got[key].(map[string]interface{})))
			default:
				if wantValue != gotValue {
					equals = false
					fmt.Fprintf(diff, "- (%T) %s: (%T) %v\n", key, key, wantValue, wantValue)
					fmt.Fprintf(diff, "+ (%T) %s: (%T) %v\n", key, key, gotValue, gotValue)
				} else {
					fmt.Fprintf(diff, "  (%T) %v: (%T) %v\n", key, key, gotValue, gotValue)
				}
			}
		} else {
			equals = false
			fmt.Fprintf(diff, "- (%T) %s: (%T) %v\n", key, key, wantValue, wantValue)
			fmt.Fprintf(diff, "+ nil\n")
		}
	}
	for key, gotValue := range got {
		if _, ok := want[key]; !ok {
			equals = false
			fmt.Fprintf(diff, "+ nil\n")
			fmt.Fprintf(diff, "- (%T) %s: (%T) %v\n", key, key, gotValue, gotValue)
		}
	}
	if !equals {
		return diff.String()
	}
	return ""
}

func stringToMap(str string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
