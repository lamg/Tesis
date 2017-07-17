package syncbd

type Originator interface {
	HaActual() bool
	Actual() interface{}
	Proxime()
	Remontar()
}

type Predicator interface {
	Ver(interface{}) bool
}

type Addetor interface {
	Adder(interface{})
}
