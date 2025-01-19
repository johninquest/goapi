// internal/routes/policy_routes.go
package routes

import (
    "goapi/internal/handlers"
    "github.com/gofiber/fiber/v2"
    "database/sql"
)

func PolicyRoutes(app *fiber.App, db *sql.DB) {
    policyHandler := handlers.NewPolicyHandler(db)
    
    // Policy route group
    policy := app.Group("/api/policy")
    
    // CRUD routes
    policy.Post("/", policyHandler.CreatePolicy)
    policy.Get("/:id", policyHandler.GetPolicy)
    policy.Put("/:id", policyHandler.UpdatePolicy)
    policy.Delete("/:id", policyHandler.DeletePolicy)
}

// Database initialization (add to your main.go or a separate db package)
func InitDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "policies.db")
    if err != nil {
        return nil, err
    }

    // Create policies table
    createTable := `
    CREATE TABLE IF NOT EXISTS policies (
        id TEXT PRIMARY KEY,
        policy_type TEXT NOT NULL,
        policy_number TEXT NOT NULL,
        insurance_provider TEXT NOT NULL,
        policy_comment TEXT,
        start_date DATETIME NOT NULL,
        end_date DATETIME NOT NULL,
        automatic_renewal BOOLEAN NOT NULL,
        created_by TEXT NOT NULL,
        premium REAL NOT NULL,
        payment_frequency INTEGER NOT NULL,
        created DATETIME NOT NULL,
        updated DATETIME NOT NULL
    );`

    _, err = db.Exec(createTable)
    if err != nil {
        return nil, err
    }

    return db, nil
}