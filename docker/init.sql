CREATE SEQUENCE coredns_records_id_seq
    INCREMENT 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    START 1
    CACHE 1;

CREATE TABLE coredns_records (
    id bigint DEFAULT nextval('coredns_records_id_seq'::regclass) NOT NULL,
    zone VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    ttl INT DEFAULT NULL,
    content TEXT,
    record_type VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
) ;

-- add test data

INSERT INTO coredns_records (zone, name, ttl, content, record_type) VALUES
('example.org.', '', 30, '{"ip": "1.1.1.1"}', 'A'),
('example.org.', '', 60, '{"ip": "1.1.1.0"}', 'A'),
('example.org.', 'test', 30, '{"text": "hello"}', 'TXT'),
('example.org.', 'mail', 30, '{"host" : "mail.example.org.","priority" : 10}', 'MX');