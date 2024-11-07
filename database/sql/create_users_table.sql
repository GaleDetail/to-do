CREATE TABLE IF NOT EXISTS users (
                                     id INT AUTO_INCREMENT,
                                     username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
    );
