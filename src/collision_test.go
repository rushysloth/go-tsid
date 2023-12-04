package main

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollision(t *testing.T) {

	t.Run("one goroutine per node", func(t *testing.T) {

		// One goroutine demonstrates one node
		goroutineCount := 10
		iterationCount := 100_000

		var collisionCounter atomic.Uint32
		var tsidMap sync.Map

		wg := &sync.WaitGroup{}

		for i := 0; i < goroutineCount; i++ {

			nodeId := i

			wg.Add(1)
			go func(nodeId int32, iterationCount int32, collisionCounter *atomic.Uint32,
				tsidMap *sync.Map, wg *sync.WaitGroup) {
				defer wg.Done()

				tsidFactory, err := TsidFactoryBuilder().
					WithNode(nodeId).
					Build()
				assert.Nil(t, err)

				for j := 0; j < int(iterationCount); j++ {
					tsid, err := tsidFactory.Generate()
					assert.Nil(t, err)

					// check if this tsid was already generated
					if _, ok := tsidMap.Load(tsid); !ok {

						// not present, store it
						tsidMap.Store(tsid, (nodeId*iterationCount)+int32(j))
						continue
					}

					// collision detected, increment counter and break out
					collisionCounter.Add(1)
					break
				}

			}(int32(nodeId), int32(iterationCount), &collisionCounter, &tsidMap, wg)

		}

		// Wait for all goroutines to complete
		wg.Wait()

		assert.Zero(t, collisionCounter.Load(), 0, "Collision detected")
	})

	t.Run("multiple goroutines per node", func(t *testing.T) {

		node := 1
		nodeBit := 1
		goroutineCount := 10
		iterationCount := 200_000

		var collisionCounter atomic.Uint32
		var tsidMap sync.Map

		wg := &sync.WaitGroup{}

		for i := 0; i < goroutineCount; i++ {

			wg.Add(1)
			go func(nodeId int32, nodeBit int32, iterationCount int32, collisionCounter *atomic.Uint32,
				tsidMap *sync.Map, wg *sync.WaitGroup) {
				defer wg.Done()

				tsidFactory, err := TsidFactoryBuilder().
					WithNodeBits(nodeBit).
					Build()
				assert.Nil(t, err)

				for j := 0; j < int(iterationCount); j++ {
					tsid, err := tsidFactory.Generate()
					assert.Nil(t, err)

					// check if this tsid was already generated
					if _, ok := tsidMap.Load(tsid); !ok {

						// not present, store it
						tsidMap.Store(tsid, (nodeId*iterationCount)+int32(j))
						continue
					}

					// collision detected, increment counter and break out
					collisionCounter.Add(1)
					break
				}

			}(int32(node), int32(nodeBit), int32(iterationCount), &collisionCounter, &tsidMap, wg)

		}

		// Wait for all goroutines to complete
		wg.Wait()

		assert.Zero(t, collisionCounter.Load(), 0, "Collision detected")
	})
}
