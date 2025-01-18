package isodisplay

const (
	TONS         = "tons"
	TONS_PER_MIN = "tons-per-min"
)

type Signal struct {
	Name          string
	RelativeValue int8 // Scale -100 to 100
	RawValue      any
	RawUnit       string
}

type Source interface {
	Output() <-chan Signal
	Close() error
}

type Sink interface {
	Input() chan<- Signal
	Close() error
}
