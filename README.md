# PokcetAnalyst

PocketAnalyst is an app meant to provide predictions on stock prices based on
historical data, technical analysis, fundamentals, and broad economic data.

## Architecture

PocketAnalyst uses MVC (Model-View-Controller) architecture.

### Database Schema

#### Companies Table

Stores basic information about companies whose stocks we're tracking.
Acts as a master reference table.

- company_id: Primary key, uniquely identifies each company
- symbol: Stock ticker symbol (e.g., "AAPL", "MSFT")
- name: Full company name
- sector: Business sector (e.g., "Technology", "Healthcare")
- industry: Specific industry within the sector
- exchange: Stock exchange where the company is listed (e.g., "NASDAQ")
- is_active: Boolean flag to indicate if we're actively tracking this company
- created_at: Timestamp when the record was created

#### Data Sources

Stores configuration information for external data providers.

- source_id: Primary key for each data source configuration
- source_name: Name of the data provider (e.g., "AlphaVantage", "YahooFinance")
- source_type: Type of data provided (e.g., "PRICE", "FUNDAMENTAL", "SENTIMENT")
- base_url: Base URL for API calls
- rate_limit_per_minute: API rate limit per minute
- rate_limit_per_day: API rate limit per day
- config_parameters: JSON containing non-sensitive configuration parameters
- is_active: Whether this data source is active
- created_at: Timestamp when the record was created

#### Data Fetch Jobs

Tracks scheduled data fetching operations.

- job_id: Primary key for each fetch job
- source_id: Foreign key to the data_sources table
- entity_type: What type of entity to fetch data for (e.g., "SYMBOL", "SECTOR", "MARKET")
- entity_value: The actual symbol, sector name, etc.
- data_type: Type of data to fetch (e.g., "PRICE", "FUNDAMENTALS", "SENTIMENT")
- frequency: How often to fetch (e.g., "daily", "weekly", "monthly")
- parameters: JSON containing additional parameters for the job
- last_execution: When job was last executed
- last_success: When job last completed successfully
- next_scheduled: When job should next run
- status: Current job status (e.g., "PENDING", "RUNNING", "SUCCESS", "FAILED")
- is_active: Whether this job is active
- created_at: Timestamp when the record was created

#### Stock Prices

Stores raw historical stock price data fetched from external APIs.

- price_id: Primary key for each daily stock record
- company_id: Foreign key linking to the companies table
- symbol: Stock ticker symbol (duplicated for query convenience)
- date: The trading date for this price record
- open_price: Opening price for the day
- high_price: Highest price during the day
- low_price: Lowest price during the day
- close_price: Closing price for the day
- adjusted_close: Closing price adjusted for splits and dividends
- volume: Number of shares traded
- dividend_amount: Amount of dividend issued on this date
- split_coefficient: Stock split factor on this date
- source_id: Foreign key linking to the data_sources table
- created_at: Timestamp when the record was created

#### Technical Indicators

Stores calculated technical indicators based on the raw stock data.
These are used as features for ML models.

- indicator_id: Primary key for each indicator record
- company_id: Foreign key linking to the companies table
- symbol: Stock ticker symbol (duplicated for query convenience)
- date: The date for this indicator value
- indicator_type: Type of indicator (e.g., "SMA", "EMA", "BOLLINGER", "RSI")
- period: Time period for the indicator (e.g., 14 days for a 14-day RSI)
- value: Primary indicator value
- upper_band: Upper band value (for indicators like Bollinger Bands)
- lower_band: Lower band value (for indicators like Bollinger Bands)
- created_at: Timestamp when this indicator was calculated

#### Fundamental Data

Stores financial data like income statements, balance sheets, and key metrics.

- fundamental_id: Primary key for each fundamental data record
- company_id: Foreign key linking to the companies table
- symbol: Stock ticker symbol (duplicated for query convenience)
- date: Date of the report/data
- report_type: Type of report (e.g., "QUARTERLY", "ANNUAL")
- data_type: Kind of data (e.g., "INCOME_STATEMENT", "BALANCE_SHEET", "CASH_FLOW", "RATIOS")
- data: JSON containing all financial metrics
- source_id: Foreign key linking to the data_sources table
- created_at: Timestamp when the record was created

#### News Events

Tracks significant news and events that might impact stock prices.

