package logl

type Config struct {
	Handler Handler
	Level   Level
	NotSafe bool // for multiple goroutine not safe
}
