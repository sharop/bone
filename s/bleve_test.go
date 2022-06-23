package s_test

import (
	"flag"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/sharop/bone/s"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"
)

var batchSize = flag.Int("batchSize", 100, "batch size for indexing")
var staticPath = flag.String("index", "nplets.bleve", "Index Path")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func TestIndex(t *testing.T) {

	IX := s.BIndex{}
	dir, err := ioutil.TempDir("", "bleve-test")
	require.NoError(t, err)
	log.Printf("Testing directory : %v", dir)
	defer os.RemoveAll(dir)
	IX.BIndex, err = s.Init(dir)
	require.NoError(t, err)

	for scenario, fn := range map[string]func(t *testing.T, ix *s.BIndex, args string){
		"Close Ix":   testCloseIx,
		"Re-open IX": testReOpenIx,
		"IndexData":  testIndex,
	} {
		t.Run(scenario, func(t *testing.T) {
			ix := &IX
			fn(t, ix, dir)
		})
	}

}

func testCloseIx(t *testing.T, ix *s.BIndex, path string) {

	err := (*ix.BIndex).Close()
	require.NoError(t, err)
	ix.Closed = err == nil

}

func testReOpenIx(t *testing.T, ix *s.BIndex, path string) {

	var err error
	ix.BIndex, err = s.Init(path)
	require.NoError(t, err)
	require.NotNil(t, (*ix.BIndex))

}
func testIndex(t *testing.T, ix *s.BIndex, path string) {
	testData := tData()

	err := ix.Index(testData[0])
	require.NoError(t, err)

	err = ix.BatchIndex(3, testData[1:])
	require.NoError(t, err)
	dc, err := (*ix.BIndex).DocCount()
	require.True(t, dc >= 9)

}
func testIndexB(t *testing.T, ix *s.BIndex, path string) {

	flag.Parse()
	log.Printf("GOMAXPROCS: %d", runtime.GOMAXPROCS(-1))

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}

	nPletIndex, err := bleve.Open(*staticPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		indexMapping, err := s.BuildIndexMapping()
		if err != nil {
			log.Fatal(err)
		}
		nPletIndex, err := bleve.New(*staticPath, indexMapping)
		if err != nil {
			log.Fatal(err)
		}

		err = indexKVPlets(nPletIndex)
		if err != nil {
			log.Fatal(err)
		}

	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Opening existing index...")
	}

	// search for some text
	query := bleve.NewMatchQuery("AB1001")
	search := bleve.NewSearchRequest(query)
	searchResults, err := nPletIndex.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)

}

func tData() []s.KVPlets {
	var NPletsTest []s.KVPlets

	NPletsTest = append(NPletsTest, s.KVPlets{Prefix: "AB1001-E", NPlet: s.NPlet{Subject: "AB1001-2001-SPACE-TEAM_VENTAS", Predicate: "", Object: ""}})
	NPletsTest = append(NPletsTest, s.KVPlets{Prefix: "AB1003-E", NPlet: s.NPlet{Subject: "AB1003-2001-PROJECT-CENTRAL_SUR", Predicate: "", Object: ""}})
	NPletsTest = append(NPletsTest, s.KVPlets{Prefix: "AB1001-T-AB1003", NPlet: s.NPlet{Subject: "AB1001-2001-SPACE-TEAM_VENTAS", Predicate: "AB1001-2001-CONTAINS", Object: "AB1003-2001-PROJECT-CENTRAL_SUR"}})
	NPletsTest = append(NPletsTest, s.KVPlets{Prefix: "AB1005-E", NPlet: s.NPlet{Subject: "AB1005-2001-SOURCE-2021", Predicate: "", Object: ""}})
	NPletsTest = append(NPletsTest, s.KVPlets{Prefix: "AB1005-T-AB1006", NPlet: s.NPlet{Subject: "AB1005-2001-SOURCE-2021", Predicate: "AB1005-2001-HAS", Object: "AB1006-2001-ASSET-2021"}})
	NPletsTest = append(NPletsTest, s.KVPlets{Prefix: "AB1006-D-AB1005", NPlet: s.NPlet{Subject: "AB1006-2001-ASSET-2021", Predicate: "AB1006-D-PARENT", Object: "AB1005-2001-SOURCE-2021"}})
	NPletsTest = append(NPletsTest, s.KVPlets{Prefix: "AB1003-T-AB1006", NPlet: s.NPlet{Subject: "AB1003-2001-PROJECT-CENTRAL_SUR", Predicate: "AB1003-2001-REFERENCE", Object: "AB1006-2001-ASSET-2021"}})
	NPletsTest = append(NPletsTest, s.KVPlets{Prefix: "AB1007-E", NPlet: s.NPlet{Subject: "AB1007-2001-SOURCE-CSV_CATALOG", Predicate: "", Object: ""}})
	NPletsTest = append(NPletsTest, s.KVPlets{Prefix: "AB1003-T-AB1007", NPlet: s.NPlet{Subject: "AB1003-2001-PROJECT-CENTRAL_SUR", Predicate: "AB1003-2001-REFERENCE", Object: "AB1007-2001-SOURCE-CSV_CATALOG"}})

	return NPletsTest
}

func indexKVPlets(i bleve.Index) error {
	//{prefix:"", "nplet":{"subject":"","predicate":"","object":""}}

	testData := tData()

	var count, batchCount int
	startTime := time.Now()
	batch := i.NewBatch()

	for _, p := range testData {
		err := batch.Index(p.Prefix, p)
		if err != nil {
			return err
		}
		batchCount++
		if batchCount >= *batchSize {
			err := i.Batch(batch)
			if err != nil {
				return err
			}
			batch = i.NewBatch()
			batchCount = 0
		}
		count++
		if count%1000 == 0 {
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(count)
			log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
		}
	}
	// flush the last batch
	if batchCount > 0 {
		err := i.Batch(batch)
		if err != nil {
			log.Fatal(err)
		}
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
	return nil

}
