package holos

// _DBName is the database name used across multiple holos components in this project
_DBName: "zitadel"

// The canonical login domain for the entire platform.  Zitadel will be active
// on a single cluster at a time, but always accessible from this domain.
_ExternalDomain: "login.\(_Platform.Model.org.domain)"
