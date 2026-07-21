CREATE TABLE IF NOT EXISTS tasks (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  status VARCHAR(50) DEFAULT 'todo',
  assignee VARCHAR(255),
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  UNIQUE KEY idx_tasks_title (title),
  INDEX idx_tasks_deleted_at (deleted_at),
  INDEX idx_tasks_status (status),
  INDEX idx_tasks_assignee (assignee)
);
