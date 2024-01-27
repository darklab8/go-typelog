package logcore

import (
	"fmt"
	"log/slog"
	"time"
)

type SlogGroup struct {
	slogs []SlogAttr
}

func (s SlogGroup) Render() []SlogAttr {
	return s.slogs
}

func (s SlogGroup) Append(params ...slog.Attr) {
	for _, param := range params {
		s.slogs = append(s.slogs, param)
	}
}

type SlogParam func(r *SlogGroup)

type SlogAttr = any

func newSlogArgs(opts ...SlogParam) []SlogAttr {
	client := &SlogGroup{}
	for _, opt := range opts {
		opt(client)
	}

	return (*client).Render()
}

func TestParam(value int) SlogParam {
	return func(c *SlogGroup) {
		c.Append(slog.Int("test_param", value))
	}
}

func Any(key string, value any) SlogParam {
	return func(c *SlogGroup) {
		c.Append(slog.String(key, fmt.Sprintf("%v", value)))
	}
}

func String(key string, value string) SlogParam {
	return func(c *SlogGroup) {
		c.Append(slog.String(key, value))
	}
}

func Int(key string, value int) SlogParam {
	return func(c *SlogGroup) {
		c.Append(slog.Int(key, value))
	}
}
func Int64(key string, value int64) SlogParam {
	return func(c *SlogGroup) {
		c.Append(slog.Int64(key, value))
	}
}
func Float32(key string, value float32) SlogParam {
	return func(c *SlogGroup) {
		c.Append(slog.Float64(key, float64(value)))
	}
}
func Time(key string, value time.Time) SlogParam {
	return func(c *SlogGroup) {
		c.Append(slog.Time(key, value))
	}
}
func Float64(key string, value float64) SlogParam {
	return func(c *SlogGroup) {
		c.Append(slog.Float64(key, value))
	}
}
func Bool(key string, value bool) SlogParam {
	return func(c *SlogGroup) {
		c.Append(slog.Bool(key, value))
	}
}

func Expected(value any) SlogParam {
	return func(c *SlogGroup) {
		c.Append(slog.String("expected", fmt.Sprintf("%v", value)))
	}
}
func Actual(value any) SlogParam {
	return func(c *SlogGroup) {
		c.Append(slog.String("actual", fmt.Sprintf("%v", value)))
	}
}

func OptError(err error) SlogParam {
	return func(c *SlogGroup) {
		c.Append(
			slog.String("error_msg", fmt.Sprintf("%v", err)),
			slog.String("error_type", fmt.Sprintf("%T", err)),
		)
	}
}

func Items[T any](value []T, item_name string) SlogParam {
	return func(c *SlogGroup) {
		sliced_string := fmt.Sprintf("%v", value)
		if len(sliced_string) > 300 {
			sliced_string = sliced_string[:300] + "...sliced string"
		}
		c.Append(slog.String(item_name, fmt.Sprintf("%v", value)))
		c.Append(slog.String(fmt.Sprintf("%s_len", item_name), fmt.Sprintf("%d", len(value))))
	}
}

func Records[T any](value []T) SlogParam {
	return Items[T](value, "records")
}

func Args(value []string) SlogParam {
	return Items[string](value, "args")
}

func Bytes(key string, value []byte) SlogParam {
	return func(c *SlogGroup) {
		c.Append(slog.String(key, string(value)))
	}
}

func Struct(value any) SlogParam {
	return func(c *SlogGroup) {
		c.Append(TurnMapToAttrs(StructToMap(value))...)
	}
}

func NestedStruct(key string, value any) SlogParam {
	return func(c *SlogGroup) {
		attrs := TurnMapToAttrs(StructToMap(value))
		c.Append(slog.Group(key, AttrsToAny(attrs)...))
	}
}

func Map(value map[string]any) SlogParam {
	return func(c *SlogGroup) {
		c.Append(TurnMapToAttrs(value)...)
	}
}
