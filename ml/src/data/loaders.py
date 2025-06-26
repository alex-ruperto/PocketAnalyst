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
        url = f"{self.api_base_url}/get"
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




