package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Cert holds the schema definition for the Cert entity.
type Cert struct {
	ent.Schema
}

// Fields of the Cert.
func (Cert) Fields() []ent.Field {
	return []ent.Field{field.String("tag").NotEmpty().MaxLen(500)}
}

// Edges of the Cert.
func (Cert) Edges() []ent.Edge {
	return nil
}
