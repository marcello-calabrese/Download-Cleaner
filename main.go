package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marce/download-cleaner/mover"
	"github.com/marce/download-cleaner/reporter"
	"github.com/marce/download-cleaner/scanner"
)

func main() {
	flagPath := flag.String("path", "", "Path to Downloads folder (default: ~/Downloads)")
	flagAge := flag.Int("age", 365, "Age threshold in days for 'old' files")
	flagDryRun := flag.Bool("dry-run", false, "Preview changes without moving any files")
	flag.Parse()

	// Resolve downloads directory
	downloadsDir := *flagPath
	if downloadsDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot determine home directory: %v\n", err)
			os.Exit(1)
		}
		downloadsDir = filepath.Join(home, "Downloads")
	}

	info, err := os.Stat(downloadsDir)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "error: %q is not a valid directory\n", downloadsDir)
		os.Exit(1)
	}

	fmt.Printf("Scanning: %s\n", downloadsDir)

	// Scan
	entries, err := scanner.Scan(downloadsDir, *flagAge)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error during scan: %v\n", err)
		os.Exit(1)
	}

	var actionable, skipped []scanner.FileEntry
	for _, e := range entries {
		if e.IsSkipped {
			skipped = append(skipped, e)
		} else {
			actionable = append(actionable, e)
		}
	}

	// Print preview table
	fmt.Println("\n=== Download Cleaner — Preview ===")

	if len(actionable) == 0 {
		fmt.Println("Nothing to move. Downloads folder is already organised.")
		if len(skipped) > 0 {
			fmt.Printf("%d file(s) with unrecognised type were left in place.\n", len(skipped))
		}
		os.Exit(0)
	}

	lines := reporter.FormatTable(actionable)
	for _, l := range lines {
		fmt.Println(l)
	}

	oldCount := 0
	for _, e := range actionable {
		if e.IsOld {
			oldCount++
		}
	}

	fmt.Printf("\n%d file(s) to move", len(actionable))
	if len(skipped) > 0 {
		fmt.Printf(", %d file(s) skipped (unrecognised type)", len(skipped))
	}
	fmt.Println(".")

	if oldCount > 0 {
		fmt.Printf("[!] %d file(s) older than %d days — will be moved to Archive/\n",
			oldCount, *flagAge)
	}

	reportPath := filepath.Join(downloadsDir, "cleaner-report.log")

	// Dry-run: write preview report and exit
	if *flagDryRun {
		fmt.Println("\n[dry-run] No files were moved.")
		if err := reporter.Write(reportPath, actionable, skipped, true); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not write report: %v\n", err)
		} else {
			fmt.Printf("Dry-run report written to: %s\n", reportPath)
		}
		os.Exit(0)
	}

	// Interactive confirmation
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nProceed? [y/N] ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\nNo input received. Aborting.")
			os.Exit(0)
		}
		input = strings.TrimSpace(strings.ToLower(input))
		switch input {
		case "y", "yes":
			fmt.Println("Great! Starting to organise your Downloads folder...")
			goto doMove
		case "n", "no", "":
			fmt.Println("No problem! Your files have been left untouched.")
			os.Exit(0)
		default:
			fmt.Printf("Please enter 'y' to proceed or 'n' to abort (got: %q).\n", input)
		}
	}

doMove:
	fmt.Println("\nMoving files...")
	moved, moveErr := mover.Move(actionable)
	fmt.Printf("All done! %d/%d file(s) successfully moved.\n", len(moved), len(actionable))

	if err := reporter.Write(reportPath, moved, skipped, false); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not write report: %v\n", err)
	} else {
		fmt.Printf("Report written to: %s\n", reportPath)
	}

	if moveErr != nil {
		os.Exit(1)
	}
}
