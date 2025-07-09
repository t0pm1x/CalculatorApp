package calculationService

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
)

type CalculationService interface {
	CreateCalculation(expression string) (Calculation, error)
	GetAllCalculations() ([]Calculation, error)
	GetCalculationById(id string) (Calculation, error)
	UpdateCalculation(id string, expression string) (Calculation, error)
	DeleteCalculation(id string) error
}

type calcService struct {
	repo CalculationRepository
}

func NewCalculationService(r CalculationRepository) CalculationService {
	return &calcService{repo: r}
}

func (s *calcService) calculateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression) //Создаем выражение
	if err != nil {

		return "", err //Передали выражение
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result), err
}

func (s *calcService) CreateCalculation(expression string) (Calculation, error) {
	result, err := s.calculateExpression(expression)
	if err != nil {
		return Calculation{}, err
	}

	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: expression,
		Result:     result,
	}

	if err := s.repo.CreateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, nil
}

func (s *calcService) GetAllCalculations() ([]Calculation, error) {
	return s.repo.GetAllCalculations()
}

func (s *calcService) GetCalculationById(id string) (Calculation, error) {
	return s.repo.GetCalculationById(id)
}

func (s *calcService) UpdateCalculation(id string, expression string) (Calculation, error) {
	calc, err := s.repo.GetCalculationById(id)
	if err != nil {
		return Calculation{}, err
	}

	result, err := s.calculateExpression(expression)
	if err != nil {
		return Calculation{}, err
	}

	calc.Expression = expression
	calc.Result = result

	if err := s.repo.UpdateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, nil
}

func (s *calcService) DeleteCalculation(id string) error {
	return s.repo.DeleteCalculation(id)
}
