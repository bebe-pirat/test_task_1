CREATE TABLE subscription(
	id SERIAL PRIMARY KEY,
	service_name VARCHAR(256) NOT NULL,
	price INT NOT NULL,
	user_id UUID NOT NULL,
	start_date date NOT NULL,
    formatted_start_date text GENERATED ALWAYS AS (to_char(start_date, 'MM-YYYY')) STORED,
	end_date DATE ,
    formatted_end_date text GENERATED ALWAYS AS (to_char(start_date, 'MM-YYYY')) STORED
);
