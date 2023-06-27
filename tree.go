package main

import (
	"fmt"
	"log"
)

type ResourceNode struct {
	resource            Resource
	collectionsChildren map[string]*ResourceNode //two fields responsible for children might reside in another struct; but for now it is okay and it might even be okay for the future
	identifierChild     *ResourceNode
}

func (node *ResourceNode) FindChildWithResourceName(nameOrId string, params map[string]string) *ResourceNode {
	if collectionChild, ok := node.collectionsChildren[nameOrId]; ok {
		return collectionChild
	}
	if node.identifierChild != nil {
		params[node.identifierChild.Name()] = nameOrId
		return node.identifierChild
	}
	return nil
}

func (node *ResourceNode) Name() string {
	return node.resource.Name()
}

type UrlsTree struct {
	root *ResourceNode
}

func NewUrlsTree(urls []Url) *UrlsTree {
	sentinel := &ResourceSentinel{name: "sentinel", setUpHandler: nil}
	nodeSentinel := ResourceNode{
		resource:            sentinel,
		collectionsChildren: map[string]*ResourceNode{},
		identifierChild:     nil,
	}
	urlsTree := &UrlsTree{root: &nodeSentinel}
	urlsTree.InsertUrls(urls)
	return urlsTree
}

// todo: maybe should use a shorter name
func (tree *UrlsTree) InsertIdentifierResourceNodeInto(currentNode *ResourceNode, currentResource Resource) *ResourceNode {
	if currentNode.identifierChild != nil && currentNode.identifierChild.Name() != currentResource.Name() {
		log.Fatal("resource identifier is already present")
	} else if currentNode.identifierChild != nil {
		return currentNode.identifierChild
	}

	node := &ResourceNode{
		resource:            currentResource, //todo: think about the possibility that the slice is gonna change
		collectionsChildren: map[string]*ResourceNode{},
		identifierChild:     nil,
	}
	currentNode.identifierChild = node
	return node
}

func (tree *UrlsTree) InsertCollectionResourceNodeInto(currentNode *ResourceNode, currentResource Resource) *ResourceNode {
	if node, ok := currentNode.collectionsChildren[currentResource.Name()]; ok {
		return node
	}

	newResourceNode := &ResourceNode{
		resource:            currentResource,
		collectionsChildren: map[string]*ResourceNode{},
		identifierChild:     nil,
	}
	currentNode.collectionsChildren[currentResource.Name()] = newResourceNode
	return newResourceNode
}

func (tree *UrlsTree) insertUrlHelper(currentNode *ResourceNode, currentUrl Url) {
	if len(currentUrl) == 0 {
		return
	}

	currentResource := currentUrl[0]
	switch currentResource.(type) {
	case *ResourceCollection:
		nextResourceNode := tree.InsertCollectionResourceNodeInto(currentNode, currentResource)
		tree.insertUrlHelper(nextResourceNode, currentUrl[1:])
	case *ResourceIdentifier:
		nextResourceNode := tree.InsertIdentifierResourceNodeInto(currentNode, currentResource)
		tree.insertUrlHelper(nextResourceNode, currentUrl[1:])
	default:
		log.Fatal("not possible")
	}
}

func (tree *UrlsTree) InsertUrl(url Url) {
	tree.insertUrlHelper(tree.root, url)
}

func (tree *UrlsTree) InsertUrls(url []Url) {
	for _, resourcesInFullUrl := range url {
		tree.InsertUrl(resourcesInFullUrl)
	}
}

func (tree *UrlsTree) printTreeHelper(currentNode *ResourceNode) {
	fmt.Println("node: ", currentNode.Name())
	collectionsChildren := currentNode.collectionsChildren
	for _, collectionsChild := range collectionsChildren {
		tree.printTreeHelper(collectionsChild)
	}
	if currentNode.identifierChild != nil {
		tree.printTreeHelper(currentNode.identifierChild)
	}
}

func (tree *UrlsTree) printTree() {
	tree.printTreeHelper(tree.root)
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
			switch currentNode.resource.(type) {
			case *ResourceSentinel: //maybe can avoid this check by storing "" as sentinel node's name
			default:
				singleChildUrl = append([]string{currentNode.Name()}, singleChildUrl...) //need to rethink how to make this more performant
			}
			allUrls = append(allUrls, singleChildUrl)
		}
	}

	if len(allChildrenUrls) == 0 {
		switch currentNode.resource.(type) {
		case *ResourceSentinel:
		default:
			allUrls = append(allUrls, []string{currentNode.Name()})
		}
	}

	return allUrls
}

func (tree *UrlsTree) String() string {
	allUrlsString := ""
	allUrls := tree.stringHelper(tree.root)
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
