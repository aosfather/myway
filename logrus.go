package main

import (
	"github.com/aosfather/myway/core"
	"github.com/sirupsen/logrus"
	"os"
)

type logrusFactory struct {
	accessLogger core.AccessLogger
	accessfile   string
}

func (this *logrusFactory) Init(conf ApplicationConfigurate) {
	//访问日志
	this.accessfile = conf.AccessLogFile
	if this.accessfile == "" {
		this.accessfile = "access_log.log"
	}

}

func (this *logrusFactory) GetAccessLogger() core.AccessLogger {
	if this.accessLogger == nil {
		l := logrus.New()
		file, err := os.OpenFile(this.accessfile, os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			l.Out = file
		} else {
			l.Info("Failed to log to file, using default stderr")
		}
		this.accessLogger = &logrusImp{l}

	}
	return this.accessLogger
}

type logrusImp struct {
	log *logrus.Logger
}

//服务访问
func (this *logrusImp) ToAccess(content *core.AccessContent) {
	entry := this.log.WithFields(logrus.Fields{"remote": content.Remote})
	entry.Info(content.Url)
}

//服务访问错误信息记录
func (this *logrusImp) ToError(e *core.ErrorContent) {

}

func (this *logrusImp) WriteTextToAccess(text string) {
	this.log.Info(text)
}
