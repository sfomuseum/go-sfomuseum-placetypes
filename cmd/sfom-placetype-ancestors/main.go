package main

import (
	"flag"
	"log"

	"github.com/sfomuseum/go-flags/multi"
	sfom_placetypes "github.com/sfomuseum/go-sfomuseum-placetypes"
	wof_placetypes "github.com/whosonfirst/go-whosonfirst-placetypes"
)

func main() {

	var roles multi.MultiString
	flag.Var(&roles, "role", "...")

	flag.Parse()

	pt_spec, err := sfom_placetypes.SFOMuseumPlacetypeSpecification()

	if err != nil {
		log.Fatal(err)
	}

	if len(roles) == 0 {
		roles = wof_placetypes.AllRoles()
	}

	for _, str_pt := range flag.Args() {

		pt, err := pt_spec.GetPlacetypeByName(str_pt)

		if err != nil {
			log.Fatal(err)
		}

		ancestors := pt_spec.AncestorsForRoles(pt, roles)

		for i, p := range ancestors {
			log.Println(i, p.Name)
		}
	}
}
