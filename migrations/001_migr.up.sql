CREATE TABLE subscription(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	service_name VARCHAR(256) NOT NULL,
	price INT NOT NULL,
	user_id UUID NOT NULL,
	start_date DATE NOT NULL,
	end_date DATE 
);