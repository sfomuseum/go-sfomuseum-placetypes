package placetypes

// This remains TBD... specifically in order to merge with WOF spec
// or even just to be able to reference "core" parents we'll need to
// maintain parity with WOFPlacetype struct

type SFOMuseumPlacetypeRecord struct {
	Id int64 `json:"sfomuseum:id"`
	Name string `json:"sfomuseum:name"`
	Role string `json:"sfomuseum:role"`
	Label string `json:"sfomuseum:label"`
	Parent []string `json:"sfomuseum:parent"`
	Concordances map[string]string `json:"sfomuseum:concordances"`
}
