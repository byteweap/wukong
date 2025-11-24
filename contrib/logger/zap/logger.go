package zap

import "github.com/byteweap/wukong/plugin/logger"

// Log TODO
type Log struct {
}

var _ logger.Logger = (*Log)(nil)

func New() *Log {
	return &Log{}
}

func (l *Log) With(s string, s2 string) logger.Logger {
	//TODO implement me
	panic("implement me")
}

func (l *Log) Debug() logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (l *Log) Info() logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (l *Log) Warn() logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (l *Log) Error() logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (l *Log) Fatal() logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (l *Log) Panic() logger.Entry {
	//TODO implement me
	panic("implement me")
}
