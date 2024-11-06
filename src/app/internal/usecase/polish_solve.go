package usecase

import (
	"container/list"
	"math"
	"strconv"

	m "APG2_SmartCalc/app/internal/models"
)

func (u *Usecase) MakePolishSolve() string {
	polishList := list.New()
	tmpList := list.New()

	for i := u.NodeList.Front(); i != nil; i = i.Next() {
		node := i.Value.(*m.Node)
		switch node.Type {
		case m.OpenBracket:
			tmpList.PushBack(node)
		case m.Xdigit,
			m.Digit:
			polishList.PushBack(node)
		case m.CloseBracket:
			u.pushBetweenBrackets(polishList, tmpList)
		default:
			if u.needPushInPolishList(tmpList, node) {
				u.pushHighPriorytyNodes(polishList, tmpList, node)
			}
			tmpList.PushBack(node)
		}
	}

	u.pushOtherInPolish(polishList, tmpList)
	u.NodeList = polishList

	return u.MakeOutput()
}

func (u *Usecase) pushOtherInPolish(polish, tmp *list.List) {
	for tmp.Len() != 0 {
		polish.PushBack(tmp.Back().Value)
		tmp.Remove(tmp.Back())
	}
}

func (u *Usecase) pushHighPriorytyNodes(polish, tmp *list.List, node *m.Node) {

	for i := tmp.Back(); i != nil; i = i.Prev() {
		lastTmpNode := i.Value.(*m.Node)
		if lastTmpNode.Type == m.OpenBracket || tmp.Len() == 0 || lastTmpNode.Priotity < node.Priotity {
			break
		}
		polish.PushBack(tmp.Back().Value)
		tmp.Remove(tmp.Back())
	}
}

func (u *Usecase) pushBetweenBrackets(polish, tmp *list.List) {
	for tmp.Len() != 0 {
		lastNode := tmp.Back().Value.(*m.Node)
		if lastNode.Type == m.OpenBracket {
			tmp.Remove(tmp.Back())
			break
		}
		polish.PushBack(lastNode)
		tmp.Remove(tmp.Back())
	}
}

func (u *Usecase) needPushInPolishList(tmp *list.List, node *m.Node) bool {
	if tmp.Len() != 0 {
		tmpBack := tmp.Back().Value.(*m.Node)
		if tmpBack.Priotity >= node.Priotity && tmpBack.Type != m.OpenBracket {
			return true
		}
	}
	return false
}

func MakeOutput1(list *list.List) string {
	var res string
	for i := list.Front(); i != nil; i = i.Next() {
		node := i.Value.(*m.Node)
		switch node.Type {
		case m.Digit:
			num := node.Value.Content.(int)
			contentNum := float64(num)
			if node.Value.TypeDigit == m.Float {
				contentNum = contentNum / math.Pow(10, float64(node.Value.Exponent))
			}

			numStr := strconv.FormatFloat(float64(contentNum), 'f', -1, 64)

			res += numStr
		default:
			content, _ := node.Value.Content.(string)
			res += content
		}
	}
	return res
}
