package examples

import (
	"testing"
	"time"

	"github.com/darklab8/go-typelog/examples/logger"
	"github.com/darklab8/go-typelog/examples/typedlogs"
	"github.com/darklab8/go-typelog/examples/types"
	"github.com/darklab8/go-typelog/typelog"
)

func TestTypedLogs(t *testing.T) {
	worker_id := types.WorkerID(5)
	logger.Log.Info("Worker was started", typedlogs.WorkerID(worker_id))

	logger := logger.Log.WithFields(typedlogs.WorkerID(worker_id), typedlogs.TaskID("abc"))
	logger.Info("Worker started task")

	logger.Info("Worker finished task")
}

func TestUsingInitialized(t *testing.T) {

	logger.Log.Debug("123")

	logger.Log.Debug("123", typelog.TestParam(456))

	logger1 := logger.Log.WithFields(typelog.Int("worker_id", 10))

	logger1.Info("Worker made action1")
	logger1.Info("Worker made action2")

	logger2 := logger.Log.WithFields(typelog.Float64("smth", 13.54))
	logger2.Debug("try now")
	logger1.Info("Worker made action1", typelog.Bool("is_check", false))
}

func TestSlogging(t *testing.T) {

	logger := typelog.NewLogger("test", typelog.WithLogLevel(typelog.LEVEL_DEBUG))
	logger.Debug("123")

	logger.Debug("123", typelog.TestParam(456))
}

func TestNested(t *testing.T) {
	logger := typelog.NewLogger("test", typelog.WithLogLevel(typelog.LEVEL_DEBUG), typelog.WithJsonFormat(true))

	logger.Debug("123", typedlogs.NestedParam("abc"))
	logger.Debug("456", typedlogs.NestedStructParam("abc"))
}

func TestCopyingLoggers(t *testing.T) {
	logger := typelog.NewLogger("test", typelog.WithLogLevel(typelog.LEVEL_DEBUG), typelog.WithJsonFormat(true))

	logger1 := logger.WithFields(typelog.String("smth", "123"))
	logger2 := logger1.WithFields(typelog.Int("smth2", 2), typelog.String("anotheparam", "abc"))
	logger3 := logger2.WithFields(typelog.Time("smth3", time.Now()))

	logger1.Info("logger1 printed")
	logger2.Info("logger2 printed")
	logger3.Info("logger3 printed")
}
