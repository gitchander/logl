package logl

type Level int

const (
	LEVEL_OFF Level = iota

	LEVEL_PANIC   // log: [ Panic ]                  - before write, flush and after call panic(message)
	LEVEL_FATAL   // log: [ Panic, Fatal ]           - before write, flush and after call os.Exit(1)
	LEVEL_ERROR   // log: [ Panic, Fatal, Error ]
	LEVEL_WARNING // log: [ Panic, Fatal, Error, Warning ]
	LEVEL_INFO    // log: [ Panic, Fatal, Error, Warning, Info ]
	LEVEL_DEBUG   // log: [ Panic, Fatal, Error, Warning, Info, Debug ]
)
