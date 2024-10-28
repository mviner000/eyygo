package admin

import (
	"github.com/mviner000/eyygo/models"
	"gorm.io/gorm"
)

var Site *AdminSite

func InitializeAdmin(db *gorm.DB) {
	Site = NewAdminSite(db)

	// Register models
	Site.Register(&models.Note{}, &ModelAdmin{
		ListFields:   []string{"ID", "Title", "Author", "IsPublished", "CreatedAt"},
		SearchFields: []string{"Title", "Content"},
		FilterFields: []string{"IsPublished", "AuthorID"},
		OrderFields:  []string{"CreatedAt", "Title"},
		FormFields:   []string{"Title", "Content", "AuthorID", "IsPublished", "Tags"},
		DB:           db,
	})

	Site.Register(&models.User{}, &ModelAdmin{
		ListFields:   []string{"ID", "Username", "Email", "IsActive", "IsStaff"},
		SearchFields: []string{"Username", "Email", "FirstName", "LastName"},
		FilterFields: []string{"IsActive", "IsStaff", "IsSuperUser"},
		OrderFields:  []string{"Username", "DateJoined"},
		FormFields:   []string{"Username", "Email", "Password", "FirstName", "LastName", "IsActive", "IsStaff"},
		DB:           db,
	})
}
