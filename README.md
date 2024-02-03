# Dcard 2024 Intern Backend Homework

## Prerequisite
* Go
* PostgreSQL

## Quick Start
Run the service:
```bash
$ go run main.go
```

linting
```bash
$ go vet ./..
```

test
```bash
$ go test ./...
ok      dcard2024       0.099s
?       dcard2024/internal/get_ads      [no test files]
?       dcard2024/internal/post_ads     [no test files]
```
## Usage
Post new advertisement
```
curl -X POST -H "Content-Type: application/json" \
  "http://localhost:8080/api/v1/ad" \
  --data '{
     "title": "AD 55",
     "startAt": "2023-12-10T03:00:00.000Z",
     "endAt": "2023-12-31T16:00:00.000Z", 
     "conditions": {
        "ageStart": 20,
        "country": ["TW", "JP"],
        "platform": ["android", "ios"]
     }
  }'
```

Get advertisements

* Query:
  ```
  curl -X GET -H "Content-Type: application/json" \
  "http://localhost:8080/api/v1/ad?offset=10&limit=2&age=24&gender=F&country=TW&platform=ios"
  ```
* Response:
  ```json
  {"items":[{"title": "AD 1","endAt" "2023-12-22T01:00:00.000Z"},{"title": "AD 31","endAt" "2023-12-30T12:00:00.000Z"}]}
  ```

## PostgreSQL Common Commands
Open PostgreSQL
```bash
$ sudo -u postgres psql
```

Add user to PostgreSQL
```
$ sudo -i -u postgres
~$ createuser --interactive
~$ Enter name of role to add: <username>
$ Shall the new role be a superuser? (y/n) y
postgres@<host>:~$ exit
```

> In this project I remove password requirement in my PostgreSQL.
> Run the following to get the path of `pg_hba.conf`.
> ```
> $ psql
> # SHOW hba_file;
> ```
> Then set authentication to `trust` in `pg_hba.conf`.  
> Ref: https://dba.stackexchange.com/questions/83164/postgresql-remove-password-requirement-for-user-postgres

Postgre CLI
```
$ psql
\c ads;    # use database ad
\dt;       # show all the table in the database
\conninfo; # get connection information
```
## SQL Common Commands
Create table
```
CREATE TABLE ad ( 
  title text,
  startAt timestamp,
  endAt timestamp,
  ageStart int,
  ageEnd int,
  gender text[],
  country text[],
  platform text[]
);
```

Insert row
```
INSERT INTO ad (title, start_at, end_at, age_start, country)
VALUES ('AD 55', '2023-12-10T03:00:00.000Z', '2024-12-31T16:00:00.000Z', 20, ARRAY ['TW']);
```

Delete every row in table
```
DELETE FROM ad;
```

Drop table
```
DROP TABLE IF EXISTS ad;
```


