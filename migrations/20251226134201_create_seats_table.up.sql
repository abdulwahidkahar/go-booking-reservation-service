CREATE TABLE seats (
  id BIGSERIAL PRIMARY KEY,
  flight_id BIGINT NOT NULL,
  seat_number VARCHAR(10) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE',
  locked_until TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_seat_flight
    FOREIGN KEY (flight_id)
    REFERENCES flights(id)
    ON DELETE CASCADE,

  CONSTRAINT unique_seat_per_flight
    UNIQUE (flight_id, seat_number),

  CONSTRAINT seat_status_check
    CHECK (status IN ('AVAILABLE', 'LOCKED', 'BOOKED'))
);
