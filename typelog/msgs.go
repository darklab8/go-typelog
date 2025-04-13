package typelog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

func (l *Logger) Debug(msg string, opts ...LogType) {
	if IsMsgEnabled(l.level_log, LEVEL_DEBUG) {
		if l.enable_scope_shown {
			opts = append(opts, String("scope", l.scope))
		}
		args := append([]SlogAttr{}, newSlogArgs(opts...)...)
		if l.enable_file_shown {
			args = append(args, logGroupFiles())
		}
		l.logger.Debug(msg, args...)
	}
}

func (l *Logger) DebugCtx(ctx context.Context, msg string, opts ...LogType) {
	l.Debug(msg, CtxOpts(ctx, opts...)...)
}

func (l *Logger) Info(msg string, opts ...LogType) {
	if IsMsgEnabled(l.level_log, LEVEL_INFO) {
		if l.enable_scope_shown {
			opts = append(opts, String("scope", l.scope))
		}
		args := append([]SlogAttr{}, newSlogArgs(opts...)...)
		if l.enable_file_shown {
			args = append(args, logGroupFiles())
		}
		l.logger.Info(msg, args...)
	}
}

func (l *Logger) InfoCtx(ctx context.Context, msg string, opts ...LogType) {
	l.Info(msg, CtxOpts(ctx, opts...)...)
}

// Warn is for just potentially bad behavior to be aware of
func (l *Logger) Warn(msg string, opts ...LogType) {
	if IsMsgEnabled(l.level_log, LEVEL_WARN) {
		if l.enable_scope_shown {
			opts = append(opts, String("scope", l.scope))
		}
		args := append([]SlogAttr{}, newSlogArgs(opts...)...)
		if l.enable_file_shown {
			args = append(args, logGroupFiles())
		}
		l.logger.Warn(msg, args...)
	}
}

func (l *Logger) WarnCtx(ctx context.Context, msg string, opts ...LogType) {
	l.Warn(msg, CtxOpts(ctx, opts...)...)
}

// Error is bad but program can recover from it
func (l *Logger) Error(msg string, opts ...LogType) {
	if IsMsgEnabled(l.level_log, LEVEL_ERROR) {
		if l.enable_scope_shown {
			opts = append(opts, String("scope", l.scope))
		}
		args := append([]SlogAttr{}, newSlogArgs(opts...)...)
		if l.enable_file_shown {
			args = append(args, logGroupFiles())
		}
		l.logger.Error(msg, args...)
	}

}

func (l *Logger) ErrorCtx(ctx context.Context, msg string, opts ...LogType) {
	l.Error(msg, CtxOpts(ctx, opts...)...)
}

// Fatal when encountered, Program is not allowed to run further with fatal
func (l *Logger) Fatal(msg string, opts ...LogType) {
	if l.enable_scope_shown {
		opts = append(opts, String("scope", l.scope))
	}
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	if l.enable_file_shown {
		args = append(args, logGroupFiles())
	}
	l.logger.Error(msg, args...)

	os.Exit(1)
}

func (l *Logger) FatalCtx(ctx context.Context, msg string, opts ...LogType) {
	l.Fatal(msg, CtxOpts(ctx, opts...)...)
}

func (l *Logger) Panic(msg string, opts ...LogType) {
	if l.enable_scope_shown {
		opts = append(opts, String("scope", l.scope))
	}
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	if l.enable_file_shown {
		args = append(args, logGroupFiles())
	}
	l.logger.Error(msg, args...)
	panic(msg)
}

func (l *Logger) PanicCtx(ctx context.Context, msg string, opts ...LogType) {
	l.Panic(msg, CtxOpts(ctx, opts...)...)
}

func (l *Logger) CheckDebug(err error, msg string, opts ...LogType) bool {
	if err == nil {
		return false
	}
	if !IsMsgEnabled(l.level_log, LEVEL_DEBUG) {
		return true
	}
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	args = append(args, slog.String("error", fmt.Sprintf("%v", err)))
	l.logger.Debug(msg, args...)
	return true
}

func (l *Logger) CheckDebugCtx(ctx context.Context, err error, msg string, opts ...LogType) {
	l.CheckDebug(err, msg, CtxOpts(ctx, opts...)...)
}

func (l *Logger) CheckWarn(err error, msg string, opts ...LogType) bool {
	if err == nil {
		return false
	}
	if !IsMsgEnabled(l.level_log, LEVEL_WARN) {
		return true
	}
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	args = append(args, slog.String("error", fmt.Sprintf("%v", err)))
	l.logger.Warn(msg, args...)
	return true
}

func (l *Logger) CheckWarnCtx(ctx context.Context, err error, msg string, opts ...LogType) {
	l.CheckWarn(err, msg, CtxOpts(ctx, opts...)...)
}

func (l *Logger) CheckError(err error, msg string, opts ...LogType) bool {
	if err == nil {
		return false
	}
	if !IsMsgEnabled(l.level_log, LEVEL_ERROR) {
		return true
	}
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	args = append(args, slog.String("error", fmt.Sprintf("%v", err)))
	l.logger.Error(msg, args...)
	return true
}

func (l *Logger) CheckErrorCtx(ctx context.Context, err error, msg string, opts ...LogType) {
	l.CheckError(err, msg, CtxOpts(ctx, opts...)...)
}

// CheckFatal has shorter error output in comparison to CheckPanic
func (l *Logger) CheckFatal(err error, msg string, opts ...LogType) {
	if err == nil {
		return
	}
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	args = append(args, slog.String("error", fmt.Sprintf("%v", err)))
	l.logger.Error(msg, args...)
	os.Exit(1)
}

func (l *Logger) CheckFatalCtx(ctx context.Context, err error, msg string, opts ...LogType) {
	l.CheckFatal(err, msg, CtxOpts(ctx, opts...)...)
}

func (l *Logger) CheckPanic(err error, msg string, opts ...LogType) {
	if err == nil {
		return
	}
	args := append([]SlogAttr{}, newSlogArgs(opts...)...)
	args = append(args,
		slog.String("err_msg", fmt.Sprintf("%v", err)),
		slog.String("err_type", fmt.Sprintf("%T", err)),
	)
	l.logger.Error(msg, args...)
	panic_logger.Error(msg, args...)
	panic(panic_str.String())
}

func (l *Logger) CheckPanicCtx(ctx context.Context, err error, msg string, opts ...LogType) {
	l.CheckPanic(err, msg, CtxOpts(ctx, opts...)...)
}
