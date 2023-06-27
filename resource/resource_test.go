package resource

import (
	"reflect"
	"testing"
)

func TestResourceGeneral(t *testing.T) {
	resourceOne := NewResourceIdentifier("username", nil)
	resourceTwo := NewResourceCollection("wordgame", nil)
	resourceThree := NewResourceCollection("stats", nil)
	url := Url{resourceOne, resourceTwo, resourceThree}
	resourceNames := []string{}
	resourceTypes := []string{}
	for _, resource := range url {
		resourceNames = append(resourceNames, resource.Name())
		switch resource.(type) {
		case *ResourceIdentifier:
			resourceTypes = append(resourceTypes, "identifier")
		case *ResourceCollection:
			resourceTypes = append(resourceTypes, "collection")
		default:
			t.Fatal("is not supposed to happen")
		}
	}

	wantedResourceNames := []string{"username", "wordgame", "stats"}
	if !reflect.DeepEqual(resourceNames, wantedResourceNames) {
		t.Fatal("the resourceNames is supposed to be: ", wantedResourceNames, " but is: ", resourceNames)
	}

	wantedResourceTypes := []string{"identifier", "collection", "collection"}
	if !reflect.DeepEqual(resourceTypes, wantedResourceTypes) {
		t.Fatal("the resourceTypes is supposed to be: ", wantedResourceTypes, " but is: ", resourceTypes)
	}
}
