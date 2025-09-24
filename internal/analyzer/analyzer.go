package analyzer

import (
	"errors"

	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	StatusOK     = "OK"
	StatusFailed = "FAILED"
)

type LogTarget struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
}

type ReportEntry struct {
	LogID        string    `json:"log_id"`
	FilePath     string    `json:"file_path"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	ErrorDetails string    `json:"error_details"`
	Err          error     `json:"-"`
	AnalyzedAt   time.Time `json:"-"`
}

func Analyze(target []LogTarget) <-chan ReportEntry {
	results := make(chan ReportEntry)
	var wg sync.WaitGroup

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	wg.Add(len(target))
	for _, t := range target {
		t := t //capture
		go func() {
			defer wg.Done()

			//Simulation d'analyse : 500ms max
			time.Sleep(time.Duration(r.Intn(500)) * time.Millisecond)

			// Vérifions que le fichier est bien lisible et existe
			stat, err := os.Stat(t.Path)
			if err != nil {
				ferr := &FileAccessError{Path: t.Path, Err: err}
				results <- ReportEntry{
					LogID:        t.ID,
					FilePath:     t.Path,
					Status:       StatusFailed,
					Message:      "Fichier introuvable ou illisible",
					ErrorDetails: ferr.Error(),
					Err:          ferr,
					AnalyzedAt:   time.Now(),
				}
				return
			}

			if stat.IsDir() {
				ferr := &FileAccessError{Path: t.Path, Err: errors.New("c'est un répertoire")}
				results <- ReportEntry{
					LogID:        t.ID,
					FilePath:     t.Path,
					Status:       StatusFailed,
					Message:      "Chemin invalide (répertoire)",
					ErrorDetails: ferr.Error(),
					Err:          ferr,
					AnalyzedAt:   time.Now(),
				}
				return
			}

			// Simulons une erreur de parsing (ParseError)
			if strings.Contains(strings.ToLower(filepath.Base(t.Path)), "corrupt") {
				perr := &ParseError{Path: t.Path, Why: "fichier corrompu détecté"}
				results <- ReportEntry{
					LogID:        t.ID,
					FilePath:     t.Path,
					Status:       StatusFailed,
					Message:      "Erreur de parsing",
					ErrorDetails: perr.Error(),
					Err:          perr,
					AnalyzedAt:   time.Now(),
				}
				return
			}
			// Si erreur de lecture
			b, rerr := os.ReadFile(t.Path)
			if rerr != nil {
				ferr := &FileAccessError{Path: t.Path, Err: rerr}
				results <- ReportEntry{
					LogID:        t.ID,
					FilePath:     t.Path,
					Status:       StatusFailed,
					Message:      "Erreur de lecture du fichier",
					ErrorDetails: ferr.Error(),
					Err:          ferr,
					AnalyzedAt:   time.Now(),
				}
				return
			}
			// Tout est OK
			msg := "Analyse terminée avec succès ! Youhou !"
			if len(b) == 0 {
				msg = "Fichier vide (Aucune donnée) - OK "
			}
			results <- ReportEntry{
				LogID:        t.ID,
				FilePath:     t.Path,
				Status:       StatusOK,
				Message:      msg,
				ErrorDetails: "",
				Err:          nil,
				AnalyzedAt:   time.Now(),
			}
		}()
	}
	// Fermeture du channel à la fin
	go func() {
		wg.Wait()
		close(results)
	}()
	return results

}
