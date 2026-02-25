package reporter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/marce/download-cleaner/scanner"
)

// Write generates a plain-text report at reportPath describing the operation.
// moved is the list of entries that were (or would be, in dry-run) moved.
// skipped is the list of entries with unrecognised extensions left in place.
// dryRun controls the mode label in the report header.
func Write(reportPath string, moved []scanner.FileEntry, skipped []scanner.FileEntry, dryRun bool) error {
	f, err := os.Create(reportPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	mode := "LIVE RUN"
	if dryRun {
		mode = "DRY RUN (no files were moved)"
	}

	sep := strings.Repeat("─", 70)

	fmt.Fprintf(w, "Download Cleaner Report\n")
	fmt.Fprintf(w, "Generated : %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "Mode      : %s\n", mode)
	fmt.Fprintf(w, "%s\n\n", sep)

	// Section 1: moved files table
	fmt.Fprintf(w, "MOVED FILES (%d total)\n\n", len(moved))
	if len(moved) > 0 {
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "FILE\tCATEGORY\tFLAGS\tDESTINATION")
		fmt.Fprintln(tw, "────\t────────\t─────\t───────────")
		for _, fe := range moved {
			flags := ""
			if fe.IsOld {
				flags = "[OLD >365d]"
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
				fe.Name, displayCategory(fe), flags, fe.DestPath)
		}
		tw.Flush()
	} else {
		fmt.Fprintln(w, "  (none)")
	}

	// Section 2: old files flagged
	var old []scanner.FileEntry
	for _, fe := range moved {
		if fe.IsOld {
			old = append(old, fe)
		}
	}
	fmt.Fprintf(w, "\n%s\n", sep)
	fmt.Fprintf(w, "OLD FILES FLAGGED (%d files older than 365 days, moved to Archive/)\n\n", len(old))
	if len(old) > 0 {
		for _, fe := range old {
			fmt.Fprintf(w, "  - %-50s  (last modified: %s)\n",
				fe.Name, fe.ModTime.Format("2006-01-02"))
		}
	} else {
		fmt.Fprintln(w, "  (none)")
	}

	// Section 3: skipped files
	fmt.Fprintf(w, "\n%s\n", sep)
	fmt.Fprintf(w, "SKIPPED FILES (%d files with unrecognised type, left in place)\n\n", len(skipped))
	if len(skipped) > 0 {
		for _, fe := range skipped {
			fmt.Fprintf(w, "  - %s\n", fe.Name)
		}
	} else {
		fmt.Fprintln(w, "  (none)")
	}

	fmt.Fprintf(w, "\n%s\n", sep)
	fmt.Fprintln(w, "End of report.")

	return nil
}

// FormatTable returns terminal-friendly lines for the preview table,
// using tabwriter for aligned columns.
func FormatTable(entries []scanner.FileEntry) []string {
	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "FILE\tCATEGORY\tFLAGS\tDESTINATION")
	fmt.Fprintln(tw, "────\t────────\t─────\t───────────")
	for _, fe := range entries {
		flags := ""
		if fe.IsOld {
			flags = "[OLD]"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			fe.Name, displayCategory(fe), flags, fe.DestPath)
	}
	tw.Flush()

	lines := strings.Split(strings.TrimRight(sb.String(), "\n"), "\n")
	return lines
}

func displayCategory(fe scanner.FileEntry) string {
	if fe.IsOld {
		return "Archive"
	}
	if fe.Category == "" {
		return "(unknown)"
	}
	return fe.Category
}
