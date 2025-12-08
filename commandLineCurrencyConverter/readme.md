# Simple Go Lang Projects

A collection of command-line utilities built with Go to practice core concepts and API integration.

---

## 📁 Projects

### 1. Command Line Currency Converter

A real-time currency converter that fetches live exchange rates and performs conversions between 160+ currencies.

**Features:**
- 💰 Convert between 160+ global currencies
- ⚡ Smart caching system with disk persistence
- 🔄 Automatic retry logic with exponential backoff
- 📊 Batch conversion mode for multiple conversions
- ⚙️ Configurable cache expiry via flags
- ✅ Input validation with smart suggestions
- 💾 Persistent cache survives program restarts

**Usage:**

```bash
# Standard single conversion
go run commandLineCurrencyConverter/commandLineCurrencyConverter.go

# Batch mode for multiple conversions
go run commandLineCurrencyConverter/commandLineCurrencyConverter.go -batch

# Custom cache expiry (in minutes)
go run commandLineCurrencyConverter/commandLineCurrencyConverter.go -cache=60

# Combine flags
go run commandLineCurrencyConverter/commandLineCurrencyConverter.go -batch -cache=120
```

**Technical Highlights:**
- RESTful API integration (ExchangeRate-API)
- JSON marshaling/unmarshaling
- File I/O for cache persistence
- Command-line flag parsing
- Error handling with retry logic
- Struct tags and custom types

---

### 2. Password Strength Calculator

*(Add description when implemented)*

---

## 🚀 Getting Started

**Prerequisites:**
- Go 1.16 or higher

**Installation:**

```bash
git clone https://github.com/SShogun/simple-Go-Lang-projs.git
cd simple-Go-Lang-projs
```

**Run a project:**

```bash
cd commandLineCurrencyConverter
go run commandLineCurrencyConverter.go
```

---

## 📚 Learning Goals

- API integration and HTTP requests
- JSON data handling
- Caching strategies
- Error handling patterns
- Command-line interfaces
- File I/O operations
- Concurrency (future)

---

## 🛠️ Technologies

- **Language:** Go 1.x
- **APIs:** ExchangeRate-API
- **Tools:** Standard library (`net/http`, `encoding/json`, `flag`)

---

## 📝 License

This is a personal learning project. Feel free to use and modify as needed.

---

## 🤝 Contributing

These are personal practice projects, but suggestions and improvements are welcome!
