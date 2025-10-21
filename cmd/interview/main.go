package main

import (
	"context"
	"log"
	"runtime"
	"sort"
	"sync"
	"time"
)

const callTimeout = 3 * time.Second

type HotelMatching struct {
	Hotel    string
	Supplier string
}

// Count of requests some of big number
func MethodForCreateWorkerPool(countOfRequests int) {
	workerPoolSize := runtime.GOMAXPROCS(0) * 2

	tasks := make(chan struct{}, workerPoolSize)

	work := func() {
		log.Printf("Do some work")
	}

	var wg sync.WaitGroup

	// Do work with work pool
	for range workerPoolSize {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _ = range tasks {
				work()
			}
		}()
	}

	// create payload for jobs
	for range countOfRequests {
		tasks <- struct{}{}
	}

	close(tasks)
	wg.Wait()

}

func main() {
	md := []HotelMatching{
		{
			Hotel:    "hilton",
			Supplier: "Supplier1",
		},
		{
			Hotel:    "hilton",
			Supplier: "Supplier2",
		},
		{
			Hotel:    "marriot",
			Supplier: "Supplier1",
		},
		{
			Hotel:    "holiday_inn",
			Supplier: "Supplier3",
		},
	}

	// TODO нйти все цены для отелей из md за callTimeout
	// Результат :: список цен со средней ценой
}

type SearchResult struct {
	Hotel string
	Price float64
}

// Дано
func SearchRPC(ctx context.Context, supplier string, hotels []string) ([]SearchResult, error) {
	// 0-2 seconds delay

	return []SearchResult{
		{Hotel: hotels[0], Price: 1200.0},
		{Hotel: hotels[1], Price: 10000},
	}, nil
}

func ProcessPricesForHotel(prices []float64) ([]float64, float64) {
	sum := float64(0)
	for _, v := range prices {
		sum += v
	}
	return prices, sum / float64(len(prices))
}

// End of conditions

func FindPrices(ctx context.Context, hotelsSupp []HotelMatching) {
	c, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	suppDict := make(map[string][]string)
	retInputCh := make(chan map[string][]string)

	// Gen func
	go func() {
		defer close(retInputCh)
		select {
		case <-c.Done():
			log.Printf("Timeout while generation")
			return
		default:
			for _, v := range hotelsSupp {
				// take all hotels for supp should find all hotels for each supp
				// SearchRPC(c, v.Hotel, hotels)
				if dv, ok := suppDict[v.Supplier]; ok {
					dv = append(dv, v.Hotel)
					suppDict[v.Supplier] = dv
				} else {
					suppDict[v.Supplier] = make([]string, 0, 10)
				}
			}
			retInputCh <- suppDict
		}
	}()

	hpCh := make(chan []SearchResult)

	// Func for get data
	go func() {
		defer close(hpCh)
		select {
		case <-c.Done():
			log.Printf("Timeout while processing gen data")
			return
		case mapa, ok := <-retInputCh:
			if ok {
				for k, hotels := range mapa {
					hotelsPrice, err := SearchRPC(c, k, hotels)
					if err != nil {
						log.Printf("Error while try to get hotels prices %v", err)
						return
					}
					hpCh <- hotelsPrice
				}
			} else {
				log.Printf("Error while try to read from gen data chan")
				return
			}
		}
	}()

	type NumResult struct {
		Prices []float64
		Sum    float64
	}
	result := make(map[string]NumResult)
	rCh := make(chan map[string]NumResult)

	type ResultHotels struct {
		Hotel    string
		Prices   []float64
		MidPrice float64
	}
	// Func for processing
	go func() {
		defer close(rCh)
		select {
		case <-c.Done():
			log.Printf("Timeout occured while proccessing sum ")
			return
		default:
			for hotelPrices := range hpCh {
				for _, hp := range hotelPrices {
					if v, ok := result[hp.Hotel]; ok {
						v.Prices = append(v.Prices, hp.Price)
						v.Sum += hp.Price
					} else {
						result[hp.Hotel] = NumResult{
							Prices: []float64{hp.Price},
							Sum:    hp.Price,
						}
					}
				}
			}
			rCh <- result
		}
	}()

}

// Work balancer

type Request interface{}

type Response interface{}

type Backend interface {
	Invoke(context.Context, Request) (Response, error)
}

func NewBackend(addr string) Backend {
	return _
}

type Balancer struct {
	// TODO
	as       []string
	rCounter []int
	status   []bool
}

var _ Backend = &Balancer{}

func (b *Balancer) SetStatus(index int, status bool) {
	b.status[index] = status
}

func (b *Balancer) Invoke(ctx context.Context, r Request) (Response, error) {

	// find addr with minimum requests
	index, addrForRequest := b.GetMinAddr()

	backend := NewBackend(addrForRequest)

	b.rCounter[index]++
	defer func() {
		b.rCounter[index]--
	}()

	response, err := backend.Invoke(ctx, r)

	if err != nil {
		log.Printf("Error while execute request :: %v, start waiting", r)
		b.SetStatus(index, false)
		<-time.After(3 * time.Second)
		response, err = backend.Invoke(ctx, r)
		if err != nil {
			return nil, err
		}
		b.SetStatus(index, true)
	}

	return response, nil
}

func (b *Balancer) GetMinAddr() (int, string) {

	min := 0
	index := 0
	// Can not find any good status
	for i, count := range b.rCounter {
		if count <= min && b.status[i] {
			min = count
			index = i
		}
	}

	return index, b.as[index]
}

func NewBalancer(addrs []string) *Balancer {
	var newBalancer Balancer = Balancer{
		as:       make([]string, len(addrs)),
		rCounter: make([]int, len(addrs)),
		status:   make([]bool, len(addrs)),
	}

	for i, addr := range addrs {
		newBalancer.as[i] = addr
		newBalancer.status[i] = true
	}
	return &newBalancer
}

var namesSell = []string{"Seller1", "Seller2", "seller4"}
var likesS = []int{18, 15, 16}

// out seller1 seller4 seller2

type Seller struct {
	Name  string
	Likes int
}

func sortSeller(names []string, likes []int) []string {
	sellersType := make([]Seller, len(names))
	for i, name := range names {
		sellersType[i] = Seller{
			Name:  name,
			Likes: likes[i],
		}
	}
	sort.Slice(sellersType, func(i, j int) bool {
		return sellersType[i].Likes > sellersType[j].Likes
	})

	result := make([]string, len(names))
	for i, seller := range sellersType {
		result[i] = seller.Name
	}
	return result
}

func QuickSortImplementation(xs []int) []int {
	if len(xs) <= 1 {
		return xs
	}
	left, right := 0, len(xs)-1
	midDot := len(xs) / 2

	xs[midDot], xs[right] = xs[right], xs[midDot]

	for i := range xs {
		if xs[i] < xs[right] {
			xs[i], xs[left] = xs[left], xs[i]
			left++
		}
	}

	xs[left], xs[right] = xs[right], xs[left]

	QuickSortImplementation(xs[:left])
	QuickSortImplementation(xs[left+1:])

	return xs
}
