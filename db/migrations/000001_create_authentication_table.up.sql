-- Create Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    verified_at TIMESTAMP
);

-- Create TempUsers table
CREATE TABLE IF NOT EXISTS temp_users (
    user_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    otp VARCHAR(255) NOT NULL,
    otp_expires TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create ForgotPassword table
CREATE TABLE IF NOT EXISTS forgot_password (
    email VARCHAR(255) NOT NULL,
    otp VARCHAR(255) NOT NULL,
    otp_expires TIMESTAMP NOT NULL,
    verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
