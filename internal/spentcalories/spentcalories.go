package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var ErrAmountOfElements = errors.New("wrong amount of elements passed")
var ErrZeroOrNegative = errors.New("value is negative or 0")
var ErrEmptyString = errors.New("empty string")
var ErrWrongActivity = errors.New("неизвестный тип тренировки")

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	splitData := strings.Split(data, ",")
	if len(splitData) != 3 {
		return 0, "", 0, ErrAmountOfElements
	}
	steps, err := strconv.Atoi(splitData[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, ErrZeroOrNegative
	}
	activity := splitData[1]
	if activity == "" {
		return 0, "", 0, ErrEmptyString
	}
	duration, err := time.ParseDuration(splitData[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, ErrZeroOrNegative
	}
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distance := stepLength * float64(steps)
	distanceInKm := distance / 1000
	return distanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)
	meanSpeed := distance / duration.Hours()
	return meanSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	durationHours := duration.Hours()
	speed := meanSpeed(steps, height, duration)
	distance := distance(steps, height)
	var calories float64
	switch activity {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", ErrWrongActivity
	}

	return fmt.Sprintf("Тип тренировки: %v\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, durationHours, distance, speed, calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, ErrZeroOrNegative
	}
	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	return (weight * speed * durationMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, ErrZeroOrNegative
	}
	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	return (weight * speed * durationMinutes * walkingCaloriesCoefficient) / minInH, nil
}
