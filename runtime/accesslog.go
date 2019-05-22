package runtime

import "github.com/aosfather/myway/core"

/*
  access 日志输出
  1、access 日志
  2、error 日志

*/

var access_log core.AccessLogger

func SetAccessLogger(log core.AccessLogger) {
	access_log = log
}

func Log(text string) {
	if access_log != nil {
		access_log.WriteTextToAccess(text)
	}

}
