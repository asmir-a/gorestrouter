package tree

import (
	"regexp"
	"testing"

	"github.com/asmir-a/gorestrouter/resource"
)

func TestInsertZeroUrls(t *testing.T) {
	urls := []resource.Url{}
	urlsTree := NewUrlsTree(urls) //urls should be passed to init the tree

	treeRepresentation := urlsTree.String()
	if treeRepresentation != "" {
		t.Fatalf(`treeRepresentation is %q and not ""`, treeRepresentation)
	}
}

func TestInsertOneUrlWithOneResource(t *testing.T) {
	url := resource.Url{resource.NewResourceIdentifier("username", nil)}
	urls := []resource.Url{url}

	urlsTree := NewUrlsTree(urls)

	treeRepresentation := urlsTree.String()
	want := regexp.MustCompile(`username`)

	if !want.MatchString(treeRepresentation) {
		t.Fatalf(`treeRepresentation after insertion is %q , want match for %#q`, treeRepresentation, want)
	}
}

func TestInsertOneUrlWithTwoResources(t *testing.T) {
	url := resource.Url{
		resource.NewResourceIdentifier("username", nil),
		resource.NewResourceCollection("wordgame", nil),
		resource.NewResourceCollection("stats", nil),
	}
	urls := []resource.Url{url}

	urlsTree := NewUrlsTree(urls)

	treeRepresentation := urlsTree.String()
	want := regexp.MustCompile(`username.*wordgame.*stats`)
	if !want.MatchString(treeRepresentation) {
		t.Fatalf(`treeRepresentation after insertion is %q , want match for %#q`, treeRepresentation, want)
	}
}

func TestInsertTwoUrlsWithTwoResources(t *testing.T) {
	urlOne := resource.Url{
		resource.NewResourceIdentifier("username", nil),
		resource.NewResourceCollection("wordgame", nil),
		resource.NewResourceCollection("stats", nil),
	}
	urlTwo := resource.Url{
		resource.NewResourceIdentifier("username", nil),
		resource.NewResourceCollection("wordgame", nil),
		resource.NewResourceCollection("words", nil),
	}
	urls := []resource.Url{urlOne, urlTwo}
	urlsTree := NewUrlsTree(urls)

	treeRepresentation := urlsTree.String()
	wantOne := regexp.MustCompile("username.*wordgame.*stats")
	wantTwo := regexp.MustCompile("username.*wordgame.*words")
	if !wantOne.MatchString(treeRepresentation) || !wantTwo.MatchString(treeRepresentation) {
		t.Fatalf(`treeRepresentation after insertion is %q, want matches for %#q and %#q`, treeRepresentation, wantOne, wantTwo)
	}
}

func TestInsertTwoUrlsWithTwoResourcesWithDifferingHeads(t *testing.T) {
	urlOne := resource.Url{
		resource.NewResourceIdentifier("username", nil),
		resource.NewResourceCollection("books", nil),
		resource.NewResourceCollection("book_id", nil),
	}
	urlTwo := resource.Url{
		resource.NewResourceCollection("books", nil),
		resource.NewResourceIdentifier("book_id", nil),
	}

	urls := []resource.Url{urlOne, urlTwo}
	urlsTree := NewUrlsTree(urls)

	treeRepresentation := urlsTree.String()
	wantOne := regexp.MustCompile("username.*books.*book_id")
	wantTwo := regexp.MustCompile("books.*book_id")
	if !wantOne.MatchString(treeRepresentation) || !wantTwo.MatchString(treeRepresentation) {
		t.Fatalf(`treeRepresentation after insertion is %q, want matches for %#q and %#q`, treeRepresentation, wantOne, wantTwo)
	}

}
