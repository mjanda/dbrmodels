# dbrmodels
Create projects for MySql database and generate table model structs for github.com/gocraft/dbr
## Getting started
Lets create MySql table
```sql
CREATE TABLE Persons
(
    PersonID    int,
    LastName    varchar(255),
    FirstName   varchar(255),
    Address     varchar(255),
    City        varchar(255)
);
```
and generate gocraft/dbr model
```go
package dbmodels
type Persons struct {
    PersonID    int64           `db:"PersonID"`
    LastName    string          `db:"LastName"`
    FirstName   string          `db:"FirstName"`
    Address     dbr.NullString  `db:"Address"`
    City        dbr.NullString  `db:"City"`
}
```
## MySql types
| MySql | GO | NULL |
|  -------------  | ------------- | ------------- |
| tinyint(1) | bool | dbr.NullBool |
| int | int64 | dbr.NullInt64 |
| float | float64 | dbr.NullFloat64 |
| * | string | dbr.NullString |

# install
```bash
go get github.com/finalist736/dbrmodels
```
## gocraft/dbr example
```go
// Get a record
var persons dbmodels.Persons
err := dbrSess.Select("*").From("persons").Where("PersonID = ?", 13).LoadStruct(&persons)
```
# Projects

Project including next data:
* Name
* DB Host
* DB Port
* DB User
* DB Password
* DB Name
* Path where to .go files located

# Using
start generate project
```bash
dbrmodels project_name
```
* list projects
```bash
dbrmodels ls
```
* create project
```bash
dbrmodels create
```
* edit project
```bash
dbrmodels edit project_name
```
* remove project
```bash
dbrmodels remove project_name
```
