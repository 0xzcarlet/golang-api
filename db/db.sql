-- Supabase PostgreSQL Schema
-- Database: golang_api

-- ============================================================
-- Drop tables (in dependency order)
-- ============================================================

DROP TABLE IF EXISTS place_category_list CASCADE;
DROP TABLE IF EXISTS place_category CASCADE;
DROP TABLE IF EXISTS place_link CASCADE;
DROP TABLE IF EXISTS place_status CASCADE;
DROP TABLE IF EXISTS place CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- ============================================================
-- Table: users
-- ============================================================

CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  name VARCHAR(120) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users (email);

-- ============================================================
-- Table: place
-- ============================================================

CREATE TABLE place (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  name VARCHAR(255),
  link VARCHAR(255),
  link_type SMALLINT,
  description TEXT,
  go_at DATE,                    -- rencana tanggal pergi ke tempat tersebut
  go_at_time TIMESTAMP,         -- jam rencana pergi
  status SMALLINT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_place_user_id ON place (user_id);

-- ============================================================
-- Table: place_category
-- ============================================================

CREATE TABLE place_category (
  id SERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  name VARCHAR(50)
);

CREATE INDEX idx_place_category_user_id ON place_category (user_id);

-- ============================================================
-- Table: place_category_list
-- ============================================================

CREATE TABLE place_category_list (
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  place_id BIGINT REFERENCES place(id) ON DELETE CASCADE ON UPDATE CASCADE,
  category_id INTEGER REFERENCES place_category(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_pcl_place_id ON place_category_list (place_id);
CREATE INDEX idx_pcl_user_id ON place_category_list (user_id);
CREATE INDEX idx_pcl_category_id ON place_category_list (category_id);

-- ============================================================
-- Table: place_link
-- ============================================================

CREATE TABLE place_link (
  id SMALLSERIAL PRIMARY KEY,
  domain VARCHAR(50),
  icon VARCHAR(255)              -- icon social media or g maps for type link
);

-- ============================================================
-- Table: place_status
-- ============================================================

CREATE TABLE place_status (
  id SMALLSERIAL PRIMARY KEY,
  name VARCHAR(50),
  background VARCHAR(8)          -- background for CSS badge with hex code
);

-- ============================================================
-- Auto-update updated_at trigger (replaces MySQL ON UPDATE)
-- ============================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_place_updated_at
  BEFORE UPDATE ON place
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
