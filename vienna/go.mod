module github.com/zoroqi/hamster/vienna

go 1.16

require (
	github.com/mattn/go-sqlite3 v1.14.7
	github.com/zoroqi/hamster/clipboard v0.0.0-incompatible
)

replace github.com/zoroqi/hamster/clipboard => ../clipboard
