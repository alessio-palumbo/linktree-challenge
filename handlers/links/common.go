package links

import (
	"encoding/json"

	"github.com/alessio-palumbo/linktree-challenge/handlers/models"
)

// addSublink unmarshal the given metadata in the correct sublink model and append it to the Link object.
// It returns the parsed model as an interface for further processing or validation.
// Note: If the sublink payload matches any, but not all the fields of the model, the matching fields
// 	 will still be parsed and the sublink considered correct
func addSublink(l *models.Link, subID string, metadata json.RawMessage) (interface{}, error) {

	// TODO this could be improved to reduce code duplication using reflection
	switch l.Type {
	case models.LinkMusic:
		sb := models.Platform{ID: subID}

		err := json.Unmarshal(metadata, &sb)
		if err != nil {
			return nil, err
		}

		l.SubLinks = append(l.SubLinks, sb)
		return sb, nil
	case models.LinkShows:
		sb := models.Show{ID: subID}

		err := json.Unmarshal(metadata, &sb)
		if err != nil {
			return nil, err
		}

		l.SubLinks = append(l.SubLinks, sb)
		return sb, nil
	}

	return nil, nil
}
