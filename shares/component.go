package shares

import "log"

type Parent struct {
}

func (p Parent) SayHello() {
	log.Println("Hello, I am " + p.Name())
}

func (p Parent) Name() string {
	return "Parent"
}

type Son struct {
	Parent
}

func (p Son) Name() string {
	return "Son"
}

func main() {
	var s Son
	// 面向对象继承的说法：Hello, I am Son
	// 但是在 GO 里面，组合：Hello, I am Parent
	s.SayHello()
}
