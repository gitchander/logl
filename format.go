package logl

type Format interface {
	Format(*Record) []byte
}
