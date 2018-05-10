// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
)

func ExampleNewPageIterator() {
	more := func(offset int) (int, interface{}, *model.AppError) {
		var batch []*model.Scheme
		var result StoreResult

		// if result = <-a.Srv.Store.Scheme().GetAllPage("", offset, batchSize); result.Err != nil {
		// 	return 0, nil, result.Err
		// }

		batch = result.Data.([]*model.Scheme)
		return len(batch), batch, nil
	}

	schemeIterator := NewPageIterator(100, more)

	for schemeIterator.HasNext {
		schemeBatch, _ := schemeIterator.Next()
		fmt.Println(schemeBatch)
	}
}
