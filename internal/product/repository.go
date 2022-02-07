package product

import (
	"context"
	"database/sql"
	"log"

	"github.com/nicoavellaneda1/storage/internal/db"
	models "github.com/nicoavellaneda1/storage/internal/models"
)

const (
	GetName       = "select id, name, type, count, price from products where name = ?"
	InsertProduct = "INSERT INTO products(name, type, count, price) VALUES( ?, ?, ?, ? )"
	GetAll        = "SELECT id, name, type, count, price FROM products"
	UpdateProduct = "UPDATE products SET name = ?, type = ?, count = ?, price = ? WHERE id = ?"
	GetId         = "select id, name, type, count, price from products where id = ?"
	DeleteProduct = "DELETE FROM products WHERE id = ?"
)

type Repository interface {
	Store(product models.Product) (models.Product, error)
	GetByName(name string) models.Product
	GetOne(id int) models.Product
	Update(product models.Product) (models.Product, error)
	GetAll() ([]models.Product, error)
	Delete(id int) error
	UpdateWithContext(ctx context.Context, product models.Product) (models.Product, error)
}
type repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetByName(name string) models.Product {
	var product models.Product
	db := db.StorageDB
	rows, err := db.Query(GetName, name)
	if err != nil {
		log.Println(err)
		return product
	}
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price); err != nil {
			log.Println(err.Error())
			return product
		}
	}
	return product
}

func (r *repository) Store(product models.Product) (models.Product, error) {
	//db := db.StorageDB                     // se inicializa la base
	stmt, err := r.db.Prepare(InsertProduct) // se prepara el SQL
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // se cierra la sentencia al terminar. Si quedan abiertas se genera consumos de memoria
	var result sql.Result
	result, err = stmt.Exec(product.Name, product.Type, product.Count, product.Price) // retorna un sql.Result y un error
	if err != nil {
		return models.Product{}, err
	}
	insertedId, _ := result.LastInsertId() // del sql.Resul devuelto en la ejecución obtenemos el Id insertado
	product.ID = int(insertedId)
	return product, nil
}

func (r *repository) GetAll() ([]models.Product, error) {
	var products []models.Product
	db := db.StorageDB
	rows, err := db.Query(GetAll)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// se recorren todas las filas
	for rows.Next() {
		// por cada fila se obtiene un objeto del tipo Product
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price); err != nil {
			log.Fatal(err)
			return nil, err
		}
		//se añade el objeto obtenido al slice products
		products = append(products, product)
	}
	return products, nil
}

func (r *repository) Update(product models.Product) (models.Product, error) {
	db := db.StorageDB                     // se inicializa la base
	stmt, err := db.Prepare(UpdateProduct) // se prepara la sentencia SQL a ejecutar
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // se cierra la sentencia al terminar. Si quedan abiertas se genera consumos de memoria
	_, err = stmt.Exec(product.Name, product.Type, product.Count, product.Price, product.ID)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (r *repository) GetOne(id int) models.Product {
	var product models.Product
	db := db.StorageDB
	rows, err := db.Query(GetId, id)
	if err != nil {
		log.Println(err)
		return product
	}
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price); err != nil {
			log.Println(err.Error())
			return product
		}
	}
	return product
}

func (r *repository) Delete(id int) error {
	db := db.StorageDB                     // se inicializa la base
	stmt, err := db.Prepare(DeleteProduct) // se prepara la sentencia SQL a ejecutar
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()     // se cierra la sentencia al terminar. Si quedan abiertas se genera consumos de memoria
	_, err = stmt.Exec(id) // retorna un sql.Result y un error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateWithContext(ctx context.Context, product models.Product) (models.Product, error) {
	db := db.StorageDB                     // se inicializa la base
	stmt, err := db.Prepare(UpdateProduct) // se prepara la sentencia SQL a ejecutar
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()                                                                                   // se cierra la sentencia al terminar. Si quedan abiertas se genera consumos de memoria
	_, err = stmt.ExecContext(ctx, product.Name, product.Type, product.Count, product.Price, product.ID) //envio el contexto
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}
