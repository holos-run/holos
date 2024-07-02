package internal

//go:generate rm -rf ent
//go:generate go run entgo.io/ent/cmd/ent generate --feature sql/upsert --target ./ent ./schema
