-- This is an example of how to insert a data source into the data_sources table.
INSERT INTO data_sources (source_name, source_type, base_url, rate_limit_per_minute, rate_limit_per_day, config_parameters, is_active) 
VALUES ('AlphaVantage', 'PRICE', 'https://www.alphavantage.co/query', 5, 500, '{}', true) 
ON CONFLICT (source_name) DO NOTHING;
