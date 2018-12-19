package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"database/sql"
	"log"
	_ "encoding/json"
	"strconv"
	"time"
	_ "os"
	_ "strings"
	"math"
)

//---------------------------------------------------------------------------
//Interface
//---------------------------------------------------------------------------
type Product interface {
	getProductType() string
	getRefundable() string
	getPrice() float64
	getTax() float64
	getAmount() float64
}

//---------------------------------------------------------------------------
//Struct - Cart and Base
//---------------------------------------------------------------------------
type Cart struct {
	ListProduct 		[]Product `json:"ListProduct"`
	TotalPrice		  float64	  `json:"TotalPrice"`
	TotalTax				float64		`json:"TotalTax"`
	TotalAmount			float64		`json:"TotalAmount"`
}

type Base struct {
  CartId 		  		string		`json:"CartId"`
  ProductName 		string		`json:"ProductName"`
  ProductTypeCode string		`json:"ProductTypeCode"`
	ProductTypeName	string		`json:"ProductTypeName"`
	Refundable			string		`json:"Refundable"`
  Price		    		string	  `json:"Price"`
	Tax							string		`json:"Tax"`
	Amount					string		`json:"Amount"`
}

//---------------------------------------------------------------------------
//Struct - Food
//---------------------------------------------------------------------------
type Food struct {
	Base
}

func (f Food) getProductType() string {
  return "Food & Beverage"
}

func (f Food) getRefundable() string {
  return "Yes"
}

func (f Food) getPrice() float64 {
	price, _ := strconv.ParseFloat(f.Price, 64)
  return price
}

func (f Food) getTax() float64 {
	price, _ := strconv.ParseFloat(f.Base.Price, 64)
  return 0.1 * price
}

func (f Food) getAmount() float64 {
	price, _ := strconv.ParseFloat(f.Base.Price, 64)
  return price + f.getTax()
}

//---------------------------------------------------------------------------
//Struct - Tobacco
//---------------------------------------------------------------------------
type Tobacco struct {
	Base
}

func (t Tobacco) getProductType() string {
  return "Tobacco"
}

func (t Tobacco) getRefundable() string {
  return "No"
}

func (t Tobacco) getPrice() float64 {
	price, _ := strconv.ParseFloat(t.Price, 64)
  return price
}


func (t Tobacco) getTax() float64 {
	price, _ := strconv.ParseFloat(t.Price, 64)
  return 10 + 0.02 * price
}

func (t Tobacco) getAmount() float64 {
	price, _ := strconv.ParseFloat(t.Price, 64)
  return price + t.getTax()
}

//---------------------------------------------------------------------------
//Struct - Entertainment
//---------------------------------------------------------------------------
type Entertainment struct {
	Base
}

func (e Entertainment) getProductType() string {
  return "Entertainment"
}

func (e Entertainment) getRefundable() string {
  return "No"
}

func (e Entertainment) getPrice() float64 {
	price, _ := strconv.ParseFloat(e.Price, 64)
  return price
}

func (e Entertainment) getTax() float64 {
	price, _ := strconv.ParseFloat(e.Price, 64)
	if (price >= 100) {
		return 0.01 * (price - 100)
	} else {
		return 0
	}
}

func (e Entertainment) getAmount() float64 {
	price, _ := strconv.ParseFloat(e.Price, 64)
  return price + e.getTax()
}

//---------------------------------------------------------------------------
//To get database connection
//---------------------------------------------------------------------------
func getDBConnection() *sql.DB {
	//var dbCredential string = "root:admin@/testdb"
	//var dbCredential string = "testuser:123@/testdb"
	var dbCredential string = "testuser:123@tcp(mysql:3306)/testdb"

	db, err := sql.Open("mysql", dbCredential)
  if err != nil {
		log.Printf("DB Connection [ERROR]")
	}else{
		log.Printf("DB Connection [SUCCESS]")
  }

	if err = db.Ping(); err != nil {
    log.Printf("Error connecting to the database: %s\n", err)
		return db
  }

	return db
}


