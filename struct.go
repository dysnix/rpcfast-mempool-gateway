package main

type AccessListItem struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

type Transaction struct {
	Type                 string           `json:"type"`
	ChainId              string           `json:"chainId"`
	Nonce                string           `json:"nonce"`
	GasPrice             string           `json:"gasPrice"`
	MaxPriorityFeePerGas string           `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string           `json:"maxFeePerGas"`
	Gas                  string           `json:"gas"`
	To                   string           `json:"to"`
	Value                string           `json:"value"`
	Input                string           `json:"input"`
	AccessList           []AccessListItem `json:"accessList"`
	V                    string           `json:"v"`
	R                    string           `json:"r"`
	S                    string           `json:"s"`
	TxHash               string           `json:"hash"`
	Peer                 string           `json:"peer"`
	//From                 string           `json:"from"`
}

type newTx struct {
	Peer        string      `json:"peer"`
	Transaction Transaction `json:"tx"`
}
