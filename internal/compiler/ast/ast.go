package ast

import (
	"fmt"

	"github.com/go-serv/internal/compiler/token"
)

type Attrib interface{}

type World struct {
	Animals []Animal
}
type Animal struct {
	Name string
}

func NewWorld(animals Attrib) (*World, error) {
	s, ok := animals.([]Animal)
	if !ok {
		return nil, fmt.Errorf("%v %v %v %v", "NewWorld", "[]Animals", "animals", animals)
	}

	return &World{Animals: s}, nil
}

func NewAnimalList() ([]Animal, error) {
	return []Animal{}, nil
}

func AppendAnimals(animals, animal Attrib) ([]Animal, error) {
	s, ok := animal.(Animal)
	if !ok {
		return nil, fmt.Errorf("%v %v %v %v", "AppendAnimals", "Animal", "animal", animal)
	}
	return append(animals.([]Animal), s), nil
}

func NewAnimal(id Attrib) (Animal, error) {
	return Animal{string(id.(*token.Token).Lit)}, nil
}

func (this *World) String() string {
	return fmt.Sprintf("%v", this.Animals)
}

func (this *Animal) String() string {
	return this.Name
}
