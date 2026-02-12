CREATE TABLE notification_subscriptions
(
    id SERIAL PRIMARY KEY,
    user_id INT,
    location_id INT,
    channel VARCHAR(20) NOT NULL,
    province_id INT(10,2),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE notification_logs
(
    id SERIAL PRIMARY KEY,
    subscription_id INT REFERENCES notification_subscriptions(id),
    location_id INT,
    water_level DECIMAL(10,2),
    message TEXT,
    status VARCHAR(20),
    channel VARCHAR(20),
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);