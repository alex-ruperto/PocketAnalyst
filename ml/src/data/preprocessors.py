"""
Data preprocessing for different ML models.
Handles scaling, sequence creation, and train/test splitting.
"""

import logging
import pandas as pd
from datetime import datetime, timedelta
from typing import List, Dict
from pocketanalyst_ml.data.pipelines import create_training_pipeline, StockDataError
from pocketanalyst_ml.features import add_technical_indicators


class TechnicalDataPreprocessor:
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
        self.end_date = datetime.now().strftime("%Y-%m-%d")
        self.start_date = (datetime.now() - timedelta(days=lookback_days)).strftime(
            "%Y-%m-%d"
        )

        self.logger.info(
            f"Training pipeline initalized for date range: {self.start_date} to {self.end_date}"
        )

    def load_training_data(
        self, symbols: List[str] = None, min_records_per_symbol: int = 100
    ) -> Dict[str, pd.DataFrame]:
        """
        Load and prepare training data for multiple stocks

        Args:
            symbols: List of symbols to load. If None, loads all available symbols
            min_records_per_symbol: Minimum number of records required per symbol

        Returns:
            Dictionary mapping symbols to DataFrames ready for training
        """
        try:
            # Get symbols to process
            if symbols is None:
                self.logger.info(
                    "No symbols specified, loading all available symbols..."
                )
                symbols = self.data_loader.get_distinct_symbols()
                self.logger.info(f"Found {len(symbols)} available symbols.")

            # Filter out symbols that are likely to cause issues
            symbols = self._filter_symbols(symbols)

            # Load data using batch processing
            self.logger.info(f"Loading data for {len(symbols)} symbols...")
            stock_data = self.data_loader.get_multiple_stocks_data(
                symbols=symbols,
                start_date=self.start_date,
                end_date=self.end_date,
                return_combined=False,
            )

            # Filter symbols with insufficient data
            filtered_data = {}
            for (
                symbol,
                df,
            ) in stock_data.items():
                if len(df) >= min_records_per_symbol:
                    filtered_data[symbol] = df
                else:
                    self.logger.warning(
                        f"Symbol {symbol} has only {len(df)} records, skipping."
                    )

            self.logger.info(
                f"Successfully loaded data for {len(filtered_data)} symbols."
            )
            return filtered_data

        except StockDataError as e:
            self.logger.error(f"Failed to load training data {e}")
            raise
        except Exception as e:
            self.logger.error(f"Unexpected error loading training data: {e}")
            raise

    def prepare_features(
        self, stock_data: Dict[str, pd.DataFrame], add_technical: bool = True
    ) -> Dict[str, pd.DataFrame]:
        """
        Add technical indicators and prepare features for all stocks

        Args:
            stock_data: Dictionary of stock DataFrames
            add_technical: Boolean deciding whether to add technical indicators. Default = True
        """
        prepared_data = {}
        for symbol, df in stock_data.items():
            try:
                self.logger.debug(f"Preparing features for {symbol}")

                # Start with a clean copy
                feature_df = df.copy()

                # Add technical indicators if requested
                if add_technical:
                    feature_df = add_technical_indicators(feature_df)

                # Drop rows with NaN values (from technical indicators)
                initial_rows = len(feature_df)
                feature_df = feature_df.dropna()
                final_rows = len(feature_df)

                if final_rows < initial_rows * 0.7:  # More than 30 % data loss
                    self.logger.warning(
                        f"{symbol}: High data loss after feature engineering "
                        f"({initial_rows} -> {final_rows} rows)"
                    )

                if final_rows > 50:  # Ensure we still have a reasonable amount of data
                    prepared_data[symbol] = feature_df
                    self.logger.debug(
                        f"Successfully prepared {final_rows} feature records for {symbol}"
                    )
                else:
                    self.logger.warning(
                        f"Too few records after feature engineering for {symbol}, skipping"
                    )

            except Exception as e:
                self.logger.error(f"Failed to prepare features for {symbol}: {e}")
                continue

        self.logger.info(
            f"Successfully prepared feature for {len(prepared_data)} symbols"
        )
        return prepared_data

    def create_training_targets(
        self,
        feature_data: Dict[str, pd.DataFrame],
        prediction_horizons: List[int] = [1, 3, 7, 30, 90],
    ) -> Dict[str, pd.DataFrame]:
        """
        Create training targets for different prediction horizons.

        Args:
            feature_data: Dictionary of feature DataFrames
            prediction_horizons: List of days ahead to predict

        Returns:
            Dictionary of DataFrames with targets added
        """
        training_data = {}

        for symbol, df in feature_data.items():
            try:
                target_df = df.copy()
                # Crete targets for prediction horizon
                for horizon in prediction_horizons:
                    """
                    Calculate the percentage return 'horizon' days in the future.
                    """
                    target_df[f"target_return_{horizon}d"] = (
                        target_df["close_price"].shift(-horizon)
                        / target_df["close_price"]
                        - 1
                    )

                    # Binary direction target (up/down)
                    target_df[f"target_direction_{horizon}d"] = (
                        target_df[f"target_return_{horizon}d"]
                    ).astype(int)

                    # Categorical return buckets for classification
                    target_df[f"target_bucket_{horizon}d"] = pd.cut(
                        target_df[f"target_return_{horizon}d"],
                        bins=[-float("inf"), -0.05, -0.02, 0.02, 0.05, float("inf")],
                        labels=["strong_down", "down", "flat", "up", "strong_up"],
                    )

                # Remove rows where we can't calculate targets (end of dataset)
                max_horizon = max(prediction_horizons)
                target_df = target_df.iloc[
                    :-max_horizon
                ]  # Remove last max_horizon rows

                # Only keep if we have enough data for training
                if len(target_df) > 100:
                    training_data[symbol] = target_df
                    self.logger.debug(
                        f"Created targets for {symbol}: {len(target_df)} training samples"
                    )
                else:
                    self.logger.warning(f"Insufficient data for targets in {symbol}")

            except Exception as e:
                self.logger.warning(f"Failed to create targets for {symbol}: {e}")
                continue

        self.logger.info(
            f"Successfully created training targets for {len(training_data)} symbols"
        )
        return training_data

    def get_combined_training_dataset(
        self,
        symbols: List[str] = None,
        prediction_horizons: List[int] = [1, 3, 7, 30, 90],
        sample_limit: int = None,
    ) -> pd.DataFrame:
        """
        Get combined training dataset with all symbols for model training

        Args:
            symbols: Specific symbols to include
            prediction_horizons: Prediction horizons to create targets for (e.g., 1d, 3d, 7d, 30d, 90d)
            sample_limit: Maximum number of samples per symbol (for balancing)

        Returns:
            Combined DataFrame ready for ML training
        """
        # Load and prepare data
        stock_data = self.load_training_data(symbols)
        feature_data = self.prepare_features(stock_data)
        training_data = self.create_training_targets(feature_data, prediction_horizons)

        if not training_data:
            raise StockDataError("No training data available after processing")

        # Apply sampling if requested
        if sample_limit:
            for symbol in training_data:
                df = training_data[symbol]
                if len(df) > sample_limit:
                    # Sample recent data preferentially
                    recent_samples = min(sample_limit // 2, len(df) // 4)
                    random_samples = sample_limit - recent_samples

                    recent_data = df.tail(recent_samples)
                    random_data = df.head(-recent_samples).sample(
                        n=random_samples, random_states=42
                    )

                    training_data[symbol] = pd.concat(
                        [random_data, recent_data]
                    ).sort_values("date")
