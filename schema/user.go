package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Unique().NotEmpty(),
		field.String("email").Unique().NotEmpty(),
		field.String("password").Sensitive().NotEmpty(),
		field.Time("date_joined").Default(time.Now),
		field.Bool("is_active").Default(true),
		field.Bool("is_staff").Default(false),
		field.Bool("is_superuser").Default(false),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
