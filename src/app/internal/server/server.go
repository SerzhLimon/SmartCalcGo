package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v4"

	"APG2_SmartCalc/app/internal/models"
	uc "APG2_SmartCalc/app/internal/usecase"
)

type Server struct {
	usecase uc.Usecase
}

func NewServer(storage uc.Strorage) *Server {
	usecase := uc.New(storage)
	return &Server{
		usecase: *usecase,
	}
}

func (s *Server) ServeHTML(c *gin.Context) {
	s.usecase.PushClear()
	c.HTML(http.StatusOK, "index.html", "")
}

func (s *Server) PushNumber(c *gin.Context) {
	var input models.PushNumberRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := null.StringFrom("")
	output, err := s.usecase.PushNumber(input.Value)
	if err != nil {
		message = null.StringFrom(err.Error())
	}

	c.JSON(http.StatusOK, models.OutputResponse{Output: output, Error: message})
}

func (s *Server) PushPoint(c *gin.Context) {
	message := null.StringFrom("")

	output, err := s.usecase.PushPoint()
	if err != nil {
		message = null.StringFrom(err.Error())
	}
	c.JSON(http.StatusOK, models.OutputResponse{Output: output, Error: message})
}

func (s *Server) PushOperator(c *gin.Context) {
	var input models.PushOperatorRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := null.StringFrom("")
	output, err := s.usecase.PushOperator(input.Value)
	if err != nil {
		message = null.StringFrom(err.Error())
	}
	c.JSON(http.StatusOK, models.OutputResponse{Output: output, Error: message})
}

func (s *Server) PushFunction(c *gin.Context) {
	var input models.PushOperatorRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output := s.usecase.PushFunction(input.Value)
	c.JSON(http.StatusOK, models.OutputResponse{Output: output, Error: null.StringFrom("")})
}

func (s *Server) PushOpenBracket(c *gin.Context) {
	output := s.usecase.PushOpenBracket()
	c.JSON(http.StatusOK, models.OutputResponse{Output: output, Error: null.StringFrom("")})
}

func (s *Server) PushCloseBracket(c *gin.Context) {
	message := null.StringFrom("")

	output, err := s.usecase.PushCloseBracket()
	if err != nil {
		message = null.StringFrom(err.Error())
	}
	c.JSON(http.StatusOK, models.OutputResponse{Output: output, Error: message})
}

func (s *Server) PushX(c *gin.Context) {
	message := null.StringFrom("")

	output, err := s.usecase.PushX()
	if err != nil {
		message = null.StringFrom(err.Error())
	}
	c.JSON(http.StatusOK, models.OutputResponse{Output: output, Error: message})
}

func (s *Server) PushChangeSign(c *gin.Context) {
	message := null.StringFrom("")

	output, err := s.usecase.PushChangeSign()
	if err != nil {
		message = null.StringFrom(err.Error())
	}
	c.JSON(http.StatusOK, models.OutputResponse{Output: output, Error: message})
}

func (s *Server) PushClear(c *gin.Context) {
	output := s.usecase.PushClear()
	c.JSON(http.StatusOK, models.OutputResponse{Output: output, Error: null.StringFrom("")})
}

func (s *Server) PushDeleteBackNode(c *gin.Context) {
	output := s.usecase.PushDeleteBackNode()
	c.JSON(http.StatusOK, models.OutputResponse{Output: output, Error: null.StringFrom("")})
}

func (s *Server) Calc(c *gin.Context) {
	var input models.MakeResultRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message := null.StringFrom("")
	output, err := s.usecase.MakeResult(input.Input)
	if err != nil {
		message = null.StringFrom(err.Error())
	}
	c.JSON(http.StatusOK, models.OutputResponse{Output: output, Error: message})
}

func (s *Server) MakePlots(c *gin.Context) {

	var input models.GraphRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := null.StringFrom("")
	img, err := s.usecase.MakeGraph(input.Input, input.Xmin, input.Xmax)
	if err != nil {
		message = null.StringFrom(err.Error())
	}

	c.JSON(http.StatusOK, models.GraphResponse{Image: img, Error: message})
}

func (s *Server) GetHistory(c *gin.Context) {

	history, _ := s.usecase.GetHistory()
	c.JSON(http.StatusOK, models.HistoryResponse{History: history})
}
