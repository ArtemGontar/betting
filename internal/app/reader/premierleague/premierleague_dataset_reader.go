package premierleague

import (
	"fmt"
	"log"

	"github.com/ArtemGontar/betting/internal/app/model"
	"github.com/PuerkitoBio/goquery"
)

func Teams() []model.Team {
	doc, err := goquery.NewDocument("https://www.premierleague.com/clubs")
	if err != nil {
		log.Fatal(err)
	}
	teams := []model.Team{}
	// Find the review items
	doc.Find(".indexItem").Each(func(i int, s *goquery.Selection) {
		badge, _ := s.Find(".badge-image").Attr("src")
		clubName := s.Find(".clubName").Text()
		stadiumName := s.Find(".stadiumName").Text()
		teams = append(teams, model.Team{
			BadgeLink:   badge,
			TeamName:    clubName,
			StadiumName: stadiumName,
		})
	})
	return teams
}

func Players() []model.Player {
	doc, err := goquery.NewDocument("https://www.premierleague.com/players")
	if err != nil {
		log.Fatal(err)
	}
	players := []model.Player{}
	// Find the review items
	doc.Find(".dataContainer").Find("tr").Each(func(i int, s *goquery.Selection) {
		playerName := s.Find(".playerName").Text()
		position := s.Find(".hide-s").First().Text()
		playerCountry := s.Find(".playerCountry").Text()
		fmt.Println(playerName, position, playerCountry)
		players = append(players, model.Player{
			Name:     playerName,
			Country:  playerCountry,
			Position: position,
		})
	})
	return players
}
