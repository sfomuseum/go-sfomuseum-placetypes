package main

import (
	"flag"
	"log"

	sfom_placetypes "github.com/sfomuseum/go-sfomuseum-placetypes"
)

func main() {

	flag.Parse()

	pt_spec, err := sfom_placetypes.SFOMuseumPlacetypeSpecification()

	if err != nil {
		log.Fatalf("Failed to load SFOM placetypes spec, %v", err)
	}

	for _, str_pt := range flag.Args() {

		pt, err := pt_spec.GetPlacetypeByName(str_pt)

		if err != nil {
			log.Fatalf("Failed to load placetype '%s', %v", str_pt, err)
		}

		log.Println(pt)
	}
}
