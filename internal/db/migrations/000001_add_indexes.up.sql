BEGIN;

CREATE TABLE IF NOT EXISTS categories (
                                          id BIGSERIAL PRIMARY KEY,
                                          name VARCHAR(255) UNIQUE NOT NULL
    );

CREATE TABLE IF NOT EXISTS companies (
                                         id BIGSERIAL PRIMARY KEY,

                                         name VARCHAR(255) NOT NULL,
    bin VARCHAR(20),
    website VARCHAR(255),
    email VARCHAR(255),
    city VARCHAR(120),
    number VARCHAR(120),
    address TEXT,
    industry VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active',

    director_name VARCHAR(255),
    director_pos VARCHAR(255),
    executive_name VARCHAR(255),
    executive_pos VARCHAR(255),
    dir_start VARCHAR(100),
    exec_start VARCHAR(100),

    linkedin VARCHAR(255),
    facebook VARCHAR(255),
    status_fb VARCHAR(100),
    status_link VARCHAR(100),
    li_last_update VARCHAR(100),
    fb_last_update VARCHAR(100),

    procurement_method VARCHAR(255),
    procurement_email VARCHAR(255),
    procurement_phone VARCHAR(120),

    hr_name VARCHAR(255),
    hr_email VARCHAR(255),
    hr_phone VARCHAR(120),

    esg_name VARCHAR(255),
    esg_email VARCHAR(255),
    esg_phone VARCHAR(120),
    esg_report_url VARCHAR(255),
    has_esg_dept BOOLEAN DEFAULT FALSE,
    last_source VARCHAR(255),
    last_parsed_at TIMESTAMP NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    category_id BIGINT,
    CONSTRAINT fk_companies_category
    FOREIGN KEY (category_id)
    REFERENCES categories(id)
    ON DELETE SET NULL
    );

CREATE INDEX IF NOT EXISTS idx_companies_bin ON companies(bin);
CREATE INDEX IF NOT EXISTS idx_companies_category_id ON companies(category_id);

CREATE TABLE IF NOT EXISTS company_logs (
                                            id BIGSERIAL PRIMARY KEY,
                                            company_id BIGINT NOT NULL,
                                            action VARCHAR(100),
    field_name VARCHAR(100),
    old_value TEXT,
    new_value TEXT,
    source VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_company_logs_company
    FOREIGN KEY (company_id)
    REFERENCES companies(id)
    ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_company_logs_company_id ON company_logs(company_id);

COMMIT;