package go_utils

import (
	"fmt"
	"os"
	"path"
	"time"
	"github.com/fsnotify/fsnotify"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	)

func newRotateLogHook(logDir, baseFilename, logLevel string, rotateDays, maxRemainCnt uint) (logrus.Hook, error)  {
	writer, err := rotatelogs.New(
		path.Join(logDir, baseFilename+"-%Y-%m-%d.log"),
		rotatelogs.WithLinkName(path.Join(logDir, baseFilename+".log")),
		rotatelogs.WithRotationTime(time.Duration(rotateDays * 24) * time.Hour),
		rotatelogs.WithRotationCount(maxRemainCnt),
		)
	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
		return nil, err
	}

	level, ok := logrus.ParseLevel(logLevel)

	if ok == nil {
		logrus.SetLevel(level)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	hook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.WarnLevel: writer,
		logrus.InfoLevel: writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors:true})

	return hook, nil
}

func loadConfig() *viper.Viper {
	logConfViper := viper.New()
	dir, _ := os.Getwd()
	configFilepath := path.Join(dir, "conf", "log.yml")
	logConfViper.SetConfigFile(configFilepath)
	err := logConfViper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("load config file: ../conf/log.yml error: %v", err)
	}

	logConfViper.WatchConfig()
	logConfViper.OnConfigChange(func(event fsnotify.Event) {
		logrus.Infof("config file: ../conf/log.yml has been changed: %s", event.Name)
	})

	return logConfViper
}

func InitLog()  {
	logConf := loadConfig()

	logConf.SetDefault("logrotate.baseFilename", "app")
	logConf.SetDefault("logrotate.logDir", "logs")
	logConf.SetDefault("logrotate.logLevel", "info")
	logConf.SetDefault("logrotate.maxRemainCnt", 30)
	logConf.SetDefault("logrotate.rotateDays", 1)

	logDir := logConf.GetString("lograotate.logDir")
	baseFilename := logConf.GetString("lograotate.baseFilename")
	logLevel := logConf.GetString("lograotate.logLevel")
	maxRemainCnt := logConf.GetUint("lograotate.maxRemainCnt")
	rotateDays := logConf.GetUint("lograotate.rotateDays")

	if logHook, err := newRotateLogHook(logDir, baseFilename, logLevel, rotateDays, maxRemainCnt); err != nil {
		panic(fmt.Errorf("newRotateLogHook error: %v.\n", err))
	} else {
		logrus.AddHook(logHook)
	}

}


