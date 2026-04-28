CREATE TABLE IF NOT EXISTS EVENTS (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    image_path VARCHAR(255),
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    location VARCHAR(255),
    delivery_method VARCHAR(255),
    created_by_user_id CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by_user_id) REFERENCES USERS(id) ON DELETE SET NULL
)
