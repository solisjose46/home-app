package dao

import (
    _ "github.com/mattn/go-sqlite3"
    "database/sql"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "home-app/app/models"
    "time"
    "github.com/solisjose46/pretty-print/debug"
)

const (
	dbPath                  = "db/home.db"
    authUserQuery           = "SELECT password FROM users WHERE username = ?"
    userIdQuery             = "SELECT userId FROM users WHERE username = ?"
    getExpenseQuery         = "SELECT e.expenseId AS ExpenseId, e.name AS Name, e.amount AS Amount, c.categoryName AS Category, u.username AS Username, e.userId AS UserId, e.createdAt AS Datetime FROM expenses e JOIN users u ON e.userId = u.userId JOIN categories c ON e.categoryId = c.categoryId WHERE e.expenseId = ?"
    addExpenseQuery         = "INSERT INTO expenses (userId, name, amount, category) VALUES (?, ?, ?, ?)"
    updateExpenseQuery      = "UPDATE expenses SET name = ?, amount = ?, category = ? WHERE expenseId = ?"
    monthlyExpensesQuery    = "SELECT e.expenseId AS ExpenseId, e.name AS Name, e.amount AS Amount, c.categoryName AS Category, u.username AS Username, e.userId AS UserId, e.createdAt AS Datetime FROM expenses e JOIN users u ON e.userId = u.userId JOIN categories c ON e.categoryId = c.categoryId WHERE e.createdAt BETWEEN ? AND ? ORDER BY e.createdAt DESC"
    monthlyCategoriesQuery  = "SELECT c.categoryName, SUM(e.amount) as balance, c.categoryLimit FROM categories c JOIN expenses e ON c.categoryId = e.categoryId WHERE e.createdAt >= ? AND e.createdAt < ? GROUP BY c.categoryName, c.categoryLimit"
)

var dao *sql.DB

func InitDB() error {
	var err error
    dao, err = sql.Open("sqlite3", dbPath)

    debug.PrintInfo(InitDB, "Connecting to db")

	if err != nil {
        debug.PrintError(InitDB, err)
		return err
	}

	return nil
}

func CloseDB() {
    if dao != nil {
        dao.Close()
    }
}

func ValidateUser(username, password string) (bool, error) {
    debug.PrintInfo(ValidateUser, "Attempting to auth user", username)

    var hashedPassword []byte
    err := dao.QueryRow(authUserQuery, username).Scan(&hashedPassword)

    if err != nil {
        debug.PrintError(ValidateUser, err)
        return false, err
    }

    err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
    
    if err == bcrypt.ErrMismatchedHashAndPassword {
        debug.PrintError(ValidateUser, err)
        return false, nil
    }
    
    if err != nil {
        debug.PrintError(ValidateUser, err)
        return false, err
    }

    debug.PrintSucc(ValidateUser, "Succ Auth", username)
    return true, nil
}

func GetUserId(username string) (string, error) {
    debug.PrintInfo(GetUserId, "getting user id")
    var userId string

    err := dao.QueryRow(userIdQuery, username).Scan(&userId)

    if err != nil {
        debug.PrintError(GetUserId, err)
        return "", err
    }
    
    debug.PrintSucc(GetUserId, "got user id")
    return userId, nil
}

func AddExpense(expense models.Expense) (bool, error) {
    debug.PrintInfo(AddExpense, "attempting to add expense")

    _, err := dao.Exec(
        addExpenseQuery, 
        expense.UserId,
        expense.Name,
        expense.Amount,
        expense.Category,
    )

    if err != nil {
        debug.PrintError(AddExpense, err)
        return false, err
    }

    debug.PrintSucc(AddExpense, "expense added!")
    return true, nil
}

