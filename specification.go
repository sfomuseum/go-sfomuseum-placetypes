package placetypes

import (
	"embed"
	"fmt"

	wof_placetypes "github.com/whosonfirst/go-whosonfirst-placetypes"
)

//go:embed placetypes.json
var FS embed.FS

func SFOMuseumPlacetypeSpecification() (*wof_placetypes.WOFPlacetypeSpecification, error) {

	r, err := FS.Open("placetypes.json")

	if err != nil {
		return nil, fmt.Errorf("Failed to open placetypes, %w", err)
	}

	defer r.Close()

	core_spec, err := wof_placetypes.DefaultWOFPlacetypeSpecification()

	if err != nil {
		return nil, fmt.Errorf("Failed to load core placetypes spec, %w", err)
	}

	sfom_spec, err := wof_placetypes.NewWOFPlacetypeSpecificationWithReader(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to load SFO Museum placetypes spec, %w", err)
	}

	err = sfom_spec.AppendPlacetypeSpecification(core_spec)

	if err != nil {
		return nil, fmt.Errorf("Failed to append SFO Museum placetypes spec to core placetypes spec, %w", err)
	}

	return sfom_spec, nil
}
