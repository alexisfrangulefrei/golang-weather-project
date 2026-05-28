# Atelier — Parsing JSON / XML de données météo en Go

**Module** : Golang — M2 Dev Manager Full Stack **Durée** : 2h30 environ (atelier du milieu de journée 2\) **Format** : solo. **Livrable** : un module Go compilable, à pousser sur le repo de la promo avant 17h10h. **Suite directe** : le TP du soir construit une API REST CRUD par-dessus ce module

---

## 1\. Contexte

Vous reprenez le code d'une équipe data d'un institut météorologique européen. Vous récupérez deux export : un fichier **JSON** produit par leur nouvelle API, et un fichier **XML** hérité de leur ancien SI qui tourne encore. Les deux fichiers décrivent les mêmes 30 stations et les mêmes 300 observations, **mais avec des schémas différents** (les équipes ne se sont jamais mises d'accord sur la modélisation).

Votre mission : produire un module Go capable de lire les deux formats et de fournir une vue **unifiée** des données à un appelant. C'est exactement le genre de travail que vous ferez en production le jour où vous devrez migrer un format vers un autre, ou consolider des sources hétérogènes (spoiler, ça arrive tout le temps).

Les fichiers fournis :

- `weather_data.json` — 30 stations, 300 observations, structure imbriquée par clés  
- `weather_data.xml` — mêmes données, structure attribut-heavy, mesures typées en éléments

---

## 2\. Objectifs pédagogiques

À la fin de l'atelier vous serez capables de :

- Décoder un fichier JSON volumineux avec `encoding/json` et des struct tags.  
- Décoder un fichier XML avec `encoding/xml`, en distinguant attributs et éléments enfants.  
- Concevoir un modèle interne unifié indépendant de la source.  
- Mapper deux schémas externes différents vers le même modèle interne.  
- Gérer proprement les champs optionnels (pointeurs, `omitempty`).  
- Exécuter quelques requêtes simples sur la collection en mémoire (filtrage, agrégation).

---

## 3\. Pré-requis (ce que vous avez vu hier)

- Slices, maps, structs (y compris struct embarqués)  
- Lecture de fichier, `os.ReadFile`  
- `encoding/json` en lecture, struct tags `json:"..."`  
- Pointeurs sur valeur scalaire pour les champs optionnels

Si l'un de ces points est flou, c'est le moment de poser la question maintenant, pas à 17h45.

---

## 4\. Mise en route

mkdir \-p weather && cd weather

go mod init github.com/efrei/weather

\# placez weather\_data.json et weather\_data.xml à la racine du module

Arborescence cible :

weather/

├── go.mod

├── weather\_data.json

├── weather\_data.xml

├── model.go        \# le modèle interne unifié

├── jsonsource.go   \# parser JSON

├── xmlsource.go    \# parser XML

├── query.go        \# quelques fonctions de requête sur la collection

└── main.go \# programme principal de démo

---

## 5\. Partie A — Exploration des deux fichiers

Ouvrez les deux fichiers dans Goland. Comparez **à l'œil** la modélisation des éléments suivants, puis remplissez le tableau dans votre `README.md` avant de coder quoi que ce soit (c'est l'étape la plus importante du TP) :

| Donnée | Comment c'est représenté en JSON ? | Comment c'est représenté en XML ? |
| :---- | :---- | :---- |
| Pays |  |  |
| Coordonnées |  |  |
| Altitude |  |  |
| Modèle de capteur |  |  |
| Température |  |  |
| Conditions ciel |  |  |
| Vent |  |  |
| Notes (optionnelles) |  |  |

Tant que ce tableau n'est pas rempli, ne commencez pas à coder.

---

## 6\. Partie B — Modèle interne unifié

Dans `model.go`, définissez les types Go qui représentent **votre vision** des données, indépendamment du format source. Quelques règles de conception :

- Les noms de champs sont en anglais, lisibles, sans suffixe d'unité (préférez `WindSpeed` avec un commentaire `// km/h` plutôt que `WindSpeedKmh`).  
- Le pays est stocké en **code ISO 2 lettres** (`FR`, `ES`, …) dans le modèle interne, peu importe ce que dit la source.  
- Les champs optionnels (les notes) sont des pointeurs (`*string`), pas des chaînes vides.  
- Le timestamp est un `time.Time`, pas un `string`.  
- Aucun struct tag JSON/XML dans `model.go`. Le modèle interne est neutre.

---

## 7\. Partie C — Parser JSON 

Dans `jsonsource.go`, écrivez :

func LoadFromJSON(path string) (\[\]Station, error)

Quelques points d'attention que vous allez devoir traiter :

- Définissez des **types intermédiaires** privés (`jsonStation`, `jsonObservation`, …) avec les tags `json:"..."` qui collent à la structure du fichier. Ne tentez pas de tagger directement votre modèle interne — vous le regretterez en partie D quand il faudra parser le XML aussi.  
- Le champ `country` du JSON est un nom complet en français (`"France"`). Votre modèle interne attend un code ISO 2 lettres. Construisez une petite `map[string]string` de conversion (une dizaine d'entrées suffit pour les pays présents).  
- Le champ `notes` est `null` dans la majorité des cas. Avec un type cible `*string` et `encoding/json`, le décodage gère ça pour vous (`nil` si null). Vérifiez que c'est bien le cas.  
- Le champ `installed_on` est un `string` au format `YYYY-MM-DD`. Vous devrez le parser avec `time.Parse("2006-01-02", s)`. Mêmes remarques pour les `timestamp` au format ISO 8601 (`time.RFC3339`).  
- Écrivez une **fonction de conversion** `(jsonStation) -> Station` qui isole le mapping. C'est elle qui fait le vrai travail. Le décodage JSON n'est qu'un transport.

Test rapide à la fin de la partie : `len(stations) == 30` et `len(stations[0].Observations) == 10`. Si ce n'est pas le cas, c'est un bug.

---

## 8\. Partie D — Parser XML 

Dans `xmlsource.go`, écrivez :

func LoadFromXML(path string) (\[\]Station, error)

Le XML est volontairement plus retors que le JSON. Quelques pièges à anticiper :

- Les attributs se taggent avec `xml:"nom,attr"` et pas `xml:"nom"`. Confondre les deux est l'erreur n°1.  
- Les mesures sont stockées en `<measure type="temperature" unit="C">22.5</measure>`. Le contenu textuel d'un élément se récupère avec un champ `Value string` taggé `xml:",chardata"`. Vous obtiendrez donc une slice `[]xmlMeasure` qu'il faudra **dispatcher** dans votre `Observation` en fonction de l'attribut `type`. C'est un mini-switch, rien de plus.  
- Les valeurs des mesures sont des strings. À convertir en `float64` ou `int` avec `strconv`.  
- L'élément `<note>` n'existe que si la note est présente. Le décodeur laisse simplement le champ à sa valeur zéro si l'élément est absent. Si vous le déclarez `*string` côté intermédiaire, ça marche tout seul.  
- Le code pays est déjà ISO ici (`country="FR"`) — pas de conversion à faire dans ce sens.  
- L'altitude est en attribut sur `<coordinates>`, pas dans un champ dédié. À remapper proprement.

Même principe qu'en partie C : un type intermédiaire `xmlStation`, et une fonction `(xmlStation) -> Station` qui fait le mapping.

Test rapide à la fin de la partie : les deux loaders renvoient des slices `[]Station` qui contiennent **les mêmes données** (mêmes IDs, mêmes valeurs numériques aux arrondis près). On va le vérifier en partie E.

---

## 9\. Partie E — Requêtes sur la collection 

Dans `query.go`, implémentez les fonctions suivantes. Toutes opèrent sur `[]Station` en mémoire, pas de prouesse algorithmique attendue.

// Renvoie les stations d'un pays donné (code ISO).

func FilterByCountry(stations \[\]Station, iso string) \[\]Station

// Renvoie la température moyenne (°C) sur l'ensemble des observations d'une station.

func AvgTemperature(s Station) float64

// Renvoie la station avec la rafale de vent la plus forte sur l'ensemble du dataset,

// avec la valeur correspondante.

func MaxWindGust(stations \[\]Station) (Station, float64)

// Renvoie un map\[code\_pays\] \-\> nombre de stations.

func CountByCountry(stations \[\]Station) map\[string\]int

Dans `cmd/main.go`, écrivez un petit programme qui :

1. Charge le JSON et le XML.  
2. Compare que les deux datasets ont le même nombre de stations et d'observations.  
3. Affiche la station la plus ventée et sa rafale.  
4. Affiche la moyenne de température de la station `FR-BOR-001` (Bordeaux Mérignac).  
5. Affiche le compte de stations par pays.

go run ./cmd/explore

Sortie attendue (format libre) :

JSON : 30 stations, 300 observations

XML  : 30 stations, 300 observations

Cohérence : OK

Station la plus ventée : XX-XXX-XXX (XX.X km/h)

Temp. moyenne Bordeaux Mérignac : XX.X °C

Stations par pays : map\[AT:1 BE:1 CH:2 CZ:1 DE:4 DK:1 ES:3 FR:8 IT:3 NL:1 NO:1 PL:1 PT:2 SE:1\]

---

## 10\. Bonus

- Ajouter une fonction `LoadFromAuto(path string) ([]Station, error)` qui détecte le format à partir de l'extension (`.json` / `.xml`) et délègue.  
- Mesurer les temps de chargement des deux formats (`time.Now()` avant/après). Lequel est le plus rapide à votre avis ? Vérifiez.  
- Écrire un test unitaire qui valide qu'un `Station` parsé depuis JSON est strictement égal à son équivalent parsé depuis XML (attention aux arrondis flottants).

---

## 11\. Ouverture sur le TP du soir

À 17h vous aurez un module Go propre qui sait charger un dataset météo depuis deux sources hétérogènes et le servir sous une forme unifiée. C'est exactement ce qu'il faut pour la suite : ce soir, vous allez exposer ces données via une **API REST CRUD** en Go pur (`net/http` \+ `encoding/json`, pas de framework). Routes prévues :

GET    /stations

GET    /stations/{id}

GET    /stations/{id}/observations

POST   /stations

PUT    /stations/{id}

DELETE /stations/{id}

Donc gardez votre code propre — ce qui est fait maintenant vous servira directement après le dîner.

---

## 12\. Livrables et notation

Vous rendrez un dossier zippé contenant le module Go, un `README.md` (avec votre tableau de la partie A \+ une phrase de description par fichier), et la sortie de `l'exécution du script`

Barème indicatif (20 points) :

- Partie A — tableau comparatif rempli et juste : **2 pts**  
- Partie B — modèle interne propre et neutre : **3 pts**  
- Partie C — parser JSON fonctionnel et conversions correctes : **4 pts**  
- Partie D — parser XML fonctionnel, attributs et chardata maîtrisés : **5 pts**  
- Partie E — requêtes correctes et programme `explore` qui tourne : **4 pts**  
- Qualité de code (nommage, séparation des responsabilités, gestion d'erreurs) : **2 pts**

---

## 13\. Pièges fréquents (à lire **avant** de déboguer pendant 45 min)

- Si votre décodage JSON renvoie une slice vide, c'est presque toujours un tag mal écrit ou un type intermédiaire mal imbriqué. Vérifiez d'abord la structure cible avant de soupçonner le fichier.  
- En XML, ne mélangez pas `xml:"name"` (élément enfant) et `xml:"name,attr"` (attribut). Relisez le fichier et catégorisez chaque donnée avant de tagger.  
- `time.Parse` est strict sur le format. Le format de référence est `2006-01-02T15:04:05Z07:00` (`time.RFC3339`) pour les timestamps, `2006-01-02` pour les dates `installed_on`. Si vous obtenez `parsing time "..." as "..."` , c'est un format de référence qui ne colle pas.  
- N'essayez pas de tagger votre struct `Station` interne pour qu'elle décode directement le JSON ET le XML. Ça finit mal. Passez par des types intermédiaires.

