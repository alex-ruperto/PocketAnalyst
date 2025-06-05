-- Database schema for analysis and prediction

-- Companies table to store company information
CREATE TABLE IF NOT EXISTS companies (
	company_id SERIAL PRIMARY KEY,
	symbol VARCHAR(20) NOT NULL UNIQUE,
	name VARCHAR(255) NOT NULL,
	sector VARCHAR(100),
	industry VARCHAR(100),
	exchange VARCHAR(50),
	is_active BOOLEAN DEFAULT TRUE,
	last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Data sources configuration table
CREATE TABLE IF NOT EXISTS data_sources (
	source_id SERIAL PRIMARY KEY,
	source_name VARCHAR(100) NOT NULL UNIQUE,	-- e.g., "AlphaVantage, YahooFinance, Google Trends"
	source_type VARCHAR(50) NOT NULL,		-- e.g., "PRICE", "FUNDAMENTAL", "SENTIMENT"
	base_url VARCHAR(255),
	rate_limit_per_minute INTEGER,
	rate_limit_per_day INTEGER,
	config_parameters JSONB,			-- Non-sensitive configuration parameters
	is_active BOOLEAN DEFAULT TRUE,
	last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Data fetch jobs for tracking what to fetch and when
CREATE TABLE IF NOT EXISTS data_fetch_jobs (
	job_id SERIAL PRIMARY KEY,
	source_id INTEGER NOT NULL REFERENCES data_sources(source_id),
	entity_type VARCHAR(50) NOT NULL,			-- SYMBOL, SECTOR, MARKET 
	entity_value VARCHAR(100) NOT NULL,		-- The actual symbol, sector name, etc.
	data_type VARCHAR(50) NOT NULL,			-- e.g., "PRICE", "FUNDAMENTALS", "SENTIMENT"
	frequency VARCHAR(20) NOT NULL,			-- "daily", "weekly", "monthly"
	parameters JSONB,				-- Additional parameters for this job
	last_execution TIMESTAMP,			-- When job was last executed
	last_success TIMESTAMP,				-- When job was last completed successfully
	next_scheduled TIMESTAMP,			-- When the job should run next
	status VARCHAR(20) DEFAULT 'PENDING',		-- "PENDING", "RUNNING", "SUCCESS", "FAILED"
	is_active BOOLEAN DEFAULT TRUE,
	last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fetch_job_unique UNIQUE (source_id, entity_type, entity_value, data_type)
);

-- Stock price data
CREATE TABLE IF NOT EXISTS stock_prices (
    price_id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES companies(company_id),
    symbol VARCHAR(20) NOT NULL,
    date DATE NOT NULL,
    open_price NUMERIC(15, 5),
    high_price NUMERIC(15, 5),
    low_price NUMERIC(15, 5),
    close_price NUMERIC(15, 5),
    adjusted_close NUMERIC(15, 5),
    volume BIGINT,
    dividend_amount NUMERIC(10, 5),
    split_coefficient NUMERIC(10, 5),
    source_id INTEGER NOT NULL REFERENCES data_sources(source_id),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT stock_price_unique UNIQUE (company_id, date, source_id)
);

-- Table for storing calculated technical indicators
CREATE TABLE IF NOT EXISTS technical_indicators (
	indicator_id SERIAL PRIMARY KEY,
	company_id INTEGER NOT NULL REFERENCES companies(company_id),
	symbol VARCHAR(20) NOT NULL,
	date DATE NOT NULL,
	indicator_type VARCHAR(50) NOT NULL,
	period INTEGER NOT NULL,
	value NUMERIC(15, 5),
	upper_band NUMERIC(15, 5),
	lower_band NUMERIC(15, 5),
	last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT technical_indicator_unique UNIQUE (company_id, date, indicator_type, period)
);

-- Fundamental data (financial statements, ratios, etc.)
CREATE TABLE IF NOT EXISTS fundamental_data (
    fundamental_id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES companies(company_id),
    symbol VARCHAR(20) NOT NULL,
    date DATE NOT NULL,                        -- Date of the report/data
    report_type VARCHAR(20) NOT NULL,          -- "QUARTERLY", "ANNUAL"
    data_type VARCHAR(50) NOT NULL,            -- "INCOME_STATEMENT", "BALANCE_SHEET", "CASH_FLOW", "RATIOS"
    data JSONB NOT NULL,                       -- Store all metrics in flexible JSON format
    source_id INTEGER NOT NULL REFERENCES data_sources(source_id),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fundamental_data_unique UNIQUE (company_id, date, report_type, data_type, source_id)
);

-- News and events data
CREATE TABLE IF NOT EXISTS news_events (
    event_id SERIAL PRIMARY KEY,
    company_id INTEGER REFERENCES companies(company_id), -- NULL means market-wide
    event_type VARCHAR(50) NOT NULL,           -- "NEWS", "EARNINGS", "DIVIDEND", "SPLIT", etc.
    event_date TIMESTAMP NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    source VARCHAR(100) NOT NULL,
    url VARCHAR(255),
    sentiment_score NUMERIC(5, 4),             -- Optional pre-calculated sentiment (-1 to 1)
    source_id INTEGER NOT NULL REFERENCES data_sources(source_id),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sentiment data from social media, news, etc.
CREATE TABLE IF NOT EXISTS sentiment_data (
    sentiment_id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES companies(company_id),
    symbol VARCHAR(20) NOT NULL,
    date DATE NOT NULL,
    source_type VARCHAR(50) NOT NULL,          -- "TWITTER", "REDDIT", "NEWS", "GOOGLE_TRENDS"
    sentiment_score NUMERIC(5, 4) NOT NULL,    -- -1.0 to 1.0
    volume INTEGER NOT NULL,                   -- Number of mentions
    trending_keywords JSONB,                   -- Keywords/phrases and their counts
    raw_data JSONB,                            -- Optional storage of source data
    source_id INTEGER NOT NULL REFERENCES data_sources(source_id),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sentiment_data_unique UNIQUE (company_id, date, source_type, source_id)
);

-- Feature sets (definitions for ML features)
CREATE TABLE IF NOT EXISTS feature_sets (
    feature_set_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    feature_definitions JSONB NOT NULL,        -- Definitions of how to calculate each feature
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Feature data (pre-computed features for ML)
CREATE TABLE IF NOT EXISTS feature_data (
    feature_data_id SERIAL PRIMARY KEY,
    feature_set_id INTEGER NOT NULL REFERENCES feature_sets(feature_set_id),
    company_id INTEGER NOT NULL REFERENCES companies(company_id),
    symbol VARCHAR(20) NOT NULL,
    date DATE NOT NULL,
    features JSONB NOT NULL,                   -- Feature name -> value mapping
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT feature_data_unique UNIQUE (feature_set_id, company_id, date)
);

-- ML models
CREATE TABLE IF NOT EXISTS ml_models (
    model_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    model_type VARCHAR(50) NOT NULL,           -- "LSTM", "RANDOM_FOREST", "XGBOOST", etc.
    target_type VARCHAR(20) NOT NULL,          -- "CLASSIFICATION", "REGRESSION"
    target_variable VARCHAR(50) NOT NULL,      -- What we're predicting: "PRICE_DIRECTION", "PRICE", "VOLATILITY"
    feature_set_id INTEGER NOT NULL REFERENCES feature_sets(feature_set_id),
    parameters JSONB NOT NULL,                 -- Hyperparameters and model configuration
    model_path VARCHAR(255),                   -- Path to stored model file/directory
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Model training history
CREATE TABLE IF NOT EXISTS model_training_history (
    training_id SERIAL PRIMARY KEY,
    model_id INTEGER NOT NULL REFERENCES ml_models(model_id),
    training_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    training_dataset_start DATE NOT NULL,
    training_dataset_end DATE NOT NULL,
    validation_dataset_start DATE NOT NULL,
    validation_dataset_end DATE NOT NULL,
    training_accuracy NUMERIC(10, 4),
    validation_accuracy NUMERIC(10, 4),
    metrics JSONB NOT NULL,                    -- Detailed metrics (precision, recall, etc.)
    training_duration_seconds INTEGER,
    notes TEXT,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Predictions generated by ML models
CREATE TABLE IF NOT EXISTS ml_predictions (
    prediction_id SERIAL PRIMARY KEY,
    model_id INTEGER NOT NULL REFERENCES ml_models(model_id),
    company_id INTEGER NOT NULL REFERENCES companies(company_id),
    symbol VARCHAR(20) NOT NULL,
    prediction_date TIMESTAMP NOT NULL,        -- When prediction was made
    target_date DATE NOT NULL,                 -- Future date being predicted
    prediction_value NUMERIC(15, 5) NOT NULL,  -- Predicted value
    prediction_confidence NUMERIC(5, 4),       -- Model confidence (0-1)
    extra_data JSONB,                          -- Additional prediction info
    actual_value NUMERIC(15, 5),               -- Actual value once known
    accuracy_metric NUMERIC(10, 4),            -- How accurate the prediction was
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT prediction_unique UNIQUE (model_id, company_id, prediction_date, target_date)
);

-- System job scheduler status tracking
CREATE TABLE IF NOT EXISTS job_execution_logs (
    log_id SERIAL PRIMARY KEY,
    job_id INTEGER REFERENCES data_fetch_jobs(job_id),
    job_type VARCHAR(50) NOT NULL,             -- "DATA_FETCH", "FEATURE_CALCULATION", "MODEL_TRAINING", etc.
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    status VARCHAR(20) NOT NULL,               -- "RUNNING", "SUCCESS", "FAILED"
    records_processed INTEGER,
    error_message TEXT,
    details JSONB,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create all necessary indexes for optimized queries
CREATE INDEX IF NOT EXISTS idx_stock_prices_symbol_date ON stock_prices(symbol, date);
CREATE INDEX IF NOT EXISTS idx_stock_prices_company_date ON stock_prices(company_id, date);
CREATE INDEX IF NOT EXISTS idx_technical_indicators_symbol_date ON technical_indicators(symbol, date);
CREATE INDEX IF NOT EXISTS idx_technical_indicators_type_period ON technical_indicators(indicator_type, period);
CREATE INDEX IF NOT EXISTS idx_fundamental_data_company_date ON fundamental_data(company_id, date);
CREATE INDEX IF NOT EXISTS idx_sentiment_data_company_date ON sentiment_data(company_id, date);
CREATE INDEX IF NOT EXISTS idx_news_events_company_date ON news_events(company_id, event_date);
CREATE INDEX IF NOT EXISTS idx_feature_data_company_date ON feature_data(company_id, date);
CREATE INDEX IF NOT EXISTS idx_ml_predictions_company_target ON ml_predictions(company_id, target_date);
CREATE INDEX IF NOT EXISTS idx_data_fetch_jobs_next_scheduled ON data_fetch_jobs(next_scheduled, is_active);
CREATE INDEX IF NOT EXISTS idx_job_execution_logs_job_id ON job_execution_logs(job_id);
