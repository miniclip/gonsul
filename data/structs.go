package data

type ConsulResult struct {
	LockIndex		int		`json:"lockIndex"`
	Key				string	`json:"Key"`
	Flags			int		`json:"Flags"`
	Value			string	`json:"Value"`
	CreateIndex		int		`json:"CreateIndex"`
	ModifyIndex		int		`json:"ModifyIndex"`
}

type Entry struct {
	KVPath 			string
	Value    		string
}