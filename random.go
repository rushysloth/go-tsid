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
	crypto_rand "crypto/rand"
	"math"
	"math/big"
	math_rand "math/rand"
	"time"
)

const (
	BYTE_SIZE        = 8
	INTEGER_SIZE_32  = 32
	INTEGER_BYTES_32 = 4
)

type Random interface {
	NextInt() (int32, error)
	NextBytes(length int32) ([]byte, error)
}

type intRandom struct {
	intSupplier     IntSupplier
	intSupplierFunc func() (int32, error)
}

func NewIntRandom(intSupplier IntSupplier) *intRandom {
	return &intRandom{
		intSupplier: intSupplier,
	}
}

func NewIntRandomWithSupplierFunc(intSupplierFunc func() (int32, error)) *intRandom {
	return &intRandom{
		intSupplierFunc: intSupplierFunc,
	}
}

func (i *intRandom) NextInt() (int32, error) {
	if i.intSupplierFunc != nil {
		return i.intSupplierFunc()
	}
	return i.intSupplier.GetInt()
}

func (i *intRandom) NextBytes(length int32) ([]byte, error) {
	bytes := make([]byte, length)

	shift := 0
	var random int32 = 0
	var err error = nil

	for j := 0; j < int(length); j++ {
		if shift < BYTE_SIZE {
			shift = INTEGER_SIZE_32

			// generate random value
			if i.intSupplierFunc != nil {
				random, err = i.intSupplierFunc()

			} else {
				random, err = i.intSupplier.GetInt()

			}
			if err != nil {
				return nil, err
			}
		}
		shift -= BYTE_SIZE
		bytes[j] = byte(uint32(random) >> shift)
	}

	return bytes, nil
}

type byteRandom struct {
	byteSupplier     ByteSupplier
	byteSupplierFunc func(length int32) ([]byte, error)
}

func NewByteRandom(randomSupplier RandomSupplier) *byteRandom {
	return &byteRandom{
		byteSupplier: randomSupplier,
	}
}

func NewByteRandomWithSupplierFunc(randomSupplierFunc func(length int32) ([]byte, error)) *byteRandom {
	return &byteRandom{
		byteSupplierFunc: randomSupplierFunc,
	}
}

func (i *byteRandom) NextInt() (int32, error) {
	var number int32 = 0
	var bytes []byte
	var err error

	if i.byteSupplierFunc != nil {
		bytes, err = i.byteSupplierFunc(INTEGER_SIZE_32)

	} else {
		bytes, err = i.byteSupplier.GetBytes(INTEGER_SIZE_32)

	}
	if err != nil {
		return int32(number), err
	}

	for j := 0; j < INTEGER_BYTES_32; j++ {
		number = int32(byte(number<<BYTE_SIZE) | (bytes[j] & 0xff))
	}
	return number, nil
}

func (i *byteRandom) NextBytes(length int32) ([]byte, error) {
	if i.byteSupplierFunc != nil {
		return i.byteSupplierFunc(length)
	}
	return i.byteSupplier.GetBytes(length)
}

// Suppliers
type RandomSupplier interface {
	GetInt() (int32, error)
	GetBytes(length int32) ([]byte, error)
}

type IntSupplier interface {
	GetInt() (int32, error)
}

type ByteSupplier interface {
	GetBytes(length int32) ([]byte, error)
}

type mathRandomSupplier struct {
}

func NewMathRandomSupplier() *mathRandomSupplier {
	return &mathRandomSupplier{}
}

func (i *mathRandomSupplier) GetInt() (int32, error) {
	rand := math_rand.New(
		math_rand.NewSource(
			time.Now().UnixNano()))

	return rand.Int31(), nil
}

func (i *mathRandomSupplier) GetBytes(length int32) ([]byte, error) {
	rand := math_rand.New(
		math_rand.NewSource(
			time.Now().UnixNano()))

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)

	return bytes, err
}

type cryptoRandomSupplier struct {
}

func NewCryptoRandomSupplier() *cryptoRandomSupplier {
	return &cryptoRandomSupplier{}
}

func (i *cryptoRandomSupplier) GetInt() (int32, error) {
	random, err := crypto_rand.Int(crypto_rand.Reader, big.NewInt(math.MaxInt32))
	return int32(random.Int64()), err
}

func (i *cryptoRandomSupplier) GetBytes(length int32) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := crypto_rand.Read(bytes)

	return bytes, err
}
