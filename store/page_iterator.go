// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	"github.com/mattermost/mattermost-server/model"
)

type PageIterator struct {
	HasNext   bool
	BatchSize int
	More      func(offset int) (int, interface{}, *model.AppError)
	offset    int
}

func (it *PageIterator) Next() (interface{}, *model.AppError) {
	batchLen, batch, err := it.More(it.offset)
	if err != nil {
		return nil, err
	}
	if batchLen < it.BatchSize {
		it.HasNext = false
	}
	it.offset += batchLen
	return batch, nil
}

func NewPageIterator(batchSize int, moreF func(offset int) (int, interface{}, *model.AppError)) *PageIterator {
	return &PageIterator{
		HasNext:   true,
		BatchSize: batchSize,
		More:      moreF,
	}
}
