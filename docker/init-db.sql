-- Create databases for our services
CREATE DATABASE productdb OWNER postgres;
CREATE DATABASE paymentdb OWNER postgres;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE productdb TO postgres;
GRANT ALL PRIVILEGES ON DATABASE paymentdb TO postgres; 