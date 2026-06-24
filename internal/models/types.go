package models

type Dependency struct {
	Name      string
	Version   string
	Ecosystem string
	Type      string // "prod" | "dev"
	DocURL    string // Optional: pre-resolved URL from lockfile
}

type DocResult struct {
	Name      string
	Version   string
	DocURL    string
	Ecosystem string
	Type      string
}
