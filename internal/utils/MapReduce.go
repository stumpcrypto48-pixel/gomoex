package utils

import (
	"math"
	"sync"
)

type Stream[T any] []chan []T

type FunctionalInterface[T, M any] interface {
	Map(func(T) M) Stream[M]
	Stream(int) []chan []T
}

type StreamType[T any] struct {
	slice []T
}

func NewStreamType[T any](inputSlice []T) *StreamType[T] {
	return &StreamType[T]{
		slice: inputSlice,
	}
}

func (s *StreamType[T]) Stream(parNum int) Stream[T] {
	if parNum == 0 {
		parNum = 1
	}

	currentLen := len(s.slice)
	elementsCount := int(math.Ceil(float64(currentLen) / float64(parNum)))
	parChainList := make([]chan []T, parNum)

	for i := range parNum {
		// log.Printf("In cycle index :: %v\n", i)
		// log.Printf("Elem count :: %v, Elem count times i :: %v, Elem count times i + 1 :: %v", elementsCount, elementsCount*i, elementsCount*(i+1))
		newChan := make(chan []T, 1)
		if elementsCount*(i+1) <= currentLen {
			newChan <- s.slice[(elementsCount * i):(elementsCount * (i + 1))]
			parChainList[i] = newChan
			close(newChan)
		} else {
			newChan <- s.slice[(elementsCount * i):]
			parChainList[i] = newChan
			close(newChan)
		}
	}
	return parChainList
}

func Map[T, R any](stream Stream[T], fn func(T) R) Stream[R] {
	var result Stream[R] = make([]chan []R, len(stream))

	var wg sync.WaitGroup
	for i, ch := range stream {

		resultChan := make(chan []R)
		result[i] = resultChan
		wg.Add(1)

		go func(channel chan []T, rch chan []R) {
			defer close(rch)
			defer wg.Done()
			items := <-channel
			funcResult := make([]R, len(items))
			for i, item := range items {
				funcResult[i] = fn(item)
			}
			rch <- funcResult
		}(ch, resultChan)

	}
	go func() {
		wg.Wait()
	}()

	return result
}

func (s Stream[T]) Collect() []T {
	result := make([][]T, len(s))
	totalLen := 0
	for i, ch := range s {
		result[i] = <-ch
		totalLen += len(result[i])
	}
	totalResult := make([]T, 0, totalLen)
	for _, slice := range result {
		totalResult = append(totalResult, slice...)
	}
	return totalResult
}

func (s *Stream[T]) Filter(fn func(T) bool) {

}
