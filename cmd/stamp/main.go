package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/toto/stamp/internal/clipboard"
	"github.com/toto/stamp/internal/config"
	"github.com/toto/stamp/internal/counter"
	"github.com/toto/stamp/internal/generator"
	"github.com/toto/stamp/internal/obsidian"
	"github.com/toto/stamp/internal/sequential"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	cfg  *config.Config
	cntr *counter.Manager
	gen  *generator.Generator

	// Flags
	flagExt            bool
	flagCopy           bool
	flagQuiet          bool
	flagAnalogCheck    bool
	flagAnalogReset    bool
	flagAnalogCounter  bool
	flagProjectCheck   bool
	flagProjectCounter bool
	flagSeqPrefix      string
	flagSeqWidth       int
	flagSeqStart       int
	flagSeqCheck       bool
	flagSeqCounter     bool
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
  - project:  PXXXX format (shorthand for seq --prefix P --width 4)
  - seq:      Custom prefix + zero-padded number (workspace scan)

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
	rootCmd.AddCommand(seqCmd)
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
		if flagAnalogCheck {
			result, err := cntr.CheckAnalog(gen.GetCurrentDate())
			if err != nil {
				return err
			}
			return outputResult(result)
		}

		if flagAnalogReset {
			if err := cntr.ResetAnalog(gen.GetCurrentDate()); err != nil {
				return err
			}
			if !flagQuiet {
				fmt.Println("Counter reset for analog notes")
			}
			return nil
		}

		if flagAnalogCounter {
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
	Long:  "Equivalent to `stamp seq --prefix P --width 4`.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSeqCommand(seqCommandOptions{
			Spec:         sequential.Spec{Prefix: "P", Width: 4, Start: 1},
			CounterLabel: "project",
			Check:        flagProjectCheck,
			Counter:      flagProjectCounter,
			TitleArgs:    args,
		})
	},
}

var seqCmd = &cobra.Command{
	Use:     "seq [title]",
	Aliases: []string{"sequential"},
	Short:   "Generate sequential codes from the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSeqCommand(seqCommandOptions{
			Spec: sequential.Spec{
				Prefix: flagSeqPrefix,
				Width:  flagSeqWidth,
				Start:  flagSeqStart,
			},
			Check:     flagSeqCheck,
			Counter:   flagSeqCounter,
			TitleArgs: args,
		})
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
	// Add counter management flags to analog, project, and seq commands
	analogCmd.Flags().BoolVar(&flagAnalogCheck, "check", false, "Check next number without incrementing")
	analogCmd.Flags().BoolVar(&flagAnalogReset, "reset", false, "Reset counter")
	analogCmd.Flags().BoolVar(&flagAnalogCounter, "counter", false, "Show current counter value")

	projectCmd.Flags().BoolVar(&flagProjectCheck, "check", false, "Check next number without incrementing")
	projectCmd.Flags().BoolVar(&flagProjectCounter, "counter", false, "Show highest existing number")

	seqCmd.Flags().StringVar(&flagSeqPrefix, "prefix", "P", "Prefix for generated code (case-insensitive match)")
	seqCmd.Flags().IntVar(&flagSeqWidth, "width", 4, "Number of digits for zero padding")
	seqCmd.Flags().IntVar(&flagSeqStart, "start", 1, "Starting number when no entries are found")
	seqCmd.Flags().BoolVar(&flagSeqCheck, "check", false, "Check next number without creating files")
	seqCmd.Flags().BoolVar(&flagSeqCounter, "counter", false, "Show highest existing number for the prefix")
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

type seqCommandOptions struct {
	Spec         sequential.Spec
	CounterLabel string
	Check        bool
	Counter      bool
	TitleArgs    []string
}

func runSeqCommand(opts seqCommandOptions) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if opts.Counter {
		highest, err := sequential.Highest(wd, opts.Spec)
		if err != nil {
			return err
		}

		label := opts.CounterLabel
		if label == "" {
			label = fmt.Sprintf("counter for prefix %s", strings.ToUpper(normalizePrefix(opts.Spec)))
		} else {
			label = fmt.Sprintf("%s counter", label)
		}

		if highest == 0 {
			fmt.Printf("Current %s: none\n", label)
		} else {
			fmt.Printf("Current %s: %d\n", label, highest)
		}
		return nil
	}

	code, _, err := sequential.Next(wd, opts.Spec)
	if err != nil {
		return err
	}

	title := strings.Join(opts.TitleArgs, " ")
	if opts.Check {
		title = ""
	}
	if title != "" {
		code += " " + title
	}

	return outputResult(code)
}

func normalizePrefix(spec sequential.Spec) string {
	if spec.Prefix == "" {
		return "P"
	}
	return spec.Prefix
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

	if wd, err := os.Getwd(); err == nil {
		if detectResult, detectErr := obsidian.Detect(wd); detectErr != nil {
			fmt.Fprintf(os.Stderr, "Obsidian detection warning: %v\n", detectErr)
			if detectResult != nil && detectResult.InVault {
				gen.ApplyLayouts(generator.LayoutOverrides{
					Default: detectResult.Layouts.Default,
					Daily:   detectResult.Layouts.Daily,
				})
			}
		} else if detectResult.InVault {
			gen.ApplyLayouts(generator.LayoutOverrides{
				Default: detectResult.Layouts.Default,
				Daily:   detectResult.Layouts.Daily,
			})
		}
	}

	// Apply default extension flag from config
	if cfg.AlwaysExtension && !rootCmd.PersistentFlags().Changed("ext") {
		flagExt = true
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
