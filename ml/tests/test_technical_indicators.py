"""
Integration test for technical indicators with Go backend API.
"""

import pytest
import requests
import logging
import pandas as pd
from pocketanalyst_ml.features import add_technical_indicators

logger = logging.getLogger(__name__)

class TestTechnicalIndicators:
    """
    Integration tests for technical indicators with Go API backend
    """

    @pytest.mark.integration
    def test_api_connection(self, api_base_url ):
        """
        Test that we can connect to the Go backend API.
        """
        try:
            response = requests.get(f"{api_base_url}/health", timeout=5)
            assert response.status_code == 200, "API health check failed"
        except requests.exceptions.ConnectionError:
            pytest.skip("Go backend not running — skipping integration tests.")

    @pytest.mark.integration
    def test_add_technical_indicators(self, api_base_url):
        """
        Test the add technical indicators function
        """
        response = requests.get(f"{api_base_url}/get?symbol=AAPL&start_date=2006-01-02&end_date=2025-01-02", 
                                timeout=5)

        if response.status_code == 200:
            data = response.json()
            df = pd.DataFrame(data)

            # Add technical indicators and store in ta_df, then store it in CSV for data analytics
            ta_df = add_technical_indicators(df)
            ta_df.to_csv('../test-reports/data-with-indicators.csv', index=False)
            logger.info(f"Data with indicators saved to test-reports/data_with_indicators.csv")
            logger.info(f"Final shape: {ta_df.shape}")
            logger.info(f"Final columns: {list(ta_df.columns)}")
        
        else:
            pytest.skip("Status code was not 200. Skipping test.")

