INSERT INTO users (
    first_name,
    last_name,
    user_active,
    access_level,
    email,
    password,
    created_at,
    updated_at
) VALUES (
    'Admin',
    'User',
    1,
    3,
    'admin@example.com',
    '$2a$12$F0IySlEPTRjXOY5l3mrl9.aEWJrpajLuVn3gKcZlXqLB9AqY0BB02',
    '2022-11-30 20:24:19',
    '2023-01-18 12:00:21'
);

INSERT INTO preferences (
    name, preference
) VALUES
('monitoring_live',       '1'),
('check_interval_amount', '3'),
('check_interval_unit',   'm'),
('notify_via_email',      '0');

-- see checker.go
INSERT INTO services (
    service_name, active, icon
) VALUES 
('HTTP', 1, '🔓'),
('HTTPS', 0, '🔒');

INSERT INTO hosts (
    host_name, canonical_name, url, ip, ipv6, location, os, active
) VALUES
('TestServer', 'TestServer', 'http://192.168.50.168:8080', '192.168.50.168', '', 'Taipei', 'Fedora', '1'),
('Router', 'Router', 'http://192.168.50.1:80', '192.168.50.1', '', 'Taipei', 'unknown', '1');

INSERT INTO host_services (
    host_id,
    service_id,
    status,
    schedule_number,
    schedule_unit
) VALUES
(1, 1, 10, 4, 's'),
(2, 1, 10, 4, 's');