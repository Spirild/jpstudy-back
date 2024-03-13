package core

import (
	"fmt"
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type RotateLogger struct {
	Logger    *lumberjack.Logger
	SplitTime time.Time
	SDuration int32
	FileName  string
}

func (rl *RotateLogger) Write(p []byte) (n int, err error) {
	if !time.Now().Before(rl.SplitTime) {
		rl.Logger.Filename = fmt.Sprintf("%s[%s].log", rl.FileName, time.Now().Format("2006_01_02_15_04_05"))
		err := rl.Logger.Rotate()
		if err != nil {
			fmt.Printf("rotate err\n")
			return 0, err
		}
		rl.SplitTime = rl.getSplitTime()
	}
	return rl.Logger.Write(p)
}

func NewLogWritter(logFile string, spliteDuration int32) *RotateLogger {
	logger := &RotateLogger{
		Logger: &lumberjack.Logger{
			Filename:   fmt.Sprintf("%s[%s].log", logFile, time.Now().Format("2006_01_02_15_04_05")),
			MaxSize:    100, // mb,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   false,
		},
		// SplitTime: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.Now().Location()),
		SDuration: spliteDuration,
		FileName:  logFile,
	}
	logger.SplitTime = logger.getSplitTime()
	return logger
}

func (rl *RotateLogger) getSplitTime() time.Time {
	dayDuration := rl.SDuration / 24
	hourDuration := rl.SDuration % 24
	current := time.Now()
	if hourDuration != 0 {
		hourDuration += int32(current.Hour())
	}
	return time.Date(current.Year(), current.Month(), current.Day()+int(dayDuration), int(hourDuration), 0, 0, 0, current.Location())
}
