package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/illenko/growth.md/internal/storage"
	"github.com/spf13/cobra"
)

var (
	cfgFile      string
	repoPath     string
	outputFormat string
	verbose      bool
)

var (
	config        *storage.Config
	skillRepo     *storage.SkillRepository
	goalRepo      *storage.GoalRepository
	pathRepo      *storage.PathRepository
	phaseRepo     *storage.PhaseRepository
	resourceRepo  *storage.ResourceRepository
	milestoneRepo *storage.MilestoneRepository
	progressRepo  *storage.ProgressLogRepository
)

var rootCmd = &cobra.Command{
	Use:   "growth",
	Short: "Git-native career development manager",
	Long: `growth.md - Track your skills, goals, and learning paths in plain Markdown files.

All your career development data is stored as human-readable Markdown files with
YAML frontmatter, versioned with Git for full history and portability.`,
	Version: "0.1.0-alpha",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeApp()
	},
	SilenceUsage: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: .growth/config.yml)")
	rootCmd.PersistentFlags().StringVar(&repoPath, "repo", "", "growth repository path (default: current directory)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "output format: table, json, yaml")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func initializeApp() error {
	if repoPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		repoPath = cwd
	}

	if cfgFile == "" {
		cfgFile = filepath.Join(repoPath, ".growth", "config.yml")
	}

	if _, err := os.Stat(cfgFile); err == nil {
		loadedConfig, err := storage.LoadConfig(cfgFile)
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "Warning: Failed to load config: %v\n", err)
			}
			config = storage.DefaultConfig()
		} else {
			config = loadedConfig
		}
	} else {
		config = storage.DefaultConfig()
	}

	if outputFormat != "" {
		config.Display.OutputFormat = outputFormat
	}

	if err := initializeRepositories(); err != nil {
		return err
	}

	return nil
}

func initializeRepositories() error {
	skillsPath := filepath.Join(repoPath, "skills")
	goalsPath := filepath.Join(repoPath, "goals")
	pathsPath := filepath.Join(repoPath, "paths")
	phasesPath := filepath.Join(repoPath, "phases")
	resourcesPath := filepath.Join(repoPath, "resources")
	milestonesPath := filepath.Join(repoPath, "milestones")
	progressPath := filepath.Join(repoPath, "progress")

	var err error

	skillRepo, err = storage.NewSkillRepository(skillsPath)
	if err != nil {
		return fmt.Errorf("failed to initialize skill repository: %w", err)
	}

	goalRepo, err = storage.NewGoalRepository(goalsPath)
	if err != nil {
		return fmt.Errorf("failed to initialize goal repository: %w", err)
	}

	pathRepo, err = storage.NewPathRepository(pathsPath)
	if err != nil {
		return fmt.Errorf("failed to initialize path repository: %w", err)
	}

	phaseRepo, err = storage.NewPhaseRepository(phasesPath)
	if err != nil {
		return fmt.Errorf("failed to initialize phase repository: %w", err)
	}

	resourceRepo, err = storage.NewResourceRepository(resourcesPath)
	if err != nil {
		return fmt.Errorf("failed to initialize resource repository: %w", err)
	}

	milestoneRepo, err = storage.NewMilestoneRepository(milestonesPath)
	if err != nil {
		return fmt.Errorf("failed to initialize milestone repository: %w", err)
	}

	progressRepo, err = storage.NewProgressLogRepository(progressPath)
	if err != nil {
		return fmt.Errorf("failed to initialize progress repository: %w", err)
	}

	return nil
}
