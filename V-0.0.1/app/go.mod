module github.com/thuatus/EnterpriteCA/tree/main/V-0.0.1/app

go 1.24.2

replace github.com/thuatus/EnterpriteCA/tree/main/V-0.0.1/ca => ../ca

require (
	github.com/gorilla/sessions v1.4.0
	github.com/thuatus/EnterpriteCA/tree/main/V-0.0.1/ca v0.0.0-00010101000000-000000000000
	github.com/thuatus/EnterpriteCA/tree/main/V-0.0.1/db v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.38.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
)

replace github.com/thuatus/EnterpriteCA/tree/main/V-0.0.1/db => ../db
