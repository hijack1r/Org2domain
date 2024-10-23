package runner

import (
	"fmt"
)

func PrintProgress(currentIndex, totalLoops int) {
	progress := float64(currentIndex) / float64(totalLoops) * 100
	fmt.Printf("\r[+] (%d/%d) Progress: %.2f%%", currentIndex, totalLoops, progress)
}
