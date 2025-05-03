package titles

import (
	"errors"

	strava "titles.run/strava/models"
	"titles.run/webhook/models"
)

func (h *TitlesCore) PostTitle(userID int64, body models.TitleRequest) error {
	user, err := h.DB.GetUserInternal(userID)
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

	isRateLimited, err := h.TTLStore.CheckRateLimit(user.ID, user.Plan)
	if err != nil {
		return errors.New("failed to check rate limit")
	}
	if isRateLimited {
		return nil
	}

	if err := h.TTLStore.IncrementRateLimit(user.ID); err != nil {
		return errors.New("failed to increment rate limit")
	}

	activity, err := h.Strava.GetActivity(user, body.ActivityID)
	if err != nil {
		return errors.New("failed to get activity")
	}

	update := strava.Update{Description: activity.Description}

	if update.Description == "" {
		update.Description = "Titled via titles․run"
	} else {
		update.Description += "\n\nTitled via titles․run"
	}

	update.Name = body.Title

	if err := h.Strava.RenameActivity(user, activity, update); err != nil {
		return errors.New("failed to rename activity")
	}

	if err := h.Ntfy.Notify(user, activity, update); err != nil {
		return errors.New("failed to send notification")
	}

	return nil
}
