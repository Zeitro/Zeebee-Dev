package gormx

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/timex"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	SlowThreshold time.Duration
	mod           string
}

func NewGormLogger(mod string) *GormLogger {
	return &GormLogger{
		SlowThreshold: 800 * time.Millisecond,
		mod:           mod,
	}
}

var _ logger.Interface = (*GormLogger)(nil)

func (l *GormLogger) LogMode(lev logger.LogLevel) logger.Interface {
	return &GormLogger{}
}
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Infof(msg, data)
}
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Errorf(msg, data)
}
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Errorf(msg, data)
}
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// duration
	elapsed := time.Since(begin)
	// sql and rows
	sql, rows := fc()
	logFields := []logx.LogField{
		logx.Field("sql", sql),
		logx.Field("duration", timex.ReprOfDuration(elapsed)),
		logx.Field("rows", rows),
	}
	subLog := logx.WithContext(ctx).WithCallerSkip(3)
	if err != nil {
		// not found record
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logFields = append(logFields, logx.Field("sql_type", "record_not_found"))
			subLog.Infow("Database ErrRecordNotFound", logFields...)
		} else {
			logFields = append(logFields, logx.Field("sql_type", "error"))
			subLog.Errorw(err.Error(), logFields...)
		}
		return
	}
	// slow log
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		logFields = append(logFields, logx.Field("sql_type", "slow_sql"))
		subLog.Sloww("Slow Sql Query", logFields...)
	}
	//
	if l.mod != service.ProMode {
		subLog.Infow("Sql Query", logFields...)
	}
}
