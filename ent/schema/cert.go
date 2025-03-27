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
	return []ent.Field{
		field.Strings("domains"),field.String("email").NotEmpty().MaxLen(200),
		field.JSON("config", map[string]string{}),field.String("certPath").Optional().MaxLen(1000),
		field.String("keyPath").Optional().MaxLen(1000),field.Bool("auto"),field.String("result").Optional(),
		
	}
}

// Edges of the Cert.
func (Cert) Edges() []ent.Edge {
	return nil
}
