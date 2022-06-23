package s

import (
	"github.com/blevesearch/bleve/v2"
	"log"
	"path"
)

var (
	DefaultFileMode = 0755
)

type BIndex struct {
	BIndex *bleve.Index
	Closed bool
}

func Init(directoryPath string) (*bleve.Index, error) {
	// open a new index

	bonePath := path.Join(directoryPath, "bone.index")
	index, err := bleve.Open(bonePath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping, err := BuildIndexMapping() //Add different kinds of mappings.
		if err != nil {
			log.Fatal(ErrorBoneIndexErrorOnInit, err)
		}

		index, err = bleve.New(bonePath, mapping)

		if err != nil {
			log.Fatal(err)
		}

	}

	return &index, nil
}

func (i BIndex) Index(item KVPlets) error {
	return (*i.BIndex).Index(item.Prefix, item)
}
func (i BIndex) BatchIndex(batchSize int, items []KVPlets) error {

	var batchCount int
	batch := (*i.BIndex).NewBatch()
	for _, item := range items {
		err := batch.Index(item.Prefix, item)
		if err != nil {
			return err
		}
		batchCount++
		if batchCount >= batchSize {
			err := (*i.BIndex).Batch(batch)
			if err != nil {
				return err
			}
			batch = (*i.BIndex).NewBatch()
			batchCount = 0
		}

	}
	// flush the last batch or batch with less batchSize
	if batchCount > 0 {
		err := (*i.BIndex).Batch(batch)
		if err != nil {
			return err
		}
	}

	return nil
}

type Indexer interface {
	Index(plets KVPlets) error
	BatchIndex(...KVPlets) error
}
