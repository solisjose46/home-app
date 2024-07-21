package dao

import (
    _ "github.com/mattn/go-sqlite3"
    "database/sql"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "home-app/app/models"
    "home-app/app/util"
    "time"
)

const (
	dbPath                  = "internal/db/sqlite/home.db"
    authUserQuery           = "SELECT password FROM users WHERE username = ?"
    userIdQuery             = "SELECT userId FROM users WHERE username = ?"
    addExpenseQuery         = "INSERT INTO expenses (userId, name, amount, category) VALUES (?, ?, ?, ?)"
    updateExpenseQuery      = "UPDATE expenses SET name = ?, amount = ?, category = ? WHERE expenseId = ?"
    monthlyExpensesQuery    = "SELECT e.expenseId AS ExpenseId, e.name AS Name, e.amount AS Amount, c.categoryName AS Category, u.username AS Username, e.userId AS UserId, e.createdAt AS Datetime FROM expenses e JOIN users u ON e.userId = u.userId JOIN categories c ON e.categoryId = c.categoryId WHERE e.createdAt BETWEEN ? AND ? ORDER BY e.createdAt DESC"
    monthlyCategoriesQuery  = "SELECT c.categoryName, SUM(e.amount) as balance, c.categoryLimit FROM categories c JOIN expenses e ON c.categoryId = e.categoryId WHERE e.createdAt >= ? AND e.createdAt < ? GROUP BY c.categoryName, c.categoryLimit"
)

var dao *sql.DB

func InitDB() error {
	var err error
    dao, err = sql.Open("sqlite3", dbPath)

    util.PrintMessage("Connecting to db")

	if err != nil {
        fmt.PrintError(err)
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
    util.PrintMessage("Attempting to auth user ", username)
    
    var hashedPassword string
    err := dao.QueryRow(authUserQuery, username).Scan(&hashedPassword)

    if err != nil {
        util.PrintError(err)
        return false, err
    }

    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        util.PrintError(err)
        return false, errors.New("password does not match")
    }

    util.PrintSuccess("Succ Auth ", useranme)
    return true, nil
}

func GetUserId(username string) (string, error) {
    util.PrintMessage("getting user id")
    var userId string

    err := dao.QueryRow(userIdQuery, username).Scan(&userId)

    if err != nil {
        util.PrintError(err)
        return "", err
    }
    
    util.PrintMessage("got user id")
    return userId, nil
}

func AddExpense(expense models.Expense) (bool, error) {
    util.PrintMessage("attempting to add expense")

    _, err := dao.Exec(
        addExpenseQuery, 
        expense.UserId,
        expense.Name,
        expense.Amount,
        expense.Category
    )

    if err != nil {
        util.PrintError(err)
        return false, err
    }

    // what if no error and does not upload
    // is that valid behavior

    util.PrintSuccess("expense added!")
    return true, nil
}

func UpdateExpense(expense models.Expense) (bool, error) {
    util.PrintMessage("Attempting to update expense. expense: ", expense.Name)

    _, err = dao.Exec(
        updateExpenseQuery, 
        expense.Name,
        expense.Amount,
        expense.Category,
        expense.ExpenseId
    )
    
    if err != nil {
        return false, err
    }

    util.PrintSuccess("expense added!")
    return true, nil
}

func GetExpensesForCurrentMonth(userId string) ([]models.Expense, error) {
    util.PrintMessage("attempting to get monthly expense")

    now := time.Now()
    firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

    nextMonth := firstOfMonth.AddDate(0, 1, 0)

    startOfMonthStr := firstOfMonth.Format("2006-01-02 15:04:05")
    nextMonthStr := nextMonth.Format("2006-01-02 15:04:05")

    util.PrintMessage("from ", startOfMonthStr, " to ", nextMonthStr)

    rows, err := dao.Query(
        monthlyExpensesQuery,
        startOfMonthStr,
        nextMonthStr
    )

    if err != nil {
        util.PrintError(err)
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
            &expense.Category
            &expense.Username,
            &expense.UserId,
            &expense.Datetime
        )

        if err != nil {
            util.PrintError(err)
            return nil, err
        }

        if expense.ExpenseId == userId {
            expense.IsOwner = true
        } else {
            expense.IsOwner = false
        }
    }

    if err = rows.Err(); err != nil {
        util.PrintError(err)
        return nil, err
    }

    util.PrintSuccess("got monthly expenses")
    return expenses, nil
}


func GetExpense(expenseId string) (models.Expense, error) {
    util.PrintMessage("getting expense w/ id: ", expenseId)
    var expense models.Expense

    err := dao.QueryRow(getExpenseQuery, expenseId).Scan(
        &expense.ExpenseId, 
        &expense.UserId, 
        &expense.Name, 
        &expense.Amount, 
        &expense.Category, 
        &expense.Datetime,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            util.PrintError(err)
            return expense, errors.New("expense not found")
        }

        util.PrintError(err)
        return expense, err
    }

    util.PrintSuccess("expense found!")
    return expense, nil
}

func GetCategoriesForCurrentMonth() ([]models.Category, error) {
    now := time.Now()
    firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
    nextMonth := firstOfMonth.AddDate(0, 1, 0)

    startOfMonthStr := firstOfMonth.Format("2006-01-02 15:04:05")
    nextMonthStr := nextMonth.Format("2006-01-02 15:04:05")



    rows, err := dao.Query(monthlyCategoriesQuery, startOfMonthStr, nextMonthStr)
    if err != nil {
        util.PrintError(err)
        return nil, err
    }
    defer rows.Close()

    var categories []models.Category
    for rows.Next() {
        var category models.Category
        err := rows.Scan(&category.Name, &category.Balance, &category.Limit)
        if err != nil {
            util.PrintError(err)
            return nil, err
        }
        categories = append(categories, category)
    }

    if err = rows.Err(); err != nil {
        util.PrintError(err)
        return nil, err
    }

    util.PrintSuccess("Got categories!")
    return categories, nil
}

