package example

const (
	Greeter_SayHello = "/helloworld.Greeter/SayHello"
	Greeter_SayBye   = "/helloworld.Greeter/SayBye"
)

var (
	FullMethods = []string{
		Greeter_SayHello,
		Greeter_SayBye,
	}
)
