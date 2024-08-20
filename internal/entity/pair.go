package entity

type Pair struct {
	Key   string
	Value interface{}
}

type EncryptedPair struct {
	Key   string
	Value string
}

type Error struct {
	Key     string
	Message string
}
