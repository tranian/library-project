CREATE USER 'library_dbuser'@'localhost' IDENTIFIED BY 'library_securepassword';
GRANT ALL PRIVILEGES ON library_db.* TO 'library_dbuser'@'localhost';
FLUSH PRIVILEGES;
