CREATE TABLE IF NOT EXISTS "order"(
    id          VARCHAR(30) PRIMARY KEY ,
    status      VARCHAR(20) NOT NULL,
    user_id     INT NOT NULL,
    accrual     INT DEFAULT 0,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now() 
);