package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

var ErrAmountOfElements = errors.New("wrong amount of elements passed")
var ErrZeroOrNegative = errors.New("value is negative or 0")
var ErrEmptyString = errors.New("empty string")
var ErrWrongActivity = errors.New("неизвестный тип тренировки")

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

func (t *Training) Parse(datastring string) (err error) {
	splitData := strings.Split(datastring, ",")
	if len(splitData) != 3 {
		return ErrAmountOfElements
	}
	steps, err := strconv.Atoi(splitData[0])
	if err != nil {
		return err
	}
	if steps <= 0 {
		return ErrZeroOrNegative
	}
	t.Steps = steps

	activity := splitData[1]
	if activity == "" {
		return ErrEmptyString
	}
	t.TrainingType = activity

	duration, err := time.ParseDuration(splitData[2])
	if err != nil {
		return err
	}
	if duration <= 0 {
		return ErrZeroOrNegative
	}
	t.Duration = duration
	return nil
}

func (t Training) ActionInfo() (string, error) {
	if t.TrainingType != "Ходьба" && t.TrainingType != "Бег" {
		return "", ErrWrongActivity
	}
	distance := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)
	var energy float64
	if t.TrainingType == "Ходьба" {
		energySpent, err := spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
		if err != nil {
			return "", err
		}
		energy = energySpent
	} else {
		energySpent, err := spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
		if err != nil {
			return "", err
		}
		energy = energySpent
	}
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", t.TrainingType, float64(t.Duration.Hours()), distance, speed, energy), nil
}
