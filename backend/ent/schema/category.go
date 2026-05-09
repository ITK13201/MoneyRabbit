package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Category holds the schema definition for the Category entity.
type Category struct{ ent.Schema }

func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.String("name"),
		field.String("color"),
		field.String("icon"),
		field.Enum("type").Values("income", "expense", "both"),
		field.Int("sort_order").Default(0),
	}
}

func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("rules", CategoryRule.Type),
		edge.To("transactions", Transaction.Type),
	}
}
