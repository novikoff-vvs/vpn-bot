package yoomoney

import "math"

func FindAmountWithCommission(requiredAmount float64) float64 {
	// Вычисляем исходную сумму
	initialAmount := requiredAmount / 0.90

	// Округляем результат до 1 знака после запятой
	roundedAmount := math.Round(initialAmount*10) / 10

	return roundedAmount
}
