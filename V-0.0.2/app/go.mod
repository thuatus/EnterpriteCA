module github.com/thuatus/EnterpriteCA/V-0.0.2/app

go 1.24.2

replace github.com/thuatus/EnterpriteCA/V-0.0.2/ca => ../ca

replace github.com/thuatus/EnterpriteCA/V-0.0.2/db => ../db

replace github.com/thuatus/EnterpriteCA/V-0.0.2/app => .

replace github.com/thuatus/EnterpriteCA/tree/main/V-0.0.2/app => .

require (
	github.com/gorilla/sessions v1.4.0
	github.com/stretchr/testify v1.10.0
	github.com/thuatus/EnterpriteCA/V-0.0.2/ca v0.0.0-00010101000000-000000000000
	github.com/thuatus/EnterpriteCA/V-0.0.2/db v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.40.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
