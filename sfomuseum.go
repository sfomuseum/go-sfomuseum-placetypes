package placetypes

import (
	"context"
	"fmt"

	wof_placetypes "github.com/whosonfirst/go-whosonfirst-placetypes"
)

type SFOMuseumDefinition struct {
	wof_placetypes.Definition
	spec *wof_placetypes.WOFPlacetypeSpecification
	prop string
	uri  string
}

func init() {
	ctx := context.Background()
	wof_placetypes.RegisterDefinition(ctx, "sfomuseum", NewSFOMuseumDefinition)
}

func NewSFOMuseumDefinition(ctx context.Context, uri string) (wof_placetypes.Definition, error) {

	spec, err := SFOMuseumPlacetypeSpecification()

	if err != nil {
		return nil, fmt.Errorf("Failed to create default WOF placetype specification, %w", err)
	}

	s := &SFOMuseumDefinition{
		spec: spec,
		prop: "sfomuseum:placetype",
		uri:  uri,
	}

	return s, nil
}

func (s *SFOMuseumDefinition) Specification() *wof_placetypes.WOFPlacetypeSpecification {
	return s.spec
}

func (s *SFOMuseumDefinition) Property() string {
	return s.prop
}

func (s *SFOMuseumDefinition) URI() string {
	return s.uri
}
