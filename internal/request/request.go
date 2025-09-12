package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	allLines, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var requestLine *RequestLine
	var errLine error
	firstLine := strings.Split(string(allLines), "\r\n")[0]
	requestLine, errLine = parseRequestLine(firstLine)
	if errLine != nil {
		return nil, errLine
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	fmt.Printf("Request for parsing :: %v \n", line)

	var requestLineResult RequestLine

	splitedRequest := strings.Split(line, " ")

	if len(splitedRequest) != 3 {
		return nil, fmt.Errorf("invalid request line")
	}
	var errV error

	requestLineResult.HttpVersion, errV = getVersionFromString(splitedRequest[2])
	if errV != nil {
		return nil, errV
	}
	requestLineResult.RequestTarget = splitedRequest[1]
	requestLineResult.Method = splitedRequest[0]

	return &requestLineResult, nil
}

func getVersionFromString(version string) (string, error) {
	allVersions := strings.Split(version, "/")
	fmt.Printf("All versions :: %v \n", allVersions)
	if len(allVersions) != 2 {
		return "", fmt.Errorf("invalid version format")
	}
	return allVersions[1], nil
}
