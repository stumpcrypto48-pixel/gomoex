package handlers

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

func TestSattelliteApp(t *testing.T) {

	// t.Run("Test handler", func(t *testing.T) {
	// 	request, _ := http.NewRequest(http.MethodGet, "/v1/getSatteliteData", nil)
	// 	response := httptest.NewRecorder()

	// 	SatteliteGETHandler(response, request)
	// 	got := response.Body.String()
	// 	want := "CALSPHERE 1"

	// 	if !strings.Contains(got, want) {
	// 		t.Errorf("Got :: %v, | want :: %v", got, want)
	// 	}
	// })

}

type TestMarshalling struct {
	A int
	B string
	C string
	D string
}

func CustomMarshalingForMyStruct(structForMarsh []TestMarshalling) []byte {

	stringBuilder := strings.Builder{}

	stringBuilder.WriteString(`[ `)
	stringBuilder.Grow(20 * len(structForMarsh))
	for _, sm := range structForMarsh {
		stringBuilder.WriteString(`{"A`)
		stringBuilder.WriteString(strconv.Itoa(sm.A))
		stringBuilder.WriteString(`, "B": `)
		stringBuilder.WriteString(sm.B)
		stringBuilder.WriteString(`, "C": `)
		stringBuilder.WriteString(sm.C)
		stringBuilder.WriteString(`, "D": `)
		stringBuilder.WriteString(sm.D)
		stringBuilder.WriteString(`},`)
		// initStr += buf.String()
		stringBuilder.Reset()

	}
	stringBuilder.WriteString(` ]`)
	return []byte(stringBuilder.String())
}

func BenchmarkMarshalilng(t *testing.B) {

	numberOfObjects := 100_000
	xs := make([]TestMarshalling, 0, numberOfObjects)

	for i := range numberOfObjects {
		xs = append(xs, TestMarshalling{
			A: i,
			B: strconv.Itoa(i),
			C: strconv.Itoa(i + 1),
			D: strconv.Itoa(i + 2),
		})

	}

	t.Run("Internal marshaller bench", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			json.Marshal(xs)
		}
	})

	// t.Run("Bench fast json", func(b *testing.B) {
	// 	var parser fastjson.Object
	// 	var resultBytes []byte
	// 	b.ResetTimer()
	// 	for i := 0; i < b.N; i++ {
	// 		parser.MarshalTo(resultBytes)
	// 	}
	// })

	t.Run("Custom marshaller ", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			CustomMarshalingForMyStruct(xs)
		}
	})

	// t.Run("Encoder marshaller", func(b *testing.B) {
	// 	var buf bytes.Buffer
	// 	encoder := json.NewEncoder(buf.)
	// 	b.ResetTimer()
	// 	for i := 0; i < b.N; i++ {
	// 		encoder.Encode(xs)
	// 	}
	// })

}
