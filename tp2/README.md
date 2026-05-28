 ## Partie A — Exploration des deux fichiers

| Donnée               | Comment c'est représenté en JSON ?                                                            | Comment c'est représenté en XML ?                                                                      |
|:---------------------|:----------------------------------------------------------------------------------------------|:-------------------------------------------------------------------------------------------------------|
| Pays                 | stations[n].country                                                                           | weather_dataset.station[@id="n"][@country="?"]                                                         |
| Coordonnées          | stations[n].location.latitude<br/>stations[n].location.longitude                              | weather_dataset.station[@id="n"].coordinates[@lat="?"][@lon="?"]                                       |
| Altitude             | stations[n].altitude_m                                                                        | weather_dataset.station[@id="n"].coordinates[@altitude="?"]                                            |
| Modèle de capteur    | stations[n].device.type                                                                       | weather_dataset.station[@id="n"].hardware[@model="?"]                                                  |
| Température          | stations[n].observations[n].temperature_celsius                                               | weather_dataset.station[@id="n"].observations.observation[@at="m"].measure[@type="temperature"].text() |
| Conditions ciel      | stations[n].observations[n].conditions                                                        | weather_dataset.station[@id="n"].observations.observation[@at="m"][@sky="?"]                           |
| Vent                 | stations[n].observations[n].wind.speed_kmh<br/>stations[n].observations[n].wind.direction_deg | weather_dataset.station[@id="n"].observations.observation[@at="m"].wind[@speed="?"][@direction="?"]    |
| Notes (optionnelles) | stations[n].observations[n].notes                                                             | weather_dataset.station[@id="n"].observations.observation[@at="m"].note.text()                         |

## Une phrase de description par fichier
- enonce_tp_meteo.md : L'énoncé du TP
- go.mod : Le module go du projet qui contient la version go utilisée.
- jsonsource.go : Le fichier qui permet d'extraire les données d'un fichier au format JSON en les convertissant dans un premier temps en un modèle d'entrée puis ensuite en un modèle unifié.
- main.go : Le fichier d'entrée du projet qui contient le programme de la partie E du TP ainsi que les tests des différentes parties commentées.
- model.go : Le fichier contient le modèle unifié du projet.
- query.go : Le fichier contient les différentes requêtes permettant d'intérroger les données formattées à partir du modèle unifié.
- README.md : Le fichier contenant le tableau de la partie A et une section intitulée "Une phrase de description par fichier".
- utils.go : Le fichier contient une fonction utilitaire utilisée dans les fichiers "jsonsource.go" et "xmlsource.go"
- weather_data.json : Le fichier contenant les données au format JSON. 
- weather_data.xml : Le fichier contenant les données au format XML.
- xmlsource.go : Le fichier qui permet d'extraire les données d'un fichier au format XML en les convertissant dans un premier temps en un modèle d'entrée puis ensuite en un modèle unifié.
- sortie_execution_programme.png : La capture d'écran montrant la sortie d'exécution du programme de la partie E dans le terminal.
