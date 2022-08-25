package dbrepo

import (
	"context"
	"log"
	"time"

	"gitlab.com/gjerry134679/vigilate/pkg/models"
)

func (m *postgresDBRepo) InsertHost(h models.Host) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	   INSERT INTO hosts (host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at)
	   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id;
	`
	var newID int
	err := m.DB.QueryRowContext(
		ctx, query,
		h.HostName, h.CanonicalName, h.URL, h.IP, h.IPv6,
		h.Location, h.OS, h.Active, time.Now(), time.Now(),
	).Scan(&newID)

	if err != nil {
		log.Println(err)
		return newID, err
	}
	return newID, nil
}

func (m *postgresDBRepo) GetHostByID(id int) (models.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
	   SELECT id, host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at
	     FROM hosts
	    WHERE id = $1;
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	h := models.Host{}
	err := row.Scan(
		&h.ID, &h.HostName, &h.CanonicalName,
		&h.URL, &h.IP, &h.IPv6, &h.Location,
		&h.OS, &h.Active, &h.CreatedAt, &h.UpdatedAt,
	)
	if err != nil {
		return h, err
	}
	return h, nil
}

func (m *postgresDBRepo) GetAllHost() ([]models.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT * FROM hosts
		ORDER BY host_name ASC;
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var hosts []models.Host
	for rows.Next() {
		var h models.Host
		err = rows.Scan(
			&h.ID, &h.HostName, &h.CanonicalName,
			&h.URL, &h.IP, &h.IPv6, &h.Location,
			&h.OS, &h.Active, &h.CreatedAt, &h.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		hosts = append(hosts, h)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return hosts, nil
}

func (m *postgresDBRepo) UpdateHost(h models.Host) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	   UPDATE hosts
	      SET host_name = $1, 
		      canonical_name = $2,
			  url = $3, 
			  ip = $4,
			  ipv6 = $5,
			  location = $6,
			  os = $7,
			  active = $8,
			  updated_at = $9
		WHERE id = $10
	`
	_, err := m.DB.ExecContext(ctx, stmt,
		h.HostName, h.CanonicalName, h.URL, h.IP, h.IPv6,
		h.Location, h.OS, h.Active, time.Now(), h.ID,
	)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
