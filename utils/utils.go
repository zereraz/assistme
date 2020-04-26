package utils

import "hash/fnv"

func GenerateHash(data string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(data))
	return hash.Sum64()
}
