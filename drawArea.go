package butch

import (
	"fmt"
)

// ProcessRenderer handles tasks to render pool invocation process
type ProcessRenderer struct {
	pool  *Pool
	lines int
}

// Paint renders process
func (pr *ProcessRenderer) Paint() {
	if pr.lines > 0 {
		for i := 0; i < pr.lines; i++ {
			fmt.Print("\033[1A") // move cursor one line up
		}
		pr.lines = 0
	}

	workers := pr.pool.Workers()

	done := 0
	for _, w := range workers {
		if w.IsDone() {
			done++
		}
	}

	pr.lines++
	fmt.Printf(
		"Batch execution. Elapsed: %0.1f sec (%.1f%%)\n",
		pr.pool.Elapsed().Seconds(),
		100.*float32(done)/float32(len(pr.pool.workers)),
	)
	fmt.Print("\033[K") // delete till end of line

	for _, w := range workers {
		pr.lines++
		if w.IsWaiting() {
			// not running
			fmt.Print("")
		} else if !w.IsDone() {
			// In progress
			fmt.Printf("  %5.1f  ", w.Elapsed().Seconds())
		} else {
			// Done
			if len(w.lastErr) > 0 {
				fmt.Print(" " + doneErr(fmt.Sprintf("[%5.1f]", w.Elapsed().Seconds())) + " ")
			} else {
				fmt.Print(" " + doneOk(fmt.Sprintf("[%5.1f]", w.Elapsed().Seconds())) + " ")
			}
		}
		if !w.IsWaiting() && !w.IsDone() {
			fmt.Print(active(w.Name), " ")
		} else {
			fmt.Print(w.Name, " ")
		}

		// Printing last message
		if len(w.lastErr) > 0 {
			fmt.Print(ferror(w.lastErr))
		} else if len(w.lastMsg) > 0 {
			fmt.Print(normal(w.lastMsg))
		}

		fmt.Println()
		fmt.Print("\033[K") // delete till end of line
	}
}

func doneOk(str string) string {
	return "\033[44m" + str + "\033[0m"
}
func doneErr(str string) string {
	return "\033[41m" + str + "\033[0m"
}

func active(str string) string {
	return "\033[1m" + str + "\033[0m"
}

func normal(str string) string {
	return "\033[1;30m" + str + "\033[0m"
}

func ferror(str string) string {
	return "\033[1;31m" + str + "\033[0m"
}
