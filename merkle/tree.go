package merkle

import "errors"

type (
	Tree interface {
		GetRootHash() string
	}

	CombineHashFn func(string, string) string

	tree struct {
		root *node
	}

	node struct {
		hash      string
		left      *node
		right     *node
		duplicate bool
	}
)

func NewTree(leafHashes []string, comboFn CombineHashFn) (Tree, error) {
	root, err := buildTree(leafHashes, comboFn)
	if err != nil {
		return nil, err
	}
	return &tree{
		root: root,
	}, nil
}

func (t *tree) GetRootHash() string {
	return t.root.hash
}

func buildTree(leafHashes []string, comboFn CombineHashFn) (*node, error) {
	if len(leafHashes) == 0 {
		return nil, errors.New("empty list of leaf node hashes provided, cannot construct tree")
	}
	currentLayer := make([]*node, len(leafHashes), len(leafHashes))
	for i := 0; i < len(leafHashes); i++ {
		currentLayer[i] = &node{
			hash:      leafHashes[i],
			left:      nil,
			right:     nil,
			duplicate: false,
		}
	}

	for len(currentLayer) > 1 {
		if len(currentLayer)%2 == 1 {
			lastNode := currentLayer[len(currentLayer)-1]
			duplicateNode := &node{
				hash:      lastNode.hash,
				left:      nil,
				right:     nil,
				duplicate: true,
			}
			currentLayer = append(currentLayer, duplicateNode)
		}

		nextLayerLen := len(currentLayer) / 2
		nextLayer := make([]*node, nextLayerLen, nextLayerLen)
		for i := 0; i < len(currentLayer); i += 2 {
			nextLayer[i/2] = &node{
				hash:      comboFn(currentLayer[i].hash, currentLayer[i+1].hash),
				left:      currentLayer[i],
				right:     currentLayer[i+1],
				duplicate: false,
			}
		}
		currentLayer = nextLayer
	}

	return currentLayer[0], nil
}
