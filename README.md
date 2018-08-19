You want to insert various fields of this same struct with package sql:

```Go
type Email struct {
    ID          int64   `col:"id"`
    UserName    string  `col:"user_name"`
    Subject     string
    Body        string  `col:"body"`
}
```

You may now name the fields that fits your SQL query so package sql can update your table:

```Go
var db *sql.DB

email := Email {
    Subject: "Hello",
    Body: "world",
}

db.Exec("UPDATE emails SET subject = ?, body = ? WHERE id = ?",
    col.Values(email, "Subject", "body", "id")...)
```

So now you can have a single 2-parameter Insert function for all your insertion needs:

```Go
import "github.com/Masterminds/squirrel"

func Insert(src Email, columns ...string) {
    squirrel.Insert("emails").
        Columns(columns...).
        Values(col.Values(src, columns...)...).
        RunWith(db).Exec()
}
```

And here is the update function:

```Go

func Update(src Email, columns ...string) {
    q := squirrel.Update("emails").Where("id = ?", src.ID)
    values := col.Values(src, columns...)

    for i, column := range columns {
        q = q.Set(column, values[i])
    }

    q.RunWith(db).Exec()
}
```
