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
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Stocks table to store raw stock price data
CREATE TABLE IF NOT EXISTS stocks (
	stock_data_id SERIAL PRIMARY KEY,
	company_id INTEGER NOT NULL REFERENCES companies(company_id),
	symbol VARCHAR(20) NOT NULL,
	date DATE NOT NULL,
	open_price NUMERIC(10, 4),
	high_price NUMERIC(10, 4),
	low_price NUMERIC(10, 4),
	close_price NUMERIC(10, 4),
	adjusted_close_price NUMERIC(10, 4),
	volume BIGINT,
	dividend_amount NUMERIC(10, 4),
	split_coefficient NUMERIC(10, 4),
	data_source VARCHAR(50) NOT NULL,
	last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT stocks_company_date_unique UNIQUE (company_id, date)
);

-- Create index for faster querying by date range and symbol
CREATE INDEX IF NOT EXISTS idx_stocks_symbol_date ON stocks(symbol, date);
CREATE INDEX IF NOT EXISTS idx_stocks_date ON stocks(date);

-- Table for storing calculated technical indicators
CREATE TABLE IF NOT EXISTS technical_indicators (
	indicator_id SERIAL PRIMARY KEY,
	stock_data_id INTEGER REFERENCES stocks(stock_data_id),
	symbol VARCHAR(20) NOT NULL,
	date DATE NOT NULL,
	indicator_type VARCHAR(50) NOT NULL,
	period INTEGER NOT NULL,
	value NUMERIC(10, 4),
	upper_band NUMERIC(10, 4),
	lower_band NUMERIC(10, 4),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT indicator_stock_type_period_unique UNIQUE (stock_data_id, indicator_type, period)
);

-- Create index for faster indicator queries
CREATE INDEX IF NOT EXISTS idx_indicators_symbol_date ON technical_indicators(symbol, date);
CREATE INDEX IF NOT EXISTS idx_indicators_type_period ON technical_indicators(indicator_type, period);

-- Table that represents ML model prediction outputs
CREATE TABLE IF NOT EXISTS ml_predictions (
	prediction_id SERIAL PRIMARY KEY,
	symbol VARCHAR(20) NOT NULL,
	target_date DATE NOT NULL,
	prediction_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
	prediction_type VARCHAR(50) NOT NULL, 
	predicted_value NUMERIC(10, 4) NOT NULL, 
	confidence NUMERIC(5, 4) NOT NULL, 
	model_version VARCHAR(50) NOT NULL, 
	features_used TEXT, 
	actual_value NUMERIC(10, 4),
	accuracy_metric NUMERIC(10, 4),
	CONSTRAINT prediction_unique UNIQUE (symbol, target_date, prediction_type, model_version)
);

-- Create index for faster prediction queries
CREATE INDEX IF NOT EXISTS idx_predictions_symbol_target ON ml_predictions(symbol, target_date);
CREATE INDEX IF NOT EXISTS idx_predictions_model ON ml_predictions(model_version);

-- Data fetch config table
CREATE TABLE IF NOT EXISTS data_fetch_configs (
    config_id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    data_source VARCHAR(50) NOT NULL,
    api_key VARCHAR(100),
    start_date DATE,
    end_date DATE,
    frequency VARCHAR(20) NOT NULL DEFAULT 'daily',
    last_fetched TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fetch_config_unique UNIQUE (symbol, data_source)
);

-- Model training history table
CREATE TABLE IF NOT EXISTS model_training_history (
    training_id SERIAL PRIMARY KEY,
    model_version VARCHAR(50) NOT NULL,
    training_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    model_type VARCHAR(50) NOT NULL,
    parameters JSONB,
    training_accuracy NUMERIC(10,4),
    validation_accuracy NUMERIC(10,4),
    features_used TEXT,
    training_duration_seconds INTEGER,
    notes TEXT
);
