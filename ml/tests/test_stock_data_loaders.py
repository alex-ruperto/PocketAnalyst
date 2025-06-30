import pytest
import requests
import logging
import pandas as pd
from pocketanalyst_ml.data.loaders import (
    StockDataPipeline,
    BatchConfig,
    StockDataError
)

logger = logging.getLogger(__name__)

class TestStockDataLoaders:
    """
    Integration tests for enhanced multi-stock data loaders.
    """

    @pytest.fixture
    def test_symbols(self):
        """
        Common test symbols
        """
        return ['AAPL', 'NVDA', 'MSFT']

    @pytest.mark.integration
    def test_api_connection(self, api_base_url):
        """
        Test that we can connect to our Go backend API.
        """
        try:
            response = requests.get(f"{api_base_url}/health", timeout=5)
            assert response.status_code == 200, "API health check failed"
            logger.info("API connection successful")
        except requests.exceptions.ConnectionError:
            pytest.skip("Go backend not running - skipping integration tests.")
