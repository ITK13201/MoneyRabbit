package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// CategoryRule holds the schema definition for the CategoryRule entity.
type CategoryRule struct{ ent.Schema }

func (CategoryRule) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.String("keyword"),
		field.Int("priority").Default(0),
	}
}

func (CategoryRule) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", Category.Type).Ref("rules").Unique().Required(),
	}
}
