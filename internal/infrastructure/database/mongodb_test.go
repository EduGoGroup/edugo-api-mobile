package database

import (
	"context"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

// TestInitMongoDB_Success verifica que la inicialización de MongoDB funcione correctamente
func TestInitMongoDB_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Levantar contenedor de MongoDB para testing
	mongoContainer, err := mongodb.RunContainer(ctx,
		testcontainers.WithImage("mongo:6"),
		mongodb.WithUsername("testuser"),
		mongodb.WithPassword("testpass"),
	)
	require.NoError(t, err)
	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Logf("Error terminando contenedor: %v", err)
		}
	}()

	// Obtener string de conexión del contenedor
	connStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Configurar config de prueba
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			MongoDB: config.MongoDBConfig{
				URI:      connStr,
				Database: "testdb",
				Timeout:  10 * time.Second,
			},
		},
	}

	// Crear logger de prueba
	testLogger := logger.NewZapLogger("debug", "json")

	// Ejecutar la función a testear
	db, err := InitMongoDB(ctx, cfg, testLogger)

	// Verificaciones
	require.NoError(t, err, "No debería haber error al inicializar MongoDB")
	require.NotNil(t, db, "La conexión DB no debería ser nil")

	// Verificar que podemos hacer operaciones
	collection := db.Collection("test")
	_, err = collection.InsertOne(ctx, bson.M{"test": "data"})
	assert.NoError(t, err, "Debería poder insertar un documento")

	// Verificar que podemos leer
	var result bson.M
	err = collection.FindOne(ctx, bson.M{"test": "data"}).Decode(&result)
	assert.NoError(t, err, "Debería poder leer el documento")
	assert.Equal(t, "data", result["test"], "El valor debería ser 'data'")

	// Limpiar
	err = collection.Drop(ctx)
	assert.NoError(t, err, "No debería haber error al eliminar la colección")
}

// TestInitMongoDB_InvalidURI verifica manejo de errores con URI inválido
func TestInitMongoDB_InvalidURI(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Config con URI inválido
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			MongoDB: config.MongoDBConfig{
				URI:      "mongodb://invalid-host-that-does-not-exist:27017",
				Database: "testdb",
				Timeout:  2 * time.Second,
			},
		},
	}

	testLogger := logger.NewZapLogger("debug", "json")

	db, err := InitMongoDB(ctx, cfg, testLogger)

	// Debería fallar al hacer ping
	assert.Error(t, err, "Debería haber error con URI inválido")
	assert.Contains(t, err.Error(), "error pinging mongodb", "El error debería mencionar el ping")
	assert.Nil(t, db, "La base de datos debería ser nil en caso de error")
}

// TestInitMongoDB_ContextTimeout verifica que respete el timeout del contexto
func TestInitMongoDB_ContextTimeout(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Config con timeout muy corto
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			MongoDB: config.MongoDBConfig{
				URI:      "mongodb://invalid-host:27017",
				Database: "testdb",
				Timeout:  1 * time.Millisecond, // Timeout muy corto
			},
		},
	}

	testLogger := logger.NewZapLogger("debug", "json")

	db, err := InitMongoDB(ctx, cfg, testLogger)

	// Debería fallar por timeout
	assert.Error(t, err, "Debería haber error por timeout")
	assert.Nil(t, db, "La base de datos debería ser nil en caso de error")
}

// TestInitMongoDB_MultipleOperations verifica que podamos realizar múltiples operaciones
func TestInitMongoDB_MultipleOperations(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Levantar contenedor de MongoDB
	mongoContainer, err := mongodb.RunContainer(ctx,
		testcontainers.WithImage("mongo:6"),
		mongodb.WithUsername("testuser"),
		mongodb.WithPassword("testpass"),
	)
	require.NoError(t, err)
	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Logf("Error terminando contenedor: %v", err)
		}
	}()

	connStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			MongoDB: config.MongoDBConfig{
				URI:      connStr,
				Database: "testdb",
				Timeout:  10 * time.Second,
			},
		},
	}

	testLogger := logger.NewZapLogger("debug", "json")

	db, err := InitMongoDB(ctx, cfg, testLogger)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Insertar múltiples documentos
	collection := db.Collection("test_multi")
	docs := []interface{}{
		bson.M{"name": "doc1", "value": 1},
		bson.M{"name": "doc2", "value": 2},
		bson.M{"name": "doc3", "value": 3},
	}

	_, err = collection.InsertMany(ctx, docs)
	assert.NoError(t, err, "Debería poder insertar múltiples documentos")

	// Contar documentos
	count, err := collection.CountDocuments(ctx, bson.M{})
	assert.NoError(t, err, "Debería poder contar documentos")
	assert.Equal(t, int64(3), count, "Debería haber 3 documentos")

	// Buscar documentos
	cursor, err := collection.Find(ctx, bson.M{})
	assert.NoError(t, err, "Debería poder buscar documentos")
	defer cursor.Close(ctx)

	var results []bson.M
	err = cursor.All(ctx, &results)
	assert.NoError(t, err, "Debería poder decodificar resultados")
	assert.Len(t, results, 3, "Debería haber 3 resultados")

	// Actualizar un documento
	updateResult, err := collection.UpdateOne(ctx,
		bson.M{"name": "doc1"},
		bson.M{"$set": bson.M{"value": 10}},
	)
	assert.NoError(t, err, "Debería poder actualizar un documento")
	assert.Equal(t, int64(1), updateResult.ModifiedCount, "Debería haber actualizado 1 documento")

	// Eliminar un documento
	deleteResult, err := collection.DeleteOne(ctx, bson.M{"name": "doc3"})
	assert.NoError(t, err, "Debería poder eliminar un documento")
	assert.Equal(t, int64(1), deleteResult.DeletedCount, "Debería haber eliminado 1 documento")

	// Verificar que ahora hay 2 documentos
	count, err = collection.CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count, "Debería haber 2 documentos restantes")

	// Limpiar
	err = collection.Drop(ctx)
	assert.NoError(t, err)
}

