package models

type NotificationProps struct {
	Message string
	Type    string
}

type PageProps struct {
	Title        string
	Path         string
	Notification NotificationProps
}
