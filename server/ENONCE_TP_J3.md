# TP J3 — API REST météo en Go (validation Postman)

**Module** : Golang — M2 Dev Manager Full Stack
**Durée** : Journée
**Format** : solo
**Rendu** : demain matin **9h00**, démo flash en début de J4
**Stack imposée** : `net/http` + `encoding/json` uniquement. **Pas de framework** (Gin, Echo, Chi).

> 🔧 **Validation Postman.** Une collection est fournie sur Teams (`EFREI_Golang_J3.postman_collection.json`). Au lieu de retaper les `curl` à la main, vous importez la collection dans Postman et vous cliquez sur chaque requête. Le statut HTTP s'affiche en couleur, vous voyez le body de la réponse, c'est plus rapide et plus visuel.

---

## Contexte

Vous avez livré hier votre module Go de parsing JSON/XML. Aujourd'hui, vous **exposez ces données** via une API REST. À 18h, vous avez 6 routes qui répondent correctement à Postman.

---

## Setup Postman (3 min, à faire en premier)

1. Si Postman n'est pas installé : `postman.com/downloads` ou via le store Windows. Compte gratuit suffit.
2. Téléchargez `EFREI_Golang_J3.postman_collection.json` depuis Teams.
3. Postman → **File → Import** → glissez le `.json`. La collection apparaît dans la sidebar.
4. Variable de collection : `baseUrl = http://localhost:8080`. Modifiable si vous changez de port.

Le dossier qui vous intéresse pour valider chaque jalon : **"Slides 26-27 — API complète (smoke tests)"**. 12 requêtes numérotées, qui matchent exactement les jalons ci-dessous.

---

## Arborescence cible

```
weather/
├── go.mod
├── weather_data.json    ← votre dataset (depuis le TP J2)
├── model.go             ← reprendre depuis votre TP J2
├── jsonsource.go        ← idem
├── store.go             ← NOUVEAU — Store en mémoire
├── handlers.go          ← NOUVEAU — les 6 handlers
└── main.go              ← bootstrap + mux + ListenAndServe
```

Vous repartez de votre repo d'hier. Pas besoin de recréer le module.

---

## Méthode d'attaque — IMPÉRATIF

**Un jalon à la fois. Validé dans Postman. Commit. Puis seulement le suivant.**


Pour CHAQUE jalon, le cycle est :

1. Écrire le code minimal du jalon
2. `go run .` (laisser tourner)
3. Postman → cliquer sur la requête de validation du jalon → **statut attendu = statut obtenu** ?
4. Si OK : `git commit -m "feat: jalon N"` et au suivant.
5. Si KO : lire le log du serveur dans le terminal + le body de la réponse Postman, corriger.

---

## Jalon 1 — Setup

**Objectif** : avoir un dossier prêt à coder.

- Dans votre repo, créez (ou réutilisez) un dossier `server/` ou un repo dédié.
- Le `go.mod` reste le même que le TP J2.
- Assurez-vous que `weather_data.json` est à la racine du module.

**Validation Postman** : aucune, c'est juste un check `go build ./...` qui passe sans erreur.

**Commit** : `chore: setup server skeleton`

---

## Jalon 2 — Hello world HTTP

**Objectif** : un serveur qui tourne et répond `ok`.

```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.ListenAndServe(":8080", mux)
}
```

**Validation Postman** : pas de requête dédiée dans la collection. Lancez `go run .` puis dans Postman, créez une requête `GET {{baseUrl}}/health` à la main, statut **200** attendu.

**Commit** : `feat: server listens on :8080, /health returns ok`

---

## Jalon 3 — Bootstrap : charger le dataset dans un Store

**Objectif** : au démarrage du serveur, vos 30 stations sont en mémoire.

Créez `store.go` :

```go
package main

type Store struct {
	stations map[string]Station
}

func NewStore() *Store { }

func (s *Store) Put(st Station)      { 
	
}

func (s *Store) Has(id string) bool  { 
	
}

func (s *Store) Get(id string) (Station, bool) {
	
}

func (s *Store) Delete(id string) bool {

}

func (s *Store) All() []Station {
	
}
```

