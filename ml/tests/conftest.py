"""
Test configuration and shared fixtures.
Path handling is done by pyproject.toml and editable installation.
"""
import pytest
import logging
import os


@pytest.fixture
def logger():
    """Provide a logger instance for tests."""
    return logging.getLogger("test_logger")


@pytest.fixture(scope="session")
def api_base_url():
    """Provide base URL for API integration tests."""
    return os.getenv("API_BASE_URL", "http://localhost:8080/api/stocks")
