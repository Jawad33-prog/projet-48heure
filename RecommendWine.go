package main

import (
	"fmt"
	"strings"
)

// Import the Wine structure from main.go// Import the Wine structure from main.go

// Fonction pour recommander un vinn
func RecommendWine(wines []Wine) {
	fmt.Println("Bienvenue dans le système de recommandation de vin !")
	fmt.Println("Répondez aux questions suivantes pour trouver un vin qui vous plaira.")

	// Question 1 : Couleur du vin
	fmt.Print("Préférez-vous un vin rouge, blanc ou rosé ? ")
	var couleur string
	fmt.Scanln(&couleur)
	couleur = strings.ToLower(strings.TrimSpace(couleur))

	// Question 2 : Goût
	fmt.Print("Préférez-vous un vin sec, doux ou fruité ? ")
	var gout string
	fmt.Scanln(&gout)
	gout = strings.ToLower(strings.TrimSpace(gout))

	// Question 3 : Occasion
	fmt.Print("Est-ce pour un repas, un apéritif ou une occasion spéciale ? ")
	var occasion string
	fmt.Scanln(&occasion)
	occasion = strings.ToLower(strings.TrimSpace(occasion))

	// Recherche dans les donnéess
	for _, wine := range wines {
		if wine.Type == couleur && wine.Taste == gout && wine.Occasion == occasion {
			fmt.Printf("\nNous vous recommandons : %s\n", wine.Title)
			return
		}
	}

	// Si aucun vin ne correspond
	fmt.Println("\nDésolé, nous n'avons pas trouvé de vin correspondant à vos préférences.")
	// Si aucun vin ne correspond
	fmt.Println("\nDésolé, nous n'avons pas trouvé de vin correspondant à vos préférences.")
}
