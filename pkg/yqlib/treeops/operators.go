package treeops

import (
	"container/list"
	"fmt"

	"gopkg.in/yaml.v3"
)

type OperatorHandler func(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error)

func UnwrapDoc(node *yaml.Node) *yaml.Node {
	if node.Kind == yaml.DocumentNode {
		return node.Content[0]
	}
	return node
}

func PipeOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	lhs, err := d.GetMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	return d.GetMatchingNodes(lhs, pathNode.Rhs)
}

func createBooleanCandidate(owner *CandidateNode, value bool) *CandidateNode {
	valString := "true"
	if !value {
		valString = "false"
	}
	node := &yaml.Node{Kind: yaml.ScalarNode, Value: valString, Tag: "!!bool"}
	return &CandidateNode{Node: node, Document: owner.Document, Path: owner.Path}
}

func nodeToMap(candidate *CandidateNode) *list.List {
	elMap := list.New()
	elMap.PushBack(candidate)
	return elMap
}

func LengthOperator(d *dataTreeNavigator, matchMap *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- lengthOperation")
	var results = list.New()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		var length int
		switch candidate.Node.Kind {
		case yaml.ScalarNode:
			length = len(candidate.Node.Value)
		case yaml.MappingNode:
			length = len(candidate.Node.Content) / 2
		case yaml.SequenceNode:
			length = len(candidate.Node.Content)
		default:
			length = 0
		}

		node := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", length), Tag: "!!int"}
		lengthCand := &CandidateNode{Node: node, Document: candidate.Document, Path: candidate.Path}
		results.PushBack(lengthCand)
	}

	return results, nil
}
