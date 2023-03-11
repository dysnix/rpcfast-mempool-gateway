// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8000", "http service address")

var RedisUrl = flag.String("redis-url", "95.217.34.222:6379", "redis address:port")
var RedisPassword = flag.String("redis-password", "afuogh3shi", "redis password")
var RedisDB = flag.Int("redis-db", 0, "redis DB")
var RedisChannel = flag.String("redis-channel", "polygon-mainnet-txs", "redis channel")

func servePeers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     *RedisUrl,
		Password: *RedisPassword,
		DB:       *RedisDB,
		PoolSize: 100,
	})

	data := rdb.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
		Key:   "peers",
		Start: 0,
		Stop:  100,
		Rev:   true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(data.Val())
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/api/peers", servePeers)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
