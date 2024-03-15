
WINDOWS : 

- **Install GOLANG :**

  -    if `go version` fails
  -   go to your env variable settings and set GOROOT to point to the directory where the wizard installed go lang (ex : set GOROOT=C:\Go)


- **Download Migrate :** 
  -    WINDOWS : Follow the steps here https://www.geeksforgeeks.org/how-to-install-golang-migrate-on-windows/
  -    MACOS   : run `brew install golang-migrate`


- **Create Migration setup :** 
  - run `mkdir db/migration` (in project folder)
  - run `migrate create -ext sql -dir db/migration -seq init_schema`


- **Database Design : https://dbdiagram.io**
    - paste exported code in up/down files accordingly 
    - note : down file must be Script that drops all tables - up file will create the tables that get generated


- **Open Terminal Ported to the PostGreSQL container :** 
  - (in vscode bash terminal) : `docker exec -it simple-bank bin/sh`
  - (in mac OS      terminal) : `docker exec -it simple-bank /bin/sh`


- **Create Database in PGSql container :** 
   - run `createdb --username=root --owner=root simple_bank`
   - run `psql simple_bank`

Create a make file to make the setup process easy between team members
    Check Makefile in project Dir
- run `docker run --name simple-bank -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres` to RUN database image using docker
- run `make dropdb` to DROP all tables
- run `make createdb` to CREATE the database 

Migrate our SQL schema to the RUNNING Database (container)
- run `migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up`

To go back to previous DB version : 
- run `migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down`

- for postgres : https://docs.sqlc.dev/en/latest/index.html || https://docs.sqlc.dev/en/latest/overview/install.html (USE IT THE DOCKER WAY ON WINDOWS  ELSE USE BREW ON MAC)
  - run `docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc init` USING CMD on WINDOWS to generate the YAML file


- For database schema changes **DURING PRODUCTION:** 
  - `migrate create -ext sql -dir db/migration -seq add_users`
  - paste the NEW lines added to the schema in there and NOT the entire schema generated


for other stuff : https://pkg.go.dev/github.com/jmoiron/sqlx

NOTE : ALL GO ENV'S USE A / AND NOT A \ AS PATH SEPERATOR
like GOBIN