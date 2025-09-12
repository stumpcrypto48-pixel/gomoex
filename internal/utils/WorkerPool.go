package utils

import (
	"context"
	"sync"
)

func WorkerPool[T, R any](c context.Context, paramPipe chan T, work func(T) R) <-chan R {
	responseChan := make(chan R, len(paramPipe))
	var wg sync.WaitGroup

	for range len(paramPipe) {

		wg.Add(1)
		select {
		case <-c.Done():
			return responseChan
		case param, ok := <-paramPipe:
			if !ok {
				return responseChan
			}
			go func() {
				defer wg.Done()
				response := work(param)
				responseChan <- response
			}()
		}
	}

	wg.Wait()
	defer close(responseChan)

	return responseChan
}
