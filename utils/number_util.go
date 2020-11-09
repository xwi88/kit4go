// Package utils common utils, number util
package utils

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
)

// Float32ToByte convert float32 to byte
func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bt := make([]byte, 4)
	binary.LittleEndian.PutUint32(bt, bits)
	return bt
}

// ByteToFloat32 convert byte to float32
func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

// Float64ToByte convert float64 to byte
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bt := make([]byte, 8)
	binary.LittleEndian.PutUint64(bt, bits)
	return bt
}

// ByteToFloat64 convert byte to float32
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func ToFixed(num float64, precision int) float64 {
	n10 := math.Pow10(precision)
	n := int(num*n10 + math.Copysign(0.5, num*n10))
	return float64(n) / n10
}

// Round round
func Round(f float64, n int) float64 {
	return ToFixed(f, n)
}

// RandomInt 随机数 int
func RandomInt(num int) int {
	return rand.Intn(65536) % num
}

// Int2Byte BigEndian
func Int2Byte(data int) ([]byte, error) {
	s1 := make([]byte, 0)
	buf := bytes.NewBuffer(s1)
	// BigEndian
	if err := binary.Write(buf, binary.BigEndian, data); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// Int2ByteLittleEndian ...
func Int2ByteLittleEndian(data int) ([]byte, error) {
	s1 := make([]byte, 0)
	buf := bytes.NewBuffer(s1)
	if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// Byte2Int BigEndian
func Byte2Int(data []byte) (int, error) {
	buf := bytes.NewBuffer(data)
	// BigEndian
	var i2 int
	if err := binary.Read(buf, binary.BigEndian, &i2); err != nil {
		return i2, err
	}
	return i2, nil
}

// Byte2IntLittleEndian ...
func Byte2IntLittleEndian(data []byte) (int, error) {
	buf := bytes.NewBuffer(data)
	var i2 int
	if err := binary.Read(buf, binary.LittleEndian, &i2); err != nil {
		return i2, err
	}
	return i2, nil
}
