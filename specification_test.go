package placetypes

import (
	_ "fmt"
	"testing"

	wof_placetypes "github.com/whosonfirst/go-whosonfirst-placetypes"	
)

func TestSFOMuseumPlacetypeSpecification(t *testing.T) {

	spec, err := SFOMuseumPlacetypeSpecification()

	if err != nil {
		t.Fatalf("Failed to create SFO Museum placetypes spec, %v", err)
	}

	placetype_names := []string{
		"gate",
		"region",
	}

	for _, name := range placetype_names {
		
		_, err = spec.GetPlacetypeByName(name)

		if err != nil {
			t.Fatalf("Failed to get '%s' placetype by name, %v", name, err)
		}
	}

	gate_pt, err := spec.GetPlacetypeByName("gate")

	if err != nil {
		t.Fatalf("Failed to load gate placetype, %v", err)
	}

	terminal_pt, err := spec.GetPlacetypeByName("terminal")

	if err != nil {
		t.Fatalf("Failed to load terminal placetype, %v", err)
	}

	county_pt, err := spec.GetPlacetypeByName("county")

	if err != nil {
		t.Fatalf("Failed to load county placetype, %v", err)
	}
	
	roles := []string{
		wof_placetypes.COMMON_ROLE,
		wof_placetypes.OPTIONAL_ROLE,
		wof_placetypes.COMMON_OPTIONAL_ROLE,
		wof_placetypes.CUSTOM_ROLE,				
	}
	
	a := spec.AncestorsForRoles(gate_pt, roles)
	count_a := len(a)

	expected_count := 24
	
	if count_a != expected_count {
		t.Fatalf("Unexpected ancestors for gate. Expected %d, but got %d", expected_count, count_a)
	}
	
	if !spec.IsAncestor(gate_pt, terminal_pt){
		t.Fatalf("Expected terminal to be ancestor of gate")
	}

	if !spec.IsAncestor(terminal_pt, county_pt){
		t.Fatalf("Expected county to be ancestor of terminal")
	}

	if !spec.IsDescendant(county_pt, terminal_pt){
		t.Fatalf("Expected terminal to be descendant of county")
	}

	roles_custom := []string{
		wof_placetypes.CUSTOM_ROLE,
	}

	planet_pt, _ := spec.GetPlacetypeByName("planet")

	if err != nil {
		t.Fatalf("Failed to get planet placetype, %v", err)
	}
		
	custom_pt := spec.DescendantsForRoles(planet_pt, roles_custom)

	expected_custom := 9

	if len(custom_pt) != expected_custom {
		t.Fatalf("Unexpected placetypes for custom role. Expected %d but got %d", expected_custom, len(custom_pt))
	}
}
