package examples

import (
	"log/slog"
	"testing"
	"time"

	"github.com/darklab8/logusgo/examples/logus"
	"github.com/darklab8/logusgo/logcore"
)

func TestUsingInitialized(t *testing.T) {

	logus.Log.Debug("123")

	logus.Log.Debug("123", logcore.TestParam(456))

	logger := logus.Log.WithFields(logcore.Int("worker_id", 10))

	logger.Info("Worker made action1")
	logger.Info("Worker made action2")

	logger2 := logus.Log.WithFields(logcore.Float64("smth", 13.54))
	logger2.Debug("try now")
	logger.Info("Worker made action1", logcore.Bool("is_check", false))
}

func TestSlogging(t *testing.T) {

	logger := logcore.NewLogger("test", logcore.WithLogLevel(logcore.LEVEL_DEBUG))
	logger.Debug("123")

	logger.Debug("123", logcore.TestParam(456))
}

func NestedParam(value string) logcore.SlogParam {
	return func(c *logcore.SlogGroup) {
		c.Append(slog.Group("nested", logcore.AttrsToAny(logcore.TurnMapToAttrs(map[string]any{
			"smth":   "abc",
			"number": 123,
		}))...))
	}
}

type Smth struct {
	Value1  string
	Number1 int
}

func NestedStructParam(value string) logcore.SlogParam {
	return func(c *logcore.SlogGroup) {
		c.Append(
			slog.Group("nested", logcore.AttrsToAny(logcore.TurnStructToAttrs(Smth{Value1: "123", Number1: 4}))...),
			slog.Int("not_nested", 345),
		)
	}
}

func TestNested(t *testing.T) {
	logger := logcore.NewLogger("test", logcore.WithLogLevel(logcore.LEVEL_DEBUG), logcore.WithJsonFormat(true))

	logger.Debug("123", NestedParam("abc"))
	logger.Debug("456", NestedStructParam("abc"))
}

func TestCopyingLoggers(t *testing.T) {
	logger := logcore.NewLogger("test", logcore.WithLogLevel(logcore.LEVEL_DEBUG), logcore.WithJsonFormat(true))

	logger1 := logger.WithFields(logcore.String("smth", "123"))
	logger2 := logger1.WithFields(logcore.Int("smth2", 2), logcore.String("anotheparam", "abc"))
	logger3 := logger2.WithFields(logcore.Time("smth3", time.Now()))

	logger1.Info("logger1 printed")
	logger2.Info("logger2 printed")
	logger3.Info("logger3 printed")
}
