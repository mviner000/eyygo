package models

// AuthPermission represents the auth_permission table
type AuthPermission struct {
	ID            uint             `germ:"primaryKey;autoIncrement"`
	Name          string           `germ:"not null"`
	ContentTypeID uint             `germ:"not null"`
	Codename      string           `germ:"not null"`
	ContentType   EyygoContentType `germ:"foreignKey:ContentTypeID"`

	// Use a composite unique index for content_type_id and codename
	UniqueIndex string `germ:"uniqueIndex:content_type_codename,priority:1"`
}
