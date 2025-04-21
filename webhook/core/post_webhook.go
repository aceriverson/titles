package titles

import (
	"errors"

	"titles.run/webhook/models"
	"titles.run/webhook/utils"
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

	if user.Plan == models.UserPlanNone {
		return nil
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

	if update.Description == "" {
		update.Description = "Titled via titles․run"
	} else {
		update.Description += "\n\nTitled via titles․run"
	}

	routeMap, err := h.Map.GenerateMap(polyline.Points)
	if err != nil {
		return errors.New("failed to generate map")
	}

	poi, err := h.Here.GetPOI(polyline.Flex, polyline.Points[len(polyline.Points)/2])
	if err != nil {
		return errors.New("failed to get poi")
	}

	title, err := h.AI.Title(user.Plan, activity, polygons, routeMap, poi)
	if err != nil {
		return errors.New("failed to get title")
	}

	update.Name = title

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
