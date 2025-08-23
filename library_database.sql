CREATE DATABASE library_db;
USE library_db;
CREATE TABLE books (
	id INT AUTO_INCREMENT PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	author VARCHAR(255),
	year INT,
	is_available BOOLEAN DEFAULT TRUE,
	checked_out_to VARCHAR(255),
	checkout_date DATETIME,
	due_date DATETIME
)
