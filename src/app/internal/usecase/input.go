package usecase

import (
	"errors"
	"math"
	"strconv"

	m "APG2_SmartCalc/app/internal/models"
)

func (u *Usecase) PushNumber(content int) (string, error) {
	node := m.NewNode(m.Digit, m.VeryHigh, m.NewNodeValue(content, 0, m.Integer))

	if u.NodeList.Back() != nil {
		lastNode := u.NodeList.Back().Value.(*m.Node)
		switch lastNode.Type {
		case m.Digit:
			num := lastNode.Value.Content.(int)
			lastNode.Value.Content = num*10 + content

			if lastNode.Value.TypeDigit == m.Float {
				lastNode.Value.Exponent += 1
			}

		case m.CloseBracket,
			m.Xdigit:

			u.PushOperator("*")
			u.NodeList.PushBack(node)
		default:
			u.NodeList.PushBack(node)
		}
	} else {
		u.NodeList.PushBack(node)
		if content == 0 {
			return u.PushPoint()
		}
	}

	return u.MakeOutput(), nil
}

func (u *Usecase) MakeOutput() string {
	var res string
	for i := u.NodeList.Front(); i != nil; i = i.Next() {
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

func (u *Usecase) PushPoint() (string, error) {
	if u.NodeList.Back() != nil {
		lastNode := u.NodeList.Back().Value.(*m.Node)
		if lastNode.Type == m.Digit && lastNode.Value.TypeDigit == m.Integer {
			lastNode.Value.TypeDigit = m.Float
			return u.MakeOutput() + ".", nil
		}
	}
	return u.MakeOutput(), errors.New("incorrect input point")
}

func (u *Usecase) PushOperator(content string) (string, error) {
	priority := m.Low
	if content == "*" || content == "/" || content == "\\%" {
		priority = m.Medium
	}
	if content == "^" {
		priority = m.High
	}
	if u.NodeList.Back() != nil {
		lastNode := u.NodeList.Back().Value.(*m.Node)
		switch lastNode.Type {
		case m.Operator:
			lastNode.Value.Content = content

		case m.OpenBracket:
			return u.MakeOutput(), errors.New("incorrect input operator")
		default:
			node := m.NewNode(m.Operator, priority, m.NewNodeValue(content, 0, 0))
			u.NodeList.PushBack(node)
		}
	} else {
		return u.MakeOutput(), errors.New("incorrect input operator")
	}

	return u.MakeOutput(), nil
}

func (u *Usecase) PushOpenBracket() string {
	if u.NodeList.Back() != nil {
		lastNode := u.NodeList.Back().Value.(*m.Node)
		switch lastNode.Type {
		case m.CloseBracket,
			m.Digit,
			m.Xdigit:
			u.PushOperator("*")
		}
	}

	node := m.NewNode(m.OpenBracket, m.VeryHigh, m.NewNodeValue("(", 0, 0))
	u.NodeList.PushBack(node)
	u.countOpenBrackets++
	return u.MakeOutput()
}

func (u *Usecase) PushCloseBracket() (string, error) {
	if u.NodeList.Back() != nil && u.countOpenBrackets > 0 {
		lastNode := u.NodeList.Back().Value.(*m.Node)
		switch lastNode.Type {
		case m.OpenBracket,
			m.Operator:
			return u.MakeOutput(), errors.New("incorrect input close bracket")
		default:
			node := m.NewNode(m.CloseBracket, m.VeryHigh, m.NewNodeValue(")", 0, 0))
			u.NodeList.PushBack(node)
			u.countOpenBrackets--
		}
	} else {
		return u.MakeOutput(), errors.New("incorrect input close bracket")
	}

	return u.MakeOutput(), nil
}

func (u *Usecase) PushFunction(content string) string {
	if u.NodeList.Back() != nil {
		lastNode := u.NodeList.Back().Value.(*m.Node)
		switch lastNode.Type {
		case m.CloseBracket,
			m.Digit,
			m.Xdigit:
			u.PushOperator("*")
		}
	}

	node := m.NewNode(m.Function, m.High, m.NewNodeValue(content, 0, 0))
	u.NodeList.PushBack(node)
	u.PushOpenBracket()
	return u.MakeOutput()
}

func (u *Usecase) PushClear() string {
	for u.NodeList.Len() != 0 {
		u.NodeList.Remove(u.NodeList.Back())
	}
	for u.NodeListForX.Len() != 0 {
		u.NodeListForX.Remove(u.NodeListForX.Back())
	}
	u.countOpenBrackets = 0
	return u.MakeOutput()
}

func (u *Usecase) PushX() (string, error) {

	if u.NodeListForX.Len() != 0 {
		err := errors.New("cant push 'X' now")
		return u.MakeOutput(), err
	}

	node := m.NewNode(m.Xdigit, m.VeryHigh, m.NewNodeValue("X", 0, 0))

	if u.NodeList.Back() != nil {
		lastNode := u.NodeList.Back().Value.(*m.Node)
		switch lastNode.Type {
		case m.CloseBracket,
			m.Digit,
			m.Xdigit:
			u.PushOperator("*")
		}
	}
	u.NodeList.PushBack(node)
	return u.MakeOutput(), nil
}

func (u *Usecase) PushChangeSign() (string, error) {
	if u.NodeList.Back() != nil {
		lastNode := u.NodeList.Back().Value.(*m.Node)
		if lastNode.Type == m.Digit {
			num := lastNode.Value.Content.(int)
			lastNode.Value.Content = -num

			u.NodeList.Remove(u.NodeList.Back())
			u.PushOpenBracket()
			u.NodeList.PushBack(lastNode)

			return u.MakeOutput(), nil
		}
	}

	return u.MakeOutput(), errors.New("incorrect input change sign")
}

func (u *Usecase) PushDeleteBackNode() string {
	if u.NodeList.Back() != nil {
		lastNode := u.NodeList.Back().Value.(*m.Node)

		switch lastNode.Type {
		case m.OpenBracket:
			u.NodeList.Remove(u.NodeList.Back())
			u.countOpenBrackets--

			lastNode := u.NodeList.Back().Value.(*m.Node)
			if lastNode.Type == m.Function {
				u.NodeList.Remove(u.NodeList.Back())
				return u.MakeOutput()
			}

		case m.CloseBracket:
			u.NodeList.Remove(u.NodeList.Back())
			u.countOpenBrackets++

		case m.Digit:
			num := lastNode.Value.Content.(int)
			lastNode.Value.Content = num / 10
			if lastNode.Value.TypeDigit == m.Float {
				lastNode.Value.Exponent -= 1
				if lastNode.Value.Exponent == 0 {
					lastNode.Value.TypeDigit = m.Integer
				}
			}

			num = lastNode.Value.Content.(int)
			if num == 0 && lastNode.Value.Exponent == 0 {
				u.NodeList.Remove(u.NodeList.Back())
				return u.MakeOutput()
			}

		default:
			u.NodeList.Remove(u.NodeList.Back())
		}
	}

	return u.MakeOutput()
}
