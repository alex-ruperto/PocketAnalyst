"""
Data preprocessing for different ML models.
Handles scaling, sequence creation, and train/test splitting.
"""

import logging
import pandas as pd
from datetime import datetime, timedelta
from typing import List, Dict, Optional
from .pipelines import (
    StockDataError,
    create_training_pipeline,
    create_inference_pipeline,
)
from ..features import add_technical_indicators


class MLTrainingPreprocessor:
    """
    Complete ML training data pipeline that efficiently loads multiple stocks 
    and prepares them for model training with multi-domain support.
    
    Supports PocketAnalyst's four-domain analysis:
    - Technical Analysis (1-7 days): Price patterns, momentum, volume
    - Sentiment Analysis (1-30 days): News, social media sentiment
    - Fundamental Analysis (30-365+ days): Financial metrics, ratios
    - Macroeconomic Analysis (7-90 days): Economic indicators, rates
    """
    
    def __init__(self, lookback_days: int = 365 * 2):
        """
        Initialize MLTrainingPreprocessor.

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
        self, symbols: Optional[List[str]] = None, min_records_per_symbol: int = 100
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
                self.logger.info("No symbols specified, loading all available symbols...")
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
            for symbol, df in stock_data.items():
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
        self, 
        stock_data: Dict[str, pd.DataFrame], 
        feature_domains: Optional[List[str]] = None,
        sentiment_data: Optional[Dict[str, pd.DataFrame]] = None,
        fundamental_data: Optional[Dict[str, pd.DataFrame]] = None,
        macro_data: Optional[pd.DataFrame] = None
    ) -> Dict[str, pd.DataFrame]:
        """
        Add multi-domain features for comprehensive analysis.

        MULTI-DOMAIN APPROACH:
        - Technical: Price patterns, momentum indicators, volume analysis
        - Sentiment: News sentiment, social media buzz, analyst opinions
        - Fundamental: Financial metrics, valuation ratios, business performance  
        - Macro: Interest rates, economic indicators, sector rotation

        Args:
            stock_data: Dictionary of stock DataFrames
            feature_domains: List of domains to include ['technical', 'sentiment', 'fundamental', 'macro']
            sentiment_data: Dictionary of sentiment DataFrames {symbol: sentiment_df}
            fundamental_data: Dictionary of fundamental DataFrames {symbol: fundamental_df}
            macro_data: DataFrame with macroeconomic indicators

        Returns:
            Dictionary of DataFrames with multi-domain features added
        """
        # Set default feature domains if none provided
        if feature_domains is None:
            feature_domains = ['technical']

        prepared_data = {}

        # Validate domain requirements
        self._validate_domain_requirements(feature_domains, sentiment_data, fundamental_data, macro_data)

        for symbol, df in stock_data.items():
            try:
                self.logger.debug(f"Preparing multi-domain features for {symbol}")

                 # Start with a clean copy
                feature_df = df.copy()

                # TECHNICAL ANALYSIS DOMAIN (1-7 days)
                if 'technical' in feature_domains:
                    self.logger.debug(f"Adding technical indicators for {symbol}")
                    feature_df = add_technical_indicators(feature_df)
                    
                    # Add additional technical features
                    feature_df = self._add_price_features(feature_df)
                    feature_df = self._add_volume_features(feature_df)
                    feature_df = self._add_temporal_features(feature_df)
                
                # SENTIMENT ANALYSIS DOMAIN (1-30 days)
                if 'sentiment' in feature_domains:
                    if sentiment_data and symbol in sentiment_data:
                        # TODO: Implement sentiment feature engineering
                        self.logger.info(f"📰 Sentiment features would be added for {symbol}")
                    else:
                        self.logger.warning(f"Sentiment data not available for {symbol}")

                # FUNDAMENTAL ANALYSIS DOMAIN (30-365+ days)
                if 'fundamental' in feature_domains:
                    if fundamental_data and symbol in fundamental_data:
                        # TODO: Implement fundamental feature engineering
                        self.logger.info(f"🏢 Fundamental features would be added for {symbol}")
                    else:
                        self.logger.warning(f"Fundamental data not available for {symbol}")

                # MACROECONOMIC ANALYSIS DOMAIN (7-90 days)
                if 'macro' in feature_domains:
                    if macro_data is not None:
                        # TODO: Implement macro feature engineering
                        self.logger.info(f"🌍 Macroeconomic features would be added for {symbol}")
                    else:
                        self.logger.warning("Macroeconomic data not available")

                # Handle missing values from feature engineering
                initial_rows = len(feature_df)
                feature_df = feature_df.dropna()
                final_rows = len(feature_df)

                # Quality check
                if final_rows < initial_rows * 0.7:  # More than 30% data loss
                    self.logger.warning(
                        f"{symbol}: High data loss after feature engineering "
                        f"({initial_rows} -> {final_rows} rows)"
                    )
                
                # Minimum samples required varies by domain
                min_samples_required = 50
                if 'fundamental' in feature_domains:
                    min_samples_required = 200  # Fundamental analysis needs more history

                if final_rows > min_samples_required:
                    prepared_data[symbol] = feature_df
                    self.logger.debug(f"Successfully prepared {final_rows} multi-domain feature records for {symbol}")
                else:
                    self.logger.warning(f"Too few records after feature engineering for {symbol}, skipping")

            except Exception as e:
                self.logger.error(f"Failed to prepare features for {symbol}: {e}")
                continue

        return prepared_data

    def create_training_targets(
        self,
        feature_data: Dict[str, pd.DataFrame],
        prediction_horizons: Optional[Dict[str, List[int]]] = None,
    ) -> Dict[str, pd.DataFrame]:
        """
        Create training targets for different prediction horizons.

        Args:
            feature_data: Dictionary of feature DataFrames
            prediction_horizons: List of days ahead to predict

        Returns:
            Dictionary of DataFrames with targets added
        """
        # Set default prediction horizons if not provided
        if prediction_horizons is None:
            prediction_horizons = {
                'technical': [1, 3, 7],        # Technical analysis: 1-7 days
                'sentiment': [1, 7, 30],       # Sentiment analysis: 1-30 days  
                'fundamental': [30, 90, 365],  # Fundamental analysis: 30-365+ days
                'macro': [7, 30, 90]          # Macroeconomic analysis: 7-90 days
            }

        training_data = {}

        for symbol, df in feature_data.items():
            try:
                target_df = df.copy()

                # Create targets for prediction horizon
                for domain, horizons in prediction_horizons.items():
                    for horizon in horizons:
                        # Future return target (regression)
                        target_col = f"target_return_{domain}_{horizon}d"
                        target_df[target_col] = (
                            target_df["close_price"].shift(-horizon) / target_df["close_price"] - 1
                        )

                        # Binary direction target (classification)
                        direction_col = f"target_direction_{domain}_{horizon}d"
                        target_df[direction_col] = (target_df[target_col] > 0).astype(int)

                        # Categorical return buckets (multi-class classification)
                        bucket_col = f"target_bucket_{domain}_{horizon}d"

                        # Use domain-specific bucket ranges
                        if domain == 'technical':
                            bins = [-float('inf'), -0.03, -0.01, 0.01, 0.03, float('inf')]
                        elif domain == 'sentiment':
                            bins = [-float('inf'), -0.05, -0.02, 0.02, 0.05, float('inf')]
                        elif domain == 'fundamental':
                            bins = [-float('inf'), -0.10, -0.05, 0.05, 0.15, float('inf')]
                        elif domain == 'macro':
                            bins = [-float('inf'), -0.08, -0.03, 0.03, 0.08, float('inf')]
                        else:
                            bins = [-float('inf'), -0.05, -0.02, 0.02, 0.05, float('inf')]

                        target_df[bucket_col] = pd.cut(
                            target_df[target_col],
                            bins=bins,
                            labels=["strong_down", "down", "flat", "up", "strong_up"],
                        )
                # Remove rows where we can't calculate targets
                max_horizon = max([h for horizons in prediction_horizons.values() for h in horizons])
                target_df = target_df.iloc[:-max_horizon]

                # Only keep if we have enough data for training
                if len(target_df) > 100:
                    training_data[symbol] = target_df
                    self.logger.debug(f"Created multi-domain targets for {symbol}: {len(target_df)} training samples")
                else:
                    self.logger.warning(f"Insufficient data for targets in {symbol}")

            except Exception as e:
                self.logger.warning(f"Failed to create targets for {symbol}: {e}")
                continue

        self.logger.info(f"Successfully created multi-domain training targets for {len(training_data)} symbols")
        return training_data

    def get_combined_training_dataset(
        self,
        symbols: Optional[List[str]] = None,
        feature_domains: Optional[List[str]] = None,
        prediction_horizons: Optional[Dict[str, List[int]]] = None,
        sample_limit: Optional[int] = None,
        sentiment_data: Optional[Dict[str, pd.DataFrame]] = None,
        fundamental_data: Optional[Dict[str, pd.DataFrame]] = None,
        macro_data: Optional[pd.DataFrame] = None
    ) -> pd.DataFrame:
        """
        Get a combined training dataset with multi-domain features for model training.
        
        This is the main public method that orchestrates the entire multi-domain pipeline:
        1. Loads data for specified symbols
        2. Prepares multi-domain features (technical, sentiment, fundamental, macro)
        3. Creates training targets for different domains and horizons
        4. Applies sampling for balanced datasets
        5. Combines all symbols into a single training DataFrame
        
        MULTI-DOMAIN APPROACH:
        - Technical: Price patterns, momentum (1-7 days)
        - Sentiment: News/social sentiment (1-30 days)
        - Fundamental: Financial metrics (30-365+ days)
        - Macro: Economic indicators (7-90 days)
        
        Args:
            symbols: Specific symbols to include. If None, uses all available symbols
            feature_domains: List of domains to include ['technical', 'sentiment', 'fundamental', 'macro']
            prediction_horizons: Dict mapping domains to prediction horizons
                               If None, uses domain-appropriate defaults
            sample_limit: Maximum number of samples per symbol (for balancing)
            sentiment_data: Dictionary of sentiment DataFrames {symbol: sentiment_df}
            fundamental_data: Dictionary of fundamental DataFrames {symbol: fundamental_df}
            macro_data: DataFrame with macroeconomic indicators
                         
        Returns:
            Combined DataFrame ready for ML training with multi-domain features and targets
        """        
        # Set defaults
        if feature_domains is None:
            feature_domains = ['technical']

        self.logger.info("Starting complete multi-domain training dataset preparation pipeline")
        self.logger.info(f"Requested domains: {feature_domains}")

        # STEP 1: Load raw stock data
        self.logger.info("Step 1/4: Loading stock data...")
        stock_data = self.load_training_data(symbols)

        # STEP 2: Prepare multi-domain features
        self.logger.info("Step 2/4: Preparing multi-domain features...")
        feature_data = self.prepare_features(
            stock_data,
            feature_domains=feature_domains,
            sentiment_data=sentiment_data,
            fundamental_data=fundamental_data,
            macro_data=macro_data
        )

        # Step 3: Create multi-domain training targets
        self.logger.info("Step 3/4: Creating multi-domain training targets...")
        training_data = self.create_training_targets(feature_data, prediction_horizons)

        if not training_data:
            raise StockDataError("No training data available after processing")

        # Step 4: Apply sampling if requested 
        if sample_limit:
            self.logger.info(f"Step 4/4: Applying sampling (limit: {sample_limit} per symbol)...")
            training_data = self._apply_sampling(training_data, sample_limit)
        else:
            self.logger.info("Step 4/4: No sampling applied (using all available data)")

        # Step 5: Combine all symbols into single DF
        combined_df = pd.concat(training_data.values(), ignore_index=True)
        combined_df = combined_df.sort_values(['symbol', 'date']).reset_index(drop=True)

        # Log results
        unique_symbols = combined_df['symbol'].nunique()
        feature_cols = [col for col in combined_df.columns if not col.startswith('target_') and col not in ['symbol', 'date']]
        target_cols = [col for col in combined_df.columns if col.startswith('target_')]

        self.logger.info(f"Multi-domain training dataset creation complete!")
        self.logger.info(f"Dataset size: {len(combined_df):,} samples from {unique_symbols} symbols")
        self.logger.info(f"Features: {len(feature_cols)} feature columns")
        self.logger.info(f"Targets: {len(target_cols)} target columns")

        return combined_df

    def _validate_domain_requirements(
        self, 
        feature_domains: List[str], 
        sentiment_data: Optional[Dict[str, pd.DataFrame]], 
        fundamental_data: Optional[Dict[str, pd.DataFrame]], 
        macro_data: Optional[pd.DataFrame]
    ) -> None:
        """Validate that required data is available for requested domains."""
        missing_requirements = []
        
        if 'sentiment' in feature_domains and sentiment_data is None:
            missing_requirements.append("sentiment_data required for sentiment domain")
        
        if 'fundamental' in feature_domains and fundamental_data is None:
            missing_requirements.append("fundamental_data required for fundamental domain")
        
        if 'macro' in feature_domains and macro_data is None:
            missing_requirements.append("macro_data required for macro domain")
        
        if missing_requirements:
            self.logger.warning("Missing data requirements:")
            for req in missing_requirements:
                self.logger.warning(f"  - {req}")

    def _apply_sampling(
        self,
        training_data: Dict[str, pd.DataFrame],
        sample_limit: int
    ) -> Dict[str, pd.DataFrame]:
        """
        Apply intelligent sampling to balance the dataset.
        
        SAMPLING STRATEGY:
        - Take 50% from recent data (more relevant to current market conditions)
        - Take 50% randomly from historical data (preserve long-term patterns)
        - Maintains chronological order for time-series consistency
        """
        sampled_data = {}
        
        for symbol, df in training_data.items():
            if len(df) > sample_limit:
                # Calculate sampling split
                recent_samples = min(sample_limit // 2, len(df) // 4)
                random_samples = sample_limit - recent_samples
                
                # Get recent data (most recent records)
                recent_data = df.tail(recent_samples)
                
                # Get random sample from historical data (excluding recent portion)
                historical_data = df.head(-recent_samples)
                if len(historical_data) >= random_samples:
                    random_data = historical_data.sample(n=random_samples, random_state=42)
                else:
                    random_data = historical_data
                
                # Combine and sort
                sampled_df = pd.concat([random_data, recent_data]).sort_values('date')
                sampled_data[symbol] = sampled_df
                
                self.logger.debug(f"Sampled {symbol}: {len(df)} -> {len(sampled_df)} records")
            else:
                sampled_data[symbol] = df
        
        return sampled_data
