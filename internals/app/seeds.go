package app

import "github.com/ballot/internals/models"

func init() {

	// add candidates
	models.NewCandidate("trump", "republican")
	models.NewCandidate("kamala", "democratic")

}
