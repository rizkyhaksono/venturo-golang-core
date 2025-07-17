## Project Capstone: The Inventory & Transaction Module

### **Overall Goal**
Build a complete inventory management and transaction system on top of our existing backend core. This will validate a user's ability to work with the entire architecture.

### **Phase 0: Prerequisite - Re-implement Authentication**
Before starting, the trainee should ensure the full authentication system, including **Access and Refresh Tokens backed by a dedicated database table**, is correctly implemented. They will need this for the protected routes in this module.

---
### **Phase 1: Foundation - Database & Models**
The first step is to define the new data structures.

1.  **Create New Database Migrations:**
    * Generate a migration to create a new `outlets` table. It should have at least an `id` and a `name`.
    * Generate another migration to create the `inventory_ledgers` table. This is the most important table. It must have columns for `id`, `item_id` (foreign key), `outlet_id` (foreign key), `transaction_id` (foreign key, nullable), `quantity_change` (an integer that can be positive for stock-in or negative for stock-out), and timestamps.
    * Modify the existing `transactions` table migration to add a required `outlet_id` foreign key.

2.  **Create the Go Models:**
    * Create a new `internal/model/outlet_model.go` file with the `Outlet` struct. Define its relationship to `Transaction` (an outlet has many transactions).
    * Create `internal/model/inventory_ledger_model.go` with the `InventoryLedger` struct and its relationships.
    * Update the existing `Transaction` model to include the `OutletID` field and its `Outlet` relationship struct.

---
### **Phase 2: Core Logic - Stock Management**
Before we can sell items, we must have them in stock.

1.  **Create an `InventoryService`:** Create a new service file for inventory logic.
2.  **Implement a "Stock In" Method:** Inside the service, create a `StockIn` method. This method will receive an `item_id`, `outlet_id`, and `quantity`. Its job is to create a new record in the `inventory_ledgers` table with a **positive** `quantity_change`.
3.  **Create the "Stock In" Route:**
    * Create a new `InventoryHandler`.
    * Define a new protected route, `POST /api/v1/inventory/stock-in`. This route should require authentication.
    * The handler should validate the incoming payload and call the `StockIn` service method.

---
### **Phase 3: The Main Feature - Transaction Flow**
Now, integrate inventory logic into the transaction process.

1.  **Modify the "Create Transaction" Service:**
    * Update the existing `CreateTransaction` service method to accept the new required `outlet_id`.

2.  **Implement Stock Validation Logic:** This is the most critical business rule.
    * Inside the `CreateTransaction` service method, **before** saving the transaction, you must check if there is enough stock.
    * For each item in the transaction, perform a query on the `inventory_ledgers` table to `SUM` the `quantity_change` for that specific `item_id` and `outlet_id`.
    * Compare the result (the current "on-hand" quantity) with the quantity being requested in the transaction.
    * If the requested quantity is greater than the on-hand quantity, the function must return an error and the transaction must fail.

---
### **Phase 4: Post-Transaction Logic**
This covers what happens after a transaction is successfully created and paid for.

1.  **Create Ledger Records:**
    * After a transaction is successfully saved, the `CreateTransaction` service must create new records in the `inventory_ledgers` table for each item sold.
    * The `quantity_change` for these records must be **negative** (e.g., if 2 items were sold, the quantity change is -2).
    * **Recommendation:** This is a perfect opportunity to reinforce asynchronous processing. The service should handle this logic inside a **goroutine** so the user doesn't have to wait for it. The `WaitGroup` we implemented should be used to track this background task.

---
### **Phase 5: Reporting**
The final step is to aggregate all the data into a useful report.

1.  **Create a `ReportService`:** Create a new service for generating reports.
2.  **Implement the Report Generation Method:**
    * This method needs to perform a complex query to generate the data for the report.
    * It must get the final "on-hand quantity" for each item at each outlet. This requires `GROUPING BY` item and outlet and using `SUM(quantity_change)`.
    * It must also fetch the individual transaction history for that item/outlet combination. This involves joining the `inventory_ledgers`, `transactions`, and `items` tables.
3.  **Create the Report Route:**
    * Create a new `ReportHandler`.
    * Define a new protected route, `GET /api/v1/reports/inventory`, that can be filtered by `item_id` and/or `outlet_id`.
    * The handler will call the `ReportService` and return the formatted report data.