package spentenergy

import (
	"errors"
	"time"
)

var ErrAmountOfElements = errors.New("wrong amount of elements passed")
var ErrZeroOrNegative = errors.New("value is negative or 0")
var ErrEmptyString = errors.New("empty string")

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, ErrZeroOrNegative
	}
	speed := MeanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	return (weight * speed * durationMinutes * walkingCaloriesCoefficient) / minInH, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, ErrZeroOrNegative
	}
	speed := MeanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	return (weight * speed * durationMinutes) / minInH, nil
}

func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := Distance(steps, height)
	meanSpeed := distance / duration.Hours()
	return meanSpeed
}

func Distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distance := stepLength * float64(steps)
	distanceInKm := distance / 1000
	return distanceInKm
}
