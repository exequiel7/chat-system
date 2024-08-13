package databases

import (
	"chat-system/cmd/api/config"
	"chat-system/cmd/api/models"
	"fmt"
	"reflect"
	"strings"

	"github.com/gocql/gocql"
	logger "github.com/sirupsen/logrus"
)

type CassandraDB struct {
	Session  *gocql.Session
	Keyspace string
}

func (db *CassandraDB) GetDbClient() *gocql.Session {
	return db.Session
}

func NewCassandraDB() *CassandraDB {
	config := config.GetConfig()
	cluster := gocql.NewCluster(config.CassandraHost)
	cluster.Port = config.CassandraPort
	cluster.Consistency = gocql.Quorum

	// Crear el keyspace y la sesión temporal
	tmpSession, err := cluster.CreateSession()
	if err != nil {
		logger.Fatalf("failed to connect to Cassandra: %v", err)
	}

	// Crear el keyspace si no existe
	err = tmpSession.Query(fmt.Sprintf(`
        CREATE KEYSPACE IF NOT EXISTS %s 
        WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
    `, config.CassandraKeyspace)).Exec()
	if err != nil {
		logger.Fatalf("failed to create keyspace: %v", err)
	}

	// Cerrar la sesión temporal
	tmpSession.Close()

	// Crear la sesión con el keyspace
	cluster.Keyspace = config.CassandraKeyspace
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Fatalf("failed to connect to Cassandra: %v", err)
	}

	logger.Info("Cassandra db connected")

	// Crear la instancia de CassandraDB
	db := &CassandraDB{
		Session:  session,
		Keyspace: config.CassandraKeyspace,
	}

	// Crear las tablas y los índices
	db.CreateSchema()

	return db
}

// Método para crear una tabla con una clave compuesta específica para "messages"
func CreateMessagesTable(session *gocql.Session, keyspace string) error {
	query := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s.messages (
            conversation_id uuid,
            id uuid,
            sender_id uuid,
            receiver_id uuid,
            content text,
            timestamp timestamp,
            PRIMARY KEY (conversation_id, timestamp, id)
        ) WITH CLUSTERING ORDER BY (timestamp ASC);
    `, keyspace)
	return session.Query(query).Exec()
}

// Método para crear tablas genéricas basado en los modelos
func CreateGenericTable(session *gocql.Session, keyspace, tableName string, model interface{}) error {
	t := reflect.TypeOf(model)
	fields := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		cqlTag := field.Tag.Get("cql")
		if cqlTag != "" {
			fields = append(fields, cqlTag)
		}
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (%s)", keyspace, tableName, strings.Join(fields, ", "))
	return session.Query(query).Exec()
}

func CreateIndex(session *gocql.Session, keyspace, tableName, columnName, indexName string) error {
	query := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s.%s (%s);", indexName, keyspace, tableName, columnName)
	return session.Query(query).Exec()
}

func (db *CassandraDB) CreateSchema() {
	// Crear la tabla de usuarios
	err := CreateGenericTable(db.Session, db.Keyspace, "users", models.User{})
	if err != nil {
		logger.Fatalf("Error creating users table: %v", err)
	}
	logger.Info("Users table created")

	// Crear índice en la tabla de usuarios
	err = CreateIndex(db.Session, db.Keyspace, "users", "username", "username_idx")
	if err != nil {
		logger.Fatalf("Error creating index on username: %v", err)
	}

	// Crear la tabla de mensajes con clave compuesta
	err = CreateMessagesTable(db.Session, db.Keyspace)
	if err != nil {
		logger.Fatalf("Error creating messages table: %v", err)
	}
	logger.Info("Messages table created")
}
