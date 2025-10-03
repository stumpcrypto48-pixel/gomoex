package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func main() {

	// workerNumber := 10
	// chanList := make([]chan int, 0, workerNumber)
	var processing Proccessing[int, int] = func(num int) int {
		return num * num
	}

	cancel := make(chan struct{})
	paramChan := make(chan int)
	resultChan := ProcessWithCancel(processing, paramChan, cancel)

	go func() {
		defer close(paramChan)
		for i := range 100 {
			if i != 0 && i%64 == 0 {
				close(cancel)
				return
			}
			paramChan <- i
		}
	}()

	for chr := range resultChan {
		log.Printf("Value result :: %v\n", chr)
	}

	resourceSemaphore := semaphore.NewWeighted(5)
	mainContext := context.Background()
	numOfStarts := 1000
	var wg sync.WaitGroup
	wg.Add(numOfStarts)
	for i := range 1000 {
		go SemaphoreFunction(mainContext, i, resourceSemaphore, &wg)
	}
	wg.Wait()
	GracefullShutdownWithGroup(mainContext)

	httpServer := &http.Server{
		Addr: "8080",
	}

	g, gCtx := errgroup.WithContext(mainContext)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		log.Printf("exit reason :: %v", err)
	}
}

func GracefullShutdownWithGroup(ctx context.Context) {
	g, gCtx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()
	g.Go(func() error {
		for {
			select {
			case <-exit:
				log.Println("Break from syscall SIGTERM")
				return nil
			case <-gCtx.Done():
				log.Println("Break the loop")
				return nil
			case <-time.After(1 * time.Second):
				log.Println("Hello in a loop")
			}
		}
	})

	g.Go(func() error {
		for {
			select {
			case <-exit:
				log.Println("Exit second")
				return nil
			case <-gCtx.Done():
				log.Println("Break the loop")
			case <-time.After(1 * time.Second):
				log.Println("Ciao in a loop")
			}
		}
	})

	err := g.Wait()
	if err != nil {
		log.Printf("Error group :: %v\n", err)
	}
	log.Println("Done gracefull")
}

func SemaphoreFunction(ctx context.Context, data int, resourceSemaphore *semaphore.Weighted, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Trying to get data gorutine num :: %v", data)
	if err := resourceSemaphore.Acquire(ctx, 1); err != nil {
		log.Printf("Can't get data")
		return
	}
	defer resourceSemaphore.Release(1)
}

type Generate[T, G any] func(T) G
type Proccessing[G, R any] func(G) R
type WorkerType[T, R any] func(<-chan T, Proccessing[T, R], int) []<-chan R

func ProcessWithCancel[T, R any](p Proccessing[T, R], paramChan <-chan T, cancelChan <-chan struct{}) <-chan R {
	rch := make(chan R)
	go func() {
		defer close(rch)
		for param := range paramChan {
			value := p(param)
			select {
			case <-cancelChan:
				log.Println("Cancel chan")
				return
			case rch <- value:
			}
		}
	}()
	return rch
}

func Worker[T, R any](param <-chan T, action Proccessing[T, R], numOfWorkers int) []<-chan R {
	rchF := make([]chan R, 0, numOfWorkers)
	for i := range numOfWorkers {
		rchF[i] = make(chan R)
	}

	for _, ch := range rchF {
		go func(channel chan R) {
			defer close(channel)
			for data := range param {
				channel <- action(data)
			}
		}(ch)
	}

	resultChan := make([]<-chan R, 0, numOfWorkers)
	for i, ch := range rchF {
		resultChan[i] = ch
	}

	return resultChan
}

func ParallelPipeline[T, G, R any](numberOfWorkers int,
	inputParamForG <-chan T,
	generate Generate[T, G],
	worker WorkerType[G, R],
	processing Proccessing[G, R]) <-chan R {

	params := make(chan G)

	go func() {
		defer close(params)
		for param := range inputParamForG {
			params <- generate(param)
		}
	}()

	workerResult := worker(params, processing, numberOfWorkers)

	return fanIn(workerResult...)
}

func Pipeline[T, G, R any](generateFunc Generate[T, G], process Proccessing[G, R], inputParamChan <-chan T) <-chan R {
	dataForProcess := make(chan G)
	rch := make(chan R)

	go func() {
		defer close(dataForProcess)
		for data := range inputParamChan {
			dataForProcess <- generateFunc(data)
		}
	}()

	go func() {
		defer close(rch)
		for data := range dataForProcess {
			rch <- process(data)
		}
	}()

	return rch
}

func Filter[T any](inputChan <-chan T, filter func(T) bool) <-chan T {
	rch := make(chan T)

	go func() {
		defer close(rch)
		for data := range inputChan {
			if filter(data) {
				rch <- data
			}
		}
	}()

	return rch
}

func Transform[T, R any](inputChan <-chan T, action func(T) R) <-chan R {
	rch := make(chan R)

	go func() {
		defer close(rch)
		for chData := range inputChan {
			rch <- action(chData)
		}
	}()

	return rch
}

func tee[T any](inputChan <-chan T) []<-chan T {
	rchF := make([]chan T, 0, 2)

	for i := range 2 {
		newChan := make(chan T)
		rchF[i] = newChan
	}

	go func() {
		for data := range inputChan {
			for i := range 2 {
				rchF[i] <- data
			}
		}

		for _, ch := range rchF {
			close(ch)
		}

	}()

	resultChan := make([]<-chan T, 0, 2)
	for i, rch := range rchF {
		resultChan[i] = rch
	}

	return resultChan

}

func fanOut[T any](numberOfOut int, inputChan <-chan T) []<-chan T {
	rchF := make([]chan T, 0, numberOfOut)

	for i := range numberOfOut {
		newChan := make(chan T)
		rchF[i] = newChan
	}

	go func() {
		idx := 0
		for data := range inputChan {
			rchF[idx] <- data
			idx = (idx + 1) % numberOfOut
		}

		for _, ch := range rchF {
			close(ch)
		}
	}()

	resultChan := make([]<-chan T, 0, numberOfOut)
	for i, rch := range rchF {
		resultChan[i] = rch
	}
	return resultChan

}

func fanIn[T any](inputChans ...<-chan T) <-chan T {
	rch := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(inputChans))

	for _, channel := range inputChans {
		go func(inputCh <-chan T) {
			defer wg.Done()
			for n := range inputCh {
				rch <- n
			}
		}(channel)
	}

	go func() {
		defer close(rch)
		wg.Wait()
	}()

	return rch
}

func worker[T, R any](numberOfWorkers int, work func(T) R, params chan T) chan R {
	var wg sync.WaitGroup

	rch := make(chan R)
	wg.Add(numberOfWorkers)

	for range numberOfWorkers {
		go func() {
			defer wg.Done()
			for param := range params {
				rch <- work(param)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(rch)
	}()

	return rch
}
