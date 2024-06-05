package concurrency

import (
	"log"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestErrGroup(t *testing.T) {
	var eg errgroup.Group
	eg.SetLimit(2) // 并发度

	for i := 0; i < 3; i++ {
		var x = i
		eg.Go(func() error {
			log.Println("working", x)
			time.Sleep(time.Second)
			return nil
		})
	}
	_ = eg.Wait()
}
