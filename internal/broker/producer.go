package broker

import (
	"encoding/json"
	"log/slog"
	"nats/api"
	"nats/internal/services"

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
	for i := 0; i < 5; i++ {

		msg, _ := services.Json_generator(i)

		err = sc.Publish(subj, msg)
		if err != nil {
			slog.Error("Error during publish: %v\n", err)
			return err
		}
		json.Unmarshal(msg, &res)
		slog.Info("Published:", res.Order_uid)
	}
	return nil

}
