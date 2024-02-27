package broker

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"nats/api"
	"nats/internal/services"
	"time"

	"github.com/nats-io/stan.go"
)

func Publish() error {

	sc, err := stan.Connect("test-cluster", "stan-pub", stan.NatsURL("nats://0.0.0.0:4222"))
	if err != nil {
		slog.Error("Can't connect a Publisher:\nMake sure a NATS Streaming Server is running at:%s", err)
	}
	defer sc.Close()
	var res api.Order
	subj := "orders"
	for {

		msg, _ := services.Json_generator(rand.IntN(100))

		err = sc.Publish(subj, msg)
		if err != nil {
			slog.Error("Error during publish: %v\n", err)
			return err
		}
		json.Unmarshal(msg, &res)
		slog.Info(fmt.Sprintf("Published: %s", res.Order_uid))
		time.Sleep(time.Second * 3)
	}

}