Modifiez `main.go` pour charger au démarrage :

```go
stations, err := LoadFromJSON("weather_data.json")
if err != nil { log.Fatal(err) }
store := NewStore()
for _, s := range stations { store.Put(s) }
log.Printf("bootstrap : %d stations chargées", len(stations))
```

**Validation** : au lancement du serveur, le log du terminal affiche `bootstrap : 30 stations chargées`.

**Commit** : `feat: bootstrap loads 30 stations into Store`

---

## Jalon 4 — `GET /stations`

**Objectif** : exposer la liste complète.

Créez `handlers.go` avec une struct `App` qui porte le Store et les helpers de réponse, puis le handler `listStations`.

```go
type App struct { store *Store }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (a *App) listStations(w http.ResponseWriter, r *http.Request) {
}
```


**Validation Postman** :
- Dossier "Slides 26-27" → requête **"01 - GET /stations (200)"** → cliquer **Send** → statut **200** + JSON array avec 30 entrées.

**Commit** : `feat: GET /stations returns the full list`

---

## Jalon 5 — `GET /stations/{id}` avec 404

**Objectif** : récupérer une station précise, ou renvoyer un 404 propre.

```go
func (a *App) getStation(w http.ResponseWriter, r *http.Request) {

}
```

Route : `mux.HandleFunc("GET /stations/{id}", app.getStation)`

**Validation Postman** :
- **"02 - GET /stations/FR-BOR-001 (200)"** → statut **200** + objet station.
- **"04 - GET /stations/XX-NOPE (404)"** → statut **404** + body `{"error":"..."}`.

**Commit** : `feat: GET /stations/{id} returns 200 or 404`

---

## Jalon 6 — `POST /stations` (201 / 400 / 409)

**Objectif** : créer une nouvelle station.

```go
func (a *App) createStation(w http.ResponseWriter, r *http.Request) {
	
}
```

Route : `mux.HandleFunc("POST /stations", app.createStation)`

