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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_IntRandom(t *testing.T) {

	t.Run("given supplier func NextInt should use supplier func to generate random values", func(t *testing.T) {
		randomValue := 5
		supplierFunc := func() (int32, error) {
			return 5, nil
		}

		intRandom := NewIntRandomWithSupplierFunc(supplierFunc)

		for i := 0; i < 20; i++ {
			value, err := intRandom.NextInt()

			assert.Nil(t, err)
			assert.Equal(t, randomValue, int(value))
		}
	})

	t.Run("given supplier func NextBytes should use supplier func to generate random values", func(t *testing.T) {
		randomValue := 10
		supplierFunc := func() (int32, error) {
			return 10, nil
		}

		intRandom := NewIntRandomWithSupplierFunc(supplierFunc)

		for i := 0; i < 20; i++ {
			var number int32 = 0

			// generate random bytes
			bytes, err := intRandom.NextBytes(INTEGER_BYTES_32)

			// converting bytes to number
			for j := 0; j < INTEGER_BYTES_32; j++ {
				number = int32(byte(number<<BYTE_SIZE) | (bytes[j] & 0xff))
			}

			assert.Nil(t, err)
			assert.Equal(t, randomValue, int(number))
		}
	})
}

func Test_ByteRandom(t *testing.T) {

	t.Run("given supplier func NextInt should use supplier func to generate random values", func(t *testing.T) {
		randomBytes := []byte{0, 0, 0, 15}
		supplierFunc := func(length int32) ([]byte, error) {
			return randomBytes, nil
		}

		byteRandom := NewByteRandomWithSupplierFunc(supplierFunc)

		for i := 0; i < 20; i++ {
			var number int32 = 0

			// generating random bytes
			bytes, err := byteRandom.NextBytes(INTEGER_BYTES_32)

			// converting bytes to number
			for j := 0; j < INTEGER_BYTES_32; j++ {
				number = int32(byte(number<<BYTE_SIZE) | (bytes[j] & 0xff))
			}

			assert.Nil(t, err)
			assert.Equal(t, number, int32(15))
		}
	})

	t.Run("given supplier func NextBytes should use supplier func to generate random values", func(t *testing.T) {
		randomBytes := []byte{0, 0, 0, 25}

		var randomValue int32 = 0
		// converting returned bytes to number
		for i := 0; i < INTEGER_BYTES_32; i++ {
			randomValue = int32(byte(randomValue<<BYTE_SIZE) | (randomBytes[i] & 0xff))
		}

		supplierFunc := func(length int32) ([]byte, error) {
			return randomBytes, nil
		}

		byteRandom := NewByteRandomWithSupplierFunc(supplierFunc)

		for i := 0; i < 20; i++ {
			var actualNumber int32 = 0

			// generate random bytes
			bytes, err := byteRandom.NextBytes(INTEGER_BYTES_32)

			// converting returned bytes to number
			for j := 0; j < INTEGER_BYTES_32; j++ {
				actualNumber = int32(byte(actualNumber<<BYTE_SIZE) | (bytes[j] & 0xff))
			}

			assert.Nil(t, err)
			assert.Equal(t, randomValue, actualNumber)
		}
	})
}

// Test_MathRandomSupplier tests the uniqueness of the random values generated
// although only 10 random values are being generated, this test may fail if the
// value generated is similar to the previous value. Failing of this test will not
// affect tsid generation logic
func Test_MathRandomSupplier(t *testing.T) {

	t.Run("GetInt should generate random values", func(t *testing.T) {
		supplier := NewMathRandomSupplier()
		var lastValue int32 = -1

		for i := 0; i < 10; i++ {

			value, err := supplier.GetInt()
			assert.Nil(t, err)
			assert.NotEqual(t, lastValue, value)

			lastValue = value

			// this will result in change of seed
			time.Sleep(time.Duration(5) * time.Millisecond)
		}
	})

	t.Run("GetBytes should generate random values", func(t *testing.T) {
		supplier := NewMathRandomSupplier()
		var lastValue int32 = -1

		for i := 0; i < 10; i++ {

			bytes, err := supplier.GetBytes(INTEGER_BYTES_32)
			assert.Nil(t, err)

			// convert bytes to number
			var value int32 = 0

			for j := 0; j < INTEGER_BYTES_32; j++ {
				value = int32(byte(value<<BYTE_SIZE) | (bytes[j] & 0xff))
			}

			assert.NotEqual(t, lastValue, value)
			lastValue = value

			// this will result in change of seed
			time.Sleep(time.Duration(5) * time.Millisecond)
		}
	})
}

// Test_CryptoRandomSupplier tests the uniqueness of the random values generated
// although only 10 random values are being generated, this test may fail if the
// value generated is similar to the previous value. Failing of this test will not
// affect tsid generation logic
func Test_CryptoRandomSupplier(t *testing.T) {

	t.Run("GetInt should generate random values", func(t *testing.T) {
		supplier := NewCryptoRandomSupplier()
		var lastValue int32 = -1

		for i := 0; i < 10; i++ {

			value, err := supplier.GetInt()
			assert.Nil(t, err)
			assert.NotEqual(t, lastValue, value)

			lastValue = value

			time.Sleep(time.Duration(5) * time.Nanosecond)
		}
	})

	t.Run("GetBytes should generate random values", func(t *testing.T) {
		supplier := NewCryptoRandomSupplier()
		var lastValue int32 = -1

		for i := 0; i < 10; i++ {

			bytes, err := supplier.GetBytes(INTEGER_BYTES_32)
			assert.Nil(t, err)

			// convert bytes to number
			var value int32 = 0

			for j := 0; j < INTEGER_BYTES_32; j++ {
				value = int32(byte(value<<BYTE_SIZE) | (bytes[j] & 0xff))
			}

			assert.NotEqual(t, lastValue, value)
			lastValue = value

			// this will result in change of seed
			time.Sleep(time.Duration(5) * time.Nanosecond)
		}
	})
}
