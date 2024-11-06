package usecase

import (
	"container/list"
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	m "APG2_SmartCalc/app/internal/models"
)

const (
	directory = "frontend/plots"
)

func (u *Usecase) xChanger(valueX float64) {
	for e := u.NodeList.Front(); e != nil; e = e.Next() {
		node := e.Value.(*m.Node)
		if node.Type == m.Xdigit {
			u.updateNode(node, valueX)
		}
	}
}

func (u *Usecase) swap(xMin, xMax *float64) {
	*xMin, *xMax = *xMax, *xMin
}

func (u *Usecase) makePoints(xMins, xMaxs string) [][][]float64 {
	xMin, _ := strconv.ParseFloat(xMins, 64)
	xMax, _ := strconv.ParseFloat(xMaxs, 64)

	if xMin > xMax {
		u.swap(&xMin, &xMax)
	}

	pointsArr := make([][][]float64, 0)
	points := make([][]float64, 0)
	step := (xMax - xMin) / 1000

	tmpList := u.copyList(u.NodeList)

	for xMin <= xMax {
		
		u.xChanger(xMin)
		if err := u.makeResult(); err != nil {

			xMin = -xMin + step


			pointsArr = append(pointsArr, points)
			points = make([][]float64, 0)
			u.NodeList = u.copyList(tmpList)

			continue
		}

		point := make([]float64, 2)
		point[0] = xMin

		node := u.NodeList.Back().Value.(*m.Node)
		num := node.Value.Content.(int)

		contentNum := float64(num)
		if node.Value.TypeDigit == m.Float {
			contentNum = contentNum / math.Pow(10, float64(node.Value.Exponent))
		}


		point[1] = contentNum
		points = append(points, point)

		xMin += step

		u.PushClear()
		u.NodeList = u.copyList(tmpList)
	}
	pointsArr = append(pointsArr, points)
	u.PushClear()

	return pointsArr
}

func (u *Usecase) copyList(old *list.List) *list.List {
	new := list.New()
	for e := old.Front(); e != nil; e = e.Next() {
		new.PushBack(e.Value.(*m.Node).Copy())
	}
	return new
}

func (u *Usecase) MakeGraph(input, xMins, xMaxs string) (string, error) {
	if xMins == "0" && xMaxs == "0" || input == "" {
		u.PushClear()
		err := errors.New("cannot make plot")
		return "", err
	}

	pointsArr := u.makePoints(xMins, xMaxs)
	p := plot.New()

	count := 0
	for _, points := range pointsArr {
		pts := make(plotter.XYs, len(points))
		for i := range pts {
			pts[i].X = float64(points[i][0])
			pts[i].Y = float64(points[i][1])
		}
		line, err := plotter.NewLine(pts)
		if err != nil {
			err = errors.Wrapf(err, "cannot make plot")
			return "", err
		}
		p.Add(line)
		count++
	}

	p.Title.Text = "Y = " + input

	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(directory, fmt.Sprintf("plot.jpg"))

	if err := p.Save(6*vg.Inch, 6*vg.Inch, filePath); err != nil {
		err = errors.Wrapf(err, "cannot make plot")
		return "", err
	}

	img, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(img)

	return encoded, nil
}
