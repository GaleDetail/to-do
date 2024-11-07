CREATE TABLE IF NOT EXISTS records (
                                       id INT AUTO_INCREMENT,
                                       user_id INT,
                                       content TEXT,
                                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                       PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
