package types

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level zapcore.Level

const TraceLevel Level = -2 // TRACE is lower than DEBUG (-1)

func (l Level) String() string {
	if l == TraceLevel {
		return "trace"
	}

	return zapcore.Level(l).String()
}

type SugarWithTrace struct {
	*zap.SugaredLogger
}

func (l *SugarWithTrace) Tracew(msg string, keysAndValues ...interface{}) {
	desugared := l.Desugar().WithOptions(zap.AddCallerSkip(1))
	if ce := desugared.Check(zapcore.Level(TraceLevel), msg); ce != nil {
		ce.Write(l.sweetenFields(keysAndValues)...)
	}
}

func (l *SugarWithTrace) sweetenFields(args []interface{}) []zap.Field {
	if len(args) == 0 {
		return nil
	}

	fields := make([]zap.Field, 0, len(args)/2)
	for i := 0; i < len(args); {
		// If odd number of args, log the last one as "UNPAIRED"
		if i == len(args)-1 {
			fields = append(fields, zap.Any("UNPAIRED", args[i]))
			break
		}

		key, ok := args[i].(string)
		if !ok {
			// if the key is not a string, fall back to fmt.Sprint
			key = fmt.Sprint(args[i])
		}

		val := args[i+1]
		// If user passed zap.Field directly
		if f, ok := val.(zap.Field); ok {
			fields = append(fields, f)
		} else {
			fields = append(fields, zap.Any(key, val))
		}
		i += 2
	}
	return fields
}

func CustomLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if l == zapcore.Level(TraceLevel) {
		enc.AppendString("TRACE")
		return
	}
	enc.AppendString(strings.ToUpper(l.String()))
}

func NewLogger() (*SugarWithTrace, error) {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.LevelKey = "level"
	encoderCfg.EncodeLevel = CustomLevelEncoder

	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	writer := zapcore.AddSync(os.Stdout)

	core := zapcore.NewCore(
		encoder,
		writer,
		zap.NewAtomicLevelAt(zapcore.Level(TraceLevel)),
	)

	logger := zap.New(core, zap.AddCaller())

	return &SugarWithTrace{logger.Sugar()}, nil
}
