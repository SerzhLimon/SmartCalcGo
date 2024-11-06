package models

type TypeNode int
type TypeDigit int
type Operators int
type Functions int
type Prioty int

type Node struct {
	Type      TypeNode
	Priotity  Prioty
	Value     *NodeContent
}

type NodeContent struct {
	Content   any
	Exponent  int
	TypeDigit TypeDigit
}

func NewNodeValue(content any, exp int, typ TypeDigit) *NodeContent {
	return &NodeContent{
		Content: content,
		Exponent: exp,
		TypeDigit: typ,
	}
}

func NewNode(types TypeNode, priority Prioty, value *NodeContent) *Node {
	return &Node{
		Type:      types,
		Priotity:  priority,
		Value:     value,
	}
}

func (n *Node) Copy() *Node {
	content := &NodeContent{
		Content: n.Value.Content,
		Exponent: n.Value.Exponent,
		TypeDigit: n.Value.TypeDigit,
	}
	return &Node{
		Type: n.Type,
		Priotity: n.Priotity,
		Value : content,
	}
}

const (
	Integer  TypeDigit = 1
	Float    TypeDigit = 2

	Digit        TypeNode = 1
	Operator     TypeNode = 2
	Function     TypeNode = 3
	OpenBracket  TypeNode = 4
	CloseBracket TypeNode = 5
	Xdigit       TypeNode = 6

	Low      Prioty = 1
	Medium   Prioty = 2
	High     Prioty = 3
	VeryHigh Prioty = 4
)
