package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage разбирает строку формата "<шаги>,<длительность>",
// например "678,0h50m", на количество шагов и time.Duration.
func parsePackage(data string) (int, time.Duration, error) {
	dataSplit := strings.Split(data, ",")

	if len(dataSplit) != 2 {
		return 0, 0, fmt.Errorf("неверный формат %q: ожидалось 2 значения, получено %d",
			data, len(dataSplit))
	}

	steps, err := strconv.Atoi(dataSplit[0])
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат количества шагов %q: %w", dataSplit[0], err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным, получено %d", steps)
	}

	duration, err := time.ParseDuration(dataSplit[1])
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат времени выполнения %q: %w", dataSplit[1], err)
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("длительность должна быть положительной, получена %v", duration)
	}

	return steps, duration, nil
}

// DayActionInfo формирует текстовый отчёт о прогулке: количество шагов,
// пройденную дистанцию в километрах и потраченные калории.
//
// Параметры:
//   - data — строка вида "678,0h50m" (шаги и длительность);
//   - weight — вес пользователя в килограммах;
//   - height — рост пользователя в метрах.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceInKm := float64(steps) * stepLength / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceInKm, calories)
}
