package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	error1    = errors.New("error1")
	error2    = errors.New("error2")
	errorList = []error{error1, error2}
)

func main() {
	ctx := context.Background()
	firstChild(ctx)

	first := ctx.Value(KeyValueType("firstChildKey"))
	second := ctx.Value(KeyValueType("secondChildKey"))
	fmt.Printf("first context value :: %v | second context value :: %v\n", first, second)

	someSlice := ComplexValueType{KeyValueType("some"), KeyValueType("asdfads"), KeyValueType("adsfasdf"), KeyValueType("adsfasdf")}
	someValue := KeyValueType("adsfasdf")
	num := badGeneric(someSlice, someValue)
	fmt.Printf("Value %v presented in slice :: %v times\n", someValue, num)

	child := ChildStruct{
		lastName: "string",
	}
	fmt.Print(child.Void())
	fmt.Print(child.Name)

	// SenderReciver()
	CallCustomError()
	apiError := ApiErrorStruct[error, error]{
		errorWrapper:    error2,
		targetErrorList: errorList,
	}

	apiError.ReturnError()

	var someType interface{}
	var pointerType *string
	fmt.Printf("String pointer eq :: %v", pointerType == nil)
	someTypeIsNil(someType)
	someTypeIsNil(pointerType)
	someType = nil
	someTypeIsNil(someType)
}

func someTypeIsNil(someType interface{}) {
	fmt.Printf("Some value is nil ? %v :: %v, %T\n", someType == nil, someType, someType)
}

type KeyValueType string

type ComplexValueType []KeyValueType
type WrongComplexValueType []string

func firstChild(c context.Context) {
	ctx := context.WithValue(c, KeyValueType("firstChildKey"), 10)
	fmt.Printf("Ctx in first child :: %v\n", ctx.Value(KeyValueType("firstChildKey")))
	secondChild(ctx)
	fmt.Printf("Ctx second in firts child :: %v\n", ctx.Value(KeyValueType("secondChildKey")))
}

func secondChild(c context.Context) {
	ctx := context.WithValue(c, KeyValueType("secondChildKey"), 20)
	fmt.Printf("Ctx in second child :: %v\n", ctx.Value(KeyValueType("secondChildKey")))
	fmt.Printf("Ctx first in second child :: %v\n", ctx.Value(KeyValueType("firstChildKey")))
}

func badGeneric[A ~[]R, R comparable](a A, r R) int {
	counter := 0
	for i, item := range a {
		fmt.Printf("Item in array index :: %v value :: %v type :: %T equal input  :: %v\n",
			i, item, item, item == r)
		if item == r {
			counter++
		}
	}
	return counter
}

type FlyMeToTheMoon interface {
	Fly() string
	Void() string
}

type ParentStruct struct {
	Pointer uintptr
	Age     int
	Name    string
}

func (p *ParentStruct) Void() string {
	return fmt.Sprintf("Fly me to the moon with type :: %T\n", *p)
}

func (p *ParentStruct) Fly() string {
	return fmt.Sprintf("Fly %v to the moon with %v age", p.Name, p.Age)
}

// struct embedding
type ChildStruct struct {
	ParentStruct
	lastName string
}

const (
	ErrorTypCustom = "FirstCustomError"
)

func SenderReciver() {
	genChan := make(chan string, 10)
	getChan := make(chan string, 10)
	closeChan := make(chan struct{})
	// closing := make(chan struct{})

	// stop := func() {
	// 	select {
	// 	case closing <- struct{}{}:
	// 		<-closeChan
	// 	case <-closeChan:
	// 	}
	// }
	// go func() {
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			return
	// 		}
	// 	}()
	// 	for {
	// 		select {
	// 		case _, ok := <-closeChan:
	// 			if !ok {
	// 				return
	// 			}
	// 		default:
	// 			genChan <- "input string"
	// 		}
	// 	}
	// }()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		for {
			select {
			case _, ok := <-closeChan:
				if !ok {
					return
				}
			default:
				genChan <- "second input string"
			}
		}
	}()

	// close manipulation
	go func() {
		defer func() {
			fmt.Println("Trying to close close chan")
			close(genChan)
		}()
		for _ = range closeChan {
			fmt.Println("Into loop of close chan")
		}
		fmt.Println("End close chan loop ")
	}()

	go func() {
		defer close(getChan)
		for data := range genChan {
			getChan <- data
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer func() {
			close(closeChan)
			wg.Done()
		}()
		for _ = range 100 {
			getData := <-getChan
			fmt.Printf("Get next data :: %v\n", getData)
		}

	}()
	wg.Wait()

}

func safeClose[T any](chanToClose chan T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Trying to close already closed chan")
			return
		}
	}()
	close(chanToClose)
}

type CustomError1 struct {
	errorType string
	err       error
}

func (e CustomError1) Error() string {
	return fmt.Sprintf("First error :: %v", e.errorType)
}

func (e CustomError1) Unwrap() error {
	return e.err
}

func CallCustomError() {
	err1 := CustomError1{
		errorType: "err1",
	}

	var errNotFound error = errors.New("not found")

	wrappedErr := fmt.Errorf("Error wrapping %v : %w", "some", err1)
	nextWrap := fmt.Errorf("Error 2 :: %v :: %w", wrappedErr, errNotFound)

	var err CustomError1
	fmt.Printf("As error err1 :: %v\n", errors.As(nextWrap, &err))
	fmt.Printf("Is error err1 :: %v\n", errors.Is(wrappedErr, &err))
}

type ApiError[T any] interface {
	ReturnError() error
}

type ApiErrorStruct[T, E error] struct {
	errorWrapper    T
	targetErrorList []E
}

func (e ApiErrorStruct[T, E]) ReturnError() {
	for _, err := range e.targetErrorList {
		if errors.Is(e.errorWrapper, err) {
			fmt.Printf("Found next error :: %v\n", err)
		}
	}

}

type A struct {
}

func (a *A) fn()  {}
func (a *A) fn2() {}

type TestInter1 interface {
	string | int
	fn()
}

type TestInter2 interface {
	string | int
	fn2()
}

func SuperFunc[T TestInter1](input T) {

}

type PipelineInterface[P, R any] interface {
	InsertParam(chan<- P) <-chan R
	Merge([]<-chan P) <-chan P
}

type PipelineStruct[P, R any] struct {
	parameters P
	resultChan chan R
}

// func (ps PipelineStruct[P, R]) InsertParam(inputChan <-chan P, operationWithParam func()) <-chan R {
// 	out := make(chan R)
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		for inputData := range inputChan {
// 		}
// 	}()

// }

func merge[T any](channels []<-chan T) <-chan T {
	out := make(chan T, len(channels))
	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, c := range channels {
		go func(<-chan T) {
			defer wg.Done()
			for val := range c {
				out <- val
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func Pipeline[T, R, O any](paramCreator func() T,
	inputParamGen func(T) <-chan T,
	outputParamWorker func(<-chan T) <-chan R,
	manipulateResults func(<-chan R) O) O {
	// create chans for input data
	inputDataChan := inputParamGen(paramCreator())
	// create chans for ouput data
	chanSize := 10
	outputDataChan := make([]<-chan R, chanSize)
	for i := range chanSize {
		outputDataChan[i] = outputParamWorker(inputDataChan)
	}
	// work with input with input func()
	resultChan := merge(outputDataChan)

	// join output of func()
	// return merged result
	return manipulateResults(resultChan)
}