**Validation Postman** (à exécuter dans l'ordre) :
- **"05 - POST création DEMO-001 (201)"** → statut **201**.
- **"06 - POST doublon DEMO-001 (409)"** → statut **409**.
- **"07 - POST JSON cassé (400)"** → statut **400**.

**Commit** : `feat: POST /stations creates a station (201/400/409)`

---

## Jalon 7 — `PUT /stations/{id}` (200 / 201)

**Objectif** : remplacer une station existante ou la créer si l'id n'existe pas.

```go
func (a *App) updateStation(w http.ResponseWriter, r *http.Request) {

}
```

Route : `mux.HandleFunc("PUT /stations/{id}", app.updateStation)`

**Validation Postman** :
- **"08 - PUT /stations/DEMO-001 remplacement (200)"** → statut **200**.
- **"09 - PUT /stations/NEW-999 création (201)"** → statut **201**.

**Commit** : `feat: PUT /stations/{id} replaces or creates`

---

## Jalon 8 — DELETE + observations + erreurs JSON typées

Trois ajouts dans ce dernier jalon.

### 8a. DELETE /stations/{id}

```go
func (a *App) deleteStation(w http.ResponseWriter, r *http.Request) {

}
```

### 8b. GET /stations/{id}/observations

```go
func (a *App) listObservations(w http.ResponseWriter, r *http.Request) {
	
}
```

### 8c. Erreurs JSON normalisées avec code interne

```go
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
}
```

Et mettre à jour les appels (`"NOT_FOUND"`, `"BAD_JSON"`, `"ID_TAKEN"`, etc.).

**Validation Postman** :
- **"03 - GET /stations/FR-BOR-001/observations (200)"** → statut **200**, array d'observations.
- **"10 - DELETE /stations/DEMO-001 (204)"** → statut **204**, pas de body.
- **"11 - DELETE /stations/DEMO-001 (404)"** → statut **404**.

**Commit** : `feat: DELETE + observations + normalized errors`

---

## Validation finale — Runner Postman

Postman a un mode "Runner" qui exécute toutes les requêtes d'un dossier dans l'ordre. Avant de rendre votre TP :

1. Lancez le serveur : `go run .`
2. Postman → clic droit sur le dossier **"Slides 26-27 — API complète"** → **Run folder**
3. Postman exécute les 12 requêtes en cascade et affiche les 12 statuts HTTP.

**Si vous avez 12/12 OK, votre TP est validé.**

Si vous avez des KO, le Runner vous indique exactement quelle requête a échoué — vous savez quoi corriger.

---

## Livrable

Avant **demain matin 9h00**, vous poussez sur votre repo :

- Le module Go complet, compilable (`go build ./...` passe sans erreur)
- Un `README.md` à la racine avec :
  - Comment lancer le serveur (`go run .`)
  - Les 6 routes et leurs codes statut attendus
  - Une mention de la collection Postman utilisée
  - Au moins une capture d'écran du Runner Postman montrant 12/12 OK (idéal)
- Vos commits doivent montrer l'avancée jalon par jalon (8 commits minimum)

Le rendu se fait par push sur votre repo Git

---

## Démo flash de demain (J4)

À 9h10 demain, **chaque étudiant a 30 secondes** pour :
1. Lancer son serveur (`go run .`)
2. Dans Postman, cliquer sur une requête de la collection et montrer le statut HTTP en vert.
3. Cliquer sur une autre (par exemple un 404 ou un 409) pour montrer une gestion d'erreur.


---

## Barème indicatif (sur 20)

| Critère | Points |
|---|---|
| Jalons 1-2 (setup + hello) | 1 |
| Jalon 3 (bootstrap) | 2 |
| Jalon 4 (GET /stations) | 2 |
| Jalon 5 (GET /{id} + 404) | 3 |
| Jalon 6 (POST + 201/400/409) | 4 |
| Jalon 7 (PUT + 200/201) | 3 |
| Jalon 8 (DELETE + observations + erreurs JSON typées) | 3 |
| Qualité de code, structure, commits propres | 2 |

Bonus +1 si Runner Postman complet capturé dans le README, ou si les erreurs ont un code interne (`NOT_FOUND`, etc.).

---

## Pièges fréquents (lisez-les **avant** de coder)

1. **`r.PathValue` ne marche pas** — vérifiez `go version`, il faut **1.22 ou plus**. Si 1.21, mettez à jour.
2. **`WriteHeader` ignoré silencieusement** — vous l'appelez après `Write` ? Inversez. Toujours l'ordre : `Header().Set` → `WriteHeader` → `Write/Encode`.
3. **`Decode` rejette les champs** — vous avez `DisallowUnknownFields`. Le payload contient un champ pas dans votre struct ? Adaptez.
4. **`time.Parse` échoue** — format de référence de Go : `2006-01-02T15:04:05Z07:00`. Pas `YYYY-MM-DD`.
5. **Le serveur ne redémarre pas après une modif** — Ctrl+C, puis `go run .`. Go ne fait pas de hot reload natif.
6. **Le port 8080 est déjà pris** — un ancien serveur tourne encore. `lsof -i :8080` (bash) ou `Get-NetTCPConnection -LocalPort 8080` (PowerShell) pour trouver le PID, puis `kill`.

---

## IAG

Pas encouragée. Vous avez tout ce qu'il faut dans :
- Le projet démo qu'on a construit ensemble en cours (10 dossiers `slide_NN_*/`)
- La collection Postman fournie
- La doc officielle : `pkg.go.dev/net/http`, `pkg.go.dev/encoding/json`
- Moi (premier appel gratuit)

Toute utilisation d'IAG non autorisée sur le TP vaudra une pénalité et des questions très salées en démo flash demain.
