package typedlogs

import (
	"log/slog"

	"github.com/darklab8/go-typelog/examples/types"
	"github.com/darklab8/go-typelog/typelog"
)

func TaskID(value types.TaskID) typelog.LogType { return typelog.String("task_id", string(value)) }

func WorkerID(value types.WorkerID) typelog.LogType {
	return typelog.Int("worker_id", int(value))
}

type Smth struct {
	Value1  string
	Number1 int
}

func NestedStructParam(value string) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(
			typelog.Group("nested", typelog.TurnStructToAttrs(Smth{Value1: value, Number1: 4})...),
			slog.Int("not_nested", 345),
		)
	}
}

func NestedParam(value string) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(typelog.Group("nested", typelog.TurnMapToAttrs(map[string]any{
			"smth":   value,
			"number": 123,
		})...))
	}
}
