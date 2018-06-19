package inject

import (
	"testing"
	"fmt"
	"encoding/json"
)

type Person struct {
	Name string
	Age  int
	Stu  *Stu `inject:""`
}
type Stu struct {
	Name   *string `inject:"stu_name"`
	Age    *uint32 `inject:"stu_age"`
	Status *bool   `inject:"stu_status"`
}

func (s *Stu) Add() {

}

func (p *Person) Add() {
	p.Age++
}
func (p *Person) AddPtr() {
	p.Age++
}

func TestInjector_Apply(t *testing.T) {
	injector := New()
	person := Person{}
	stu := Stu{}
	injector.Apply(&stu, &person)
	name := "XXXX"
	var age uint32 = 45
	status := true
	injector.MapString("stu_name", &name)
	injector.MapUint32("stu_age", &age)
	injector.MapBool("stu_status", &status)
	injector.Inject()
	p := injector.Service(&Person{}).(Person)
	jsonStr,_:=json.Marshal(p)
	fmt.Println(string(jsonStr))
}
