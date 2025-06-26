# PocketAnalyst
*The honest approach to ML in finance.*

## What is PocketAnalyst?

PocketAnalyst is a probabilistic stock analysis tool that tells you what you actually need to know: **the probability of different price movements at different time horizons**, based on historical patterns.

Instead of claiming to predict the future with false precision, PocketAnalyst gives you honest, probabilistic insights:

- *"Based on technical patterns, there's a 35% chance AAPL gains 5%+ in the next 7 days"*
- *"Based on fundamental metrics, there's a 60% chance of 15%+ gains over the next 6 months"*
- *"Based on recent sentiment, there's a 25% chance of 3%+ downside in the next week"*

## The Philosophy

### What We Believe
- **Markets are uncertain** - anyone claiming high accuracy is probably lying
- **Different factors matter at different timeframes** - technical analysis for days, fundamentals for months
- **Probabilities are more honest than predictions** - 30% chance of gain is better than "will go up 2.3%"
- **Humans should make decisions** - we provide information, you decide what to do

### What We Don't Do
- ❌ Automate trading decisions
- ❌ Promise stunning accuracy  
- ❌ Replace human judgment
- ❌ Guarantee profits

### What We Do
- ✅ Analyze historical patterns across multiple domains
- ✅ Provide probabilistic insights for different time horizons
- ✅ Acknowledge uncertainty honestly
- ✅ Help you make better-informed decisions

## How It Works

### Multi-Domain Analysis

PocketAnalyst analyzes four different domains of market information:

#### 🔧 Technical Analysis (1-7 days)
- **What**: Price patterns, momentum indicators, volume analysis
- **Good for**: Short-term trading opportunities, entry/exit timing
- **Example**: *"Chart patterns suggest 40% chance of +3% move in next 3 days"*

#### 📰 Sentiment Analysis (1-30 days)  
- **What**: News sentiment, social media buzz, analyst opinions
- **Good for**: Event-driven moves, earnings reactions
- **Example**: *"Positive news flow indicates 45% chance of +5% gain over next 2 weeks"*

#### 🏢 Fundamental Analysis (30-365+ days)
- **What**: Financial metrics, valuation ratios, business performance  
- **Good for**: Long-term investment decisions, value identification
- **Example**: *"Valuation metrics suggest 65% chance of +20% gain over next year"*

#### 🌍 Macroeconomic Analysis (7-90 days)
- **What**: Interest rates, economic indicators, sector rotation
- **Good for**: Market timing, sector allocation, risk management
- **Example**: *"Economic conditions favor 55% chance of +10% sector outperformance over next quarter"*

### Nested Probability Targets

For each domain and time horizon, we calculate probabilities of hitting different return thresholds:

```
7-Day Technical Analysis for AAPL:
├── 📈 Upside Probabilities
│   ├── +1% gain: 45% chance
│   ├── +3% gain: 25% chance  
│   ├── +5% gain: 12% chance
│   └── +10% gain: 3% chance
└── 📉 Downside Probabilities
    ├── -1% loss: 35% chance
    ├── -3% loss: 18% chance
    ├── -5% loss: 8% chance
    └── -10% loss: 2% chance
```

### Adaptive Learning

- **Weekly retraining** keeps models current with market conditions
- **Multiple time windows** balance adaptation with stability  
- **Performance tracking** monitors prediction calibration over time
- **Regime detection** flags when markets behave unusually

## Technology Stack

### Backend (Go)
- **REST API** for data retrieval and storage
- **PostgreSQL** for historical stock data and predictions
- **Multi-provider data ingestion** (FMP, Alpha Vantage, others)
- **Robust error handling** and rate limiting

### ML Pipeline (Python)
- **Multi-domain feature engineering** using pandas and pandas_ta
- **Ensemble models** combining Random Forest, XGBoost, and linear methods
- **Probabilistic predictions** with proper uncertainty quantification
- **Automated retraining** with performance validation

### Data Sources
- **Stock prices**: Financial Modeling Prep, Alpha Vantage
- **Fundamentals**: Quarterly earnings and financial statements  
- **Sentiment**: News APIs and social media feeds (planned)
- **Macro**: Economic indicators and Fed data (planned)

## Quick Start

### Prerequisites
- Go 1.24+
- Python 3.11+
- PostgreSQL 12+
- 16GB+ RAM (48GB recommended for training)

### Setup
```bash
# Clone the repository
git clone https://github.com/yourusername/PocketAnalyst.git
cd PocketAnalyst

# Set up the database
psql -f database/schema.sql

# Configure environment variables
cp .env.example .env
# Edit .env with your API keys and database URL

# Start the Go API
cd api
go run server/main.go

# Set up Python environment
cd ../ml
pip install -r requirements.txt

# Run initial training
python -m training.train_models
```

### Basic Usage
```bash
# Fetch data for a symbol
curl "http://localhost:8080/api/stocks/fetch?symbol=AAPL"

# Get predictions
curl "http://localhost:8080/api/ml/predictions?symbol=AAPL"

# Trigger model retraining
python -m training.retrain_models
```

## Example Output

```json
{
  "symbol": "AAPL",
  "as_of": "2025-01-20T10:00:00Z",
  "predictions": {
    "technical_3d": {
      "upside_probabilities": {
        "gain_1pct": 0.42,
        "gain_3pct": 0.23,
        "gain_5pct": 0.11
      },
      "downside_probabilities": {
        "loss_1pct": 0.31,
        "loss_3pct": 0.16,
        "loss_5pct": 0.07
      },
      "confidence": 0.68
    },
    "fundamental_90d": {
      "upside_probabilities": {
        "gain_5pct": 0.58,
        "gain_10pct": 0.34,
        "gain_20pct": 0.19
      },
      "downside_probabilities": {
        "loss_5pct": 0.22,
        "loss_10pct": 0.09,
        "loss_20pct": 0.04
      },
      "confidence": 0.71
    }
  }
}
```

## Roadmap

### Phase 1: Core Technical Analysis ✅
- Multi-stock data pipeline
- Technical indicator features
- Ensemble prediction models
- Basic probability outputs

### Phase 2: Multi-Domain Expansion 🚧  
- Fundamental analysis integration
- Sentiment analysis pipeline
- Macroeconomic indicators
- Domain-specific model optimization

### Phase 3: Advanced Features 📋
- Real-time prediction serving
- Portfolio-level analysis
- Advanced visualization dashboard
- Mobile app for predictions

### Phase 4: Scale & Polish 📋
- Cloud deployment infrastructure
- High-frequency retraining
- Advanced ensemble techniques
- Professional API for institutions

## Contributing

We welcome contributions that align with our philosophy of honest, probabilistic analysis:

- **Bug fixes** and performance improvements
- **New data sources** for any of the four domains
- **Model improvements** that enhance calibration
- **Visualization tools** for probability displays
- **Documentation** and educational content

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Disclaimer

PocketAnalyst is a research and educational tool. All predictions are probabilistic and based on historical patterns that may not repeat. Past performance does not guarantee future results. Always do your own research and consider your risk tolerance before making investment decisions.

**This is not financial advice.**

---

*Built with the belief that honest uncertainty is better than false confidence.*
