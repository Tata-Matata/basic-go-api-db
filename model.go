package main

import (
	"database/sql"
	"fmt"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {
	query := "SELECT id, name, quantity, price from products"
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	products := []product{}

	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		} else {
			products = append(products, p)
		}
	}

	return products, nil
}

func getProduct(db *sql.DB, id int) (*product, error) {
	query := fmt.Sprintf("SELECT name, quantity, price from products WHERE id=%v", id)
	row := db.QueryRow(query)

	product := product{ID: id}
	err := row.Scan(&product.Name, &product.Quantity, &product.Price)

	if err != nil {
		return nil, err
	}

	return &product, nil
}
