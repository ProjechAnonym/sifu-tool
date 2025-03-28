// Code generated by ent, DO NOT EDIT.

package cert

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the cert type in the database.
	Label = "cert"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldDomains holds the string denoting the domains field in the database.
	FieldDomains = "domains"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldConfig holds the string denoting the config field in the database.
	FieldConfig = "config"
	// FieldCertPath holds the string denoting the certpath field in the database.
	FieldCertPath = "cert_path"
	// FieldKeyPath holds the string denoting the keypath field in the database.
	FieldKeyPath = "key_path"
	// FieldAuto holds the string denoting the auto field in the database.
	FieldAuto = "auto"
	// FieldResult holds the string denoting the result field in the database.
	FieldResult = "result"
	// Table holds the table name of the cert in the database.
	Table = "certs"
)

// Columns holds all SQL columns for cert fields.
var Columns = []string{
	FieldID,
	FieldDomains,
	FieldEmail,
	FieldConfig,
	FieldCertPath,
	FieldKeyPath,
	FieldAuto,
	FieldResult,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// EmailValidator is a validator for the "email" field. It is called by the builders before save.
	EmailValidator func(string) error
	// CertPathValidator is a validator for the "certPath" field. It is called by the builders before save.
	CertPathValidator func(string) error
	// KeyPathValidator is a validator for the "keyPath" field. It is called by the builders before save.
	KeyPathValidator func(string) error
)

// OrderOption defines the ordering options for the Cert queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByEmail orders the results by the email field.
func ByEmail(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEmail, opts...).ToFunc()
}

// ByCertPath orders the results by the certPath field.
func ByCertPath(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCertPath, opts...).ToFunc()
}

// ByKeyPath orders the results by the keyPath field.
func ByKeyPath(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldKeyPath, opts...).ToFunc()
}

// ByAuto orders the results by the auto field.
func ByAuto(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAuto, opts...).ToFunc()
}

// ByResult orders the results by the result field.
func ByResult(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldResult, opts...).ToFunc()
}
