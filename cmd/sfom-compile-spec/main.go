package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	wof_placetypes "github.com/whosonfirst/go-whosonfirst-placetypes"
	sfom_placetypes "github.com/sfomuseum/go-sfomuseum-placetypes"	
)

func main() {

	iterator_uri := flag.String("iterator-uri", "directory://", "")
	iterator_source := flag.String("iterator-source", "/usr/local/sfomuseum/sfomuseum-placetypes/placetypes", "")

	flag.Parse()

	ctx := context.Background()

	mu := new(sync.RWMutex)
	parent_map := new(sync.Map)

	// TO DO: pull in "core" placetypes spec and add parent pointers
	
	
	custom_placetypes := make([]*sfom_placetypes.SFOMuseumPlacetypeDefinition, 0)

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		if filepath.Ext(path) != ".json" {
			return nil
		}

		var pt *sfom_placetypes.SFOMuseumPlacetypeDefinition

		dec := json.NewDecoder(r)
		err := dec.Decode(&pt)

		if err != nil {
			return fmt.Errorf("Failed to decode %s, %w", path, err)
		}

		parent_map.Store(pt.Name, pt.Id)

		mu.Lock()
		defer mu.Unlock()

		custom_placetypes = append(custom_placetypes, pt)
		return nil
	}

	iter, err := iterator.NewIterator(ctx, *iterator_uri, iter_cb)

	if err != nil {
		log.Fatalf("Failed to create iterator, %v", err)
	}

	err = iter.IterateURIs(ctx, *iterator_source)

	if err != nil {
		log.Fatalf("Failed to iterate URIs, %v", err)
	}

	// START OF... not sure...

	spec := make(map[string]*wof_placetypes.WOFPlacetype)

	for _, pt := range custom_placetypes {

		// Legacy stuff, oh well...
		str_id := strconv.FormatInt(pt.Id, 10)
		
		parents := pt.Parent
		parent_ids := make([]int64, len(parents))

		for idx, p := range parents {

			p_id, ok := parent_map.Load(p)

			if !ok {
				log.Fatalf("Unable to derive parent ID for %s", p)
			}

			parent_ids[idx] = p_id.(int64)
		}

		spec[str_id] = &wof_placetypes.WOFPlacetype{
			Id: pt.Id,
			Role:   pt.Role,
			Name:   pt.Name,
			Parent: parent_ids,
		}
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(spec)

	if err != nil {
		log.Fatalf("Failed to encode spec, %v", err)
	}
}
