package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("a", "b", "ex@gmal.com", "pwd")
	assert.Nil(t, err)

	fmt.Printf("%+v\n", user)
}

func TestNewTask(t *testing.T) {
	task, err := NewTask("tilte", "desc")
	assert.Nil(t, err)

	fmt.Printf("%+v\n", task)
}
