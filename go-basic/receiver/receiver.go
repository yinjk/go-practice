package receiver

import "fmt"

type Person struct {
	name  string
	child *Person
}

func NewPerson(name string) *Person {
	return &Person{
		name: name,
		child: &Person{
			name:  "zhangshan",
			child: nil,
		},
	}
}

func (p *Person) Name() string {
	fmt.Println(p)
	return p.name
}

func (p *Person) Child() *Person {
	return p.child
}
