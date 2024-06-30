package dto

import (
	"homework/pkg/hash"
)

type (
	IdsWithHashes struct {
		Ids    []string
		Hashes []string
	}
)

func NewIdsWithHashes(ids []string, hashes []string) (IdsWithHashes, error) {
	if len(ids) == len(hashes) {
		return IdsWithHashes{Ids: ids, Hashes: hashes}, nil
	}
	return IdsWithHashes{}, ErrListWithHashesDifferentLength
}

func GenHashes(strings []string) (IdsWithHashes, error) {
	var hashes []string
	for i := 0; i < len(strings); i++ {
		hashes = append(hashes, hash.GenerateHash())
	}
	return NewIdsWithHashes(strings, hashes)
}
