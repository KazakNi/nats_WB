package broker

import (
	"log/slog"
	"time"

	"github.com/nats-io/stan.go"
)

func Publish() {

	// connect to NATS-server
	/* nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		slog.Error("Can't pub:", err)
	}
	defer nc.Close()
	*/
	// connect to STAN-streaming "nats-streaming-server-v0.3.8"
	sc, err := stan.Connect("test-cluster", "stan-pub", stan.NatsURL("nats://0.0.0.0:4222"))
	if err != nil {
		slog.Error("Can't connect a Publisher:\nMake sure a NATS Streaming Server is running at:%s", err)
	}
	defer sc.Close()

	subj, msg := "orders", []byte("Healthcheck")
	err = sc.Publish(subj, msg)
	if err != nil {
		slog.Error("Error during publish: %v\n", err)
	}
	slog.Info("Published [%s] : '%s'\n", subj, msg)
	<-time.After(5 * time.Second)

}
