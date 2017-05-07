package db

var (
	DriverProstgres postgres
	DriverMongo     mongo
)

type driver interface {
	Name() string
}

// Postgres driver
type postgres struct{ driver }

func (p postgres) Name() string { return "postgres" }

// Mongo driver
type mongo struct{ driver }

func (m mongo) Name() string { return "mongo" }

type DB interface {
	Close() error
	Driver() driver
	Schema() string
	CollectionName(string) string
	Connection() interface{}
	HasTable(string) (bool, error)
	CreateTable(string) error
	ClearTable(string) error
}
