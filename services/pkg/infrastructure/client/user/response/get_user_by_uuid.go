package response

import "time"

type GetUserByUUIDResponse struct {
	Result struct {
		User struct {
			UUID         string `json:"UUID"`
			Email        string `json:"Email"`
			Subscription struct {
				ID       int         `json:"ID"`
				UserUUID string      `json:"UserUUID"`
				User     interface{} `json:"User"`
				PlanCode string      `json:"PlanCode"`
				Plan     struct {
					Code         string    `json:"Code"`
					Name         string    `json:"name"`
					Price        int       `json:"price"`
					DurationDays int       `json:"duration_days"`
					MaxDevices   int       `json:"max_devices"`
					Description  string    `json:"description"`
					CreatedAt    time.Time `json:"created_at"`
				} `json:"Plan"`
				Type      string      `json:"Type"`
				StartedAt time.Time   `json:"StartedAt"`
				ExpiresAt time.Time   `json:"ExpiresAt"`
				IsActive  bool        `json:"IsActive"`
				AutoRenew bool        `json:"AutoRenew"`
				CreatedAt time.Time   `json:"CreatedAt"`
				UpdatedAt time.Time   `json:"UpdatedAt"`
				DeletedAt interface{} `json:"DeletedAt"`
			} `json:"Subscription"`
			ChatId    int         `json:"ChatId"`
			CreatedAt time.Time   `json:"CreatedAt"`
			IsActive  bool        `json:"IsActive"`
			UpdatedAt time.Time   `json:"UpdatedAt"`
			DeletedAt interface{} `json:"DeletedAt"`
		} `json:"user"`
	} `json:"result"`
}
