# Dcard 2024 Intern Backend Homework

## Prerequisite
* Go
* PostgreSQL


## Commands
Run the service:
```
go run main.go
```

linting
```
go vet ./..
```

postgreSQL
```
# sudo -u postgres psql
sudo -i -u postgres
~$ createuser --interactive
~$ Enter name of role to add: daniel
$ Shall the new role be a superuser? (y/n) y
postgres@b09902053-host:~$ exit


psql
\c ads; # use database ad
\dt; # show all the table in the database
\conninfo; # get connection information

select

# create table
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

# insert rows
INSERT INTO ad (title, start_at, end_at, age_start, country)
VALUES ('AD 55', '2023-12-10T03:00:00.000Z', '2024-12-31T16:00:00.000Z', 20, ARRAY ['TW']);

# delete every row in table
DELETE FROM tourneys;

# drop table
DROP TABLE IF EXISTS tourneys;



remove postgreSQL password

```
