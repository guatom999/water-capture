CREATE TABLE notification_subscriptions  (
    id SERIAL PRIMARY KEY,
    user_id INT,
    location_id INT,
    channel VARCHAR(20) NOT NULL, -- 'email', 'line', 'sms'
    target VARCHAR(255) NOT NULL, -- email address, LINE user ID, phone
    threshold_level DECIMAL(10,2), -- alert when level exceeds this
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE notification_logs (
    id SERIAL PRIMARY KEY,
    subscription_id INT REFERENCES notification_subscriptions(id),
    location_id INT,
    water_level DECIMAL(10,2),
    message TEXT,
    channel VARCHAR(20),
    status VARCHAR(20), -- 'pending', 'sent', 'failed'
    sent_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);