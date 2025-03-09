/*
Copyright (c) 2023

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tsid

import (
	"sync/atomic"
	"time"
)

const (
	TSID_EPOCH     int64 = 1672531200000 // 2023-01-01T00:00:00.000Z
	TSID_BYTES     int32 = 8
	TSID_CHARS     int32 = 13 // ToString returns a string of length 13
	RANDOM_BITS    int32 = 22
	RANDOM_MASK    int32 = 0x003fffff
	NODE_BITS_1024 int32 = 10
)

var ALPHABET_UPPERCASE []rune = []rune("0123456789ABCDEFGHJKMNPQRSTVWXYZ")
var ALPHABET_LOWERCASE []rune = []rune("0123456789abcdefghjkmnpqrstvwxyz")
var ALPHABET_VALUES []int64

var atomicCounter atomic.Uint32

func init() {
	ALPHABET_VALUES = make([]int64, 128)
	for i := 0; i < len(ALPHABET_VALUES); i++ {
		ALPHABET_VALUES[i] = -1
	}

	// Numbers
	ALPHABET_VALUES['0'] = 0x00
	ALPHABET_VALUES['1'] = 0x01
	ALPHABET_VALUES['2'] = 0x02
	ALPHABET_VALUES['3'] = 0x03
	ALPHABET_VALUES['4'] = 0x04
	ALPHABET_VALUES['5'] = 0x05
	ALPHABET_VALUES['6'] = 0x06
	ALPHABET_VALUES['7'] = 0x07
	ALPHABET_VALUES['8'] = 0x08
	ALPHABET_VALUES['9'] = 0x09

	ALPHABET_VALUES['a'] = 0x0a
	ALPHABET_VALUES['b'] = 0x0b
	ALPHABET_VALUES['c'] = 0x0c
	ALPHABET_VALUES['d'] = 0x0d
	ALPHABET_VALUES['e'] = 0x0e
	ALPHABET_VALUES['f'] = 0x0f
	ALPHABET_VALUES['g'] = 0x10
	ALPHABET_VALUES['h'] = 0x11
	ALPHABET_VALUES['j'] = 0x12
	ALPHABET_VALUES['k'] = 0x13
	ALPHABET_VALUES['m'] = 0x14
	ALPHABET_VALUES['n'] = 0x15
	ALPHABET_VALUES['p'] = 0x16
	ALPHABET_VALUES['q'] = 0x17
	ALPHABET_VALUES['r'] = 0x18
	ALPHABET_VALUES['s'] = 0x19
	ALPHABET_VALUES['t'] = 0x1a
	ALPHABET_VALUES['v'] = 0x1b
	ALPHABET_VALUES['w'] = 0x1c
	ALPHABET_VALUES['x'] = 0x1d
	ALPHABET_VALUES['y'] = 0x1e
	ALPHABET_VALUES['z'] = 0x1f

	ALPHABET_VALUES['i'] = 0x01
	ALPHABET_VALUES['l'] = 0x01
	ALPHABET_VALUES['o'] = 0x00

	ALPHABET_VALUES['A'] = 0x0a
	ALPHABET_VALUES['B'] = 0x0b
	ALPHABET_VALUES['C'] = 0x0c
	ALPHABET_VALUES['D'] = 0x0d
	ALPHABET_VALUES['E'] = 0x0e
	ALPHABET_VALUES['F'] = 0x0f
	ALPHABET_VALUES['G'] = 0x10
	ALPHABET_VALUES['H'] = 0x11
	ALPHABET_VALUES['J'] = 0x12
	ALPHABET_VALUES['K'] = 0x13
	ALPHABET_VALUES['M'] = 0x14
	ALPHABET_VALUES['N'] = 0x15
	ALPHABET_VALUES['P'] = 0x16
	ALPHABET_VALUES['Q'] = 0x17
	ALPHABET_VALUES['R'] = 0x18
	ALPHABET_VALUES['S'] = 0x19
	ALPHABET_VALUES['T'] = 0x1a
	ALPHABET_VALUES['V'] = 0x1b
	ALPHABET_VALUES['W'] = 0x1c
	ALPHABET_VALUES['X'] = 0x1d
	ALPHABET_VALUES['Y'] = 0x1e
	ALPHABET_VALUES['Z'] = 0x1f

	ALPHABET_VALUES['I'] = 0x01
	ALPHABET_VALUES['L'] = 0x01
	ALPHABET_VALUES['O'] = 0x00
}

type Tsid struct {
	number int64
}

// NewTsid returns pointer to new tsid
func NewTsid(number int64) *Tsid {
	return &Tsid{
		number: number,
	}
}

// Fast returns a pointer to new random tsid
func Fast() *Tsid {
	// Incrementing before using it
	cnt := atomicCounter.Add(1)

	time := (time.Now().UnixMilli() - TSID_EPOCH) << RANDOM_BITS
	tail := cnt & uint32(RANDOM_MASK)

	return NewTsid(time | int64(tail))
}

// FromNumber returns pointer to tsid using the given number
func FromNumber(number int64) *Tsid {
	return NewTsid(number)
}

// FromBytes returns pointer to tsid by converting the given bytes to
// number
func FromBytes(bytes []byte) *Tsid {

	// TODO: Add validation

	var number int64 = 0

	number |= int64(bytes[0]&0xff) << 56
	number |= int64(bytes[1]&0xff) << 48
	number |= int64(bytes[2]&0xff) << 40
	number |= int64(bytes[3]&0xff) << 32
	number |= int64(bytes[4]&0xff) << 24
	number |= int64(bytes[5]&0xff) << 16
	number |= int64(bytes[6]&0xff) << 8
	number |= int64(bytes[7]) & 0xff

	return NewTsid(int64(number))
}

// FromString returns pointer to tsid by converting the given string to
// number. It validates the string before conversion.
func FromString(str string) *Tsid {
	arr := ToRuneArray(str)

	var number int64 = 0

	number |= ALPHABET_VALUES[arr[0]] << 60
	number |= ALPHABET_VALUES[arr[1]] << 55
	number |= ALPHABET_VALUES[arr[2]] << 50
	number |= ALPHABET_VALUES[arr[3]] << 45
	number |= ALPHABET_VALUES[arr[4]] << 40
	number |= ALPHABET_VALUES[arr[5]] << 35
	number |= ALPHABET_VALUES[arr[6]] << 30
	number |= ALPHABET_VALUES[arr[7]] << 25
	number |= ALPHABET_VALUES[arr[8]] << 20
	number |= ALPHABET_VALUES[arr[9]] << 15
	number |= ALPHABET_VALUES[arr[10]] << 10
	number |= ALPHABET_VALUES[arr[11]] << 5
	number |= ALPHABET_VALUES[arr[12]]

	return NewTsid(int64(number))
}

// ToRuneArray converts the given string to rune array. It also performs
// validations on the rune array
func ToRuneArray(str string) []rune {
	arr := []rune(str)

	if !IsValidRuneArray(arr) {
		return nil // TODO: Throw error
	}
	return arr
}

// IsValidRuneArray validates the rune array.
func IsValidRuneArray(arr []rune) bool {

	if arr == nil || len(arr) != int(TSID_CHARS) {
		return false
	}

	if (ALPHABET_VALUES[arr[0]] & 0b10000) != 0 {
		return false
	}

	for i := 0; i < len(arr); i++ {
		if ALPHABET_VALUES[arr[i]] == -1 {
			return false
		}
	}
	return true
}

// ToNumber returns the numerical component of the tsid
func (t *Tsid) ToNumber() int64 {
	return t.number
}

// ToBytes converts the number to bytes and returns the byte array
func (t *Tsid) ToBytes() []byte {
	bytes := make([]byte, TSID_BYTES)

	bytes[0] = byte(uint64(t.number) >> 56)
	bytes[1] = byte(uint64(t.number) >> 48)
	bytes[2] = byte(uint64(t.number) >> 40)
	bytes[3] = byte(uint64(t.number) >> 32)
	bytes[4] = byte(uint64(t.number) >> 24)
	bytes[5] = byte(uint64(t.number) >> 16)
	bytes[6] = byte(uint64(t.number) >> 8)
	bytes[7] = byte(t.number)

	return bytes
}

// ToString converts the number to a canonical string.
// The output is 13 characters long and only contains characters from
// Crockford's base32 alphabets
func (t *Tsid) ToString() string {
	return t.ToStringWithAlphabets(ALPHABET_UPPERCASE)
}

// ToLowerCase converts the number to a canonical string in lower case.
// The output is 13 characters long and only contains characters from
// Crockford's base32 alphabets
func (t *Tsid) ToLowerCase() string {
	return t.ToStringWithAlphabets(ALPHABET_LOWERCASE)
}

// ToStringWithAlphabets converts the number to string using the given alphabets and returns it
func (t *Tsid) ToStringWithAlphabets(alphabets []rune) string {
	chars := make([]rune, TSID_CHARS)

	chars[0] = alphabets[((uint64(t.number) >> 60) & 0b11111)]
	chars[1] = alphabets[((uint64(t.number) >> 55) & 0b11111)]
	chars[2] = alphabets[((uint64(t.number) >> 50) & 0b11111)]
	chars[3] = alphabets[((uint64(t.number) >> 45) & 0b11111)]
	chars[4] = alphabets[((uint64(t.number) >> 40) & 0b11111)]
	chars[5] = alphabets[((uint64(t.number) >> 35) & 0b11111)]
	chars[6] = alphabets[((uint64(t.number) >> 30) & 0b11111)]
	chars[7] = alphabets[((uint64(t.number) >> 25) & 0b11111)]
	chars[8] = alphabets[((uint64(t.number) >> 20) & 0b11111)]
	chars[9] = alphabets[((uint64(t.number) >> 15) & 0b11111)]
	chars[10] = alphabets[((uint64(t.number) >> 10) & 0b11111)]
	chars[11] = alphabets[((uint64(t.number) >> 5) & 0b11111)]
	chars[12] = alphabets[(uint64(t.number) & 0b11111)]

	return string(chars)
}

// IsValid checks if the given tsid string is valid or not
func (t *Tsid) IsValid(str string) bool {
	return len(str) != 0 && IsValidRuneArray([]rune(str))
}

// GetRandom returns random component (node + counter) of the tsid
func (t *Tsid) GetRandom() int64 {
	return t.number & int64(RANDOM_MASK)
}

// GetUnixMillis returns time of creation in millis since 1970-01-01
func (t *Tsid) GetUnixMillis() int64 {
	return t.getTime() + TSID_EPOCH
}

// GetUnixMillis returns time of creation in millis since 1970-01-01
func (t *Tsid) GetUnixMillisWithCustomEpoch(epoch int64) int64 {
	return t.getTime() + epoch
}

// getTime returns the time component
func (t *Tsid) getTime() int64 {
	return int64(uint64(t.number) >> int64(RANDOM_BITS))
}
