package filter

/**
  lua filter implements
  腳本filter，對腳本提供操作header及parameter和設置的能力

*/
type LuaFilter struct {
	r       EntityReader
	w       EntityWriter
	context map[string]interface{}
}
