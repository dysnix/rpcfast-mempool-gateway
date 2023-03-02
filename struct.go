package main

type AccessListItem struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

type Transaction struct {
	ChainId              string           `json:"chainId"`
	Nonce                string           `json:"nonce"`
	MaxPriorityFeePerGas string           `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string           `json:"maxFeePerGas"`
	GasLimit             string           `json:"gasLimit"`
	To                   string           `json:"to"`
	Value                string           `json:"value"`
	Data                 string           `json:"data"`
	AccessList           []AccessListItem `json:"accessList"`
	V                    string           `json:"v"`
	R                    string           `json:"r"`
	S                    string           `json:"s"`
	From                 string           `json:"from"`
	TxHash               string           `json:"txhash"`
	Peer                 string           `json:"peer"`
}

type PeerInfo struct {
	Ip    string `redis:"ip"`
	Port  int    `redis:"port"`
	Score int    `redis:"score"`
}
