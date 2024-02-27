package broker

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"nats/api"
	"nats/internal/db"
	"nats/internal/storage"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

const subj = "orders"

func natsConnect() (*nats.Conn, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		slog.Error("Consumer can't connect", err)
		return nil, err
	}
	slog.Info("Sub is connecting")
	return nc, nil
}

func getSubConnection(nc *nats.Conn) stan.Conn {
	sc, err := stan.Connect("test-cluster", "stan-sub", stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			slog.Error("Connection sub lost, reason: %v", reason)
		}))
	if err != nil {
		slog.Error("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, nats.DefaultURL)
	}

	slog.Info(fmt.Sprintf("Connected to: %s", nats.DefaultURL))
	return sc
}

func Subscribe_to_channel() {
	var chunk api.Order

	nc, err := natsConnect()
	if err != nil {
		panic("Coud not connect to NATS-server")
	}

	defer nc.Close()
	sc := getSubConnection(nc)
	defer sc.Close()

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

	for {
		select {
		case m := <-ch:
			fmt.Println("Message has arrived!: \n ", string(m.Data))
			err = json.Unmarshal(m.Data, &chunk)

			if err != nil {
				slog.Error("Error while validating sub chunk", err)
			}

			err = chunk.Validate()

			if err != nil {
				slog.Error(fmt.Sprintf("Error while validating chunk id #%s: err - %s", chunk.Order_uid, err))
			} else {
				storage.AppCache.Set(chunk.Order_uid, m.Data)
				db.InsertItem(db.DBConnection, chunk)
			}
			m.Ack()

		case <-time.After(15 * time.Second):
			return
		}
	}

}
