package setting

import "time"

type JarvisSetting struct {
	// 채널 ID.
	ID string `json:"id"`

	// 생성 일시.
	CreatedAt time.Time `json:"created_at"`
	// 마지막 수정 일시.
	UpdatedAt time.Time `json:"updated_at"`

	// 스크럼 알림 여부.
	ScrumNotificationEnabled bool `json:"scrum_notification_enabled"`
	// 스크럼 알림 시간.
	ScrumNotificationTime string `json:"scrum_notification_time"`
}
