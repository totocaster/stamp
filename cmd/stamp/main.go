package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/toto/stamp/internal/config"
	"github.com/toto/stamp/internal/counter"
	"github.com/toto/stamp/internal/generator"
	"github.com/toto/stamp/internal/clipboard"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	cfg       *config.Config
	cntr      *counter.Manager
	gen       *generator.Generator

	// Flags
	flagExt    bool
	flagCopy   bool
	flagQuiet  bool
	flagCheck  bool
	flagReset  bool
	flagSet    int
	flagCounter bool
)

var rootCmd = &cobra.Command{
	Use:   "stamp [type]",
	Short: "Generate note filenames based on date/time",
	Long: `Stamp is a CLI tool for generating note filenames following Toto's note naming conventions.

Available note types:
  - daily:    YYYY-MM-DD format
  - fleeting: YYYY-MM-DD-FHHMMSS format
  - voice:    YYYY-MM-DD-VTHHMMSS format
  - analog:   YYYY-MM-DD-AN format (sequential per day)
  - monthly:  YYYY-MM format
  - yearly:   YYYY format
  - project:  PXXXX format (auto-incrementing)

Default (no type): YYYY-MM-DD-HHMM format`,
	RunE: runDefault,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&flagExt, "ext", false, "Add .md extension to output")
	rootCmd.PersistentFlags().BoolVar(&flagCopy, "copy", false, "Copy to clipboard (macOS only)")
	rootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "Quiet mode (no additional output)")

	// Add subcommands
	rootCmd.AddCommand(dailyCmd)
	rootCmd.AddCommand(fleetingCmd)
	rootCmd.AddCommand(voiceCmd)
	rootCmd.AddCommand(analogCmd)
	rootCmd.AddCommand(monthlyCmd)
	rootCmd.AddCommand(yearlyCmd)
	rootCmd.AddCommand(projectCmd)
	rootCmd.AddCommand(versionCmd)
}

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Generate daily note filename (YYYY-MM-DD)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return outputResult(gen.Daily())
	},
}

var fleetingCmd = &cobra.Command{
	Use:   "fleeting",
	Short: "Generate fleeting note filename (YYYY-MM-DD-FHHMMSS)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return outputResult(gen.Fleeting())
	},
}

var voiceCmd = &cobra.Command{
	Use:   "voice",
	Short: "Generate voice transcript filename (YYYY-MM-DD-VTHHMMSS)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return outputResult(gen.Voice())
	},
}

var analogCmd = &cobra.Command{
	Use:   "analog",
	Short: "Generate analog/slipbox note filename (YYYY-MM-DD-AN)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagCheck {
			result, err := cntr.CheckAnalog(gen.GetCurrentDate())
			if err != nil {
				return err
			}
			return outputResult(result)
		}

		if flagReset {
			if err := cntr.ResetAnalog(gen.GetCurrentDate()); err != nil {
				return err
			}
			if !flagQuiet {
				fmt.Println("Counter reset for analog notes")
			}
			return nil
		}

		if flagCounter {
			count, err := cntr.GetAnalogCounter(gen.GetCurrentDate())
			if err != nil {
				return err
			}
			fmt.Printf("Current analog counter for %s: %d\n", gen.GetCurrentDate(), count)
			return nil
		}

		result, err := cntr.NextAnalog(gen.GetCurrentDate())
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var monthlyCmd = &cobra.Command{
	Use:   "monthly",
	Short: "Generate monthly review filename (YYYY-MM)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return outputResult(gen.Monthly())
	},
}

var yearlyCmd = &cobra.Command{
	Use:   "yearly",
	Short: "Generate yearly review filename (YYYY)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return outputResult(gen.Yearly())
	},
}

var projectCmd = &cobra.Command{
	Use:   "project [title]",
	Short: "Generate project number (PXXXX)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagCheck {
			result, err := cntr.CheckProject()
			if err != nil {
				return err
			}
			return outputResult(result)
		}

		if flagReset {
			if err := cntr.ResetProject(); err != nil {
				return err
			}
			if !flagQuiet {
				fmt.Println("Counter reset for project")
			}
			return nil
		}

		if flagSet > 0 {
			if err := cntr.SetProject(flagSet); err != nil {
				return err
			}
			if !flagQuiet {
				fmt.Printf("Project counter set to %d\n", flagSet)
			}
			return nil
		}

		if flagCounter {
			count, err := cntr.GetProjectCounter()
			if err != nil {
				return err
			}
			fmt.Printf("Current project counter: %d\n", count)
			return nil
		}

		var title string
		if len(args) > 0 {
			title = args[0]
		}

		result, err := cntr.NextProject(title)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("stamp version %s\ncommit: %s\nbuilt: %s\n", version, commit, date)
	},
}

func init() {
	// Add counter management flags to analog and project commands
	analogCmd.Flags().BoolVar(&flagCheck, "check", false, "Check next number without incrementing")
	analogCmd.Flags().BoolVar(&flagReset, "reset", false, "Reset counter")
	analogCmd.Flags().BoolVar(&flagCounter, "counter", false, "Show current counter value")

	projectCmd.Flags().BoolVar(&flagCheck, "check", false, "Check next number without incrementing")
	projectCmd.Flags().BoolVar(&flagReset, "reset", false, "Reset counter")
	projectCmd.Flags().IntVar(&flagSet, "set", 0, "Set counter to specific value")
	projectCmd.Flags().BoolVar(&flagCounter, "counter", false, "Show current counter value")
}

func runDefault(cmd *cobra.Command, args []string) error {
	// If an argument is provided, treat it as a subcommand
	if len(args) > 0 {
		// Try to find and execute the subcommand
		for _, subcmd := range cmd.Commands() {
			if subcmd.Name() == args[0] {
				return subcmd.RunE(subcmd, args[1:])
			}
		}
		return fmt.Errorf("unknown note type: %s", args[0])
	}

	// Default behavior: output timestamp
	return outputResult(gen.Default())
}

func outputResult(result string) error {
	if flagExt {
		result += ".md"
	}

	if flagCopy {
		if err := clipboard.Copy(result); err != nil {
			// Fall back to stdout if clipboard fails
			fmt.Println(result)
			return fmt.Errorf("clipboard error: %w", err)
		}
		if !flagQuiet {
			fmt.Println(result)
			fmt.Println("Copied to clipboard!")
		}
	} else {
		fmt.Println(result)
	}

	return nil
}

func main() {
	var err error

	// Load configuration
	cfg, err = config.Load()
	if err != nil {
		// Use defaults if config loading fails
		cfg = config.Default()
	}

	// Initialize counter manager
	cntr, err = counter.New(cfg.CounterFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing counter: %v\n", err)
		os.Exit(1)
	}

	// Initialize generator
	gen, err = generator.New(cfg.Timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing generator: %v\n", err)
		os.Exit(1)
	}

	// Apply default extension flag from config
	if cfg.AlwaysExtension && !rootCmd.PersistentFlags().Changed("ext") {
		flagExt = true
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}