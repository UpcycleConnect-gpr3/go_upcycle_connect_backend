CREATE TABLE IF NOT EXISTS OBJECT_DELIVERY_METHOD (
    object_id INT NOT NULL,
    delivery_method_id INT NOT NULL,
    PRIMARY KEY (object_id, delivery_method_id)
);
CREATE TABLE IF NOT EXISTS OBJECT_PROJECT (
    object_id INT NOT NULL,
    project_id INT NOT NULL,
    PRIMARY KEY (object_id, project_id)
);
CREATE TABLE IF NOT EXISTS OBJECT_USER (
    object_id INT NOT NULL,
    user_id CHAR(36) NOT NULL,
    PRIMARY KEY (object_id, user_id)
);
CREATE TABLE IF NOT EXISTS PROJECT_STEP (
    project_id INT NOT NULL,
    step_id INT NOT NULL,
    PRIMARY KEY (project_id, step_id)
);