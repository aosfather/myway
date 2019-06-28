package core

//应用接口
type Application interface {
	ReloadAuth()     //重新加载权限
	Restart()        //重启
	ShutdownNotify() //关闭通知
	PauseService()
	RestoreService()
	ReloadConfig()
	ReloadApis()
}
