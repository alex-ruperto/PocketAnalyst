import pytest
import requests
import logging
import pandas as pd
from pocketanalyst_ml.data.pipelines import (
    StockDataPipeline,
    BatchConfig,
    StockDataError,
    create_inference_pipeline,
    create_training_pipeline,
)

logger = logging.getLogger(__name__)


class TestStockDataLoaders:
    """
    Integration tests for enhanced multi-stock data loaders.
    """

    @pytest.fixture
    def training_pipeline(self):
        """Create a training pipeline for tests."""
        return create_training_pipeline()

    @pytest.fixture
    def inference_pipeline(self):
        """Create an inference pipeline for tests."""
        return create_inference_pipeline()

    @pytest.fixture
    def test_symbols(self):
        """
        Common test symbols
        """
        return ["AAPL", "NVDA", "MSFT"]

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

    @pytest.mark.integration
    def test_get_distinct_symbols(self, training_pipeline):
        """
        Test fetching all available symbols
        """
        try:
            symbols = training_pipeline.get_distinct_symbols()

            assert isinstance(symbols, list), "Symbols should be returned as a list"
            assert len(symbols) > 0, "Should have at least some symbols available"

            # Check that symbols look reasonable
            for symbol in symbols[:5]:  # Check first few
                assert isinstance(symbol, str), (
                    f"Symbol should be string, got {type(symbol)}"
                )
                assert len(symbol) <= 5, f"Symbol {symbol} seems too long"
                assert symbol.replace(".", "").isalnum(), (
                    f"Symbol {symbol} contains invalid characters"
                )

            logger.info(f"✓ Successfully fetched {len(symbols)} symbols")
            logger.info(f"Symbols: {symbols}")

        except StockDataError as e:
            pytest.fail(f"Failed to fetch symbols: {e}")
