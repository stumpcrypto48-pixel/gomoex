package utils

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig[T, R any] struct {
	testCase       T
	testCaseResult R
}

func TestStream(t *testing.T) {

	t.Log("Create test configuration for stream testing")
	streamTestConfig := TestConfig[[]int, int]{
		testCase:       []int{1, 2, 3, 4, 5, 6, 7},
		testCaseResult: 2,
	}
	t.Logf("Created new test case obj :: %+v\n", streamTestConfig)

	t.Run("Test stream creation ", func(t *testing.T) {
		t.Log("Trying to create new stream object")
		newStream := NewStreamType[int](streamTestConfig.testCase)
		t.Logf("Create new stream obj :: %+v\n", newStream)

		t.Log("Create new streams from object sized == 2")
		chanStreams := newStream.Stream(2)
		t.Logf("Have created new streams :: %+v", chanStreams)

		var wg sync.WaitGroup

		for _, channel := range chanStreams {
			if channel != nil {
				wg.Add(1)
				go func(ch chan []int) {
					defer wg.Done()
					data := <-ch
					fmt.Printf("Channel content :: %v\n", data)
				}(channel)
			}
		}
		t.Log("Waiting til end")
		wg.Wait()
		assert.Equal(t, streamTestConfig.testCaseResult, len(chanStreams))
	})

	t.Run("Test Map function", func(t *testing.T) {
		stream := NewStreamType(streamTestConfig.testCase).Stream(2)
		mapperResult := Map[int, int](stream, SomeTestingFuncForMapper)
		for _, chanResult := range mapperResult {
			for c := range chanResult {
				t.Logf("Channel mapper values :: %v ", c)
			}
		}
		assert.Equal(t, streamTestConfig.testCaseResult, len(mapperResult))
	})

	t.Run("Test map collect function", func(t *testing.T) {
		stream := NewStreamType(streamTestConfig.testCase).Stream(2)
		mapperResult := Map[int, int](stream, SomeTestingFuncForMapper).Collect()
		t.Logf("Mapper result :: %v", mapperResult)
		assert.Equal(t, []int{1, 4, 9, 16, 25, 36, 49}, mapperResult)
	})
}

func SomeTestingFuncForMapper(item int) int {
	return item * item
}
