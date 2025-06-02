package db

import (
	"database/sql"
	"fmt"
	"github.com/ademaxweb/mfa-go-core/pkg/data"
	_ "github.com/lib/pq"
	"log"
	"reflect"
	"strings"
)

type PgsqlDB struct {
	conn *sql.DB
}

func NewPgsqlDB(connStr string) (*PgsqlDB, error) {
	instance := &PgsqlDB{}

	err := instance.Connect(connStr)
	if err != nil {
		return nil, err
	}

	err = instance.CreateUsersTable()
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (p *PgsqlDB) Connect(connStr string) error {
	var err error
	if p.conn != nil {
		_ = p.conn.Close()
	}
	p.conn, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = p.conn.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (p *PgsqlDB) CreateUser(d data.User) (int, error) {
	query := `insert into users (name, email, password) values ($1, $2, $3) returning id`
	row := p.conn.QueryRow(query, d.Name, d.Email, d.Password)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *PgsqlDB) DeleteUser(id int) error {
	query := `delete from users where id = $1`
	_, err := p.conn.Exec(query, id)
	return err
}

func (p *PgsqlDB) UpdateUser(id int, d data.User) error {
	query := `update users set `
	var args []interface{}
	argPos := 1

	v := reflect.ValueOf(d)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			query += fmt.Sprintf("%s = $%d, ", strings.ToLower(t.Field(i).Name), argPos)
			args = append(args, field.Interface())
			argPos++
		}
	}

	query = query[:len(query)-2] + fmt.Sprintf(", updated_at = NOW() WHERE id = $%d", argPos)
	args = append(args, id)

	_, err := p.conn.Exec(query, args...)
	return err
}

func (p *PgsqlDB) GetUser(id int) (*data.User, error) {
	query := `select * from users where id = $1`
	row := p.conn.QueryRow(query, id)
	var u data.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.ModifiedAt)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, NotFound
	}
	return &u, nil
}

func (p *PgsqlDB) GetAllUsers() ([]data.User, error) {
	query := `select * from users`
	rows, err := p.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	users := make([]data.User, 0)

	for rows.Next() {
		var u data.User
		err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.ModifiedAt)
		users = append(users, u)
	}

	return users, nil
}

func (p *PgsqlDB) GetUserByEmail(email string) (*data.User, error) {
	query := `select * from users where email = $1`
	row := p.conn.QueryRow(query, email)
	var u data.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.ModifiedAt)
	if err != nil {
		return nil, NotFound
	}
	return &u, nil
}

func (p *PgsqlDB) CreateUsersTable() error {
	query := `create table if not exists users (
    	"id" serial primary key,
    	"name" varchar(255) not null,
    	"email" varchar(255) not null,
    	"password" varchar(255) not null,
    	"created_at" timestamp not null default current_timestamp,
    	"updated_at" timestamp not null default current_timestamp
    )`

	_, err := p.conn.Exec(query)
	return err
}
