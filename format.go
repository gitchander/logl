package logl

type Formatter interface {
	Format(*Record) []byte
}
