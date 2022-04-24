package merkle

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed real_leaf_node_hashes.txt
var realLeafNodeHashes string

func TestNewTree_CommaAppend(t *testing.T) {
	testCases := []struct {
		leafNodeHashes []string
		expectErr      bool
		expectRootHash string
	}{
		{
			leafNodeHashes: nil,
			expectErr:      true,
		},
		{
			leafNodeHashes: []string{"A"},
			expectErr:      false,
			expectRootHash: "A",
		},
		{
			leafNodeHashes: []string{"A", "B"},
			expectErr:      false,
			expectRootHash: "A, B",
		},
		{
			leafNodeHashes: []string{"A", "B", "C", "D", "E", "F"},
			expectErr:      false,
			expectRootHash: "A, B, C, D, E, F, E, F",
		},
		{
			leafNodeHashes: []string{"A", "B", "C", "D", "E", "F", "G", "H"},
			expectErr:      false,
			expectRootHash: "A, B, C, D, E, F, G, H",
		},
	}

	for _, tc := range testCases {
		tree, err := NewTree(tc.leafNodeHashes, CommaAppendCombineFn)
		if tc.expectErr {
			assert.Error(t, err)
			assert.Nil(t, tree)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expectRootHash, tree.GetRootHash())
		}
	}
}

func TestNewTree_RealLeafNodeHashes(t *testing.T) {
	hashes := strings.Split(realLeafNodeHashes, "\n")
	assert.Equal(t, "93a9b791b88a42c3f5c61b89a7cdc42c3f09843c0f982fdf2df4f90d0254c57b", hashes[0])
	assert.Equal(t, "7cc686e0dcb94e021c3c386ab21a2c37b45ff0268bbe49f7069b335b43a964c4", hashes[len(hashes)-1])
	assert.Equal(t, 1577, len(hashes))
	tree, err := NewTree(hashes, Sha256CombineHashFn)
	assert.NoError(t, err)
	assert.Equal(t, "a9eebd039eaefae9e31893e315a31f55e5188c4ee659025eb3ba6e1df2d5ed06", tree.GetRootHash())
}
