# RPCFast Mempool Gateway

## TODO
* [x] Listen redis pub/sub
* [x] Save and order peers by score (based on first tx receive)
* [x] API for get peers ordered by score
* [x] WS: add ws server
* [ ] JsonRPC: implement eth_sendRawTransaction with broadcast to top-10 peers by score