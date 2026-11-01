package main

type Stats struct {
	Count         int
	Sum, Min, Max int64
}

func Calc(nums []int64) Stats {
	if len(nums) < 2 {
		return Stats{}
	}

	// посчитаем первую разницу вне цикла для задания начальных значений
	diff := nums[1] - nums[0]
	stats := Stats{Count: 1, Sum: diff, Min: diff, Max: diff}

	// и начнём цикл со 2-й разницы
	prev := nums[1]
	for _, next := range nums[2:] {
		diff = next - prev
		switch {
		case diff < stats.Min:
			stats.Min = diff
		case diff > stats.Max:
			stats.Max = diff
		}

		stats.Sum += diff
		stats.Count++
		prev = next
	}

	return stats
}
