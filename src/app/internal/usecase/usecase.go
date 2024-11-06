package usecase

import (
	"container/list"
	"errors"
	"math"
	"strconv"

	"APG2_SmartCalc/app/internal/models"
)

type Strorage interface {
	Get() ([][]string, error)
	Set(key, value string) error
}

type Usecase struct {
	NodeList          *list.List
	NodeListForX      *list.List
	countOpenBrackets int
	operations        map[string]calcFunc
	storage           Strorage
}

func New(storage Strorage) *Usecase {
	nodeList := list.New()
	nodeListForX := list.New()
	operations := bindType()
	return &Usecase{
		NodeList:     nodeList,
		NodeListForX: nodeListForX,
		operations:   operations,
		storage:      storage,
	}
}

func bindType() map[string]calcFunc {
	operations := map[string]calcFunc{
		"+":    calcPlus,
		"-":    calcMinus,
		"*":    calcMult,
		"/":    calcDiv,
		"%":    calcMod,
		"^":    calcPow,
		"log":  calcLog,
		"ln":   calcLn,
		"sin":  calcSin,
		"asin": calcAsin,
		"cos":  calcCos,
		"acos": calcAcos,
		"tan":  calcTan,
		"atan": calcAtan,
		"sqrt": calcSqrt,
	}
	return operations
}

func (u *Usecase) MakeResult(input string) (string, error) {

	mustSave := true
	if u.haveXsInExpression() {
		u.NodeListForX = u.copyList(u.NodeList)
		for u.NodeList.Len() != 0 {
			u.NodeList.Remove(u.NodeList.Back())
		}
		err := errors.New("you must input x")
		return u.MakeOutput(), err
	}
	if u.NodeListForX.Len() != 0 {
		mustSave = false
		u.insertXsValues()
	}
	if err := u.makeResult(); err != nil {
		return u.MakeOutput(), err
	}

	node := u.NodeList.Back().Value.(*models.Node)
	num := node.Value.Content.(int)

	contentNum := float64(num)
	contentNum = contentNum / math.Pow(10, float64(node.Value.Exponent))

	man, exp := math.Modf(contentNum)
	if exp == 0 {
		node.Value.Content = int(man)
		node.Value.Exponent = 0
		node.Value.TypeDigit = models.Integer
	}
	if mustSave {
		stringValue := strconv.FormatFloat(contentNum, 'f', -1, 64)
		u.storage.Set(input, stringValue)
	}

	return u.MakeOutput(), nil
}

func (u *Usecase) makeResult() error {
	if err := u.validation(); err != nil {
		return err
	}
	u.MakePolishSolve()
	if err := u.Calculate(); err != nil {
		u.PushClear()
		return err
	}

	u.countOpenBrackets = 0
	return nil
}

func (u *Usecase) needCalcOperation(el *list.Element) bool {
	v1 := el.Prev().Value.(*models.Node)
	v2 := el.Prev().Prev().Value.(*models.Node)
	if v1.Type == models.Digit && v2.Type == models.Digit {
		return true
	}
	return false
}

func (u *Usecase) needCalcFunction(el *list.Element) bool {
	v1 := el.Prev().Value.(*models.Node)
	if v1.Type == models.Digit {
		return true
	}
	return false
}

func (u *Usecase) calcOpeartion(el *list.Element) error {
	v2 := el.Prev().Value.(*models.Node)
	v1 := el.Prev().Prev().Value.(*models.Node)
	resNode := el.Value.(*models.Node)

	function := u.operations[resNode.Value.Content.(string)]

	num1 := v1.Value.Content.(int)
	contentNumV1 := float64(num1)
	if v1.Value.TypeDigit == models.Float {
		contentNumV1 = contentNumV1 / math.Pow(10, float64(v1.Value.Exponent))
	}

	num2 := v2.Value.Content.(int)
	contentNumV2 := float64(num2)
	if v2.Value.TypeDigit == models.Float {
		contentNumV2 = contentNumV2 / math.Pow(10, float64(v2.Value.Exponent))
	}

	if math.IsNaN(contentNumV1) || math.IsNaN(contentNumV2) {
		err := errors.New("result is not a number")
		return err
	}

	result := function(contentNumV1, contentNumV2)
	if math.IsNaN(result) || math.IsInf(result, 0) || math.IsInf(result, 1) {
		err := errors.New("result is not a number")
		return err
	}
	u.updateNode(resNode, result)

	u.NodeList.Remove(el.Prev())
	u.NodeList.Remove(el.Prev())
	return nil
}

func (u *Usecase) calcFunction(el *list.Element) error {
	v := el.Prev().Value.(*models.Node)
	resNode := el.Value.(*models.Node)

	function := u.operations[resNode.Value.Content.(string)]

	num := v.Value.Content.(int)
	contentNum := float64(num)
	if v.Value.TypeDigit == models.Float {
		contentNum = contentNum / math.Pow(10, float64(v.Value.Exponent))
	}

	result := function(contentNum, 0)
	if math.IsNaN(result) || math.IsInf(result, 0) || math.IsInf(result, 1) {
		err := errors.New("result is not a number")
		return err
	}
	u.updateNode(resNode, result)

	u.NodeList.Remove(el.Prev())
	return nil
}

func (u *Usecase) Calculate() error {
	el := u.NodeList.Back()
	for u.NodeList.Len() != 1 {
		node := el.Value.(*models.Node)
		switch node.Type {
		case models.Operator:
			if u.needCalcOperation(el) {
				if err := u.calcOpeartion(el); err != nil {
					return err
				}
			}
		case models.Function:
			if u.needCalcFunction(el) {
				if err := u.calcFunction(el); err != nil {
					return err
				}
			}
		}
		el = el.Prev()
		if el == u.NodeList.Front() || el == nil {
			el = u.NodeList.Back()
		}
	}

	return nil
}

func (u *Usecase) makeContent(a float64) int {
	mantisa := a * math.Pow(10, 7)
	return int(mantisa)
}

func (u *Usecase) updateNode(resNode *models.Node, result float64) {
	resNode.Value.Content = u.makeContent(result)
	resNode.Value.Exponent = 7
	resNode.Type = models.Digit
	resNode.Value.TypeDigit = models.Float
}

func (u *Usecase) GetHistory() ([][]string, error) {
	history, err := u.storage.Get()
	if err != nil {
		return nil, err
	}
	return history, err
}

func (u *Usecase) insertXsValues() {

	open := models.NewNode(models.OpenBracket, models.VeryHigh, models.NewNodeValue("(", 0, 0))
	close := models.NewNode(models.CloseBracket, models.VeryHigh, models.NewNodeValue(")", 0, 0))
	u.NodeList.PushFront(open)
	u.NodeList.PushBack(close)

	for e := u.NodeListForX.Front(); e != nil; e = e.Next() {
		node := e.Value.(*models.Node)
		if node.Type == models.Xdigit {
			for el := u.NodeList.Front(); el != nil; el = el.Next() {
				nodePush := el.Value.(*models.Node)
				u.NodeListForX.InsertBefore(nodePush, e)
			}
			u.NodeListForX.Remove(e)
		}
	}

	u.NodeList = u.copyList(u.NodeListForX)
	for u.NodeListForX.Len() != 0 {
		u.NodeListForX.Remove(u.NodeListForX.Back())
	}
}
