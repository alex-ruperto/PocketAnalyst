"""
Data preprocessing for different ML models.
Handles scaling, sequence creation, and train/test splitting.
"""

import logging
import pandas as pd
from datetime import datetime, timedelta
from typing import List, Dict
from pocketanalyst_ml.data.loaders import create_training_pipeline, StockDataError
from pocketanalyst_ml.features import add_technical_indicators


class TechnicalDataPreprocessor
"""
Complete ML training data pipeline that efficiently loads multiple stocks and prepares them for model training.
"""

def __init__(self, lookback_days: int = 365 * 2):
    """
    Initialize TechnicalDataPreprocessor.

    Args:
        lookback_days: Number of days of historical data for training
    """
    self.lookback_days = lookback_days
    self.data_loader = create_training_pipeline()
    self.logger = logging.getLogger(__name__)

    # Calculate data range for training data
    self.end_date = datetime.now().strftime('%Y-%m-%d')
    self.start_date = (datetime.now() - timedelta(days=lookback_days)).strftime('%Y-%m-%d')


