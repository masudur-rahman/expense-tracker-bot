package configs

import "fmt"

var PurrfectConfig PawsitiveConfiguration

type PawsitiveConfiguration struct {
	Server   ServerConfig   `json:"server" yaml:"server"`
	GRPC     GRPCConfig     `json:"grpc" yaml:"grpc"`
	Database DatabaseConfig `json:"database" yaml:"database"`
	Session  SessionConfig  `json:"session" yaml:"session"`
}

type ServerConfig struct {
	Host   string `json:"host" yaml:"host"`
	Port   int    `json:"port" yaml:"port"`
	Domain string `json:"domain" yaml:"domain"`
}

type GRPCConfig struct {
	ServerHost string `json:"serverHost" yaml:"serverHost"`
	ClientHost string `json:"clientHost" yaml:"clientHost"`
	Port       int    `json:"port" yaml:"port"`
}

type DatabaseConfig struct {
	Type     DatabaseType     `json:"type" yaml:"type"`
	ArangoDB DBConfigArangoDB `json:"arangodb" yaml:"arangodb"`
	Postgres DBConfigPostgres `json:"postgres" yaml:"postgres"`
}

type DatabaseType string

const (
	DatabaseArangoDB DatabaseType = "arangodb"
	DatabasePostgres DatabaseType = "postgres"
)

type DBConfigArangoDB struct {
	Name     string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
}

type DBConfigPostgres struct {
	Name     string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	SSLMode  string `json:"sslmode" yaml:"sslmode"`
}

func (cp DBConfigPostgres) String() string {
	return fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v sslmode=%v", cp.User, cp.Password, cp.Name, cp.Host, cp.Port, cp.SSLMode)
}

type SessionConfig struct {
	Name       string `json:"name" yaml:"name"`
	HttpOnly   bool   `json:"httpOnly" yaml:"httpOnly"`
	CSRFSecret string `json:"csrfSecret" yaml:"csrfSecret"`
	CSRFHeader string `json:"csrfHeader" yaml:"csrfHeader"`
	CSRFForm   string `json:"csrfForm" yaml:"csrfForm"`
	SessionKey string `json:"sessionKey" yaml:"sessionKey"`
}
