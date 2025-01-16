package services

import (
	"database/sql"
	"fmt"

	"go-clickhouse-example/models"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type DBService struct {
	conn *sql.DB
}

func (db *DBService) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(query, args...)
}

func NewDBService(clickhouseURL string) *DBService {
	conn, err := sql.Open("clickhouse", clickhouseURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to ClickHouse: %v", err))
	}
	return &DBService{conn: conn}
}
func (db *DBService) CreateTable() {
	// Create items table
	itemTableQuery := `
	CREATE TABLE IF NOT EXISTS items (
		id UInt64,
		name String,
		price Float64
	) ENGINE = MergeTree()
	ORDER BY id
	`
	if _, err := db.conn.Exec(itemTableQuery); err != nil {
		panic(fmt.Sprintf("Failed to create items table: %v", err))
	}

	// Create sequence table for items
	sequenceTableQuery := `
	CREATE TABLE IF NOT EXISTS item_sequence (
		last_id UInt64
	) ENGINE = TinyLog
	`
	if _, err := db.conn.Exec(sequenceTableQuery); err != nil {
		panic(fmt.Sprintf("Failed to create sequence table: %v", err))
	}

	_, err := db.conn.Exec("INSERT INTO item_sequence (last_id) VALUES (0)")
	if err != nil {
		fmt.Println("Sequence table already initialized")
	}

	// Create users table
	usersTableQuery := `
CREATE TABLE  IF NOT EXISTS users (
    user_id UInt64 PRIMARY KEY,
    username String,
    password String,
    role String
) ENGINE = MergeTree()
ORDER BY user_id;
	`
	if _, err := db.conn.Exec(usersTableQuery); err != nil {
		panic(fmt.Sprintf("Failed to create users table: %v", err))
	}
}

func (db *DBService) SaveItem(item *models.ItemResponse) error {
	var nextID uint64
	err := db.conn.QueryRow("SELECT COALESCE(MAX(last_id), 0) + 1 AS next_id FROM item_sequence").Scan(&nextID)
	if err != nil {
		return fmt.Errorf("failed to fetch next ID: %w", err)
	}

	_, err = db.conn.Exec("INSERT INTO item_sequence (last_id) VALUES (?)", nextID)
	if err != nil {
		return fmt.Errorf("failed to update sequence table: %w", err)
	}

	query := `INSERT INTO items (id, name, price) VALUES (?, ?, ?)`
	_, err = db.conn.Exec(query, nextID, item.Name, item.Price)
	if err != nil {
		return fmt.Errorf("failed to insert item into database: %w", err)
	}

	item.ID = nextID
	return nil
}

func (db *DBService) GetItemByID(id uint64) (models.ItemResponse, error) {
	query := `SELECT id, name, price FROM items WHERE id = ?`
	row := db.conn.QueryRow(query, id)

	var item models.ItemResponse
	err := row.Scan(&item.ID, &item.Name, &item.Price)
	if err != nil {
		return models.ItemResponse{}, err
	}

	return item, nil
}
func (db *DBService) UpdateItem(id uint64, item models.ItemResponse) error {
	query := `ALTER TABLE items UPDATE name = ?, price = ? WHERE id = ?`
	_, err := db.conn.Exec(query, item.Name, item.Price, id)
	return err
}

func (db *DBService) DeleteItem(id uint64) error {
	query := `DELETE FROM items WHERE id = ?`
	_, err := db.conn.Exec(query, id)
	return err
}

func (db *DBService) SaveUser(user *models.User) error {
	// Fetch the next available user_id from the user_sequence table
	var nextUserID uint64
	err := db.conn.QueryRow("SELECT COALESCE(MAX(last_user_id), 0) + 1 AS next_user_id FROM user_sequence").Scan(&nextUserID)
	if err != nil {
		return fmt.Errorf("failed to fetch next user_id: %w", err)
	}

	// Update the user_sequence table with the new user_id
	_, err = db.conn.Exec("INSERT INTO user_sequence (last_user_id) VALUES (?)", nextUserID)
	if err != nil {
		return fmt.Errorf("failed to update user_sequence table: %w", err)
	}

	// Insert the new user with the generated user_id
	query := `INSERT INTO users (user_id, username, password, role) VALUES (?, ?, ?, ?)`
	_, err = db.conn.Exec(query, nextUserID, user.Username, user.Password, user.Role)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// Assign the generated user_id to the user model
	user.ID = nextUserID
	return nil
}

// GetUserByUsername retrieves a user by their username
func (db *DBService) GetUserByUsername(username string) (models.UserResponse, error) {
	query := `SELECT user_id, username, password, role FROM users WHERE username = ?`
	row := db.conn.QueryRow(query, username)

	var user models.UserResponse
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return models.UserResponse{}, err
	}
	return user, nil
}

func (db *DBService) GetAllItems() ([]models.ItemResponse, error) {
	// Query to get all items
	query := `SELECT id, name, price FROM items`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}
	defer rows.Close()

	var items []models.ItemResponse
	// Iterate through the rows and append each item to the items slice
	for rows.Next() {
		var item models.ItemResponse
		err := rows.Scan(&item.ID, &item.Name, &item.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while fetching items: %w", err)
	}

	return items, nil
}
