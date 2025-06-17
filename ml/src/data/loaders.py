"""
Loading of data via API requests.
Handles GET requests from PostgreSQL database/Go API, retrieving data in JSON format.
"""

import pytest
import requests
import logging
import pandas as pd
import os

class StockDataError(Exception):
    """
    Custom exception for stock data retrieval errors.
    """
    pass

class StockDataLoader:
    def __init__(self, api_base_url: str, timeout: int = 5):
        self.api_base_url = self.api_base_url
        self.timeout = timeout
        self.logger = logging.getLogger(__name__)

    def api_base_url(self):
        """
        Provides the base URL for making requests.
        """
        return os.getenv("API_BASE_URL", "http://localhost:8080/api/stocks")

    def GetStockData(self, symbol, start_date, end_date) -> pd.DataFrame:
        url = f"{self.api_base_url}/get?symbol={symbol}&start_date={start_date}&end_date={end_date}"

        try:
            response = requests.get(url, timeout=self.timeout)
            response.raise_for_status() # Raise HTTPError for bad status codes

            data = response.json()
            if not data:
                raise StockDataError(f"No data returned for symbol {symbol}")

            df = pd.DataFrame(data)
            if df.empty:
                raise StockDataError(f"Empty dataset returned for symbol {symbol}")

            return df

        except requests.exceptions.Timeout:
            self.logger.error(f"Timeout occured while fetching data for {symbol}")
            raise StockDataError(f"Request timeout for symbol {symbol}")

        except requests.exceptions.ConnectionError:
            self.logger.error(f"Connection error while fetching data for {symbol}")
            raise StockDataError(f"Connection error for symbol {symbol}")
            
        except requests.exceptions.HTTPError as e:
            self.logger.error(f"HTTP error {e.response.status_code} for {symbol}")
            raise StockDataError(f"API returned error {e.response.status_code} for symbol {symbol}")
            
        except (ValueError, KeyError) as e:
            self.logger.error(f"Invalid JSON response for {symbol}: {e}")
            raise StockDataError(f"Invalid response format for symbol {symbol}")

