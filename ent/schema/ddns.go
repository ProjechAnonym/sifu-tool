package schema

import (
	"sifu-tool/models"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// DDNS holds the schema definition for the DDNS entity.
type DDNS struct {
	ent.Schema
}

// Fields of the DDNS.
func (DDNS) Fields() []ent.Field {
	return []ent.Field{
		field.Int("v4method").Optional(),field.Int("v6method").Optional(),
		field.String("ipv6").Optional().MaxLen(128),field.String("rev6").Optional().MaxLen(100),
		field.String("ipv4").Optional().MaxLen(16),field.String("rev4").Optional().MaxLen(100),
		field.String("v4script").Optional().MaxLen(1000),field.String("v4interface").Optional().MaxLen(100),
		field.String("v6script").Optional().MaxLen(1000),field.String("v6interface").Optional().MaxLen(100),
		field.JSON("domains", []models.Domain{}),field.JSON("config", map[string]string{}),
		field.JSON("webhook", map[string]string{}).Optional(),
	}
}

// Edges of the DDNS.
func (DDNS) Edges() []ent.Edge {
	return nil
}
