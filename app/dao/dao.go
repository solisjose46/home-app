package dao

import (
    _ "github.com/mattn/go-sqlite3"
    "database/sql"
    "errors"
    "github.com/solisjose46/pretty-print/debug"
    "golang.org/x/crypto/bcrypt"
    "home-app/app/models"
    "time"
)

type Dao struct{
    conn *sql.DB
}

func (dao *Dao) InitDB() error {
    debug.PrintInfo(dao.InitDB, "Initializing db")
	var err error
    dao.conn, err = sql.Open("sqlite3", dbPath)

    debug.PrintInfo(dao.InitDB, "Connecting to db")

	if err != nil {
        debug.PrintError(dao.InitDB, err)
		return err
	}

	return nil
}

func (dao *Dao) CloseDB() {
    if dao.conn != nil {
        dao.conn.Close()
    }
}

func (dao *Dao) ValidateUser(username, password string) (bool, error) {
    debug.PrintInfo(dao.ValidateUser, "Attempting to auth user", username)

    var hashedPassword []byte
    err := dao.conn.QueryRow(authUserQuery, username).Scan(&hashedPassword)


    if err == sql.ErrNoRows{
        debug.PrintInfo(dao.ValidateUser, "User not found")
        return false, nil
    }

    if err != nil {
        debug.PrintError(dao.ValidateUser, err)
        return false, errors.New("Error querying user")
    }

    err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
    
    if err == bcrypt.ErrMismatchedHashAndPassword {
        debug.PrintInfo(dao.ValidateUser, "Invalid password for user")
        return false, nil
    }
    
    if err != nil {
        debug.PrintError(dao.ValidateUser, err)
        return false, err
    }

    debug.PrintSucc(dao.ValidateUser, "Succ Auth", username)
    return true, nil
}

func (dao *Dao) GetUserId(username string) (string, error) {
    debug.PrintInfo(dao.GetUserId, "getting user id")
    var userId string

    err := dao.conn.QueryRow(userIdQuery, username).Scan(&userId)

    if err != nil {
        debug.PrintError(dao.GetUserId, err)
        return "", err
    }
    
    debug.PrintSucc(dao.GetUserId, "got user id")
    return userId, nil
}

func (dao *Dao) AddExpense(expense *models.Expense) (bool, error) {
    debug.PrintInfo(dao.AddExpense, "attempting to add expense")

    _, err := dao.conn.Exec(
        addExpenseQuery, 
        expense.UserId,
        expense.Name,
        expense.Amount,
        expense.Category,
    )

    if err != nil {
        debug.PrintError(dao.AddExpense, err)
        return false, err
    }

    debug.PrintSucc(dao.AddExpense, "expense added!")
    return true, nil
}

func (dao *Dao) UpdateExpense(expense *models.Expense) (bool, error) {
    debug.PrintInfo(dao.UpdateExpense, "Attempting to update expense. expense:", expense.Name)

    debug.PrintInfo(dao.UpdateExpense, "CATEGORY:", expense.Category)

    _, err := dao.conn.Exec(
        updateExpenseQuery,
        expense.Category,
        expense.Name,
        expense.Amount,
        expense.ExpenseId,
    )
    
    if err != nil {
        debug.PrintError(dao.UpdateExpense, err)
        return false, err
    }

    debug.PrintSucc(dao.UpdateExpense, "expense added!")
    return true, nil
}

func (dao *Dao) GetExpensesForCurrentMonth(userId string) (*[]models.Expense, error) {
    debug.PrintInfo(dao.GetExpensesForCurrentMonth, "attempting to get monthly expense")

    now := time.Now()
    firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

    nextMonth := firstOfMonth.AddDate(0, 1, 0)

    startOfMonthStr := firstOfMonth.Format("2006-01-02 15:04:05")
    nextMonthStr := nextMonth.Format("2006-01-02 15:04:05")

    debug.PrintInfo(dao.GetExpensesForCurrentMonth, "from ", startOfMonthStr, " to ", nextMonthStr)

    rows, err := dao.conn.Query(
        monthlyExpensesQuery,
        startOfMonthStr,
        nextMonthStr,
    )

    if err != nil {
        debug.PrintError(dao.GetExpensesForCurrentMonth, err)
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
            debug.PrintError(dao.GetExpensesForCurrentMonth, err)
            return nil, err
        }

        expense.IsOwner = expense.UserId == userId

        expenses = append(expenses, expense)
        debug.PrintInfo(dao.GetExpensesForCurrentMonth, "expense:", expense.Name)
    }

    if err = rows.Err(); err != nil {
        debug.PrintError(dao.GetExpensesForCurrentMonth, err)
        return nil, err
    }

    debug.PrintSucc(dao.GetExpensesForCurrentMonth, "got monthly expenses")
    return &expenses, nil
}

func (dao *Dao) GetExpense(expenseId string) (*models.Expense, error) {
    debug.PrintInfo(dao.GetExpense, "getting expense w/ id: ", expenseId)
    var expense models.Expense

    err := dao.conn.QueryRow(getExpenseQuery, expenseId).Scan(
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
            debug.PrintError(dao.GetExpense, err)
            return nil, errors.New("expense not found")
        }

        debug.PrintError(dao.GetExpense, err)
        return nil, err
    }

    debug.PrintSucc(dao.GetExpense, "expense found!")
    return &expense, nil
}

func (dao *Dao) GetCategoriesForCurrentMonth() (*[]models.Category, error) {
    debug.PrintInfo(dao.GetCategoriesForCurrentMonth, "Getting categories")
    now := time.Now()
    firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
    nextMonth := firstOfMonth.AddDate(0, 1, 0)

    startOfMonthStr := firstOfMonth.Format("2006-01-02 15:04:05")
    nextMonthStr := nextMonth.Format("2006-01-02 15:04:05")

    rows, err := dao.conn.Query(monthlyCategoriesQuery, startOfMonthStr, nextMonthStr)
    if err != nil {
        debug.PrintError(dao.GetCategoriesForCurrentMonth, err)
        return nil, err
    }
    defer rows.Close()

    var categories []models.Category
    for rows.Next() {
        var category models.Category
        err := rows.Scan(&category.Name, &category.Balance, &category.Limit)
        if err != nil {
            debug.PrintError(dao.GetCategoriesForCurrentMonth, err)
            return nil, err
        }
        categories = append(categories, category)
    }

    if err = rows.Err(); err != nil {
        debug.PrintError(dao.GetCategoriesForCurrentMonth, err)
        return nil, err
    }

    debug.PrintSucc(dao.GetCategoriesForCurrentMonth, "Got categories!")
    return &categories, nil
}