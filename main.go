package main

import (
	"fmt"
	"time"

	"github.com/vrachtwagen/anatomy-of-go/context"
)

func main() {
	fmt.Println(context.Background())
	fmt.Println(context.TODO())

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		fmt.Println("canceled:", ctx.Err())
	}()
	cancel() // your call
	time.Sleep(time.Duration(10 * time.Millisecond))

}
