package process

type trieNode struct {
	children map[string]*trieNode
	proc     *process
}

type trie struct {
	node *trieNode
}

func newNode() *trieNode {
	return &trieNode{
		children: make(map[string]*trieNode),
	}
}

func (n *trieNode) isLeaf() bool {
	return n.proc != nil
}
