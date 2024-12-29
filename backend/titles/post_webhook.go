package titles

import (
	"errors"

	"titles.run/titles/models"
	"titles.run/titles/utils"
)

func (h *TitlesCore) PostWebhook(webhook models.Webhook) error {
	is_duplicate, err := h.Dedupe.DedupeActivity(webhook.ObjectID)
	if err != nil {
		return errors.New("failed to dedupe activity")
	}
	if is_duplicate {
		return nil
	}

	user, err := h.DB.GetUserInternal(webhook.OwnerID)
	if err != nil {
		return errors.New("user not found")
	}

	user, err = h.Strava.RefreshUser(user)
	if err != nil {
		return errors.New("failed to refresh user")
	}

	if err := h.DB.UpdateUser(user); err != nil {
		return errors.New("failed to update user")
	}

	activity, err := h.Strava.GetActivity(user, webhook.ObjectID)
	if err != nil {
		return errors.New("failed to get activity")
	}

	if activity.Map.Polyline == "" {
		return nil
	}

	polyline := utils.PolylineFromGoogle(activity.Map.Polyline)
	if err := polyline.Validate(); err != nil {
		return errors.New("invalid polyline")
	}

	polygons, err := h.DB.GetIntersectingPolygons(user.ID, polyline.Points)
	if err != nil {
		return errors.New("failed to get intersecting polygons")
	}

	update := models.Update{Description: activity.Description}

	if !user.AI && len(polygons) > 0 {
		update.Name = polygons[0].Name

		if update.Description == "" {
			update.Description = "Titled via titles․run"
		} else {
			update.Description += "\n\nTitled via titles․run"
		}
	} else if user.AI {
		if update.Description == "" {
			update.Description = "Titled via titles․run/ai"
		} else {
			update.Description += "\n\nTitled via titles․run/ai"
		}

		routeMap, err := h.Map.GenerateMap(polyline.Points)
		if err != nil {
			return errors.New("failed to generate map")
		}

		poi, err := h.Here.GetPOI(polyline.Flex, polyline.Points[len(polyline.Points)/2])
		if err != nil {
			return errors.New("failed to get poi")
		}

		title, err := h.AI.Title(activity, polygons, routeMap, poi)
		if err != nil {
			return errors.New("failed to get title")
		}

		update.Name = title
	} else {
		return nil
	}

	if err := h.Strava.RenameActivity(user, activity, update); err != nil {
		return errors.New("failed to rename activity")
	}

	if err := h.Dedupe.AddActivity(webhook.ObjectID); err != nil {
		return errors.New("failed to add activity")
	}

	if err := h.Ntfy.Notify(user, activity, update); err != nil {
		return errors.New("failed to send notification")
	}

	return nil
}

func (h *TitlesCore) UnauthorizeUser(userID int64) error {
	return h.DB.UnauthorizeUser(userID)
}
