"""
Loading of data via API requests.
Handles GET requests from PostgreSQL database/Go API, retrieving data in JSON format.
"""

import requests
import logging
import time
import pandas as pd
import os
from typing import List, Optional 

class StockDataError(Exception):
    """
    Custom exception for stock data retrieval errors.
    """
    pass

class StockDataPipeline:
    """"""
    def __init__(self, api_base_url: Optional[str] = None, timeout: int = 30, max_retries: int = 3, rate_limit_delay: float = 0.1):
        self.api_base_url = api_base_url or os.getenv("API_BASE_URL", "http://localhost:8080/api/stocks")
        self.timeout = timeout
        self.max_retries = max_retries
        self.rate_limit_delay = rate_limit_delay
        self.logger = logging.getLogger(__name__)

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



