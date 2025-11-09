package app

func Greeting(name string) string {
	if name == "" {
		return "hello, world"
	}
	return "hello, " + name
}