//---------------------------------------------------------------------------
//Fetch Product From Cart
//API URL: http://127.0.0.1:3001/cart/getProducts?page=1&per_page=20
//---------------------------------------------------------------------------
func GetProductFromCart(c *gin.Context) {
	t := time.Now()

	log.Printf("Called Service GetProductFromCart at %v\n", t.Format(time.RFC3339))

	pageIdx, _ := strconv.Atoi(c.Query("page"))
	limitPerPage, _ := strconv.Atoi(c.Query("per_page"))

	db := getDBConnection()

	getProductUrl := "http://127.0.0.1:3001/cart/getProducts"
	select_query := "SELECT CART_ID, PRODUCT_NAME, TAX_CODE, PRICE FROM CART ORDER BY CART_ID asc"
	count_query := "SELECT COUNT(1) FROM (" + select_query + ") count"

	//Count Data
	rows1,err := db.Query(count_query)
    if err != nil {
			log.Printf("Error fetching cart: %v\n", err)
      return
    }

	var rowCount int
	for rows1.Next() {
	    err = rows1.Scan(&rowCount)
        if err != nil {
            log.Printf("rows.Scan(...) failed.\n\t%s\n", err.Error())
            return
        }
	}

	var pageCount string
	var fromRowNo int
	var toRowNo int

	if rowCount < limitPerPage || limitPerPage == 0 {
		pageCount = "1"
		toRowNo = rowCount
	} else{
		num := float64(rowCount) / float64(limitPerPage)
		pageCount = strconv.Itoa(int(math.Ceil(num)))
		toRowNo = pageIdx * limitPerPage
	}
	fromRowNo = ((pageIdx - 1) * limitPerPage)

	//Paging navigation data
	var nextPageNav string
	var prevPageNav string

	pageCountInt, _ := strconv.Atoi(pageCount)
	if (pageCountInt != 1 && pageIdx < pageCountInt) {
		nextPageNav = getProductUrl + "?page=" + strconv.Itoa(pageIdx+1) + "&per_page=" + strconv.Itoa(limitPerPage)
	}

	if pageIdx > 1 {
		prevPageNav = getProductUrl + "?page=" + strconv.Itoa(pageIdx-1) + "&per_page=" + strconv.Itoa(limitPerPage)
	}

	select_paging_query := select_query + " LIMIT " + strconv.Itoa(fromRowNo) + "," + strconv.Itoa(limitPerPage)
	log.Printf("Select paging query: %v\n", select_paging_query)

	var arrLimit int

	if (pageIdx != 0 && limitPerPage != 0) {
		select_query = select_paging_query

		if (pageIdx == 1 && rowCount < limitPerPage) {
			arrLimit = rowCount
		} else if (pageIdx != 1 && pageIdx == pageCountInt) {
			arrLimit = rowCount - ((pageIdx - 1) * limitPerPage)
		} else {
			arrLimit = limitPerPage
		}
	} else{
		arrLimit = rowCount
	}

	//Fetch Data
	arrayOfCart := make([]Product, arrLimit)

	rows, err := db.Query(select_query)
    if err != nil {
			log.Printf("Error fetching cart: %v\n", err)
      return
    }

	iCounter := 0

	for rows.Next() {
		var cartId					string
		var productName 		string
		var productTypeCode	string
    var price 					string

		err = rows.Scan(&cartId, &productName, &productTypeCode, &price)
    if err != nil {
    	log.Printf("rows.Scan(...) failed.\n\t%s\n", err.Error())
      return
    }

		base := Base{CartId: cartId, ProductName: productName, ProductTypeCode: productTypeCode, Price: price}

		switch productTypeCode {
			case "1":
				item := Food{base}
				item.Base.ProductTypeName = item.getProductType()
				item.Base.Refundable = item.getRefundable()
				item.Base.Tax = strconv.FormatFloat(item.getTax(), 'f', 2, 64)
				item.Base.Amount = strconv.FormatFloat(item.getAmount(), 'f', 2, 64)
				arrayOfCart[iCounter] = item
			case "2":
				item := Tobacco{base}
				item.ProductTypeName = item.getProductType()
				item.Refundable = item.getRefundable()
				item.Tax = strconv.FormatFloat(item.getTax(), 'f', 2, 64)
				item.Amount = strconv.FormatFloat(item.getAmount(), 'f', 2, 64)
				arrayOfCart[iCounter] = item
			case "3":
				item := Entertainment{base}
				item.ProductTypeName = item.getProductType()
				item.Refundable = item.getRefundable()
				item.Tax = strconv.FormatFloat(item.getTax(), 'f', 2, 64)
				item.Amount = strconv.FormatFloat(item.getAmount(), 'f', 2, 64)
				arrayOfCart[iCounter] = item
		}

		iCounter = iCounter + 1
  }

	rows.Close()
  db.Close()

	//Calculate cart summary
	var totalPrice, totalTax, totalAmount float64

	for _, v := range arrayOfCart {
			totalPrice = totalPrice + v.getPrice()
	}

	for _, v := range arrayOfCart {
			totalTax = totalTax + v.getTax()
	}

	for _, v := range arrayOfCart {
			totalAmount = totalAmount + v.getAmount()
	}


	//Define return struct
	var cart Cart
	cart.ListProduct = arrayOfCart
	cart.TotalPrice = totalPrice
	cart.TotalTax = totalTax
	cart.TotalAmount = totalAmount

	//Constract response message
	var result gin.H
	if err != nil {
		result = gin.H {
			"response_status" : "500",
			"response_message" : err.Error(),
			"total" : 0,
			"per_page": 0,
			"current_page": 0,
			"last_page": 0,
			"next_page_url": nil,
			"prev_page_url": nil,
			"from": 0,
			"to": 0,
			"data" : nil,
		}
	} else {
		result = gin.H {
			"response_status" : "200",
			"response_message" : "ok",
			"total": rowCount,
			"per_page": limitPerPage,
			"current_page": pageIdx,
			"last_page": pageCount,
			"next_page_url": nextPageNav,
			"prev_page_url": prevPageNav,
			"from": fromRowNo,
			"to": toRowNo,
			"data" : cart,
		}
	}

	log.Printf("Response message: %v\n", result)
	c.JSON(http.StatusOK, result)
}



