package dao

const (
	dbPath                  = "db/home.db"
    authUserQuery           = "SELECT password FROM users WHERE username = ?"
    userIdQuery             = "SELECT userId FROM users WHERE username = ?"
    getExpenseQuery         = "SELECT e.expenseId AS ExpenseId, e.name AS Name, e.amount AS Amount, c.categoryName AS Category, u.username AS Username, e.userId AS UserId, e.createdAt AS Datetime FROM expenses e JOIN users u ON e.userId = u.userId JOIN categories c ON e.categoryId = c.categoryId WHERE e.expenseId = ?"
    addExpenseQuery         = "INSERT INTO expenses (userId, name, amount, categoryId, createdAt) SELECT ?, ?, ?, categoryId, CURRENT_TIMESTAMP FROM categories WHERE categoryName = ?"
    updateExpenseQuery      = "WITH CategoryID AS (SELECT categoryId FROM categories WHERE categoryName = ?) UPDATE expenses SET name = ?, amount = ?, categoryId = (SELECT categoryId FROM CategoryID) WHERE expenseId = ?"
    monthlyExpensesQuery    = "SELECT e.expenseId AS ExpenseId, e.name AS Name, e.amount AS Amount, c.categoryName AS Category, u.username AS Username, e.userId AS UserId, e.createdAt AS Datetime FROM expenses e JOIN users u ON e.userId = u.userId JOIN categories c ON e.categoryId = c.categoryId WHERE e.createdAt BETWEEN ? AND ? ORDER BY e.createdAt DESC"
    monthlyCategoriesQuery  = "SELECT c.categoryName, SUM(e.amount) as balance, c.categoryLimit FROM categories c JOIN expenses e ON c.categoryId = e.categoryId WHERE e.createdAt >= ? AND e.createdAt < ? GROUP BY c.categoryName, c.categoryLimit"
)