package main

import (
	"context"
	"fmt"
	"hello/pkg/helloservice"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	Port        = ":3001"
	connections = 100000
	avg         = 5
)

func main() {
	latencyAvg := 0.0
	rateAvg := 0.0

	for j := 0; j < avg; j++ {
		var wg sync.WaitGroup
		var errCount int
		var ts []float64

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
		defer cancel()

		conn, err := grpc.DialContext(ctx, Port,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			fmt.Println("error", err)
			return
		}
		defer conn.Close()

		client := helloservice.NewHelloServiceClient(conn)

		for i := 0; i < connections; i++ {
			wg.Add(1)

			go func(wg *sync.WaitGroup) {
				t1 := time.Now()

				defer wg.Done()

				resp, err := client.Echo(ctx, &helloservice.Request{Message: "hello"})
				if err != nil {
					errCount += 1
				}
				_ = resp

				t2 := time.Since(t1)
				ts = append(ts, t2.Seconds())
			}(&wg)
		}

		wg.Wait()
		// fmt.Println("errors:", errCount)
		// fmt.Println("total:", connections)
		// fmt.Printf("rate: %.2f %%\n", float64(errCount)/connections*100)

		t := 0.0
		for _, val := range ts {
			t += val
		}

		t = t / float64(len(ts))
		latencyAvg += t
		rateAvg += float64(errCount) / connections * 100
	}
	latencyAvg /= avg
	rateAvg /= avg
	fmt.Printf("%.2f %.2f\n", rateAvg, latencyAvg)
}
