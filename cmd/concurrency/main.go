package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"httpfromtcp/rootmod/internal/model"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func main() {

	httpServer := &http.Server{
		Addr: ":8080",
	}

	mainCtx := context.Background()

	g, gCtx := errgroup.WithContext(mainCtx)

	g.Go(func() error {
		return httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(mainCtx)
	})

	http.HandleFunc("/v1/getTrades", func(response http.ResponseWriter, request *http.Request) {
		rCtx := request.Context()
		ticker := request.URL.Query().Get("ticker")
		log.Printf("Input ticker :: %v\n", ticker)
		responseChan := make(chan []byte)

		go func() {
			defer close(responseChan)
			resp, err := GetPolyDataTrades(rCtx, ticker)
			if err != nil {
				response.Write([]byte(fmt.Sprintf("Error while try to get data :: %v", err)))
				return
			}
			responseChan <- resp
		}()

		// resp := Worker(responseChan, CalcMidPrice, 4)
		// var totalResponse []byte

		// outCh := fanIn(resp)
		// for ch := range outCh {
		// 	jsonResponse, err := json.Marshal(ch)
		// 	if err != nil {
		// 		response.Write([]byte(fmt.Sprintf("Error while marshall data into json :: %v", err)))
		// 		return
		// 	}
		// 	totalResponse = jsonResponse
		// }
		// response.Write(totalResponse)
		for respData := range responseChan {
			// var encodedData model.TradeModel
			// err := xml.Unmarshal(respData, &encodedData)
			// if err != nil {
			// 	response.Write([]byte(fmt.Sprintf("Error while try to unmarshal data :: %v", err)))
			// }
			// log.Printf("All data response :: %+v \n", encodedData.Data.Rows.Items)

			midPrice := CalcMidPrice(respData)

			jsonResponse, err := json.Marshal(midPrice)
			if err != nil {
				response.Write([]byte(fmt.Sprintf("Error while try to encode data into json :: %v", err)))
				return
			}
			response.Write(jsonResponse)
		}

	})

	http.HandleFunc("/v1/pipeline", func(response http.ResponseWriter, request *http.Request) {
		rCtx := request.Context()
		paramChan := make(chan int)
		go func() {
			for i := range 100 {
				paramChan <- i
			}
			close(paramChan)
		}()

		var generatorFunc Generate[int, int] = func(param int) int {
			return param * param
		}
		var processFunc Proccessing[int, string] = func(input int) string {
			if input > 500 {
				return "Ok :: " + strconv.Itoa(input)
			} else {
				return "Not enougth " + strconv.Itoa(input)
			}
		}
		rCh := GracefullShutDownPipeline(rCtx, paramChan, generatorFunc, processFunc, Worker)

		for _, ch := range rCh {
			for data := range ch {
				log.Printf("Chan data :: %v \n", data)
				response.Write([]byte(data + " \n"))
			}
		}
		log.Println("End work")
	})

	if err := g.Wait(); err != nil {
		log.Printf("Exit reason :: %v \n", err)
	}
}

func StartApp(ctx context.Context, handlers http.HandlerFunc) {
	mainCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: handlers,
		BaseContext: func(_ net.Listener) context.Context {
			return mainCtx
		},
	}

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		fmt.Println("Exit reason ")
	}
}

func GetPolyDataTrades(ctx context.Context, ticker string) ([]byte, error) {

	cont, cancel := context.WithTimeout(ctx, 3*time.Second)
	respChan := make(chan []byte)
	defer cancel()

	go func() {
		defer close(respChan)
		moexAddres := fmt.Sprintf("http://iss.moex.com/iss/engines/stock/markets/shares/securities/%v/candles.xml", ticker)
		req, err := http.NewRequestWithContext(cont, "GET", moexAddres, nil)
		if err != nil {
			fmt.Errorf("Error while try to get request to :: %v \n | err :: %w", moexAddres, err)
		}

		q := req.URL.Query()
		q.Add("from", "2023-05-25")
		q.Add("till", "2025-09-01")
		q.Add("interval", "24")

		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error while request")
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error while read data from response body")
			return
		}
		respChan <- body
	}()

	select {
	case <-ctx.Done():
		log.Printf("Timeout exit")
		return nil, fmt.Errorf("Timeout")
	case respData, ok := <-respChan:
		if ok {
			return respData, nil
		} else {
			return nil, fmt.Errorf("Error while try to get data from channel")
		}
	}

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

func GracefullShutDownPipeline[T, G, R any](c context.Context,
	paramsForG <-chan T,
	generator Generate[T, G],
	processFunc Proccessing[G, R],
	worker WorkerType[G, R]) []<-chan R {

	parNum := 4
	ctx, stop := signal.NotifyContext(c, syscall.SIGTERM)

	var wg sync.WaitGroup
	genParam := make(chan G)

	wg.Add(1)
	go func() {
		defer stop()
		defer wg.Done()
		defer close(genParam)
		for {
			select {
			case <-ctx.Done():
				log.Println("End work after stop by signal")
				return
			case paramForG, ok := <-paramsForG:
				if ok {
					genParam <- generator(paramForG)
				} else {
					log.Println("Input param channel closed stop generator goroutine")
					return
				}
			}
		}
	}()

	return worker(genParam, processFunc, parNum)
}

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
	for range numOfWorkers {
		newCh := make(chan R)
		rchF = append(rchF, newCh)
	}

	for i, ch := range rchF {

		go func(idx int, channel chan R) {
			defer log.Printf("Closed worker chan :: %v \n", idx)
			defer close(channel)
			for data := range param {
				channel <- action(data)
			}
		}(i, ch)
	}

	resultChan := make([]<-chan R, 0, numOfWorkers)
	for _, ch := range rchF {
		resultChan = append(resultChan, ch)
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

	return fanIn(workerResult)
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

func fanIn[T any](inputChans []<-chan T) <-chan T {
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

type TotalResponseTrades struct {
	MidPrices []ResponseTrades `json:"trades"`
	MeanPrice float64          `json:"meanPrice"`
	Variance  float64          `json:"variance"`
}

type ResponseTrades struct {
	MidPrice   float64 `json:"midPrice"`
	SquareDiff float64 `json:"squareDiff"`
	Begin      string  `json:"Begin"`
	End        string  `json:"End"`
}

func CalcVariance(mean float64, xs []float64) float64 {
	allDist := 0.0
	for _, v := range xs {
		allDist += (mean - v) * (mean - v)
	}
	return allDist / float64(len(xs)-1)
}

func CalcMidPrice(response []byte) TotalResponseTrades {
	var tradeModel model.TradeModel
	err := xml.Unmarshal(response, &tradeModel)
	if err != nil {
		fmt.Errorf("Err while try to unmarshall response")
		return TotalResponseTrades{}
	}
	rows := tradeModel.Data.Rows.Items
	rs := make([]ResponseTrades, 0, len(rows))
	midPrice := make([]float64, 0, len(rows))
	mean := 0.0

	for i, row := range rows {
		midPrice = append(midPrice, (row.LowPrice+row.HighPrice)/2)
		mean += float64(midPrice[i])
	}
	mean = mean / float64(len(midPrice))

	for i, row := range rows {
		rs = append(rs, ResponseTrades{
			MidPrice:   midPrice[i],
			SquareDiff: (midPrice[i] - mean) * (midPrice[i] - mean),
			Begin:      row.StartDate,
			End:        row.EndDate,
		})
	}

	return TotalResponseTrades{
		MeanPrice: mean,
		Variance:  CalcVariance(mean, midPrice),
		MidPrices: rs,
	}
}
