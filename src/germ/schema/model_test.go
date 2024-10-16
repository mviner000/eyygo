package schema_test

import (
	"database/sql"
	"time"

	"github.com/mviner000/eyygo/src/germ"
	"github.com/mviner000/eyygo/src/germ/utils/tests"
)

type User struct {
	*germ.Model
	Name      *string
	Age       *uint
	Birthday  *time.Time
	Account   *tests.Account
	Pets      []*tests.Pet
	Toys      []*tests.Toy `germ:"polymorphic:Owner"`
	CompanyID *int
	Company   *tests.Company
	ManagerID *uint
	Manager   *User
	Team      []*User           `germ:"foreignkey:ManagerID"`
	Languages []*tests.Language `germ:"many2many:UserSpeak"`
	Friends   []*User           `germ:"many2many:user_friends"`
	Active    *bool
}

type (
	mytime time.Time
	myint  int
	mybool = bool
)

type AdvancedDataTypeUser struct {
	ID           sql.NullInt64
	Name         *sql.NullString
	Birthday     sql.NullTime
	RegisteredAt mytime
	DeletedAt    *mytime
	Active       mybool
	Admin        *mybool
}

type BaseModel struct {
	ID        uint
	CreatedAt time.Time
	CreatedBy *int
	Created   *VersionUser `germ:"foreignKey:CreatedBy"`
	UpdatedAt time.Time
	DeletedAt germ.DeletedAt `germ:"index"`
}

type VersionModel struct {
	BaseModel
	Version int
}

type VersionUser struct {
	VersionModel
	Name     string
	Age      uint
	Birthday *time.Time
}
