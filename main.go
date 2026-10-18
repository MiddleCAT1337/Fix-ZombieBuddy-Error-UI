package main

import (
	"fmt"
	"os"

	"github.com/whitehorizon/zb-patch/internal/apply"
	"github.com/whitehorizon/zb-patch/internal/discover"
	"github.com/whitehorizon/zb-patch/internal/elevation"
	"github.com/whitehorizon/zb-patch/internal/plan"
	"github.com/whitehorizon/zb-patch/internal/steamcheck"
)

const toolName = "zb-patch"

func main() {
	os.Exit(run())
}

func run() int {
	if err := elevation.EnsureAdmin(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", toolName, err)
		return 1
	}

	paths, err := discover.Resolve(discover.Overrides{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", toolName, err)
		return 1
	}

	fmt.Printf("Steam:      %s\n", paths.Steam)
	fmt.Printf("PZ:         %s\n", paths.PZ)
	fmt.Printf("Workshop:   %s\n", paths.WorkshopZB)
	fmt.Println()

	if steamcheck.Running() {
		fmt.Println("Note: Restart Project Zomboid after patching to load the Java agent.")
		fmt.Println()
	}

	cs, err := plan.Build(paths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", toolName, err)
		return 1
	}

	for _, line := range cs.PreviewLines {
		fmt.Println("  -", line)
	}
	fmt.Println()

	if err := apply.Run(cs); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", toolName, err)
		return 1
	}

	fmt.Println("Done. Restart Project Zomboid from Steam.")
	return 0
}
