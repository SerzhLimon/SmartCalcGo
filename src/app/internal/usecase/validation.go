package usecase

import (
	"errors"

	m "APG2_SmartCalc/app/internal/models"
)

func (u *Usecase) validation() error {
	if u.NodeList.Len() == 0 {
		return errors.New("empty input")
	}
	if u.countOpenBrackets != 0 {
		return errors.New("invalid count brackets")
	}
	lastNode := u.NodeList.Back().Value.(*m.Node)
	switch lastNode.Type {
		case m.Operator,
			m.OpenBracket:
			return errors.New("invalid input")
	}

	return nil
}

func (u *Usecase) haveXsInExpression() bool {
	for e := u.NodeList.Front(); e != nil; e = e.Next() {
		node := e.Value.(*m.Node)
		if node.Type == m.Xdigit {
			return true
		}
	}
	return false
}