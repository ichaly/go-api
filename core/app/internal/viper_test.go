package internal

import (
	"fmt"
	"testing"
)

func TestNewViper(t *testing.T) {
	v, _ := NewViper()
	fmt.Println(v.AllSettings())
}
