package schema

import "homework/internal/model"

type (
	IdsWithHashes struct {
		Ids    []string
		Hashes []string
	}

	PageParam struct {
		Size uint
		Page uint
	}

	GetParam struct {
		Ids         []string
		Status      model.Status
		Order       string
		Limit       uint
		RecipientId string
		Offset      uint
	}
)

func NewIdsWithHashes(ids []string, hashes []string) (IdsWithHashes, error) {
	if len(ids) == len(hashes) {
		return IdsWithHashes{Ids: ids, Hashes: hashes}, nil
	}
	return IdsWithHashes{}, ErrListWithHashesDifferentLength
}
