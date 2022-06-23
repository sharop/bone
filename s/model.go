package s

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

type NPlet struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
}

type KVPlets struct {
	Prefix string `json:"prefix"`
	NPlet  NPlet  `json:"nplet"`
}

func BuildIndexMapping() (mapping.IndexMapping, error) {

	prefixKey := bleve.NewTextFieldMapping()
	subject := bleve.NewTextFieldMapping()
	predicate := bleve.NewTextFieldMapping()
	object := bleve.NewTextFieldMapping()

	//Creating a document in order to manage the triplets
	nPlet := bleve.NewDocumentMapping()
	nPlet.AddFieldMappingsAt("subject", subject)
	nPlet.AddFieldMappingsAt("predicate", predicate)
	nPlet.AddFieldMappingsAt("object", object)

	keyMapping := bleve.NewDocumentMapping()
	keyMapping.AddSubDocumentMapping("nplets", nPlet)
	keyMapping.AddFieldMappingsAt("prefix", prefixKey)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("key", keyMapping)

	//NPLET {"subject":"","predicate":"","object":""}
	//PREFIX (TEXT BIN)

	return indexMapping, nil
}
