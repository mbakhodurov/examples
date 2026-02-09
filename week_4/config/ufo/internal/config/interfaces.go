package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type UFOGRPCCONFIG interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}
