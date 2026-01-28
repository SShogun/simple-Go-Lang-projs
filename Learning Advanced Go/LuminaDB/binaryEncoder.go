package main

import (
	"encoding/binary"
	"fmt"
)

func Encoder(timestamp int64, key string) []byte {
	keyBytes := []byte(key)
	keyLen := uint32(len(keyBytes))

	data := make([]byte, 8+4+len(keyBytes))

	binary.BigEndian.PutUint64(data[0:8], uint64(timestamp))
	binary.BigEndian.PutUint32(data[8:12], keyLen)

	copy(data[12:], keyBytes)

	return data
}

func Decoder(data []byte) (int64, string, error) {
	if len(data) < 12 {
		return 0, "", fmt.Errorf("data too short to decode")
	}

	timestamp := binary.BigEndian.Uint64(data[0:8])
	keyLen := binary.BigEndian.Uint32(data[8:12])

	if len(data) < 12+int(keyLen) {
		return 0, "", fmt.Errorf("Insufficient data for key payload")
	}
	key := string(data[12 : 12+keyLen])
	return int64(timestamp), key, nil
}
