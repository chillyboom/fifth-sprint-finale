package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

var errAmountOfElements = errors.New("wrong amount of elements passed")
var errZeroOrNegative = errors.New("value is negative or 0")
var errEmptyString = errors.New("empty string")

type DaySteps struct {
	Steps    int
	Duration time.Duration
	personaldata.Personal
}

func (ds *DaySteps) Parse(datastring string) (err error) {
	splitData := strings.Split(datastring, ",")
	if len(splitData) != 2 {
		return errAmountOfElements
	}
	steps, err := strconv.Atoi(splitData[0])
	if err != nil {
		return err
	}
	if steps <= 0 {
		return errZeroOrNegative
	}
	ds.Steps = steps

	duration, err := time.ParseDuration(splitData[1])
	if err != nil {
		return err
	}
	if duration <= 0 {
		return errZeroOrNegative
	}
	ds.Duration = duration
	return nil
}

func (ds DaySteps) ActionInfo() (string, error) {
	distance := spentenergy.Distance(ds.Steps, ds.Height)
	energy, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", ds.Steps, distance, energy), nil
}
