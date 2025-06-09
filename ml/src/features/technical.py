import pandas as pd
import pandas_ta as ta
import numpy as np
import requests
from typing import List

def add_technical_indicators(df: pd.DataFrame) -> pd.DataFrame:
    # Work with a copy of the original dataframe
    result_df = df.copy()
    
    if 'open_price' in df.columns:
        column_mapping = {
            'open_price': 'open',
            'high_price': 'high',
            'low_price': 'low',
            'close_price': 'close'
        }
        result_df = result_df.rename(columns=column_mapping)

    # Check if the required columns are in the df.
    required = ['open', 'high', 'low', 'close', 'volume']
    missing = [col for col in required if col not in result_df.columns]
    if missing:
        raise ValueError(f"Missing required columns: {missing}")

    # Technical Indicators

    """
    1. Moving Averages
    EMA = Exponential Moving Average
    SMA = Simple Moving Average
    """
    result_df['ema_12'] = ta.ema(result_df['close'], length=20)
    result_df['ema_26'] = ta.ema(result_df['close'], length=26)
    result_df['sma_20'] = ta.sma(result_df['close'], length=20)
    result_df['sma_50'] = ta.sma(result_df['close'], length=50)

    """
    2. Momentum Inidicators
    RSI = Relative Strength Index
    """
    result_df['rsi_14'] = ta.rsi(result_df['close'], length=14)
    result_df['rsi_30'] = ta.rsi(result_df['close'], length=30)

    """
    3. Volume Indiciators
    """
    result_df['volume_sma'] = ta.sma(result_df['volume'], length=20)
    result_df['volume_ratio'] = result_df['volume'] / result_df['volume_sma']

    """
    4. MACD
    MACD = Moving Average Convergence Divergence
    This is a momentum indicator uses two EMAs (typically a 12-day and a 26-day) to identify price trends and potential reversals.
    MACD Line is the difference between the 12-day EMA and the 26-day EMA. Above 0 = bullish, below 0 = bearish.
    Signal line represents buy and sell signals.
    """
    macd_data = ta.macd(result_df['close'])
    if macd_data is not None:
        result_df['macd'] = macd_data['MACD_12_26_9']
        result_df['macd_signal'] = macd_data['MACDs_12_26_9']
        result_df['macd_histogram'] = macd_data['MACDh_12_26_9']

    """
    5. Bollinger Bands
    """
    bb_data = ta.bbands(result_df['close'], length=20)
    if bb_data is not None:
        result_df['bb_lower'] = bb_data['BBL_20_2.0']
        result_df['bb_middle'] = bb_data['BBM_20_2.0']
        result_df['bb_upper'] = bb_data['BBU_20_2.0']
        result_df['bb_width'] = (bb_data['BBU_20_2.0'] - bb_data['BBL_20_2.0']) / bb_data['BBM_20_2.0']
        result_df['bb_position'] = (result_df['close'] - bb_data['BBL_20_2.0']) / (bb_data['BBU_20_2.0'] - bb_data['BBL_20_2.0'])

    """
    6. Volatility
    """
    result_df['atr'] = ta.atr(result_df['high'], result_df['low'], result_df['close'], length=14)

    """
    7. Price returns
    """
    result_df['returns_1d'] = result_df['close'].pct_change()
    result_df['returns_5d'] = result_df['close'].pct_change(5)
    
    return result_df
