package dao

import (
    "database/sql"
    "time"
    "errors"
    "golang.org/x/crypto/bcrypt"
    _ "github.com/mattn/go-sqlite3"
    "home-app/internal/models"
)

const (
	dbPath = "internal/db/sqlite/home.db"
)

var dao *sql.DB

func InitDB() error {
	var err error
    dao, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() {
    if dao != nil {
        dao.Close()
    }
}

func ValidateUser(userId int, username, password string) (bool, error) {
    var hashedPassword string

    err := dao.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
    if err != nil {
        return false, err
    }

    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        return false, errors.New("password does not match")
    }

    return true, nil
}

func AddExpense(expense templates.Expense) (bool, error) {
    _, err := dao.Exec("INSERT INTO expenses (userId, name, amount, category, createdAt) VALUES (?, ?, ?, ?, ?)", 
        expense.UserId, expense.Name, expense.Amount, expense.Category, expense.Datetime)
    if err != nil {
        return false, err
    }

    return true, nil
}

func UpdateExpense(expense templates.Expense) (bool, error) {
    var userId string

    err := dao.QueryRow("SELECT userId FROM expenses WHERE expenseId = ?", expense.ExpenseId).Scan(&userId)
    if err != nil {
        return false, err
    }

    if userId != expense.UserId {
        return false, errors.New("user ID does not match")
    }

    _, err = dao.Exec("UPDATE expenses SET name = ?, amount = ?, category = ?, createdAt = ? WHERE expenseId = ?", 
        expense.Name, expense.Amount, expense.Category, expense.Datetime, expense.ExpenseId)
    if err != nil {
        return false, err
    }

    return true, nil
}

func GetExpensesForCurrentMonth(userId int) ([]templates.Expense, error) {
    currentMonth := time.Now().Format("2006-01")
    rows, err := dao.Query("SELECT expenseId, userId, name, amount, category, createdAt FROM expenses WHERE strftime('%Y-%m', createdAt) = ? AND userId = ? ORDER BY createdAt DESC", currentMonth, userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var expenses []templates.Expense
    for rows.Next() {
        var expense templates.Expense
        err := rows.Scan(&expense.ExpenseId, &expense.UserId, &expense.Name, &expense.Amount, &expense.Category, &expense.Datetime)
        if err != nil {
            return nil, err
        }
        expenses = append(expenses, expense)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return expenses, nil
}