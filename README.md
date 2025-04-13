# Typelog - Static typed structured logging

## Description

This project araised from the need to log backend applications, aws lambdas and other stuff in modern cloud ecosystem. Logging systems today are able easily parsing JSON format out of the box.
Static typing approach brings here consistent way to define key values to final msg, as well as easier following Domain Driven Design, where logs consistently describe what they log. Static typed logging brings easy refactoring to any present logs.

## Features

- Accepts static typed components as optional params
  - it will not accept `any` options as slog.
  - has shortcut WithFields, to make clone of the logger with default logging fields
- Easy to turn on/off parameters by environment variables
  - Ability to define different log levels for different created loggers
- Easier turning complex objects into structured logging
  - accepts maps and structs as its params. It will parse them on their own.
- has ability to show "scope" of current logging (from which module the logs are)
- has ability to grab otlp span id/trace id if they are present

[See folder examples](./examples)

## Alternative Versions

- [Version in python](https://github.com/darklab8/py-typelog)

## How to use

install with `go get github.com/darklab8/go-typelog`

examples/logger/main.go
```go
package logger

import "github.com/darklab8/go-typelog/typelog"

var Log *typelog.Logger = typelog.NewLogger("typelog")
```

examples/params_test.go
```go
package examples

import (
	"log/slog"
	"testing"
	"time"

	"github.com/darklab8/go-typelog/examples/logger"
	"github.com/darklab8/go-typelog/typelog"
)

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

func NestedParam(value string) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(typelog.Group("nested", typelog.TurnMapToAttrs(map[string]any{
			"smth":   "abc",
			"number": 123,
		})...))
	}
}

type Smth struct {
	Value1  string
	Number1 int
}

func NestedStructParam(value string) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(
			typelog.Group("nested", typelog.TurnStructToAttrs(Smth{Value1: "123", Number1: 4})...),
			slog.Int("not_nested", 345),
		)
	}
}

func TestNested(t *testing.T) {
	logger := typelog.NewLogger("test", typelog.WithLogLevel(typelog.LEVEL_DEBUG), typelog.WithJsonFormat(true))

	logger.Debug("123", NestedParam("abc"))
	logger.Debug("456", NestedStructParam("abc"))
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
```
