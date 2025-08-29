module github.com/thuatus/EnterpriteCA/V-0.0.2/user

go 1.24.2

replace github.com/thuatus/EnterpriteCA/V-0.0.2/user => ../user

replace github.com/thuatus/EnterpriteCA/V-0.0.2/db => ../db

require (
	github.com/pquerna/otp v1.5.0
	github.com/thuatus/EnterpriteCA/V-0.0.2/db v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.41.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
)
