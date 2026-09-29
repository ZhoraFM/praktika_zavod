package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// ---------- МОДЕЛИ (структуры данных) ----------

// Сотрудник (работник завода)
type Employee struct {
	ID             int    `json:"id" db:"id"`
	FullName       string `json:"full_name" db:"full_name"`
	Phone          string `json:"phone" db:"phone"`
	Address        string `json:"address" db:"address"`
	GraduationYear int    `json:"graduation_year" db:"graduation_year"`
	Position       string `json:"position" db:"position"` // из таблицы Квалификация
}

// Заказчик
type Customer struct {
	ID              int      `json:"id" db:"id"`
	FullName        string   `json:"full_name" db:"full_name"`
	Address         string   `json:"address" db:"address"`
	RegularDiscount *float64 `json:"regular_discount" db:"regular_discount"`
}

// Продукция
type Product struct {
	ID          int    `json:"id" db:"id"`
	ProductName string `json:"product_name" db:"product_name"`
	Cost        int    `json:"cost" db:"cost"`
}

// Заказ
type Order struct {
	ID          int       `json:"id" db:"id"`
	CustomerID  int       `json:"customer_id" db:"customer_id"`
	Description string    `json:"description" db:"description"`
	OrderDate   time.Time `json:"order_date" db:"order_date"`
}

// Партия (выпуск продукции)
type Production struct {
	ID             int       `json:"id" db:"id"`
	OrderID        int       `json:"order_id" db:"order_id"`
	EmployeeID     int       `json:"employee_id" db:"employee_id"`
	Description    string    `json:"description" db:"description"`
	ProductionDate time.Time `json:"production_date" db:"production_date"`
}

// Связь партии с продукцией
type ProductionProduct struct {
	ProductionID int `json:"production_id" db:"production_id"`
	ProductID    int `json:"product_id" db:"product_id"`
}

// ---------- РЕПОЗИТОРИЙ (работа с БД) ----------

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// ----- СОТРУДНИКИ -----
func (r *Repository) GetEmployees() ([]Employee, error) {
	query := `
		SELECT e.*, q.position 
		FROM Employees e 
		LEFT JOIN Qualifications q ON e.id = q.employee_id
	`
	var employees []Employee
	err := r.db.Select(&employees, query)
	return employees, err
}

