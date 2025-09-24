package cmd

import (
	"errors"
	"fmt"
	"loganalyzer/internal/analyzer"
	"loganalyzer/internal/config"
	"loganalyzer/internal/reporter"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	cfgPath      string
	output       string
	filterStatus string
)
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyse des logs listés dans un fichier de config JSON",
	Long:  `Analyse des logs listés dans un fichier de config JSON et exporte les résultats dans un fichier JSON.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1) Charger la config
		targets, err := config.LoadConfig(cfgPath)
		if err != nil {
			return fmt.Errorf("charger config: %w", err)
		}
		if len(targets) == 0 {
			fmt.Println("Aucun log à analyser (config vide)")
			return nil
		}
		// 2) lancer l'analyse concurente
		resultsCh := analyzer.Analyze(targets)

		// 3) Collecter + afficher un résumé
		var entries []analyzer.ReportEntry
		for e := range resultsCh {
			if e.Err != nil {
				switch {
				case errors.Is(e.Err, analyzer.ErrFileAccess):
					fmt.Printf("❌ [%s] %s | %s | %s\n", e.LogID, e.FilePath, e.Status, e.Message)
				case errors.Is(e.Err, analyzer.ErrParse):
					fmt.Printf("❌ [%s] %s | %s | %s\n", e.LogID, e.FilePath, e.Status, e.Message)
				default:
					fmt.Printf("❌ [%s] %s | %s | %s\n", e.LogID, e.FilePath, e.Status, e.Message)
				}
				fmt.Println("    →", e.ErrorDetails)

			} else {
				fmt.Printf("✅ [%s] %s | %s\n", e.LogID, e.FilePath, e.Status)
			}
			entries = append(entries, e)

		}

		filtered := entries
		if filterStatus == analyzer.StatusOK || filterStatus == analyzer.StatusFailed {
			var tmp []analyzer.ReportEntry
			for _, e := range entries {
				if e.Status == filterStatus {
					tmp = append(tmp, e)
				}
			}
			filtered = tmp
		}

		// Exporter JSON si demandé
		if output != "" {
			// Bonus 2 : Ajouter un timestamp au nom de fichier si c'est un .json
			if strings.HasSuffix(strings.ToLower(output), ".json") {
				stamp := time.Now().Format("060102") // YYMMDD
				dir, base := filepath.Split(output)
				output = filepath.Join(dir, stamp+"_"+base)
			}
			if err := reporter.SaveJSON(output, filtered); err != nil {
				return err
			}

			fmt.Printf("✅ Résultats exportés dans %s\n", output)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVarP(&cfgPath, "config", "c", "config.json", "Chemin vers le fichier de configuration JSON")
	analyzeCmd.Flags().StringVarP(&output, "output", "o", "", "Chemin vers le fichier de sortie JSON (optionnel)")
	// Bonus 3 : Filtrer les résultats par status
	analyzeCmd.Flags().StringVarP(&filterStatus, "status", "s", "", "Filtrer les résultats par status (OK, FAILED)")
}
