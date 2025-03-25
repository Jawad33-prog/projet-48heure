package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Country to Emoji mapping
var countryToEmoji = map[string]string{
	"US":             "🇺🇸",
	"France":         "🇫🇷",
	"Italy":          "🇮🇹",
	"Spain":          "🇪🇸",
	"Portugal":       "🇵🇹",
	"Argentina":      "🇦🇷",
	"Chile":          "🇨🇱",
	"Australia":      "🇦🇺",
	"New Zealand":    "🇳🇿",
	"South Africa":   "🇿🇦",
	"Germany":        "🇩🇪",
	"Austria":        "🇦🇹",
	"Greece":         "🇬🇷",
	"Canada":         "🇨🇦",
	"Brazil":         "🇧🇷",
	"Bulgaria":       "🇧🇬",
	"Hungary":        "🇭🇺",
	"Slovenia":       "🇸🇮",
	"Romania":        "🇷🇴",
	"Croatia":        "🇭🇷",
	"Georgia":        "🇬🇪",
	"Mexico":         "🇲🇽",
	"Turkey":         "🇹🇷",
	"Israel":         "🇮🇱",
	"Ukraine":        "🇺🇦",
	"Uruguay":        "🇺🇾",
	"Lebanon":        "🇱🇧",
	"Moldova":        "🇲🇩",
	"Czech Republic": "🇨🇿",
	"Serbia":         "🇷🇸",
	"India":          "🇮🇳",
	"China":          "🇨🇳",
	"England":        "🏴󠁧󠁢󠁥󠁮󠁧󠁿",
}

// Wine struct represents the wine data model
type Wine struct {
	Points        int     `json:"points,omitempty"`
	Title         string  `json:"title,omitempty"`
	Description   string  `json:"description,omitempty"`
	TasterName    string  `json:"taster_name,omitempty"`
	TasterTwitter string  `json:"taster_twitter_handle,omitempty"`
	Price         float64 `json:"price,omitempty"`
	Designation   string  `json:"designation,omitempty"`
	Variety       string  `json:"variety,omitempty"`
	Region        string  `json:"region_1,omitempty"`
	Province      string  `json:"province,omitempty"`
	Country       string  `json:"country,omitempty"`
	CountryEmoji  string  `json:"-"`
	Winery        string  `json:"winery,omitempty"`
}

// WineFilter represents the filtering options
type WineFilter struct {
	MinPoints int
	MaxPrice  float64
	Country   string
	Variety   string
}

// Global variable to store wines
var wines []Wine

// loadWinesFromJSON loads wine data from a JSON file
func loadWinesFromJSON(filename string) error {
	// Open the JSON file
	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Attempt to parse as an array of wines
	var wineArray []Wine
	err = json.Unmarshal(file, &wineArray)
	if err == nil {
		// Add country emojis
		for i := range wineArray {
			wineArray[i].CountryEmoji = countryToEmoji[wineArray[i].Country]
		}
		wines = wineArray
		return nil
	}

	// If array parsing fails, try parsing as a single object
	var singleWine Wine
	err = json.Unmarshal(file, &singleWine)
	if err == nil {
		singleWine.CountryEmoji = countryToEmoji[singleWine.Country]
		wines = []Wine{singleWine}
		return nil
	}

	// If both parsing methods fail, try cleaning the JSON
	cleanedJSON := strings.ReplaceAll(string(file), "\n", "")
	cleanedJSON = strings.ReplaceAll(cleanedJSON, "\r", "")
	cleanedJSON = strings.TrimSpace(cleanedJSON)

	// Remove potential BOM (Byte Order Mark)
	if strings.HasPrefix(cleanedJSON, "\xef\xbb\xbf") {
		cleanedJSON = strings.TrimPrefix(cleanedJSON, "\xef\xbb\xbf")
	}

	// Try parsing the cleaned JSON
	err = json.Unmarshal([]byte(cleanedJSON), &wineArray)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	// Add country emojis
	for i := range wineArray {
		wineArray[i].CountryEmoji = countryToEmoji[wineArray[i].Country]
	}
	wines = wineArray
	return nil
}

