package tree

import (
	"fmt"
	"log"

	"github.com/asmir-a/gorestrouter/resource"
)

type UrlsTree struct {
	Root *ResourceNode
}

func NewUrlsTree(urls []resource.Url) *UrlsTree {
	sentinel := resource.NewResourceSentinel()
	nodeSentinel := ResourceNode{
		Resource:            sentinel,
		collectionsChildren: map[string]*ResourceNode{},
		identifierChild:     nil,
	}
	urlsTree := &UrlsTree{Root: &nodeSentinel}
	urlsTree.InsertUrls(urls)
	return urlsTree
}

// todo: maybe should use a shorter name
func (tree *UrlsTree) InsertIdentifierResourceNodeInto(
	currentNode *ResourceNode,
	currentResource resource.Resource,
) *ResourceNode {
	if currentNode.identifierChild != nil && currentNode.identifierChild.Resource.Name() != currentResource.Name() {
		log.Fatal("resource identifier is already present")
	} else if currentNode.identifierChild != nil {
		return currentNode.identifierChild
	}

	node := &ResourceNode{
		Resource:            currentResource, //todo: think about the possibility that the slice is gonna change
		collectionsChildren: map[string]*ResourceNode{},
		identifierChild:     nil,
	}
	currentNode.identifierChild = node
	return node
}

func (tree *UrlsTree) InsertCollectionResourceNodeInto(
	currentNode *ResourceNode,
	currentResource resource.Resource,
) *ResourceNode {
	if node, ok := currentNode.collectionsChildren[currentResource.Name()]; ok {
		return node
	}

	newResourceNode := &ResourceNode{
		Resource:            currentResource,
		collectionsChildren: map[string]*ResourceNode{},
		identifierChild:     nil,
	}
	currentNode.collectionsChildren[currentResource.Name()] = newResourceNode
	return newResourceNode
}

func (tree *UrlsTree) insertUrlHelper(currentNode *ResourceNode, currentUrl resource.Url) {
	if len(currentUrl) == 0 {
		return
	}

	currentResource := currentUrl[0]
	switch currentResource.(type) {
	case *resource.ResourceCollection:
		nextResourceNode := tree.InsertCollectionResourceNodeInto(currentNode, currentResource)
		tree.insertUrlHelper(nextResourceNode, currentUrl[1:])
	case *resource.ResourceIdentifier:
		nextResourceNode := tree.InsertIdentifierResourceNodeInto(currentNode, currentResource)
		tree.insertUrlHelper(nextResourceNode, currentUrl[1:])
	default:
		log.Fatal("not possible")
	}
}

func (tree *UrlsTree) InsertUrl(url resource.Url) {
	tree.insertUrlHelper(tree.Root, url)
}

func (tree *UrlsTree) InsertUrls(url []resource.Url) {
	for _, resourcesInFullUrl := range url {
		tree.InsertUrl(resourcesInFullUrl)
	}
}

func (tree *UrlsTree) printTreeHelper(currentNode *ResourceNode) {
	fmt.Println("node: ", currentNode.Resource.Name())
	collectionsChildren := currentNode.collectionsChildren
	for _, collectionsChild := range collectionsChildren {
		tree.printTreeHelper(collectionsChild)
	}
	if currentNode.identifierChild != nil {
		tree.printTreeHelper(currentNode.identifierChild)
	}
}

func (tree *UrlsTree) printTree() {
	tree.printTreeHelper(tree.Root)
}

func (tree *UrlsTree) stringHelper(currentNode *ResourceNode) [][]string {
	allUrls := [][]string{}

	allChildrenUrls := [][][]string{}
	for _, collectionNode := range currentNode.collectionsChildren {
		allChildrenUrls = append(allChildrenUrls, tree.stringHelper(collectionNode))
	}
	if currentNode.identifierChild != nil {
		allChildrenUrls = append(allChildrenUrls, tree.stringHelper(currentNode.identifierChild))
	}

	for _, singleChildUrls := range allChildrenUrls {
		for _, singleChildUrl := range singleChildUrls {
			switch currentNode.Resource.(type) {
			case *resource.ResourceSentinel: //maybe can avoid this check by storing "" as sentinel node's name
			default:
				singleChildUrl = append([]string{currentNode.Resource.Name()}, singleChildUrl...) //need to rethink how to make this more performant
			}
			allUrls = append(allUrls, singleChildUrl)
		}
	}

	if len(allChildrenUrls) == 0 {
		switch currentNode.Resource.(type) {
		case *resource.ResourceSentinel:
		default:
			allUrls = append(allUrls, []string{currentNode.Resource.Name()})
		}
	}

	return allUrls
}

func (tree *UrlsTree) String() string {
	allUrlsString := ""
	allUrls := tree.stringHelper(tree.Root)
	for _, singleUrl := range allUrls {
		singleUrlString := ""
		for _, resourceNodeName := range singleUrl {
			singleUrlString += resourceNodeName + " "
		}
		singleUrlString += "\n"
		allUrlsString += singleUrlString
	}
	return allUrlsString
}
