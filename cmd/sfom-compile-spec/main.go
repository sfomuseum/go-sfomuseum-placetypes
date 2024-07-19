package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	_ "slices"
	"strconv"
	"sync"

	sfom_placetypes "github.com/sfomuseum/go-sfomuseum-placetypes"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	wof_placetypes "github.com/whosonfirst/go-whosonfirst-placetypes"
)

func main() {

	iterator_uri := flag.String("iterator-uri", "directory://", "")
	iterator_source := flag.String("iterator-source", "/usr/local/sfomuseum/sfomuseum-placetypes/placetypes", "")

	flag.Parse()

	ctx := context.Background()

	mu := new(sync.RWMutex)
	parent_map := new(sync.Map)

	sfom_spec, err := sfom_placetypes.SFOMuseumPlacetypeSpecification()

	if err != nil {
		log.Fatalf("Failed to open latest spec, %v", err)
	}

	for str_id, pt := range sfom_spec.Catalog() {

		id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			log.Fatalf("Failed to parse ID (%s) for %s, %v", str_id, pt.Name, err)
		}

		parent_map.Store(pt.Name, id)
	}

	// END OF pull in "core" placetypes spec and add parent pointers

	custom_placetypes := make([]*sfom_placetypes.SFOMuseumPlacetypeRecord, 0)

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		if filepath.Ext(path) != ".json" {
			return nil
		}

		var pt *sfom_placetypes.SFOMuseumPlacetypeRecord

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

	new_spec := make(map[string]*wof_placetypes.WOFPlacetype)

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

		/*
			wof_concordance, exists := pt.Concordances["wof:placetype"]

			if exists {

				slog.Info("Placetype has wof:placetype concordance", "placetype", pt.Name, "wof placetype", wof_concordance)

				wof_pt, err := wof_placetypes.GetPlacetypeByName(wof_concordance)

				if err != nil {
					slog.Warn("Failed to load wof:placetype, skipping", "placetype", wof_concordance, "error", err)
				} else {

					for _, pid := range wof_pt.Parent {

						if slices.Contains(parent_ids, pid){
							continue
						}

						slog.Info("Add parent for WOF placetype", "placetype", pt.Name, "wof placetype", wof_pt.Name, "parent", pid)
						parent_ids = append(parent_ids, pid)
					}
				}
			}
		*/

		slog.Info("Add placetype", "id", pt.Id, "name", pt.Name, "parents", parent_ids)

		new_spec[str_id] = &wof_placetypes.WOFPlacetype{
			Id:     pt.Id,
			Role:   pt.Role,
			Name:   pt.Name,
			Parent: parent_ids,
		}
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(new_spec)

	if err != nil {
		log.Fatalf("Failed to encode spec, %v", err)
	}
}
