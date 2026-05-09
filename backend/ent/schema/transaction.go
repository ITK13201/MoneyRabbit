package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Transaction holds the schema definition for the Transaction entity.
type Transaction struct{ ent.Schema }

func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.Time("date"),
		field.String("description"),
		field.Int("amount"),
		field.Enum("import_format_id").Values("smbc_bank", "smbc_card"),
		field.Time("imported_at").Default(time.Now).Immutable(),
	}
}

func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", Category.Type).Ref("transactions").Unique(),
	}
}
