CREATE TABLE  IF NOT EXISTS  "withdraw"(
    id          VARCHAR(30) PRIMARY KEY ,
    sum     DOUBLE PRECISION NOT NULL,
    user_id     INT NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now() 
);