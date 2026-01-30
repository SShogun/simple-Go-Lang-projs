package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type RespParser struct {
	reader *bufio.Reader
}

func NewRespParser(rd io.Reader) *RespParser {
	return &RespParser{reader: bufio.NewReader(rd)}
}

func (p *RespParser) Parse() ([]string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if line[0] != '*' {
		return nil, fmt.Errorf("expect '*' at start of array")
	}

	count, _ := strconv.Atoi(line[1 : len(line)-2])
	args := make([]string, count)

	for i := 0; i < count; i++ {
		args[i], err = p.readBulkString()
		if err != nil {
			return nil, err
		}
	}

	return args, nil
}

func (p *RespParser) readBulkString() (string, error) {
	line, _ := p.reader.ReadString('\n') // Read the "$3\r\n"
	if line[0] != '$' {
		return "", fmt.Errorf("expected '$' for bulk string")
	}

	size, _ := strconv.Atoi(line[1 : len(line)-2])

	// Read exactly 'size' bytes + the 2 bytes for \r\n
	data := make([]byte, size+2)
	io.ReadFull(p.reader, data)

	return string(data[:size]), nil
}
