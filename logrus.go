package main

import "github.com/aosfather/myway/core"

type logrusFactory struct {
}

func (this *logrusFactory) GetAccessLogger() core.AccessLogger {

	return nil
}
