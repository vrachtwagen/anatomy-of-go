package main

import (
	"fmt"
	"time"

	"github.com/vrachtwagen/anatomy-of-go/context"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	select {
	case <-ctx.Done():
		fmt.Println("timed out:", ctx.Err())
	case <-time.After(100 * time.Millisecond):
		fmt.Println("didn't time out")
	}
}
