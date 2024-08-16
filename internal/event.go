package internal

import (
	"encoding/json"
	"wikitracker/pkg/models"
)

func JsonToEvent(jsonStr string) (*models.Event, error) {
	var event models.Event
	err := json.Unmarshal([]byte(jsonStr), &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
