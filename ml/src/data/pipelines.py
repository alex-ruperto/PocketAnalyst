"""
Loading of data via API requests.
Handles GET requests from PostgreSQL database/Go API, retrieving data in JSON format.
"""

import requests
import logging
import time
import pandas as pd
import os
from typing import List, Optional, Dict, Union
from concurrent.futures import ThreadPoolExecutor, as_completed
import threading
from dataclasses import dataclass 

class StockDataError(Exception):
    """
    Custom exception for stock data retrieval errors.
    """
    pass

@dataclass
class BatchConfig:
    """
    Configuration for batch processing operations
    """
    batch_size: int = 50 # Max symbols per API call
    max_workers: int = 4 # Number of concurrent requests
    retry_failed: bool = True # Retry failed symbols individually
    validate_data: bool = True # Validate data quality after loading

class StockDataPipeline:
    """"""
    def __init__(
        self, 
        api_base_url: Optional[str] = None, 
        timeout: int = 60, 
        max_retries: int = 3, 
        rate_limit_delay: float = 0.1,
        batch_config: Optional[BatchConfig] = None
    ):
        self.api_base_url = api_base_url or os.getenv("API_BASE_URL", "http://localhost:8080/api/stocks")
        self.timeout = timeout
        self.max_retries = max_retries
        self.rate_limit_delay = rate_limit_delay
        self.batch_config = batch_config or BatchConfig()
        self.logger = logging.getLogger(__name__)

        # Thread-safe counters for monitoring
        self._lock = threading.Lock()
        self._stats = {
            'symbols_requested': 0,
            'symbols_successful': 0,
            'symbols_failed': 0,
            'total_records_loaded': 0,
            'api_calls_made': 0
        }

    def get_stock_data(self, symbol: str, start_date: str, end_date: str) -> pd.DataFrame:
        """
        Fetch stock data for a single symbol with retry logic.

        Args:
            symbol: Stock ticker symbol
            start_date: Start date in YYYY-MM-DD format
            end_date: End date in YYYY-MM-DD format

        Returns:
            DataFrame containing stock price data

        Raises:
            StockDataError: If data retrieval fails after all retries
        """
        url = f"{self.api_base_url}/get-stock"
        params = {
            "symbol": symbol,
            "start_date": start_date,
            "end_date": end_date
        }

        for attempt in range(self.max_retries):
            try:
                self.logger.debug(f"Fetching data for {symbol}, attempt {attempt + 1}")

                response = requests.get(url, params=params, timeout=self.timeout)
                response.raise_for_status()

                data = response.json()
                if not data:
                    raise StockDataError(f"No data returned for symbol: {symbol}")

                df = pd.DataFrame(data)
                if df.empty:
                    raise StockDataError(f"Empty dataset returned for symbol: {symbol}")

                # Convert date column to datetime and ensure proper format
                df['date'] = pd.to_datetime(df['date'])

                # Sort by date to ensure chronological order
                df = df.sort_values('date').reset_index(drop=True)

                # Apply rate limiting after successful request
                time.sleep(self.rate_limit_delay)

                self.logger.debug(f"Successfully fetched {len(df)} records for {symbol}")
                return df

            except requests.exceptions.Timeout:
                self.logger.warning(f"Timeout for {symbol}, attempt {attempt + 1}")
                if attempt == self.max_retries - 1:
                    raise StockDataError(f"Request timeout for symbol {symbol} after {self.max_retries} attempts")

            except requests.exceptions.ConnectionError:
                self.logger.warning(f"Connection error for {symbol}, attempt {attempt + 1}")
                if attempt == self.max_retries - 1:
                    raise StockDataError(f"Connection error for symbol {symbol} after {self.max_retries} attempts")

            except requests.exceptions.HTTPError as e:
                self.logger.error(f"HTTP error {e.response.status_code} for {symbol}")
                if e.response.status_code in [404, 400]:  # Don't retry client errors
                    raise StockDataError(f"API returned error {e.response.status_code} for symbol {symbol}")
                if attempt == self.max_retries - 1:
                    raise StockDataError(f"HTTP error for symbol {symbol} after {self.max_retries} attempts")

            except (ValueError, KeyError) as e:
                self.logger.error(f"Invalid JSON response for {symbol}: {e}")
                raise StockDataError(f"Invalid response format for symbol {symbol}")

            # Wait before retrying with exponential backoff
            if attempt < self.max_retries - 1:
                wait_time = 2 ** attempt
                self.logger.warning(f"Retrying {symbol} in {wait_time}s...")
                time.sleep(wait_time)

        # If all retries fail, raise StockDataError
        raise StockDataError(f"All {self.max_retries} attempts failed for symbol: {symbol}")


    def get_multiple_stocks_data(
        self,
        symbols: List[str],
        start_date: str,
        end_date: str,
        return_combined: bool = False
    ) -> Union[Dict[str, pd.DataFrame], pd.DataFrame]:
        """
        Fetch stock data for multiple symbols using batch processing.

        Args: 
            symbols: List of stock ticker symbols
            start_date: In YYYY-MM-DD format
            end_date: In YYYY-MM-DD format
            return_combined: If true, return a single DataFrame with all symbols

        Returns:
            Dictionary mapping symbols to DataFrames, or a single combined DataFrame

        Raises:
            StockDataError: If critical batch processing fails
        """
        if not symbols:
            raise StockDataError("No symbols were provided for batch processing")

        self.logger.info(f"Starting batch processing for {len(symbols)} symbols")

        # Reset stats for this operation
        with self._lock:
            self._stats = {k: 0 for k in self._stats}
            self._stats['symbols_requested'] = len(symbols)

        # Split symbols into batches to respect API limits
        symbol_batches = self._create_batches(symbols)

        # Process batches concurrently
        all_stock_data = {}
        failed_symbols = []

        with ThreadPoolExecutor(max_workers=self.batch_config.max_workers) as executor:
            # Submit batch processing tasks
            future_to_batch = {
                executor.submit(
                    self._process_symbol_batch,
                    batch,
                    start_date,
                    end_date,
                ): batch
                for batch in symbol_batches
            }
            
            # Collect results as they complete
            for future in as_completed(future_to_batch):
                batch = future_to_batch[future]
                try:
                    batch_data, batch_failed = future.result()
                    all_stock_data.update(batch_data)
                    failed_symbols.extend(batch_failed)

                except Exception as e:
                    self.logger.error(f"Batch processing failed for symbols {batch}: {e}")
                    failed_symbols.extend(batch)

        # Retry failed symbols individually if configured
        if failed_symbols and self.batch_config.retry_failed:
            self.logger.info(f"Retrying {len(failed_symbols)} failed symbols individually")
            retry_data = self._retry_failed_symbols(failed_symbols, start_date, end_date)
            all_stock_data.update(retry_data)

        # Validate and clean data if configured
        if self.batch_config.validate_data:
            all_stock_data = self._validate_and_clean_data(all_stock_data)

        # Log final statistics
        self._log_batch_statistics(failed_symbols)

        # Return in requested format
        if return_combined:
            return self._combine_stock_dataframes(all_stock_data)
        else:
            return all_stock_data

    def get_distinct_symbols(self) -> List[str]:
        """
        Fetch all distinct symbols available in the database.

        Returns:
            List of stock symbols

        Raises:
            StockDataError: If symbols cannot be retrieved
        """
        url = f"{self.api_base_url}/get-distinct-symbols"

        for attempt in range(self.max_retries):
            try:
                self.logger.debug(f"Fetching avalable symbols, attempt {attempt + 1}")

                response = requests.get(url, timeout=self.timeout)
                response.raise_for_status()

                symbols = response.json()
                if not symbols or not isinstance(symbols, list):
                    raise StockDataError("No symbols returned or invalid format")

                self.logger.info(f"Found {len(symbols)} available symbols")
                return symbols

            except requests.exceptions.Timeout:
                self.logger.warning(f"Timeout fetching symbols, attempt {attempt + 1}")
                if attempt == self.max_retries - 1:
                    raise StockDataError(f"Request timeout fetching symbols after {self.max_retries} attempts")

            except requests.exceptions.ConnectionError:
                self.logger.warning(f"Connection error fetching symbols, attempt {attempt + 1}")
                if attempt == self.max_retries - 1:
                    raise StockDataError(f"Connection error fetching symbols after {self.max_retries} attempts")

            except requests.exceptions.HTTPError as e:
                self.logger.error(f"HTTP error {e.response.status_code} fetching symbols")
                raise StockDataError(f"API returned error {e.response.status_code} fetching symbols")

            except (ValueError, KeyError) as e:
                self.logger.error(f"Invalid JSON response fetching symbols: {e}")
                raise StockDataError(f"Invalid response format fetching symbols")

            # Wait before retrying
            if attempt < self.max_retries - 1:
                time.sleep(2 ** attempt)

        # If all retries fail, raise StockDataError
        raise StockDataError(f"All {self.max_retries} attempts failed to fetch all distinct symbols.")

    def _create_batches(self, symbols: List[str]) -> List[List[str]]:
        """
        Split symbols into batches for API processing.
        """
        batches = []
        for i in range(0, len(symbols), self.batch_config.batch_size):
            batch = symbols[i:i + self.batch_config.batch_size]
            batches.append(batch)

        self.logger.debug(f"Created {len(batches)} batches with max size {self.batch_config.batch_size}")
        return batches

    def _process_symbol_batch(
        self,
        symbols: List[str],
        start_date: str,
        end_date: str
    ) -> tuple[Dict[str, pd.DataFrame], List[str]]:
        """
        Process a single batch of symbols using the multi-stock endpoint.

        Returns:
            Tuple of (successful_data_dict, failed_symbols_list)
        """
        url = f"{self.api_base_url}/get-multiple"
        params = {
            "symbols": ",".join(symbols),
            "start_date": start_date,
            "end_date": end_date
        }

        # Track API call
        with self._lock:
            self._stats['api_calls_made'] += 1

        for attempt in range(self.max_retries):
            try:
                self.logger.debug(f"Processing batch {symbols[:3]}... (attempt {attempt + 1})")

                response = requests.get(url, params= params, timeout=self.timeout)
                response.raise_for_status

                # Parse multi-stock JSON response
                stock_data_dict = response.json()

                # Convert each symbol's data to a DataFrame
                successful_data = {}
                failed_symbols = []

                for symbol in symbols:
                    try:
                        symbol_data = stock_data_dict.get(symbol, [])

                        if not symbol_data:
                            self.logger.warning(f"No data returned for symbol: {symbol}")
                            failed_symbols.append(symbol)
                            continue

                        # Convert to DataFrame and process
                        df = pd.DataFrame(symbol_data)
                        df = self._process_dataframe(df)

                        successful_data[symbol] = df

                        # Update statistics
                        with self._lock:
                            self._stats['symbols_successful'] += 1
                            self._stats['total_records_loaded'] += len(df)

                        self.logger.debug(f"Successfully processed {len(df)} recrods for {symbol}")

                    except Exception as e:
                        self.logger.error(f"Failed to process data for {symbol}: {e}")
                        failed_symbols.append(symbol)
                        with self._lock:
                            self._stats['symbols_failed'] += 1

                # Apply rate limiting
                time.sleep(self.rate_limit_delay)

                return successful_data, failed_symbols

            except requests.exceptions.Timeout as e:
                self.logger.warning(f"Timeout for batch {symbols[:3]}..., attempt {attempt + 1}")
                if attempt == self.max_retries - 1:
                    with self._lock:
                        self._stats['symbols_failed'] += len(symbols)
                    return {}, symbols

            except requests.exceptions.RequestException as e:
                self.logger.error(f"Request error for batch {symbols[:3]}...: {e}")
                if attempt == self.max_retries - 1:
                    with self._lock:
                        self._stats['symbols_failed'] += len(symbols)
                    return {}, symbols

            except Exception as e:
                self.logger.error(f"Unexpected error processing batch {symbols[:3]}...:{e}")
                if attempt == self.max_retries - 1:
                    with self._lock:
                        self._stats['symbols_failed'] += len(symbols)
                    return {}, symbols

            # Wait before retrying with exponential backoff
            if attempt < self.max_retries - 1:
                wait_time = 2 ** attempt
                self.logger.warning(f"Retrying batch in {wait_time}s...")
                time.sleep(wait_time)

        # Safety net, generally should not be reached
        return {}, symbols

    def _retry_failed_symbols(
        self,
        failed_symbols: List[str],
        start_date: str,
        end_date: str
    ) -> Dict[str, pd.DataFrame]:
        """
        Retry failed symbols individually using the single-stock endpoint.
        """
        retry_data = {}

        for symbol in failed_symbols:
            try:
                df = self.get_stock_data(symbol, start_date, end_date)
                retry_data[symbol] = df

                with self._lock:
                    self._stats['symbols_successful'] += 1
                    self._stats['symbols_failed'] -= 1
                    self._stats['total_records_loaded'] += len(df)

                self.logger.info(f"Successfully retrieved {symbol}")

            except Exception as e:
                self.logger.error(f"Final retry failed for {symbol}: {e}")

        return retry_data

    def _process_dataframe(self, df: pd.DataFrame) -> pd.DataFrame:
        """
        Common DataFrame processing logic.
        """
        if df.empty:
            return df

        # Convert date column to datetime
        df['date'] = pd.to_datetime(df['date'])

        # Sort by date to ensure chronological order
        df = df.sort_values('date').reset_index(drop=True)

        # Ensure numeric columns are properly typed
        numeric_columns = [
            'open_price', "high_price", "low_price", "close_price",
            'adjusted_close', 'volume', 'dividend amount', 'split_coefficient'
        ]

        for col in numeric_columns:
            if col in df.columns:
                df[col] = pd.to_numeric(df[col], errors='coerce')

        return df

    def _validate_and_clean_data(self, stock_data: Dict[str,pd.DataFrame]) -> Dict[str, pd.DataFrame]:
        """
        Validate data quality and remove problematic records.
        """

        cleaned_data = {}

        for symbol, df in stock_data.items():
            if df.empty:
                self.logger.warning(f"Empty dataframe for {symbol}")
                continue

            # Check for required columns
            required_cols = ['date', 'open_price', 'high_price', 'low_price', 'close_price', 'volume']
            missing_cols = [col for col in required_cols if col not in df.columns]

            if missing_cols:
                self.logger.error(f"Missing rqeuired columns for {symbol}: {missing_cols}")
                continue

            # Remove rows with invalid price data
            initial_count = len(df)
            df = df.dropna(subset=['open_price', 'high_price', 'low_price', 'close_price'])
            df = df[df['volume'] >= 0] # Volume should not be negative

            # Price Validation
            df = df[(df['high_price'] >= df['low_price']) & (df['open_price'] > 0 ) & (df['close_price'] > 0)]

            cleaned_count = len(df)
            if cleaned_count != initial_count:
                self.logger.info(f"Cleaned {symbol}: {initial_count} -> {cleaned_count} records")

            if cleaned_count > 0:
                cleaned_data[symbol] = df
            else:
                self.logger.warning(f"All data filtered out for {symbol}")

        return cleaned_data
    
    def _combine_stock_dataframes(self, stock_data: Dict[str, pd.DataFrame]) -> pd.DataFrame:
        """
        Combine multiple stock DataFrames into a single DataFrame.
        """
        if not stock_data:
            return pd.DataFrame()

        # Concatenate all DataFrames
        combined_df = pd.concat(stock_data.values(), ignore_index=True)
        # Sort by symbol and date for consistency
        combined_df = combined_df.sort_values(['symbol', 'date']).reset_index(drop=True)

        self.logger.info(f"Combined data: {len(combined_df)} total records for {len(stock_data)} symbols")

        return combined_df

    def _log_batch_statistics(self, failed_symbols :List[str]):
        """
        Log comprehensive statistics about the batch operation.
        """
        stats = self._stats.copy()

        success_rate = (stats['symbols_successful'] / stats['symbols_requested']) * 100

        self.logger.info(f"""
        Batch Processing Complete:
        ========================
        Symbols Requested: {stats['symbol_requested']}
        Symbols Successful: {stats['symbols_successful']}
        Symbols Failed: {stats['symbols_failed']}
        Success Rate: {success_rate:..1f}%
        API Calls Made: {stats['api_calls_made']}
        Total Records Loaded: {stats['total_records_loaded']}
        """)

        if failed_symbols:
            self.logger.warning(f"Failed symbols: {failed_symbols}")

# Factory functions for different use cases
def create_training_pipeline() -> StockDataPipeline:
    """
    Create a pipeline optimized for ML training data loading.
    """
    config = BatchConfig(
        batch_size=25, # Conservative for training
        max_workers=2, # Try not to overwhelm the API during training.
        retry_failed=True,
        validate_data=True
    )

    return StockDataPipeline(
        timeout=120, # Longer timeout for training
        batch_config = config
    )

def create_inference_pipeline() -> StockDataPipeline:
    """
    Create a pipeline optimized for real-time inference data loading.
    """
    config = BatchConfig(
        batch_size=10, # Smaller batches for faste response
        max_workers=3, # More concurrent for speed
        retry_failed=True,
        validate_data=False # Skip validation for speed
    )

    return StockDataPipeline(
        timeout=30, # Shorter timeout for inference
        rate_limit_delay=0.05, # Faster rate for inference
        batch_config=config
    )
