# TP3 — GoLog Analyzer (loganalyzer)

Description  
GoLog Analyzer est un outil CLI écrit en Go pour analyser plusieurs fichiers de logs listés dans un fichier de configuration JSON, exécuter les analyses en parallèle, gérer proprement les erreurs et exporter un rapport JSON. Ce README explique l'utilisation (avec commandes `go run main.go` en premier — utile si `./loganalyzer` n'est pas exécutable sur Windows) et l'équivalent binaire (`./loganalyzer` / `.\loganalyzer.exe`).

Résumé des fonctionnalités

- Charger une liste de logs depuis un fichier JSON (`--config` / `-c`).
- Lancer l'analyse de chaque fichier en parallèle (une goroutine par cible).
- Détecter et reporter les erreurs (fichier introuvable, parsing corrompu, etc.).
- Archiver en mémoire chaque résultat (status, message, details, timestamp).
- Exporter le rapport final au format JSON (`--output` / `-o`).
- Options bonus : création automatique des dossiers d'export, horodatage du nom du fichier, filtrage par status, commande `add-log`.

Structure du projet

- `cmd/` : commandes Cobra (`root.go`, `analyze.go`, ...).
- `internal/config` : lecture et validation du fichier de config JSON.
- `internal/analyzer` : logique d'analyse, erreurs personnalisées.
- `internal/reporter` : écriture/export JSON (création dossiers si nécessaire).
- `fixtures/` et `test_logs/` : exemples de logs pour tests.
- `main.go` / `go.mod` : point d'entrée et gestion des dépendances.

Installation (rapide)

- Ouvrir PowerShell et se placer dans le dossier du projet TP3 :

```powershell
cd "C:\Users\Dev_Note\Desktop\Dev-Student\M2\Golang\Cours 1\exemple\TP3"
```

- Installer les dépendances (si besoin) :

```powershell
go mod tidy
```

Commandes importantes — exemples concrets (double : `go run main.go` d'abord, puis l'équivalent binaire)

- Lancer l'analyse (exemple concret avec `config2.json` et sortie dans `out/reports/result.json`):

```powershell
# exécution sans build (recommandé si ./loganalyzer ne s'exécute pas)
go run main.go analyze -c config2.json -o out/reports/result.json

# équivalent si vous avez construit l'exécutable
# sous Windows
.\loganalyzer.exe analyze -c config2.json -o out/reports/result.json
# ou (posix)
./loganalyzer analyze -c config2.json -o out/reports/result.json
```

- Lancer l'analyse et laisser Cobra générer le nom daté du fichier (bonus horodatage) :

```powershell
go run main.go analyze -c config.json -o reports/report.json --dated
.\loganalyzer.exe analyze -c config.json -o reports/report.json --dated
```

- Afficher l'aide de la commande `analyze` :

```powershell
go run main.go analyze --help
.\loganalyzer.exe analyze --help
```

- Ajouter un log dans un fichier de config (sous-commande `add-log`) :

```powershell
go run main.go add-log --id=my-app --path=test_logs/app.log --type=custom-app --file=config.json
.\loganalyzer.exe add-log --id=my-app --path=test_logs/app.log --type=custom-app --file=config.json
```

- Construire un exécutable :

```powershell
go build -o loganalyzer.exe .
# puis
.\loganalyzer.exe analyze -c config.json -o report.json
```

Flags et options courants

- `--config, -c <path>` : chemin vers le fichier JSON de configuration (obligatoire pour `analyze`).
- `--output, -o <path>` : chemin de sortie pour le rapport JSON (optionnel ; sinon `report.json`).
- `--status <OK|FAILED>` : filtrer l'affichage / export selon le statut (bonus).
- `--dated` : si présent, préfixe le fichier de sortie avec la date (AAMMJJ).
- `--help` : aide par commande (Cobra fournit automatiquement `--help`).

Exemple de config.json (rappel)

```json
[
  {
    "id": "web-server-1",
    "path": "test_logs/access.log",
    "type": "nginx-access"
  },
  {
    "id": "app-backend-2",
    "path": "test_logs/errors.log",
    "type": "custom-app"
  }
]
```

- Remarque : les chemins sont relatifs au répertoire courant d'exécution. Pour éviter les erreurs, lancez la commande depuis le dossier `TP3` ou utilisez des chemins absolus.

Que fait chaque package (bref)

- internal/config : lit le JSON de configuration, normalise les chemins et valide la présence du fichier `path`.
- internal/analyzer : contient les analyseurs (ex: nginx-access, custom-app). Gère erreurs personnalisées (fichier introuvable, parse error) et exécute la simulation d’analyse.
- internal/reporter : sérialise le tableau de résultats et écrit le fichier JSON ; crée les dossiers si nécessaires (`os.MkdirAll`).
- cmd/ : implemente l'interface utilisateur CLI (Cobra), parsing des flags et orchestration générale.

Conseils / points d'attention

- Toujours exécuter depuis le répertoire racine `TP3` pour que les chemins relatifs (ex: `test_logs/...`) soient résolus correctement.
- Si une entrée du `config.json` pointe vers un chemin non existant, l'analyse enregistrera un résultat avec `status: FAILED` et `error_details`. C'est normal pour les cas de test.
- Pour debug : lancer `go run main.go analyze -c config.json -o report.json` et vérifier les DEBUG prints si vous avez ajouté des logs temporaires dans `cmd/analyze.go`.
- Pour détecter des problèmes de concurrence (si vous activez l'analyse parallèle), exécutez avec le détecteur de race :

```powershell
go run -race main.go analyze -c config.json -o report.json
```

Format du rapport (`report.json`)

- Chaque entrée contient :
  - `log_id` : id depuis la config
  - `file_path` : chemin analysé
  - `status` : `OK` ou `FAILED`
  - `message` : résumé lisible
  - `error_details` : texte d'erreur en cas d'échec

Installer / publier sur GitHub (rapide)

- Initialiser le repo dans TP3 (si nécessaire) :

```powershell
git init
git add .
git commit -m "TP3 - GoLog Analyzer"
# ajouter remote et push...
```

Prochaines améliorations suggérées

- Faire l’analyse réellement concurrente (goroutines + mutex ou channel pour le Storer).
- Persister l’historique en JSON et ajouter option de lecture.
- Tests unitaires pour les analyzeurs et la gestion d’erreurs personnalisées.
- Ajout d’une interface web minimale pour visualiser les rapports.

A propos de l'auteur  
TP3 réalisé dans le cadre de l'initiation à Go. Réalisé par Fairytale-Dev(Farid-Efrei) (étudiant — Alternant à l'Efrei).

🦋 Paix à tous 🦋
