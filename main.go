package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var products []Product

func main() {
	// Cargar productos desde el archivo JSON al iniciar
	loadProducts()

	// Crear router de Gin
	router := gin.Default()

	// Rutas
	router.GET("/ping", pingHandler)
	router.GET("/products", getProductsHandler)
	router.GET("/products/:id", getProductByIDHandler)
	router.GET("/products/search", searchProductsByPriceHandler)

	// Iniciar servidor en puerto 8080
	fmt.Println("Servidor iniciado en http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}

// loadProducts carga los productos desde el archivo JSON
func loadProducts() {
	file, err := ioutil.ReadFile("products.json")
	if err != nil {
		log.Fatal("Error al leer el archivo products.json:", err)
	}

	err = json.Unmarshal(file, &products)
	if err != nil {
		log.Fatal("Error al parsear el JSON:", err)
	}

	fmt.Printf("Se cargaron %d productos desde products.json\n", len(products))
}

// pingHandler responde con "pong" y status 200
func pingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// getProductsHandler devuelve todos los productos
func getProductsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

// getProductByIDHandler devuelve un producto por su ID
func getProductByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	for _, product := range products {
		if product.ID == id {
			c.JSON(http.StatusOK, product)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
}

// searchProductsByPriceHandler busca productos con precio mayor a priceGt
func searchProductsByPriceHandler(c *gin.Context) {
	priceGtStr := c.Query("priceGt")
	if priceGtStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro priceGt es requerido"})
		return
	}

	priceGt, err := strconv.ParseFloat(priceGtStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "priceGt debe ser un número válido"})
		return
	}

	var filteredProducts []Product
	for _, product := range products {
		if product.Price > priceGt {
			filteredProducts = append(filteredProducts, product)
		}
	}

	c.JSON(http.StatusOK, filteredProducts)
}