func (r *Repository) GetEmployeeByID(id int) (*Employee, error) {
	query := `
		SELECT e.*, q.position 
		FROM Employees e 
		LEFT JOIN Qualifications q ON e.id = q.employee_id 
		WHERE e.id = ?
	`
	var employee Employee
	err := r.db.Get(&employee, query, id)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *Repository) CreateEmployee(employee *Employee) error {
	query := `
		INSERT INTO Employees (full_name, phone, address, graduation_year) 
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, employee.FullName, employee.Phone, employee.Address, employee.GraduationYear)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	employee.ID = int(id)

	if employee.Position != "" {
		queryQual := "INSERT INTO Qualifications (employee_id, position) VALUES (?, ?)"
		_, err = r.db.Exec(queryQual, employee.ID, employee.Position)
	}
	return err
}

func (r *Repository) UpdateEmployee(employee *Employee) error {
	query := `
		UPDATE Employees 
		SET full_name = ?, phone = ?, address = ?, graduation_year = ? 
		WHERE id = ?
	`
	_, err := r.db.Exec(query, employee.FullName, employee.Phone, employee.Address, employee.GraduationYear, employee.ID)
	if err != nil {
		return err
	}

	if employee.Position != "" {
		queryQual := `
			INSERT INTO Qualifications (employee_id, position) 
			VALUES (?, ?) 
			ON DUPLICATE KEY UPDATE position = ?
		`
		_, err = r.db.Exec(queryQual, employee.ID, employee.Position, employee.Position)
	}
	return err
}

func (r *Repository) DeleteEmployee(id int) error {
	_, err := r.db.Exec("DELETE FROM Qualifications WHERE employee_id = ?", id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM Employees WHERE id = ?", id)
	return err
}

// ----- ЗАКАЗЧИКИ -----
func (r *Repository) GetCustomers() ([]Customer, error) {
	var customers []Customer
	err := r.db.Select(&customers, "SELECT * FROM Customers")
	return customers, err
}

func (r *Repository) GetCustomerByID(id int) (*Customer, error) {
	var customer Customer
	err := r.db.Get(&customer, "SELECT * FROM Customers WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *Repository) CreateCustomer(customer *Customer) error {
	query := "INSERT INTO Customers (full_name, address, regular_discount) VALUES (?, ?, ?)"
	result, err := r.db.Exec(query, customer.FullName, customer.Address, customer.RegularDiscount)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	customer.ID = int(id)
	return nil
}

func (r *Repository) UpdateCustomer(customer *Customer) error {
	query := "UPDATE Customers SET full_name = ?, address = ?, regular_discount = ? WHERE id = ?"
	_, err := r.db.Exec(query, customer.FullName, customer.Address, customer.RegularDiscount, customer.ID)
	return err
}

func (r *Repository) DeleteCustomer(id int) error {
	_, err := r.db.Exec("DELETE FROM Customers WHERE id = ?", id)
	return err
}

// ----- ПРОДУКЦИЯ -----
func (r *Repository) GetProducts() ([]Product, error) {
	var products []Product
	err := r.db.Select(&products, "SELECT * FROM Products")
	return products, err
}

func (r *Repository) GetProductByID(id int) (*Product, error) {
	var product Product
	err := r.db.Get(&product, "SELECT * FROM Products WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Repository) CreateProduct(product *Product) error {
	query := "INSERT INTO Products (product_name, cost) VALUES (?, ?)"
	result, err := r.db.Exec(query, product.ProductName, product.Cost)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = int(id)
	return nil
}

func (r *Repository) UpdateProduct(product *Product) error {
	query := "UPDATE Products SET product_name = ?, cost = ? WHERE id = ?"
	_, err := r.db.Exec(query, product.ProductName, product.Cost, product.ID)
	return err
}

func (r *Repository) DeleteProduct(id int) error {
	_, err := r.db.Exec("DELETE FROM Products WHERE id = ?", id)
	return err
}

// ----- ЗАКАЗЫ -----
func (r *Repository) GetOrders() ([]Order, error) {
	var orders []Order
	err := r.db.Select(&orders, "SELECT * FROM Orders")
	return orders, err
}

func (r *Repository) GetOrderByID(id int) (*Order, error) {
	var order Order
	err := r.db.Get(&order, "SELECT * FROM Orders WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *Repository) GetOrdersByCustomer(customerID int) ([]Order, error) {
	var orders []Order
	err := r.db.Select(&orders, "SELECT * FROM Orders WHERE customer_id = ?", customerID)
	return orders, err
}

func (r *Repository) CreateOrder(order *Order) error {
	query := "INSERT INTO Orders (customer_id, description, order_date) VALUES (?, ?, ?)"
	result, err := r.db.Exec(query, order.CustomerID, order.Description, order.OrderDate)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	order.ID = int(id)
	return nil
}

func (r *Repository) UpdateOrder(order *Order) error {
	query := "UPDATE Orders SET customer_id = ?, description = ?, order_date = ? WHERE id = ?"
	_, err := r.db.Exec(query, order.CustomerID, order.Description, order.OrderDate, order.ID)
	return err
}

func (r *Repository) DeleteOrder(id int) error {
	_, err := r.db.Exec("DELETE FROM Orders WHERE id = ?", id)
	return err
}

// ----- ПАРТИИ (ВЫПУСК ПРОДУКЦИИ) -----
func (r *Repository) GetProduction() ([]Production, error) {
	var apps []Production
	err := r.db.Select(&apps, "SELECT * FROM Production")
	return apps, err
}

func (r *Repository) GetProductionByID(id int) (*Production, error) {
	var app Production
	err := r.db.Get(&app, "SELECT * FROM Production WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *Repository) GetProductionByOrder(orderID int) ([]Production, error) {
	var apps []Production
	err := r.db.Select(&apps, "SELECT * FROM Production WHERE order_id = ?", orderID)
	return apps, err
}

func (r *Repository) CreateProduction(app *Production) error {
	query := "INSERT INTO Production (order_id, employee_id, description, production_date) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(query, app.OrderID, app.EmployeeID, app.Description, app.ProductionDate)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	app.ID = int(id)
	return nil
}

func (r *Repository) UpdateProduction(app *Production) error {
	query := "UPDATE Production SET order_id = ?, employee_id = ?, description = ?, production_date = ? WHERE id = ?"
	_, err := r.db.Exec(query, app.OrderID, app.EmployeeID, app.Description, app.ProductionDate, app.ID)
	return err
}

func (r *Repository) DeleteProduction(id int) error {
	_, err := r.db.Exec("DELETE FROM ProductionProducts WHERE production_id = ?", id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM Production WHERE id = ?", id)
	return err
}

// ----- СВЯЗИ ПАРТИЙ С ПРОДУКЦИЕЙ -----
func (r *Repository) AddProductToProduction(productionID, productID int) error {
	query := "INSERT INTO ProductionProducts (production_id, product_id) VALUES (?, ?)"
	_, err := r.db.Exec(query, productionID, productID)
	return err
}

func (r *Repository) RemoveProductFromProduction(productionID, productID int) error {
	query := "DELETE FROM ProductionProducts WHERE production_id = ? AND product_id = ?"
	_, err := r.db.Exec(query, productionID, productID)
	return err
}

func (r *Repository) GetProductionProducts(productionID int) ([]Product, error) {
	query := `
		SELECT p.* 
		FROM Products p
		JOIN ProductionProducts pp ON p.id = pp.product_id
		WHERE pp.production_id = ?
	`
	var products []Product
	err := r.db.Select(&products, query, productionID)
	return products, err
}

func (r *Repository) GetProductionWithProducts(productionID int) (*Production, []Product, error) {
	app, err := r.GetProductionByID(productionID)
	if err != nil {
		return nil, nil, err
	}
	products, err := r.GetProductionProducts(productionID)
	if err != nil {
		return app, nil, err
	}
	return app, products, nil
}

// ---------- ОБРАБОТЧИКИ (HTTP Handlers) ----------

// ----- СОТРУДНИКИ -----
func GetEmployees(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employees, err := repo.GetEmployees()
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка: %v", err), http.StatusInternalServerError)
			return
		}
		if len(employees) == 0 {
			w.Write([]byte("Сотрудников нет"))
			return
		}
		for _, e := range employees {
			w.Write([]byte(fmt.Sprintf("ID: %d, ФИО: %s, Должность: %s, Телефон: %s\n",
				e.ID, e.FullName, e.Position, e.Phone)))
		}
	}
}

func GetEmployeeByID(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		employee, err := repo.GetEmployeeByID(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Сотрудник с ID %d не найден", id), http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("Сотрудник: %s, Должность: %s, Телефон: %s, Адрес: %s, Год окончания: %d",
			employee.FullName, employee.Position, employee.Phone, employee.Address, employee.GraduationYear)))
	}
}

func CreateEmployee(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employee Employee
		err := json.NewDecoder(r.Body).Decode(&employee)
		if err != nil {
			http.Error(w, fmt.Sprintf("Некорректный JSON: %v", err), http.StatusBadRequest)
			return
		}
		err = repo.CreateEmployee(&employee)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка создания сотрудника: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(employee)
	}
}

func UpdateEmployee(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		var employee Employee
		err = json.NewDecoder(r.Body).Decode(&employee)
		if err != nil {
			http.Error(w, fmt.Sprintf("Некорректный JSON: %v", err), http.StatusBadRequest)
			return
		}
		employee.ID = id
		err = repo.UpdateEmployee(&employee)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка обновления: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Сотрудник обновлён!"))
	}
}

func DeleteEmployee(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		err = repo.DeleteEmployee(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка удаления: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Сотрудник удалён!"))
	}
}

// ----- ЗАКАЗЧИКИ -----
func GetCustomers(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customers, err := repo.GetCustomers()
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка: %v", err), http.StatusInternalServerError)
			return
		}
		if len(customers) == 0 {
			w.Write([]byte("Заказчиков нет"))
			return
		}
		for _, c := range customers {
			discount := 0.0
			if c.RegularDiscount != nil {
				discount = *c.RegularDiscount
			}
			w.Write([]byte(fmt.Sprintf("ID: %d, Наименование: %s, Адрес: %s, Скидка: %.0f%%\n",
				c.ID, c.FullName, c.Address, discount*100)))
		}
	}
}

func GetCustomerByID(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		customer, err := repo.GetCustomerByID(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Заказчик с ID %d не найден", id), http.StatusNotFound)
			return
		}
		discount := 0.0
		if customer.RegularDiscount != nil {
			discount = *customer.RegularDiscount
		}
		w.Write([]byte(fmt.Sprintf("Заказчик: %s, Адрес: %s, Скидка: %.0f%%",
			customer.FullName, customer.Address, discount*100)))
	}
}

func CreateCustomer(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customer Customer
		err := json.NewDecoder(r.Body).Decode(&customer)
		if err != nil {
			http.Error(w, fmt.Sprintf("Некорректный JSON: %v", err), http.StatusBadRequest)
			return
		}
		err = repo.CreateCustomer(&customer)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка создания заказчика: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(customer)
	}
}

func DeleteCustomer(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		err = repo.DeleteCustomer(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка удаления: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Заказчик удалён!"))
	}
}

// ----- ПРОДУКЦИЯ -----
func GetProducts(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := repo.GetProducts()
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка: %v", err), http.StatusInternalServerError)
			return
		}
		if len(products) == 0 {
			w.Write([]byte("Продукции нет"))
			return
		}
		for _, p := range products {
			w.Write([]byte(fmt.Sprintf("ID: %d, Продукция: %s, Стоимость: %d руб.\n", p.ID, p.ProductName, p.Cost)))
		}
	}
}

func GetProductByID(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		product, err := repo.GetProductByID(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Продукция с ID %d не найдена", id), http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("Продукция: %s, Стоимость: %d руб.", product.ProductName, product.Cost)))
	}
}

func CreateProduct(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, fmt.Sprintf("Некорректный JSON: %v", err), http.StatusBadRequest)
			return
		}
		err = repo.CreateProduct(&product)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка создания продукции: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}

func DeleteProduct(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		err = repo.DeleteProduct(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка удаления: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Продукция удалена!"))
	}
}

// ----- ЗАКАЗЫ -----
func GetOrders(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := repo.GetOrders()
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка: %v", err), http.StatusInternalServerError)
			return
		}
		if len(orders) == 0 {
			w.Write([]byte("Заказов нет"))
			return
		}
		for _, o := range orders {
			customer, _ := repo.GetCustomerByID(o.CustomerID)
			customerName := "неизвестно"
			if customer != nil {
				customerName = customer.FullName
			}
			w.Write([]byte(fmt.Sprintf("ID: %d, Заказчик: %s, Описание: %s, Дата: %s\n",
				o.ID, customerName, o.Description, o.OrderDate.Format("2006-01-02"))))
		}
	}
}

func GetOrderByID(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		order, err := repo.GetOrderByID(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Заказ с ID %d не найден", id), http.StatusNotFound)
			return
		}
		customer, _ := repo.GetCustomerByID(order.CustomerID)
		customerName := "неизвестно"
		if customer != nil {
			customerName = customer.FullName
		}
		w.Write([]byte(fmt.Sprintf("Заказ: %d\nЗаказчик: %s\nОписание: %s\nДата: %s",
			order.ID, customerName, order.Description, order.OrderDate.Format("2006-01-02"))))
	}
}

func CreateOrder(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order Order
		err := json.NewDecoder(r.Body).Decode(&order)
		if err != nil {
			http.Error(w, fmt.Sprintf("Некорректный JSON: %v", err), http.StatusBadRequest)
			return
		}
		err = repo.CreateOrder(&order)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка создания заказа: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	}
}

func DeleteOrder(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		err = repo.DeleteOrder(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка удаления: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Заказ удалён!"))
	}
}

// ----- ПАРТИИ (ВЫПУСК ПРОДУКЦИИ) -----
func GetProductionHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apps, err := repo.GetProduction()
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка: %v", err), http.StatusInternalServerError)
			return
		}
		if len(apps) == 0 {
			w.Write([]byte("Партий нет"))
			return
		}
		for _, a := range apps {
			order, _ := repo.GetOrderByID(a.OrderID)
			employee, _ := repo.GetEmployeeByID(a.EmployeeID)
			orderInfo := "неизвестно"
			if order != nil {
				orderInfo = fmt.Sprintf("Заказ #%d", order.ID)
			}
			employeeName := "неизвестно"
			if employee != nil {
				employeeName = employee.FullName
			}
			products, _ := repo.GetProductionProducts(a.ID)
			productNames := ""
			for i, p := range products {
				if i > 0 {
					productNames += ", "
				}
				productNames += p.ProductName
			}
			w.Write([]byte(fmt.Sprintf("ID: %d, Дата: %s, Описание: %s, %s, Сотрудник: %s, Продукция: [%s]\n",
				a.ID, a.ProductionDate.Format("2006-01-02"), a.Description, orderInfo, employeeName, productNames)))
		}
	}
}

func GetProductionByIDHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		app, products, err := repo.GetProductionWithProducts(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Партия с ID %d не найдена", id), http.StatusNotFound)
			return
		}
		order, _ := repo.GetOrderByID(app.OrderID)
		employee, _ := repo.GetEmployeeByID(app.EmployeeID)

		orderInfo := "неизвестно"
		if order != nil {
			orderInfo = fmt.Sprintf("Заказ #%d", order.ID)
		}
		employeeName := "неизвестно"
		if employee != nil {
			employeeName = employee.FullName
		}

		productNames := ""
		for i, p := range products {
			if i > 0 {
				productNames += ", "
			}
			productNames += fmt.Sprintf("%s (%d руб.)", p.ProductName, p.Cost)
		}

		w.Write([]byte(fmt.Sprintf("Партия ID: %d\nДата: %s\nОписание: %s\n%s\nСотрудник: %s\nПродукция: [%s]",
			app.ID, app.ProductionDate.Format("2006-01-02"), app.Description, orderInfo, employeeName, productNames)))
	}
}

func CreateProductionHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			OrderID        int    `json:"order_id"`
			EmployeeID     int    `json:"employee_id"`
			Description    string `json:"description"`
			ProductionDate string `json:"production_date"`
			ProductIDs     []int  `json:"product_ids"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, fmt.Sprintf("Некорректный JSON: %v", err), http.StatusBadRequest)
			return
		}

		if input.OrderID == 0 {
			http.Error(w, "order_id обязательно", http.StatusBadRequest)
			return
		}
		if input.EmployeeID == 0 {
			http.Error(w, "employee_id обязательно", http.StatusBadRequest)
			return
		}
		if input.Description == "" {
			http.Error(w, "description обязательно", http.StatusBadRequest)
			return
		}

		var productionDate time.Time
		if input.ProductionDate == "" {
			productionDate = time.Now()
		} else {
			productionDate, err = time.Parse("2006-01-02", input.ProductionDate)
			if err != nil {
				productionDate, err = time.Parse("2006-01-02T15:04:05Z", input.ProductionDate)
				if err != nil {
					http.Error(w, "Неверный формат даты. Используйте YYYY-MM-DD или YYYY-MM-DDTHH:MM:SSZ", http.StatusBadRequest)
					return
				}
			}
		}

		app := &Production{
			OrderID:        input.OrderID,
			EmployeeID:     input.EmployeeID,
			Description:    input.Description,
			ProductionDate: productionDate,
		}

		err = repo.CreateProduction(app)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка создания партии: %v", err), http.StatusInternalServerError)
			return
		}

		for _, productID := range input.ProductIDs {
			err = repo.AddProductToProduction(app.ID, productID)
			if err != nil {
				log.Printf("Ошибка добавления продукции %d: %v", productID, err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(app)
	}
}

func UpdateProductionHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		var production Production
		err = json.NewDecoder(r.Body).Decode(&production)
		if err != nil {
			http.Error(w, fmt.Sprintf("Некорректный JSON: %v", err), http.StatusBadRequest)
			return
		}
		production.ID = id

		err = repo.UpdateProduction(&production)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка обновления: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Партия обновлена!"))
	}
}

func DeleteProductionHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		err = repo.DeleteProduction(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка удаления: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Партия удалена!"))
	}
}

// Добавление продукции к партии
func AddProductToProductionHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		productionID, err := strconv.Atoi(vars["productionId"])
		if err != nil {
			http.Error(w, "Неверный ID партии", http.StatusBadRequest)
			return
		}

		var request struct {
			ProductID int `json:"product_id"`
		}
		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Некорректный JSON: %v", err), http.StatusBadRequest)
			return
		}

		err = repo.AddProductToProduction(productionID, request.ProductID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка добавления продукции: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Продукция добавлена к партии!"))
	}
}

// Удаление продукции из партии
func RemoveProductFromProductionHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		productionID, err := strconv.Atoi(vars["productionId"])
		if err != nil {
			http.Error(w, "Неверный ID партии", http.StatusBadRequest)
			return
		}
		productID, err := strconv.Atoi(vars["productId"])
		if err != nil {
			http.Error(w, "Неверный ID продукции", http.StatusBadRequest)
			return
		}

		err = repo.RemoveProductFromProduction(productionID, productID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка удаления продукции: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Продукция удалена из партии!"))
	}
}

// Получение продукции для партии
func GetProductionProductsHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		productionID, err := strconv.Atoi(vars["productionId"])
		if err != nil {
			http.Error(w, "Неверный ID партии", http.StatusBadRequest)
			return
		}

		products, err := repo.GetProductionProducts(productionID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка: %v", err), http.StatusInternalServerError)
			return
		}

		if len(products) == 0 {
			w.Write([]byte("Продукции для этой партии нет"))
			return
		}

		for _, p := range products {
			w.Write([]byte(fmt.Sprintf("Продукция: %s, Стоимость: %d руб.\n", p.ProductName, p.Cost)))
		}
	}
}

