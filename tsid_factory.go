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
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// Lock will be used to control access for creating tsidFactory instance
var lock = &sync.Mutex{}

// Lock will be used to control synchronize access to random value generator
var rLock = &sync.Mutex{}

// Only a single instance of tsidFactory will be used per node
var tsidFactoryInstance *TsidFactory

// tsidFactory is a singleton which
// should be used to generate random tsid
type TsidFactory struct {
	node        int32
	nodeBits    int32
	nodeMask    int32
	counter     int32
	counterBits int32
	counterMask int32
	lastTime    int64
	customEpoch int64
	clock       Clock
	random      Random
	randomBytes int32
}

func newTsidFactory(builder *tsidFactoryBuilder) (*TsidFactory, error) {

	// properties from builder
	tsidFactory := &TsidFactory{
		customEpoch: builder.GetCustomEpoch(),
		clock:       builder.GetClock(),
		random:      builder.GetRandom(),
	}

	// get node bits
	nodeBits, err := builder.GetNodeBits()
	if err != nil {
		log.Print(err.Error())
		return nil, errors.New("failed to initialize tsid factory")
	}
	tsidFactory.nodeBits = nodeBits

	// properties to be calculated
	tsidFactory.counterBits = int32(RANDOM_BITS) - builder.nodeBits
	tsidFactory.counterMask = int32(uint32(RANDOM_MASK) >> builder.nodeBits)
	tsidFactory.nodeMask = int32(uint32(RANDOM_MASK) >> tsidFactory.counterBits)

	tsidFactory.randomBytes = ((tsidFactory.counterBits - 1) / 8) + 1

	// get node id
	node, err := builder.GetNode()
	if err != nil {
		log.Print(err.Error())
		return nil, errors.New("failed to initialize tsid factory")
	}
	tsidFactory.node = node & int32(tsidFactory.nodeMask)

	tsidFactory.lastTime = tsidFactory.clock.UnixMilli()
	randomNumber, err := tsidFactory.getRandomValue()
	if err != nil {
		return nil, err
	}

	tsidFactory.counter = randomNumber
	return tsidFactory, nil
}

// Generate will return a tsid with random number
func (factory *TsidFactory) Generate() (*Tsid, error) {
	time, err := factory.getTime()
	if err != nil {
		return nil, err
	}

	time = time << RANDOM_BITS
	node := factory.node << factory.counterBits
	counter := factory.counter & factory.counterMask

	tsidNumber := int64(time | int64(node) | int64(counter))
	return NewTsid(tsidNumber), nil
}

func (factory *TsidFactory) getTime() (int64, error) {
	time := factory.clock.UnixMilli()
	if time <= factory.lastTime {
		lock.Lock()
		factory.counter++
		carry := uint32(factory.counter) >> factory.counterBits
		factory.counter = factory.counter & factory.counterMask
		time = factory.lastTime + int64(carry)
		lock.Unlock()

	} else {
		value, err := factory.getRandomValue()
		if err != nil {
			return 0, err
		}
		lock.Lock()
		factory.counter = value
		lock.Unlock()
	}
	factory.lastTime = time
	return (time - factory.customEpoch), nil
}

func (factory *TsidFactory) getRandomValue() (int32, error) {
	return factory.getRandomCounter()
}

func (factory *TsidFactory) getRandomCounter() (int32, error) {
	switch factory.random.(type) {
	case *byteRandom:
		{
			rLock.Lock()
			bytes, err := factory.random.NextBytes(factory.randomBytes)
			rLock.Unlock()

			if err != nil {
				return 0, err
			}

			switch len(bytes) {
			case 1:
				return int32((bytes[0] & 0xff) & byte(factory.counterMask)), nil
			case 2:
				return ((int32(bytes[0]&0xff) << 8) | int32(bytes[1]&0xff)) & factory.counterMask, nil
			case 3:
				return ((int32(bytes[0]&0xff) << 16) | (int32(bytes[1]&0xff) << 8) |
					int32(bytes[2]&0xff)) & factory.counterMask, nil
			}
		}
	case *intRandom:
		{
			rLock.Lock()
			value, err := factory.random.NextInt()
			rLock.Unlock()
			if err != nil {
				return 0, err
			}

			return int32(value & factory.counterMask), nil
		}
	}

	return 0, errors.New("invalid random")
}

type tsidFactoryBuilder struct {
	node        int32
	nodeBits    int32
	customEpoch int64
	clock       Clock
	random      Random
}

// TsidFactoryBuilder should be used to get instance of tsidFactory
func TsidFactoryBuilder() *tsidFactoryBuilder {
	return &tsidFactoryBuilder{}
}

func (builder *tsidFactoryBuilder) WithNode(node int32) *tsidFactoryBuilder {
	builder.node = node
	return builder
}

func (builder *tsidFactoryBuilder) WithNodeBits(nodeBits int32) *tsidFactoryBuilder {
	builder.nodeBits = nodeBits
	return builder
}

func (builder *tsidFactoryBuilder) WithCustomEpoch(customEpoch int64) *tsidFactoryBuilder {
	builder.customEpoch = customEpoch
	return builder
}

func (builder *tsidFactoryBuilder) WithClock(clock Clock) *tsidFactoryBuilder {
	builder.clock = clock
	return builder
}

func (builder *tsidFactoryBuilder) WithRandom(random Random) *tsidFactoryBuilder {
	builder.random = random
	return builder
}

// GetNode returns the provided node id. Default is zero.
func (builder *tsidFactoryBuilder) GetNode() (int32, error) {
	if builder.nodeBits <= 0 {
		return 0, nil
	}
	max := int32(1<<builder.nodeBits) - 1

	if builder.node < 0 || builder.node > max {
		err := fmt.Sprintf("node id out of range [0, %d]: %d", max, builder.node)
		return 0, errors.New(err)
	}
	return builder.node, nil
}

// GetNodeBits returns the provided node bits. Default is zero.
// Range: [0, 20]
func (builder *tsidFactoryBuilder) GetNodeBits() (int32, error) {
	max := 20

	if builder.nodeBits < 0 || builder.nodeBits > 20 {
		err := fmt.Sprintf("node bits out of range [0, %d]: %d", max, builder.nodeBits)
		return 0, errors.New(err)
	}
	return builder.nodeBits, nil
}

func (builder *tsidFactoryBuilder) GetClock() Clock {
	if builder.clock == nil {
		builder.clock = time.Now().UTC()
	}
	return builder.clock
}

func (builder *tsidFactoryBuilder) GetRandom() Random {
	if builder.random == nil {
		randomSupplier := NewMathRandomSupplier()
		builder.random = NewIntRandom(randomSupplier)
	}

	return builder.random
}

func (builder *tsidFactoryBuilder) GetCustomEpoch() int64 {
	if builder.customEpoch == 0 {
		builder.customEpoch = TSID_EPOCH
	}
	return builder.customEpoch
}

func (builder *tsidFactoryBuilder) Build() (*TsidFactory, error) {
	if tsidFactoryInstance != nil {
		return tsidFactoryInstance, nil
	}

	lock.Lock()
	defer lock.Unlock()

	var err error = nil
	if tsidFactoryInstance == nil {
		tsidFactoryInstance, err = newTsidFactory(builder)
		if err != nil {
			return nil, err
		}
	}
	return tsidFactoryInstance, nil
}

func (builder *tsidFactoryBuilder) NewInstance() (*TsidFactory, error) {
	return newTsidFactory(builder)
}
