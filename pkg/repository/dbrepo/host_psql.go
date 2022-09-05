package dbrepo

import (
	"context"
	"fmt"
	"log"
	"strconv"
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

	stmt := `
	INSERT INTO host_services (host_id, service_id, active, schedule_number, schedule_unit,
	status, created_at, updated_at)
	VALUES ($1, 1, 0, 3, 'm', $2, $3, $4);
	`
	_, err = m.DB.ExecContext(ctx, stmt, newID, models.ServiceStatusPending, time.Now(), time.Now())
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

	query = `
	SELECT hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number,
		   hs.schedule_unit, hs.last_check, hs.status, hs.created_at, hs.updated_at,
		   s.id, s.service_name, s.active, s.icon, s.created_at, s.updated_at
	  FROM host_services AS hs
      LEFT JOIN services AS s
	    ON (s.id = hs.service_id)
	 WHERE host_id = $1;
	`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		fmt.Println(err)
		return h, err
	}
	defer rows.Close()

	var HostService []models.HostService
	for rows.Next() {
		hs := models.HostService{}
		err = rows.Scan(
			&hs.ID, &hs.HostID, &hs.ServiceID, &hs.Active, &hs.ScheduleNumber,
			&hs.ScheduleUnit, &hs.LastCheck, &hs.Status, &hs.CreatedAt, &hs.UpdatedAt,
			&hs.Service.ID, &hs.Service.ServiceName, &hs.Service.Active, &hs.Service.Icon,
			&hs.Service.CreatedAt, &hs.Service.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return h, err
		}
		HostService = append(HostService, hs)
	}
	h.HostService = HostService
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
		HostService := make([]models.HostService, 0)
		err = rows.Scan(
			&h.ID, &h.HostName, &h.CanonicalName,
			&h.URL, &h.IP, &h.IPv6, &h.Location,
			&h.OS, &h.Active, &h.CreatedAt, &h.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		hs_query := `
		SELECT hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number,
			   hs.schedule_unit, hs.last_check, hs.status, hs.created_at, hs.updated_at,
		       s.id, s.service_name, s.active, s.icon, s.created_at, s.updated_at
	      FROM host_services AS hs
   		  LEFT JOIN services AS s
	 	    ON (s.id = hs.service_id)
  		 WHERE host_id = $1;`

		hs_rows, err := m.DB.QueryContext(ctx, hs_query, h.ID)
		if err != nil {
			return nil, err
		}
		// defer hs_rows.Close()
		for hs_rows.Next() {
			var hs models.HostService

			err = hs_rows.Scan(
				&hs.ID, &hs.HostID, &hs.ServiceID, &hs.Active, &hs.ScheduleNumber,
				&hs.ScheduleUnit, &hs.LastCheck, &hs.Status, &hs.CreatedAt, &hs.UpdatedAt,
				&hs.Service.ID, &hs.Service.ServiceName, &hs.Service.Active, &hs.Service.Icon,
				&hs.Service.CreatedAt, &hs.Service.UpdatedAt,
			)
			if err != nil {
				log.Println(err)
				return hosts, err
			}
			// log.Println("append service:", hs.Service.ServiceName)
			HostService = append(HostService, hs)
		}
		h.HostService = HostService
		hosts = append(hosts, h)
		hs_rows.Close()
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return hosts, nil
}

func (m *postgresDBRepo) GetAllServiceStatusCounts() (map[models.ServiceStatus]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ssc := map[models.ServiceStatus]int{
		models.ServiceStatusHealthy: 0,
		models.ServiceStatusProblem: 0,
		models.ServiceStatusPending: 0,
		models.ServiceStatusWarning: 0,
	}

	qry := `
	SELECT status, COUNT (id) 
 	  FROM host_services
     WHERE active = 1
     GROUP BY status;`

	rows, err := m.DB.QueryContext(ctx, qry)
	if err != nil {
		fmt.Println(err)
		return ssc, err
	}
	defer rows.Close()

	for rows.Next() {
		var ss models.ServiceStatus
		var c int

		err = rows.Scan(&ss, &c)
		if err != nil {
			log.Println(err)
			return ssc, err
		}
		ssc[ss] = c
	}
	return ssc, nil
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

func (m *postgresDBRepo) UpdateHostServiceStatusByID(hostID, serviceID, active int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	UPDATE public.host_services
	   SET active = $1
     WHERE host_id = $2 AND service_id = $3;
	`
	_, err := m.DB.ExecContext(ctx, stmt,
		active, hostID, serviceID,
	)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetServiceByStatus(status models.ServiceStatus) ([][3]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	qry := `
	SELECT h.id, h.host_name, s.service_name
      FROM host_services as hs
      LEFT JOIN services as s
        ON hs.service_id = s.id
      LEFT JOIN hosts as h
        ON hs.host_id = h.id
     WHERE hs.status = $1 AND hs.active = 1
	 ORDER By h.host_name, s.service_name;`

	hostServiceNamePair := make([][3]string, 0)
	rows, err := m.DB.QueryContext(ctx, qry, int(status))
	if err != nil {
		fmt.Println(err)
		return hostServiceNamePair, err
	}
	defer rows.Close()

	for rows.Next() {
		var pair = [3]string{}
		var i int
		err = rows.Scan(&i, &pair[1], &pair[2])
		pair[0] = strconv.Itoa(i)
		if err != nil {
			log.Println(err)
			return hostServiceNamePair, err
		}
		hostServiceNamePair = append(hostServiceNamePair, pair)
	}
	return hostServiceNamePair, nil
}