// ---------- ГЛАВНАЯ ФУНКЦИЯ ----------

func main() {
	// Подключение к БД (измените на свою)
	dsn := "root:root@tcp(127.0.0.1:3307)/factory?parseTime=true"

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("БД не отвечает:", err)
	}
	log.Println("✅ Подключено к MySQL!")

	repo := NewRepository(db)
	r := mux.NewRouter()

	// === СОТРУДНИКИ ===
	r.HandleFunc("/employees", GetEmployees(repo)).Methods("GET")
	r.HandleFunc("/employees/{id}", GetEmployeeByID(repo)).Methods("GET")
	r.HandleFunc("/employees", CreateEmployee(repo)).Methods("POST")
	r.HandleFunc("/employees/{id}", UpdateEmployee(repo)).Methods("PUT")
	r.HandleFunc("/employees/{id}", DeleteEmployee(repo)).Methods("DELETE")

	// === ЗАКАЗЧИКИ ===
	r.HandleFunc("/customers", GetCustomers(repo)).Methods("GET")
	r.HandleFunc("/customers/{id}", GetCustomerByID(repo)).Methods("GET")
	r.HandleFunc("/customers", CreateCustomer(repo)).Methods("POST")
	r.HandleFunc("/customers/{id}", DeleteCustomer(repo)).Methods("DELETE")

	// === ПРОДУКЦИЯ ===
	r.HandleFunc("/products", GetProducts(repo)).Methods("GET")
	r.HandleFunc("/products/{id}", GetProductByID(repo)).Methods("GET")
	r.HandleFunc("/products", CreateProduct(repo)).Methods("POST")
	r.HandleFunc("/products/{id}", DeleteProduct(repo)).Methods("DELETE")

	// === ЗАКАЗЫ ===
	r.HandleFunc("/orders", GetOrders(repo)).Methods("GET")
	r.HandleFunc("/orders/{id}", GetOrderByID(repo)).Methods("GET")
	r.HandleFunc("/orders", CreateOrder(repo)).Methods("POST")
	r.HandleFunc("/orders/{id}", DeleteOrder(repo)).Methods("DELETE")

	// === ПАРТИИ (ВЫПУСК ПРОДУКЦИИ) ===
	r.HandleFunc("/production", GetProductionHandler(repo)).Methods("GET")
	r.HandleFunc("/production/{id}", GetProductionByIDHandler(repo)).Methods("GET")
	r.HandleFunc("/production", CreateProductionHandler(repo)).Methods("POST")
	r.HandleFunc("/production/{id}", UpdateProductionHandler(repo)).Methods("PUT")
	r.HandleFunc("/production/{id}", DeleteProductionHandler(repo)).Methods("DELETE")

	// === СВЯЗИ ПАРТИЙ С ПРОДУКЦИЕЙ ===
	r.HandleFunc("/production/{productionId}/products", GetProductionProductsHandler(repo)).Methods("GET")
	r.HandleFunc("/production/{productionId}/products", AddProductToProductionHandler(repo)).Methods("POST")
	r.HandleFunc("/production/{productionId}/products/{productId}", RemoveProductFromProductionHandler(repo)).Methods("DELETE")

	log.Println("🚀 Сервер запущен на порту :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
