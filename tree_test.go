package main

import (
	"regexp"
	"testing"
)

func TestInsertZeroUrls(t *testing.T) {
	urls := []Url{}
	urlsTree := NewUrlsTree() //urls should be passed to init the tree
	urlsTree.InsertUrls(urls)

	treeRepresentation := urlsTree.String()
	if treeRepresentation != "" {
		t.Fatalf(`treeRepresentation is %q and not ""`, treeRepresentation)
	}
}

func TestInsertOneUrlWithOneResource(t *testing.T) {
	url := Url{&ResourceIdentifier{name: "username"}}
	urls := []Url{url}

	urlsTree := NewUrlsTree()
	urlsTree.InsertUrls(urls)

	treeRepresentation := urlsTree.String()
	want := regexp.MustCompile(`username`)

	if !want.MatchString(treeRepresentation) {
		t.Fatalf(`treeRepresentation after insertion is %q , want match for %#q`, treeRepresentation, want)
	}
}

func TestInsertOneUrlWithTwoResources(t *testing.T) {
	url := Url{
		&ResourceIdentifier{name: "username"},
		&ResourceCollection{name: "wordgame"},
		&ResourceCollection{name: "stats"},
	}
	urls := []Url{url}

	urlsTree := NewUrlsTree()
	urlsTree.InsertUrls(urls)

	treeRepresentation := urlsTree.String()
	want := regexp.MustCompile(`username.*wordgame.*stats`)
	if !want.MatchString(treeRepresentation) {
		t.Fatalf(`treeRepresentation after insertion is %q , want match for %#q`, treeRepresentation, want)
	}
}

func TestInsertTwoUrlsWithTwoResources(t *testing.T) {
	urlOne := Url{
		&ResourceIdentifier{name: "username"},
		&ResourceCollection{name: "wordgame"},
		&ResourceCollection{name: "stats"},
	}
	urlTwo := Url{
		&ResourceIdentifier{name: "username"},
		&ResourceCollection{name: "wordgame"},
		&ResourceCollection{name: "words"},
	}
	urls := []Url{urlOne, urlTwo}

	urlsTree := NewUrlsTree()
	urlsTree.InsertUrls(urls)

	treeRepresentation := urlsTree.String()
	wantOne := regexp.MustCompile("username.*wordgame.*stats")
	wantTwo := regexp.MustCompile("username.*wordgame.*words")
	if !wantOne.MatchString(treeRepresentation) || !wantTwo.MatchString(treeRepresentation) {
		t.Fatalf(`treeRepresentation after insertion is %q, want matches for %#q and %#q`, treeRepresentation, wantOne, wantTwo)
	}
}

func TestInsertTwoUrlsWithTwoResourcesWithDifferingHeads(t *testing.T) {
	urlOne := Url{
		&ResourceIdentifier{name: "username"},
		&ResourceCollection{name: "books"},
		&ResourceIdentifier{name: "book_id"},
	}
	urlTwo := Url{
		&ResourceCollection{name: "books"},
		&ResourceIdentifier{name: "book_id"},
	}

	urls := []Url{urlOne, urlTwo}
	urlsTree := NewUrlsTree()
	urlsTree.InsertUrls(urls)

	treeRepresentation := urlsTree.String()
	wantOne := regexp.MustCompile("username.*books.*book_id")
	wantTwo := regexp.MustCompile("books.*book_id")
	if !wantOne.MatchString(treeRepresentation) || !wantTwo.MatchString(treeRepresentation) {
		t.Fatalf(`treeRepresentation after insertion is %q, want matches for %#q and %#q`, treeRepresentation, wantOne, wantTwo)
	}

}
