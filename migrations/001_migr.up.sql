CREATE TABLE subscription(
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(256) NOT NULL,
    price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    formatted_start_date TEXT,
    end_date DATE,
    formatted_end_date TEXT
);

CREATE OR REPLACE FUNCTION update_formatted_dates()
RETURNS TRIGGER AS $$
BEGIN
    NEW.formatted_start_date := to_char(NEW.start_date, 'MM-YYYY');
    NEW.formatted_end_date := CASE 
        WHEN NEW.end_date IS NOT NULL THEN to_char(NEW.end_date, 'MM-YYYY')
        ELSE NULL 
    END;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_dates
    BEFORE INSERT OR UPDATE ON subscription
    FOR EACH ROW
    EXECUTE FUNCTION update_formatted_dates();