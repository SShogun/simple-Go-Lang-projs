package main

import "encoding/binary"

func Encoder(timestamp int64, key string) []byte {
	keyBytes := []byte(key)
	keyLen := uint32(len(keyBytes))

	data := make([]byte, 8+4+len(keyBytes))

	binary.BigEndian.PutUint64(data[0:8], uint64(timestamp))
	binary.BigEndian.PutUint32(data[8:12], keyLen)

	copy(data[12:], keyBytes)

	return data
}
