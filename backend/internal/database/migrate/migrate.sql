CREATE TABLE IF NOT EXISTS locations (
    id           BIGSERIAL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,        
    description  TEXT,                         
    latitude     DECIMAL(9,6) NOT NULL,        
    longitude    DECIMAL(9,6) NOT NULL,        
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE danger_level AS ENUM ('SAFE', 'WATCH', 'DANGER', 'CRITICAL');

CREATE TYPE water_status AS ENUM ('ACTIVE', 'PENDING_DELETION', 'DELETED');

CREATE TABLE IF NOT EXISTS water_levels (
    id            BIGSERIAL PRIMARY KEY,
    location_id   BIGINT NOT NULL,
    level_cm      NUMERIC(10,2) NOT NULL,      
    image         VARCHAR(255),                 
    danger        danger_level NOT NULL,       
    is_flooded    BOOLEAN NOT NULL,           
    source        VARCHAR(50),                
    measured_at   TIMESTAMPTZ NOT NULL,       
    note          TEXT,
    status        water_status NOT NULL,
    deleted_at    TIMESTAMPTZ,                
    scheduled_delete_at TIMESTAMPTZ                   
);