package entities

// A consul GET response
type ConsulResult struct {
	LockIndex   int    `json:"lockIndex"`
	Key         string `json:"Key"`
	Flags       int    `json:"Flags"`
	Value       string `json:"Value"`
	CreateIndex int    `json:"CreateIndex"`
	ModifyIndex int    `json:"ModifyIndex"`
}

// A consul Transaction payload
type ConsulTxn struct {
	KV ConsulTxnKV `json:"KV"`
}

// A consul KV Transaction payload
type ConsulTxnKV struct {
	Verb  *string `json:"Verb"`
	Key   *string `json:"Key"`
	Value *string `json:"Value,omitempty"`
}


