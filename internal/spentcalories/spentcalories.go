package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining разбирает строку формата "<шаги>,<активность>,<длительность>",
// например "3456,Ходьба,3h00m", на количество шагов, вид активности
// и time.Duration.
func parseTraining(data string) (int, string, time.Duration, error) {
	dataSplit := strings.Split(data, ",")

	if len(dataSplit) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат %q: ожидалось 3 значения, получено %d",
			data, len(dataSplit))
	}

	steps, err := strconv.Atoi(dataSplit[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат количества шагов %q: %w", dataSplit[0], err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("ноличество шагов должно быть положительным, получено %d", steps)
	}

	duration, err := time.ParseDuration(dataSplit[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат времени выполнения %q: %w", dataSplit[2], err)
	}

	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("длительность должна быть положительной, получена %v", duration)
	}

	return steps, dataSplit[1], duration, nil
}

// distance вычисляет пройденную дистанцию в километрах по количеству шагов
// и росту пользователя в метрах.
func distance(steps int, height float64) float64 {
	return (float64(steps) * (height * stepLengthCoefficient)) / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч по количеству шагов,
// росту пользователя в метрах и продолжительности активности.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

// TrainingInfo формирует итоговый отчёт о тренировке.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	switch activityType {
	case "Ходьба":
		dist := distance(steps, height)
		speed := meanSpeed(steps, height, duration)
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("расчёт калорий для ходьбы: %w", err)
		}
		return fmt.Sprintf(
				"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
				activityType, duration.Hours(), dist, speed, calories,
			),
			nil

	case "Бег":
		dist := distance(steps, height)
		speed := meanSpeed(steps, height, duration)
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("расчёт калорий для бега: %w", err)
		}
		return fmt.Sprintf(
				"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
				activityType, duration.Hours(), dist, speed, calories,
			),
			nil

	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %q", activityType)
	}
}

// RunningSpentCalories вычисляет количество калорий, потраченных при беге,
// по количеству шагов, весу и росту пользователя и продолжительности бега.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным, получено %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным, получен %v", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным, получен %v", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("длительность должна быть положительной, получена %v", duration)
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := weight * speed * durationInMinutes / minInH

	return calories, nil

}

// WalkingSpentCalories вычисляет количество калорий, потраченных при ходьбе,
// по количеству шагов, весу и росту пользователя и продолжительности ходьбы.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным, получено %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным, получен %v", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным, получен %v", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("длительность должна быть положительной, получена %v", duration)
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes / minInH) * walkingCaloriesCoefficient

	return calories, nil
}
