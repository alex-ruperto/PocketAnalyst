"""
Integration test for technical indicators with Go backend API.
"""

import pytest
import requests
import pandas as pd
from src.features import add_technical_indicators

class TestTechnicalIndicators:
    """
    Integration tests for technical indicators with Go API backend
    """

    BASE_URL = "http://localhost:8080/api/stocks"

    def test_api_connection(self):
        """
        Test that we can connect to the Go backend API.
        """
        try:
            response = requests.get(f"{self.BASE_URL}/health", timeout=5)
            assert response.status_code == 200, "API health check failed"
        except requests.exceptions.ConnectionError:
            pytest.skip("Go backend not running — skipping integration tests.")