- event_id: Primary key for each event
- company_id: Foreign key linking to the companies table (NULL for market-wide events)
- event_type: Type of event (e.g., "NEWS", "EARNINGS", "DIVIDEND", "SPLIT")
- event_date: When the event occurred
- title: Event title or headline
- content: Full event content or description
- source: Where the event information came from
- url: Link to original content
- sentiment_score: Pre-calculated sentiment (-1 to 1)
- source_id: Foreign key linking to the data_sources table
- created_at: Timestamp when the record was created

#### Sentiment Data

Stores sentiment analysis from social media, news, etc.

- sentiment_id: Primary key for each sentiment record
- company_id: Foreign key linking to the companies table
- symbol: Stock ticker symbol (duplicated for query convenience)
- date: Date of the sentiment data
- source_type: Where sentiment was measured (e.g., "TWITTER", "REDDIT", "NEWS", "GOOGLE_TRENDS")
- sentiment_score: Sentiment rating (-1.0 to 1.0)
- volume: Number of mentions/posts
- trending_keywords: JSON containing keywords and their frequencies
- raw_data: Optional JSON storage for source data
- source_id: Foreign key linking to the data_sources table
- created_at: Timestamp when the record was created

#### Feature Sets

Defines collections of features for use in ML models.

- feature_set_id: Primary key for each feature set
- name: Name of the feature set
- description: Description of what this feature set is used for
- feature_definitions: JSON defining how to calculate each feature
- created_at: Timestamp when the record was created
- updated_at: Timestamp when the record was last updated

#### Feature Data

Stores pre-computed features for ML models.

- feature_data_id: Primary key for each feature data record
- feature_set_id: Foreign key linking to the feature_sets table
- company_id: Foreign key linking to the companies table
- symbol: Stock ticker symbol (duplicated for query convenience)
- date: Date for these feature values
- features: JSON mapping feature names to values
- created_at: Timestamp when the record was created

#### ML Models

Stores information about machine learning models.

- model_id: Primary key for each model
- name: Model name
- model_type: Type of model (e.g., "LSTM", "RANDOM_FOREST", "XGBOOST")
- target_type: Type of prediction (e.g., "CLASSIFICATION", "REGRESSION")
- target_variable: What we're predicting (e.g., "PRICE_DIRECTION", "PRICE", "VOLATILITY")
- feature_set_id: Foreign key linking to the feature_sets table
- parameters: JSON containing model hyperparameters and configuration
- model_path: Path to stored model file/directory
- is_active: Whether this model is active
- created_at: Timestamp when the record was created
- updated_at: Timestamp when the record was last updated

#### Model Training History

Tracks the history of machine learning model training sessions.

- training_id: Primary key for each training session
- model_id: Foreign key linking to the ml_models table
- training_date: When the model was trained
- training_dataset_start: Start date of training data
- training_dataset_end: End date of training data
- validation_dataset_start: Start date of validation data
- validation_dataset_end: End date of validation data
- training_accuracy: Accuracy on the training dataset
- validation_accuracy: Accuracy on the validation dataset
- metrics: JSON containing detailed metrics (precision, recall, etc.)
- training_duration_seconds: How long the training took
- notes: Any additional notes about this training session
- created_at: Timestamp when the record was created

#### ML Predictions

Stores predictions generated by machine learning models.

- prediction_id: Primary key for each prediction
- model_id: Foreign key linking to the ml_models table
- company_id: Foreign key linking to the companies table
- symbol: Stock ticker symbol (duplicated for query convenience)
- prediction_date: When the prediction was made
- target_date: Future date being predicted
- prediction_value: The predicted value
- prediction_confidence: Model confidence (0-1)
- extra_data: Additional prediction information
- actual_value: The actual value once the target date arrives
- accuracy_metric: How accurate the prediction was
- created_at: Timestamp when the record was created

#### Job Execution Logs

Tracks execution of system jobs for monitoring and debugging.

- log_id: Primary key for each log entry
- job_id: Foreign key linking to the data_fetch_jobs table
- job_type: Type of job (e.g., "DATA_FETCH", "FEATURE_CALCULATION", "MODEL_TRAINING")
- start_time: When the job started
- end_time: When the job ended
- status: Job status (e.g., "RUNNING", "SUCCESS", "FAILED")
- records_processed: Number of records processed
- error_message: Error message if job failed
- details: Additional job details
- created_at: Timestamp when the record was created

## Testing

- Get an Alpha Vantage API Key from [here](https://www.alphavantage.co/)
- Run export ALPHA_VANTAGE_API_KEY=yourapikey
- Run chmod +x scripts/run-tests.sh in terminal
- Use one of the following make commands
    - make test-all
    - make test-clients
    - make test-repositories
    - make test-core
    - make test-clean
