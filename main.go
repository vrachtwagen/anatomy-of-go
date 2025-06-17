package main

import (
	"fmt"

	"github.com/vrachtwagen/anatomy-of-go/context"
)

func main() {
	// Start with a cancelable parent
	parent, cancel := context.WithCancel(context.Background())
	child := context.WithoutCancel(parent)

	// Cancel the parent
	cancel()

	// child.Done() should still be nil:
	fmt.Println("child done channel:", child.Done() == nil) // true
	fmt.Println("child err:", child.Err())                  // <nil>

	// But values still flow:
	type key string
	ctx2 := context.WithValue(child, key("x"), 42)
	fmt.Println("value:", ctx2.Value(key("x")))
}
