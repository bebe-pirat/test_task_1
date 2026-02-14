CREATE TABLE subscription(
	id INT PRIMARY KEY,
	service_name VARCHAR(256) NOT NULL,
	price INT NOT NULL,
	user_id INT NOT NULL,
	start_date DATE
);