// filterWines applies the filters to the wine list
func filterWines(filter WineFilter) []Wine {
	var filteredWines []Wine

	for _, wine := range wines {
		// Apply filters
		if wine.Points >= filter.MinPoints &&
			wine.Price <= filter.MaxPrice &&
			(filter.Country == "" || strings.EqualFold(wine.Country, filter.Country)) &&
			(filter.Variety == "" || strings.EqualFold(wine.Variety, filter.Variety)) {
			filteredWines = append(filteredWines, wine)
		}
	}

	return filteredWines
}

// getUniqueValues extracts unique values for a specific field
func getUniqueValues(field string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueValues []string

	for _, wine := range wines {
		var value string
		switch field {
		case "country":
			value = wine.Country
		case "variety":
			value = wine.Variety
		}

		if value != "" && !uniqueMap[value] {
			uniqueMap[value] = true
			uniqueValues = append(uniqueValues, value)
		}
	}

	return uniqueValues
}

// Handler for the main page
func wineMarketplaceHandler(w http.ResponseWriter, r *http.Request) {
	// Parse filter parameters
	minPoints := parseIntParam(r, "minPoints", 0)
	maxPrice := parseFloatParam(r, "maxPrice", 1000)
	country := r.URL.Query().Get("country")
	variety := r.URL.Query().Get("variety")

	// Create filter
	filter := WineFilter{
		MinPoints: minPoints,
		MaxPrice:  maxPrice,
		Country:   country,
		Variety:   variety,
	}

	// Filter wines
	filteredWines := filterWines(filter)

	// Prepare template data
	data := struct {
		Wines         []Wine
		Filter        WineFilter
		Countries     []string
		Varieties     []string
		TotalWines    int
		CountryEmojis map[string]string
	}{
		Wines:         filteredWines,
		Filter:        filter,
		Countries:     getUniqueValues("country"),
		Varieties:     getUniqueValues("variety"),
		TotalWines:    len(filteredWines),
		CountryEmojis: countryToEmoji,
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("templates/marketplace.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Helper function to parse integer parameters
func parseIntParam(r *http.Request, param string, defaultValue int) int {
	valueStr := r.URL.Query().Get(param)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// Helper function to parse float parameters
func parseFloatParam(r *http.Request, param string, defaultValue float64) float64 {
	valueStr := r.URL.Query().Get(param)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

// [Conservez les structs et fonctions existantes comme Wine, countryToEmoji, etc.]

// Nouvelle fonction pour obtenir un pays aléatoire
func getRandomCountry() string {
	countries := getUniqueCountries()
	if len(countries) == 0 {
		return ""
	}
	return countries[rand.Intn(len(countries))]
}

// Nouvelle fonction pour obtenir une région aléatoire pour un pays
func getRandomRegionForCountry(country string) string {
	regions := getUniqueRegionsForCountry(country)
	if len(regions) == 0 {
		return ""
	}
	return regions[rand.Intn(len(regions))]
}

// Nouvelle fonction pour obtenir une variété aléatoire pour une région
func getRandomVarietyForRegion(country, region string) string {
	varieties := getUniqueVarietiesForRegion(country, region)
	if len(varieties) == 0 {
		return ""
	}
	return varieties[rand.Intn(len(varieties))]
}

// Gestionnaire de route pour la sélection aléatoire
func randomWineSelectionHandler(w http.ResponseWriter, r *http.Request) {
	// Initialiser le générateur de nombres aléatoires
	rand.Seed(time.Now().UnixNano())

	// Sélection aléatoire des étapes
	randomCountry := getRandomCountry()
	randomRegion := getRandomRegionForCountry(randomCountry)
	randomVariety := getRandomVarietyForRegion(randomCountry, randomRegion)

	// Filtrer les vins
	selectedWines := filterWinesBySelection(randomCountry, randomRegion, randomVariety)

	// Données pour le template
	data := struct {
		RandomCountry string
		RandomRegion  string
		RandomVariety string
		SelectedWines []Wine
		CountryEmojis map[string]string
	}{
		RandomCountry: randomCountry,
		RandomRegion:  randomRegion,
		RandomVariety: randomVariety,
		SelectedWines: selectedWines,
		CountryEmojis: countryToEmoji,
	}

	// Parser et exécuter le template
	tmpl, err := template.ParseFiles("templates/random-wine-selection.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Nouvelle fonction pour obtenir les pays uniques
func getUniqueCountries() []string {
	uniqueCountries := make(map[string]bool)
	var countries []string

	for _, wine := range wines {
		if wine.Country != "" && !uniqueCountries[wine.Country] {
			uniqueCountries[wine.Country] = true
			countries = append(countries, wine.Country)
		}
	}

	return countries
}

// Nouvelle fonction pour obtenir les régions d'un pays spécifique
func getUniqueRegionsForCountry(country string) []string {
	uniqueRegions := make(map[string]bool)
	var regions []string

	for _, wine := range wines {
		if strings.EqualFold(wine.Country, country) && wine.Region != "" && !uniqueRegions[wine.Region] {
			uniqueRegions[wine.Region] = true
			regions = append(regions, wine.Region)
		}
	}

	return regions
}

// Nouvelle fonction pour obtenir les variétés d'une région spécifique
func getUniqueVarietiesForRegion(country, region string) []string {
	uniqueVarieties := make(map[string]bool)
	var varieties []string

	for _, wine := range wines {
		if strings.EqualFold(wine.Country, country) &&
			strings.EqualFold(wine.Region, region) &&
			wine.Variety != "" &&
			!uniqueVarieties[wine.Variety] {
			uniqueVarieties[wine.Variety] = true
			varieties = append(varieties, wine.Variety)
		}
	}

	return varieties
}

// Nouvelle fonction pour filtrer les vins selon pays, région et variété
func filterWinesBySelection(country, region, variety string) []Wine {
	var filteredWines []Wine

	for _, wine := range wines {
		if strings.EqualFold(wine.Country, country) &&
			(region == "" || strings.EqualFold(wine.Region, region)) &&
			(variety == "" || strings.EqualFold(wine.Variety, variety)) {
			filteredWines = append(filteredWines, wine)
		}
	}

	return filteredWines
}

// Gestionnaire de route pour la nouvelle approche de sélection
func wineSelectionHandler(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	region := r.URL.Query().Get("region")
	variety := r.URL.Query().Get("variety")

	// Données pour le template
	data := struct {
		Countries       []string
		Regions         []string
		Varieties       []string
		SelectedWines   []Wine
		CountryEmojis   map[string]string
		SelectedCountry string
		SelectedRegion  string
		SelectedVariety string
	}{
		Countries:       getUniqueCountries(),
		CountryEmojis:   countryToEmoji,
		SelectedCountry: country,
		SelectedRegion:  region,
		SelectedVariety: variety,
	}

	// Si un pays est sélectionné, récupérer ses régions
	if country != "" {
		data.Regions = getUniqueRegionsForCountry(country)
	}

	// Si une région est sélectionnée, récupérer ses variétés
	if country != "" && region != "" {
		data.Varieties = getUniqueVarietiesForRegion(country, region)
	}

	// Si tous les filtres sont définis, sélectionner les vins
	if country != "" && region != "" && variety != "" {
		data.SelectedWines = filterWinesBySelection(country, region, variety)
	}

	// Parser et exécuter le template
	tmpl, err := template.ParseFiles("templates/wine-selection.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Charger les vins
	err := loadWinesFromJSON("wine-data-set.json")
	if err != nil {
		fmt.Printf("Erreur de chargement des données de vin : %v\n", err)
		os.Exit(1)
	}

	// Servir les fichiers statiques
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route pour la sélection aléatoire
	http.HandleFunc("/", randomWineSelectionHandler)

	// Démarrer le serveur
	fmt.Println("Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
