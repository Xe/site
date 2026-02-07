package templates

import "fmt"

// formatDollars converts cents to a dollar string representation.
func formatDollars(cents int) string {
	return fmt.Sprintf("%.2f", float64(cents)/100)
}
