package core

//应用接口
type Application interface {
	ShutDown()
	PauseService()
	RestoreService()
	ReloadConfig()
	ReloadApis()
}
