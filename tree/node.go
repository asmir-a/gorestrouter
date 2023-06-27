package tree

import "github.com/asmir-a/gorestrouter/resource"

type ResourceNode struct {
	Resource            resource.Resource
	collectionsChildren map[string]*ResourceNode //two fields responsible for children might reside in another struct; but for now it is okay and it might even be okay for the future
	identifierChild     *ResourceNode
}

func (node *ResourceNode) FindChildWithResourceName(nameOrId string, params map[string]string) *ResourceNode {
	if collectionChild, ok := node.collectionsChildren[nameOrId]; ok {
		return collectionChild
	}
	if node.identifierChild != nil {
		params[node.identifierChild.Resource.Name()] = nameOrId
		return node.identifierChild
	}
	return nil
}