// TestInitMongoDB_DatabaseName verifica que use el nombre de base de datos correcto
func TestInitMongoDB_DatabaseName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mongoContainer, err := mongodb.RunContainer(ctx,
		testcontainers.WithImage("mongo:6"),
		mongodb.WithUsername("testuser"),
		mongodb.WithPassword("testpass"),
	)
	require.NoError(t, err)
	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Logf("Error terminando contenedor: %v", err)
		}
	}()

	connStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	expectedDBName := "my_custom_database"
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			MongoDB: config.MongoDBConfig{
				URI:      connStr,
				Database: expectedDBName,
				Timeout:  10 * time.Second,
			},
		},
	}

	testLogger := logger.NewZapLogger("debug", "json")

	db, err := InitMongoDB(ctx, cfg, testLogger)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Verificar que el nombre de la base de datos es el esperado
	assert.Equal(t, expectedDBName, db.Name(), "El nombre de la base de datos debería coincidir")
}

// TestInitMongoDB_PingFailure simula un fallo en el ping después de conectar
func TestInitMongoDB_PingFailure(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Levantar contenedor y luego detenerlo inmediatamente para simular fallo
	mongoContainer, err := mongodb.RunContainer(ctx,
		testcontainers.WithImage("mongo:6"),
		mongodb.WithUsername("testuser"),
		mongodb.WithPassword("testpass"),
	)
	require.NoError(t, err)

	connStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Terminar el contenedor antes de intentar conectar
	err = mongoContainer.Terminate(ctx)
	require.NoError(t, err)

	// Esperar un poco para asegurar que el contenedor esté terminado
	time.Sleep(500 * time.Millisecond)

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			MongoDB: config.MongoDBConfig{
				URI:      connStr,
				Database: "testdb",
				Timeout:  2 * time.Second,
			},
		},
	}

	testLogger := logger.NewZapLogger("debug", "json")

	db, err := InitMongoDB(ctx, cfg, testLogger)

	// Debería fallar porque el contenedor ya no existe
	assert.Error(t, err, "Debería haber error porque el servidor no está disponible")
	assert.Nil(t, db, "La base de datos debería ser nil")
}

// TestInitMongoDB_ConcurrentConnections verifica conexiones concurrentes
func TestInitMongoDB_ConcurrentConnections(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mongoContainer, err := mongodb.RunContainer(ctx,
		testcontainers.WithImage("mongo:6"),
		mongodb.WithUsername("testuser"),
		mongodb.WithPassword("testpass"),
	)
	require.NoError(t, err)
	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Logf("Error terminando contenedor: %v", err)
		}
	}()

	connStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			MongoDB: config.MongoDBConfig{
				URI:      connStr,
				Database: "testdb",
				Timeout:  10 * time.Second,
			},
		},
	}

	testLogger := logger.NewZapLogger("debug", "json")

	db, err := InitMongoDB(ctx, cfg, testLogger)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Ejecutar operaciones concurrentes
	collection := db.Collection("concurrent_test")
	done := make(chan bool)

	for i := 0; i < 20; i++ {
		go func(id int) {
			_, err := collection.InsertOne(ctx, bson.M{"id": id, "data": "test"})
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Esperar a que todas las goroutines terminen
	for i := 0; i < 20; i++ {
		<-done
	}

	// Verificar que se insertaron todos los documentos
	count, err := collection.CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Equal(t, int64(20), count, "Debería haber 20 documentos insertados")

	// Limpiar
	err = collection.Drop(ctx)
	assert.NoError(t, err)
}
