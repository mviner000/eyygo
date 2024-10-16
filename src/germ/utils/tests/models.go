package tests

import (
	"database/sql"
	"time"

	"github.com/mviner000/eyygo/src/germ"
)

// User has one `Account` (has one), many `Pets` (has many) and `Toys` (has many - polymorphic)
// He works in a Company (belongs to), he has a Manager (belongs to - single-table), and also managed a Team (has many - single-table)
// He speaks many languages (many to many) and has many friends (many to many - single-table)
// His pet also has one Toy (has one - polymorphic)
// NamedPet is a reference to a named `Pet` (has one)
type User struct {
	germ.Model
	Name      string
	Age       uint
	Birthday  *time.Time
	Account   Account
	Pets      []*Pet
	NamedPet  *Pet
	Toys      []Toy   `germ:"polymorphic:Owner"`
	Tools     []Tools `germ:"polymorphicType:Type;polymorphicId:CustomID"`
	CompanyID *int
	Company   Company
	ManagerID *uint
	Manager   *User
	Team      []User     `germ:"foreignkey:ManagerID"`
	Languages []Language `germ:"many2many:UserSpeak;"`
	Friends   []*User    `germ:"many2many:user_friends;"`
	Active    bool
}

type Account struct {
	germ.Model
	UserID sql.NullInt64
	Number string
}

type Pet struct {
	germ.Model
	UserID *uint
	Name   string
	Toy    Toy `germ:"polymorphic:Owner;"`
}

type Toy struct {
	germ.Model
	Name      string
	OwnerID   string
	OwnerType string
}

type Tools struct {
	germ.Model
	Name     string
	CustomID string
	Type     string
}

type Company struct {
	ID   int
	Name string
}

type Language struct {
	Code string `germ:"primarykey"`
	Name string
}

type Coupon struct {
	ID               int              `germ:"primarykey; size:255"`
	AppliesToProduct []*CouponProduct `germ:"foreignKey:CouponId;constraint:OnDelete:CASCADE"`
	AmountOff        uint32           `germ:"column:amount_off"`
	PercentOff       float32          `germ:"column:percent_off"`
}

type CouponProduct struct {
	CouponId  int    `germ:"primarykey;size:255"`
	ProductId string `germ:"primarykey;size:255"`
	Desc      string
}

type Order struct {
	germ.Model
	Num      string
	Coupon   *Coupon
	CouponID string
}

type Parent struct {
	germ.Model
	FavChildID uint
	FavChild   *Child
	Children   []*Child
}

type Child struct {
	germ.Model
	Name     string
	ParentID *uint
	Parent   *Parent
}