//---------------------------------------------------------------------------
//Add Product To Cart
//API URL: http://localhost:3001/cart/add
//---------------------------------------------------------------------------
func AddProductToCart(c *gin.Context) {
	t := time.Now()

	log.Printf("Called Service AddProductToCart at %v\n", t.Format(time.RFC3339))

	var product Base
	c.BindJSON(&product)

	log.Printf("Incoming post request with data ProductName: %s, ProductTypeCode: %s, Price: %s", product.ProductName, product.ProductTypeCode, product.Price)

	insert_query := fmt.Sprintf("INSERT INTO CART (PRODUCT_NAME, TAX_CODE, PRICE) VALUES ('%s', '%s', '%s')", product.ProductName, product.ProductTypeCode, product.Price)

	log.Printf("Insert query: %v\n", insert_query)

	db := getDBConnection()

	_, err := db.Exec(insert_query)
  if err != nil {
		log.Printf("Error inserting product to cart: %v\n", err)
    return
  }

	db.Close()

	//Constract response message
	var result gin.H
	if err != nil {
		result = gin.H {
			"response_status" : "500",
			"response_message" : err.Error(),
		}
	} else {
		result = gin.H {
			"response_status" : "200",
			"response_message" : "ok",
		}
	}

	log.Printf("Response message: %v\n", result)
	c.JSON(http.StatusOK, result)

}
