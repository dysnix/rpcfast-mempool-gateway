// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"

	"github.com/fasthttp/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	var txCounter = 0
	var prevTxCounter = 0

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     *RedisUrl,
		Password: *RedisPassword,
		DB:       *RedisDB,
		PoolSize: 100,
	})

	// There is no error because go-redis automatically reconnects on error.
	pubsub := rdb.Subscribe(ctx, *RedisChannel)

	// Close the subscription when we are done.
	defer pubsub.Close()

	ch := pubsub.Channel()

	// Periodic task
	ticker := time.NewTicker(time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				diff := txCounter - prevTxCounter
				prevTxCounter = txCounter
				log.Printf("TXs processed: %d [%d tx / sec]", txCounter, diff)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for msg := range ch {
		tx := newTx{}
		err := json.Unmarshal([]byte(msg.Payload), &tx)
		if err != nil {
			log.Println("read:", err)
			return
		}

		enode, err := rdb.Get(ctx, tx.Transaction.TxHash).Result()
		switch {
		case err == redis.Nil || enode == "":
			// Broadcast tx to clients
			c.hub.broadcast <- []byte(msg.Payload)

			go func() {
				// Save TX in cache
				rdb.Set(ctx, tx.Transaction.TxHash, tx.Peer, 0)

				// Add peer with default score
				rdb.ZAdd(ctx, "peers", redis.Z{
					Score:  float64(1),
					Member: tx.Peer,
				})
			}()
			continue
		case err != nil:
			log.Println("Get failed", err)
			continue
		}

		rdb.ZIncrBy(ctx, "peers", float64(1), enode)
		txCounter += 1
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 1024)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
