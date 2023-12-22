package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User) error
	GetUsers() ([]*User, error)
	GetUserByEmail(string) (*User, error)
	GetUserById(int) (*User, error)
	CreateTask(*Task) error
	UpdateTask(*Task) error
	DeleteTask(int) error
	GetTasks() ([]*Task, error)
	GetTaskById(int) (*Task, error)
	//GetTaskByStatus(Status) ([]*Task, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	_, err := s.createUserTable(), s.createTaskTable()
	return err
}

func (s *PostgresStore) createTaskTable() error {
	query := `
		create table if not exists task (
		id serial primary key,
		title varchar(50),
		description varchar(155),
		status int
		)
	`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createUserTable() error {
	query :=
		`create table if not exists users (
			id serial primary key,
			first_name varchar(100),
			last_name varchar(100),
			email varchar(50),
			encrypted_password varchar(100),
			created_at timestamp
		)
	`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateUser(user *User) error {
	query := `insert into users 
	(first_name, last_name, email, encrypted_password, created_at)
	values ($1, $2, $3, $4, $5)`

	_, err := s.db.Query(
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.EncryptedPassword,
		user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CreateTask(task *Task) error {
	query := `insert into "task" 
	(title, description)
	values ($1, $2)`

	_, err := s.db.Query(
		query,
		task.Title,
		task.Description,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateTask(*Task) error {
	return nil
}

func (s *PostgresStore) DeleteTask(id int) error {
	_, err := s.db.Query("delete from task where id = $1", id)
	return err
}

func (s *PostgresStore) GetTaskById(id int) (*Task, error) {
	rows, err := s.db.Query("select * from task where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoTask(rows)
	}

	return nil, fmt.Errorf("task %d not found", id)
}

func (s *PostgresStore) GetUserByEmail(email string) (*User, error) {
	rows, err := s.db.Query("select * from users where email = $1", email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user %s not found with email ", email)
}

func (s *PostgresStore) GetUserById(id int) (*User, error) {
	rows, err := s.db.Query("select * from users where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user %d not found with id ", id)
}

func (s *PostgresStore) GetTasks() ([]*Task, error) {
	rows, err := s.db.Query("select * from task")
	if err != nil {
		return nil, err
	}

	tasks := []*Task{}
	for rows.Next() {
		task, err := scanIntoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *PostgresStore) GetUsers() ([]*User, error) {
	rows, err := s.db.Query("select * from users")
	if err != nil {
		return nil, err
	}

	users := []*User{}
	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func scanIntoTask(rows *sql.Rows) (*Task, error) {
	task := new(Task)
	err := rows.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
	)
	return task, err
}

func scanIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.EncryptedPassword,
		&user.CreatedAt,
	)
	return user, err

}
