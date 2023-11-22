[![Go Report Card](https://goreportcard.com/badge/github.com/tlkamp/go-semaphore)](https://goreportcard.com/report/github.com/tlkamp/go-semaphore)
![testworkflow](https://github.com/tlkamp/go-semaphore/actions/workflows/test.yaml/badge.svg?branch=main)

# Resizable Semaphore

A resizable semaphore for Golang.

Inspired by and possible due to https://github.com/eapache/channels. 

## Feature Overview
- Conforms to typical semaphore interfaces
- Semaphore can be resized, up or down, at runtime

## Example
```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tlkamp/go-semaphore"
)

func main() {
	sem := semaphore.New(2)
	wg := &sync.WaitGroup{}

	for i := 0; i <= 10; i++ {
		err := sem.Acquire(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		wg.Add(1)
		go func(n int) {
			defer sem.Release()
			defer wg.Done()

			time.Sleep(3 * time.Second)
			fmt.Printf("Processed: %d\n", n)
		}(i)
	}

	wg.Wait()
	fmt.Println("Complete")
}
```

## Contribute
Contributions are accepted in the form of issues and PRs.

PRs must have:
- Test cases
- Documentation
- Example (if applicable)
