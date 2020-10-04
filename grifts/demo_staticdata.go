package grifts

import (
	"fmt"
	"math/rand"
)

func randActivityTitle() string {
	actFirst := []string{
		"Skøre rekorder",
		"Marsmændene kommer",
		"Tør du?",
		"Bålmærke",
		"Banan",
		"Agent",
		"På skovtur",
		"Bygge huler",
		"Fødselsdag",
		"Lege",
		"Lejr for en aften",
		"Lygten i spanden og ræveleg",
		"Pynt bussen",
		"Trylleri",
		"Masterchef Hovedret",
		"Musik",
		"Foto",
		"Blepaintball",
	}

	actPostfix := []string{"møde", "løb", "fest", "mærke"}

	return fmt.Sprintf("%s-%s", actFirst[rand.Intn(len(actFirst))], actPostfix[rand.Intn(len(actPostfix))])
}
