package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	typeConstructors map[string]func() any
}

func NewContainer() *Container {
	return &Container{
		typeConstructors: make(map[string]func() any, 64),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	if _, exist := c.typeConstructors[name]; exist {
		return
	}
	typeConstructor, correct := constructor.(func() any) // а что если для конструктора нужны параметры?
	if !correct {
		return
	}
	c.typeConstructors[name] = typeConstructor
}

func (c *Container) Resolve(name string) (interface{}, error) {
	constructor, exist := c.typeConstructors[name]
	if !exist {
		return nil, fmt.Errorf("no constructor registered for %s", name)
	}

	return constructor(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
