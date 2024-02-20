package broker

import (
	"fmt"
	"log/slog"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func Subscribe() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		slog.Error("Consumer can't connect", err)
	}
	slog.Info("Sub is connecting")
	defer nc.Close()

	sc, err := stan.Connect("test-cluster", "stan-sub", stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			slog.Error("Connection sub lost, reason: %v", reason)
		}))
	if err != nil {
		slog.Error("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, nats.DefaultURL)
	}

	slog.Info("Connected to: %", nats.DefaultURL)

	defer sc.Close()
	subj := "orders"
	ch := make(chan *stan.Msg, 1024)
	mcb := func(msg *stan.Msg) {
		ch <- msg
	}
	sub, err := sc.Subscribe(subj,
		mcb, stan.DeliverAllAvailable())
	if err != nil {
		slog.Error("Consumer can't subscribe %s", err)
	}

	defer sub.Close()
	defer sub.Unsubscribe()

	select {
	case m := <-ch:
		fmt.Printf("Message has arrived!: \n %s", m.Data)
		m.Ack()
	case <-time.After(5 * time.Second):
		break
	}
}
