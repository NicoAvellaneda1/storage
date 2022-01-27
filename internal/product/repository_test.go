package product

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/nicoavellaneda1/storage/internal/db"
	models "github.com/nicoavellaneda1/storage/internal/models"

	//"github.com/nicoavellaneda1/storage/internal/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetByName(t *testing.T) {
	repo := NewRepo(db.StorageDB)

	expectedProduct := models.Product{ID: 6, Name: "TV", Type: "Smart", Count: 1, Price: 10.0}

	result := repo.GetByName(expectedProduct.Name)

	assert.Equal(t, expectedProduct, result)
}

func TestStore(t *testing.T) {
	repo := NewRepo(db.StorageDB)

	expectedProduct := models.Product{ID: 13, Name: "HH", Type: "Smart", Count: 1, Price: 20.0}

	product := models.Product{Name: "HH", Type: "Smart", Count: 1, Price: 20.0}
	result, err := repo.Store(product)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, expectedProduct, result)
	assert.NoError(t, err)
}

func TestGetAll(t *testing.T) {
	repo := NewRepo(db.StorageDB)

	result, err := repo.GetAll()
	if err != nil {
		t.Fail()
	}

	assert.Len(t, result, 5)
	assert.NoError(t, err)
}

func TestUpdateWithContext(t *testing.T) {
	repo := NewRepo(db.StorageDB)

	expectedProduct := models.Product{ID: 4, Name: "PC", Type: "Smart", Count: 1, Price: 30.0}

	product := models.Product{ID: 4, Name: "PC", Type: "Smart", Count: 1, Price: 30.0}
	//se define un context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := repo.UpdateWithContext(ctx, product)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, expectedProduct, result)
	assert.NoError(t, err)
}

//Test con go-txdb
func Test_sqlRepository_Store(t *testing.T) {
	db, err := InitDb()
	assert.NoError(t, err)
	repository := NewRepo(db)

	expectedProduct := models.Product{ID: 15, Name: "Aire", Type: "Smart", Count: 1, Price: 20.0}

	product := models.Product{Name: "Aire", Type: "Smart", Count: 1, Price: 20.0}

	result, err := repository.Store(product)
	assert.NoError(t, err)
	getResult := repository.GetOne(9)
	assert.NotNil(t, getResult)
	assert.Equal(t, expectedProduct, result)
}

func init() {
	txdb.Register("txdb", "mysql", "meli_sprint_user:Meli_Sprint#123@/storage")
}

func InitDb() (*sql.DB, error) {
	db, err := sql.Open("txdb", uuid.New().String())
	if err == nil {
		return db, db.Ping()
	}
	return db, err
}
