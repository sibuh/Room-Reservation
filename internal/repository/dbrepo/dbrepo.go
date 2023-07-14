package dbrepo

import (
	"booking/internal/pkg/config"
	"booking/internal/pkg/models"
	"booking/internal/repository"
	"booking/platform/pass"
	"booking/platform/token"
	"context"
	"database/sql"
	"errors"
	"time"
)

type postgresDbRepo struct {
	app *config.AppConfig
	DB  *sql.DB
}

func NewPostgresDbRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDbRepo{
		app: a,
		DB:  conn,
	}
}
func (p *postgresDbRepo) MakeReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stmt := `insert into reservations(first_name,last_name,email,phone,start_date,end_date,room_id,
		created_at,updated_at)
		values($1,$2,$3,$4,$5,$6,$7,$8,$9)returning id`
	var id int
	err := p.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
func (p *postgresDbRepo) InsertRoomRestriction(restrict models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stmt := `insert into room_restrictions(room_id,restriction_id,reservation_id,start_date,end_date,created_at,updated_at)
			values($1,$2,$3,$4,$5,$6,$7)`
	_, err := p.DB.ExecContext(ctx, stmt,
		restrict.RoomID,
		restrict.RestrictionID,
		restrict.ReservationID,
		restrict.StartDate,
		restrict.EndDate,
		time.Now(),
		time.Now())
	if err != nil {
		return err
	}
	return nil
}
func (p *postgresDbRepo) SearchAvailabilityByRoomID(roomID int, startDate, endDate time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var numRooms int
	stmt := `select
				count(id)
			from 
				room_restrictions
			where 
				room_id=$1
			and $2 > start_date 
			and $3 < end_date
				`
	row := p.DB.QueryRowContext(ctx, stmt, roomID, endDate, startDate)
	err := row.Scan(&numRooms)
	if err != nil {
		return false, err
	}
	if numRooms == 0 {
		return false, nil
	}
	return true, nil
}
func (p *postgresDbRepo) SearchAvailableRooms(startDate, endDate time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stmt := `select
				r.id,r.room_name 
			from
				 rooms r
			where 
				r.id not in (select rr.room_id from room_restrictions rr where $1 > start_date or $2 < end_date)
				`
	rows, err := p.DB.QueryContext(ctx, stmt, endDate, startDate)
	if err != nil {
		return nil, err
	}
	var rooms []models.Room
	for rows.Next() {
		var r models.Room
		rows.Scan(&r.ID, &r.RoomName)
		rooms = append(rooms, r)
	}

	return rooms, nil
}
func (p *postgresDbRepo) InsertRooms(arg models.AddRoomRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stmt := `insert into rooms (room_number,created_at,updated_at,room_type) values($1,$2,$3,$4)`
	_, err := p.DB.ExecContext(ctx, stmt,
		arg.RoomNumber,
		time.Now(),
		time.Now(),
		arg.RoomType)
	if err != nil {
		return err
	}
	return nil
}
func (p *postgresDbRepo) Login(arg models.LoginRequest) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stmt := `select * from users where email=$1`
	var user models.User
	row := p.DB.QueryRowContext(ctx, stmt, arg.Email)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
	)
	if err != nil {
		return "", err
	}
	match := pass.CheckPasswordHash(arg.Password, user.PasswordHash)
	if match {
		jwt := token.NewJwtMaker("1234567890123456789012")
		payload := token.Payload{
			UserName:  user.FirstName,
			CreatedAt: time.Now(),
			Duration:  time.Hour * 24,
		}
		tokenString, err := jwt.CreateToken(payload)
		if err != nil {
			return "", err
		}
		return tokenString, nil
	}
	return "", errors.New("password didnot match")
}