func UpdateExpense(expense models.Expense) (bool, error) {
    debug.PrintInfo(UpdateExpense, "Attempting to update expense. expense:", expense.Name)

    _, err := dao.Exec(
        updateExpenseQuery, 
        expense.Name,
        expense.Amount,
        expense.Category,
        expense.ExpenseId,
    )
    
    if err != nil {
        debug.PrintError(UpdateExpense, err)
        return false, err
    }

    debug.PrintSucc(UpdateExpense, "expense added!")
    return true, nil
}

func GetExpensesForCurrentMonth(userId string) ([]models.Expense, error) {
    debug.PrintInfo(GetExpensesForCurrentMonth, "attempting to get monthly expense")

    now := time.Now()
    firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

    nextMonth := firstOfMonth.AddDate(0, 1, 0)

    startOfMonthStr := firstOfMonth.Format("2006-01-02 15:04:05")
    nextMonthStr := nextMonth.Format("2006-01-02 15:04:05")

    debug.PrintInfo(GetExpensesForCurrentMonth, "from ", startOfMonthStr, " to ", nextMonthStr)

    rows, err := dao.Query(
        monthlyExpensesQuery,
        startOfMonthStr,
        nextMonthStr,
    )

    if err != nil {
        debug.PrintError(GetExpensesForCurrentMonth, err)
        return nil, err
    }
    defer rows.Close()

    var expenses []models.Expense
    for rows.Next() {
        var expense models.Expense
        err := rows.Scan(
            &expense.ExpenseId,
            &expense.Name,
            &expense.Amount,
            &expense.Category,
            &expense.Username,
            &expense.UserId,
            &expense.Datetime,
        )

        if err != nil {
            debug.PrintError(GetExpensesForCurrentMonth, err)
            return nil, err
        }

        // this seems messy
        if expense.ExpenseId == userId {
            expense.IsOwner = true
        } else {
            expense.IsOwner = false
        }
    }

    if err = rows.Err(); err != nil {
        debug.PrintError(GetExpensesForCurrentMonth, err)
        return nil, err
    }

    debug.PrintSucc(GetExpensesForCurrentMonth, "got monthly expenses")
    return expenses, nil
}


func GetExpense(expenseId string) (models.Expense, error) {
    debug.PrintInfo(GetExpense, "getting expense w/ id: ", expenseId)
    var expense models.Expense

    err := dao.QueryRow(getExpenseQuery, expenseId).Scan(
        &expense.ExpenseId, 
        &expense.Name, 
        &expense.Amount, 
        &expense.Category, 
        &expense.Username,
        &expense.UserId, 
        &expense.Datetime,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            debug.PrintError(GetExpense, err)
            return expense, errors.New("expense not found")
        }

        debug.PrintError(GetExpense, err)
        return expense, err
    }

    debug.PrintSucc(GetExpense, "expense found!")
    return expense, nil
}

func GetCategoriesForCurrentMonth() ([]models.Category, error) {
    debug.PrintInfo(GetCategoriesForCurrentMonth, "Getting categories")
    now := time.Now()
    firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
    nextMonth := firstOfMonth.AddDate(0, 1, 0)

    startOfMonthStr := firstOfMonth.Format("2006-01-02 15:04:05")
    nextMonthStr := nextMonth.Format("2006-01-02 15:04:05")

    rows, err := dao.Query(monthlyCategoriesQuery, startOfMonthStr, nextMonthStr)
    if err != nil {
        debug.PrintError(GetCategoriesForCurrentMonth, err)
        return nil, err
    }
    defer rows.Close()

    var categories []models.Category
    for rows.Next() {
        var category models.Category
        err := rows.Scan(&category.Name, &category.Balance, &category.Limit)
        if err != nil {
            debug.PrintError(GetCategoriesForCurrentMonth, err)
            return nil, err
        }
        categories = append(categories, category)
    }

    if err = rows.Err(); err != nil {
        debug.PrintError(GetCategoriesForCurrentMonth, err)
        return nil, err
    }

    debug.PrintSucc(GetCategoriesForCurrentMonth, "Got categories!")
    return categories, nil
}