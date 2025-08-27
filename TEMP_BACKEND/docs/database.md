# Database ER diagram

This ER diagram is generated from the `User` model in `main.go` and shows the database table and columns as created by GORM (SQLite).

```mermaid
erDiagram
    USERS {
        int ID PK "primary key"
        string EMAIL "unique, not null"
        string PASSWORD
        string FIRST_NAME
        string LAST_NAME
        string PHONE
        date BIRTHDAY
        int POINTS
        datetime CREATED_AT
    }

    %% No additional tables in this simple example. Keep relationships here when you add more models.
```

Notes

- GORM default pluralizes struct names; `User` becomes `users` table.
- The `email` field has a unique index (`gorm:"uniqueIndex"`).
- `temp_backend.db` is created in the project root and GORM auto-migrates the `users` table on startup.
