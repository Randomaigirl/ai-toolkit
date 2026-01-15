# 🔥 AI Toolkit - Multi-LLM Development Suite

<div align="center">

[![Made with Love](https://img.shields.io/badge/Made%20with-💜-purple)]()
[![Python](https://img.shields.io/badge/Python-3.11+-blue.svg)](https://www.python.org/)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6.svg)](https://www.typescriptlang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**A comprehensive toolkit for building AI applications across multiple LLM providers**

*Built by Revy (˃ᆺ˂) 💜*

</div>

---

## 🌟 Features

- ⚡ **Multi-LLM Comparison** - Query GPT-4, Claude, Gemini, Llama, DeepSeek, Mistral simultaneously
- 🔄 **Streaming Responses** - Real-time token-by-token streaming
- 💾 **Smart Caching** - Reduce costs with intelligent response caching
- 📊 **Token Tracking** - Accurate token counting and cost estimation
- 🎨 **Beautiful UI** - Cyberpunk-themed interfaces with animations
- 🔐 **Rate Limiting** - Built-in protection against abuse
- 📈 **Metrics & Analytics** - Track usage, costs, and performance

---

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/Randomaigirl/ai-toolkit.git
cd ai-toolkit

# Python setup
pip install -r requirements.txt

# Try the demo
python llm_comparator.py
```

---

## 💻 What's Included

### Python: Async LLM Comparator
```python
from llm_comparator import LLMComparator

comparator = LLMComparator()
results = await comparator.compare_models(
    models=["gpt-4o", "claude-sonnet-4-5"],
    prompt="Explain quantum computing"
)
```

### Go: High-Performance API Gateway
```bash
go run gateway.go
# Server with caching, rate limiting, and metrics
```

### TypeScript/React: Chat Component
```tsx
<AIChat config={{
  apiEndpoint: '/api/llm',
  model: 'gpt-4o',
  streamEnabled: true
}} />
```

### HTML: Standalone Interface
Beautiful cyberpunk-themed chat interface with no build process needed!

---

## 📊 Supported Models

| Provider | Models | Speed | Cost |
|----------|--------|-------|------|
| OpenAI | GPT-4o, GPT-4 | ⚡⚡⚡ | 💰💰💰 |
| Anthropic | Claude Sonnet 4.5 | ⚡⚡⚡ | 💰💰 |
| Google | Gemini 2.0 Flash | ⚡⚡⚡⚡ | 💰 |
| Meta | Llama 3.3 70B | ⚡⚡⚡ | 💰 |
| DeepSeek | DeepSeek V3 | ⚡⚡⚡⚡ | 💰 |
| Mistral | Mistral Large | ⚡⚡⚡ | 💰💰 |

---

## 🎯 Use Cases

1. **Model Selection** - Compare LLMs to choose the best for your needs
2. **A/B Testing** - Validate outputs across multiple models
3. **Cost Optimization** - Track and minimize API costs
4. **Development** - Prototype without vendor lock-in
5. **Production** - Deploy with high-performance caching

---

## 🤝 Contributing

Contributions welcome! Check out our [contributing guidelines](CONTRIBUTING.md).

---

## 📄 License

MIT License - see [LICENSE](LICENSE)

---

<div align="center">

**Made with 💜 by [@Randomaigirl](https://github.com/Randomaigirl)**

Building AI that feels alive, not just intelligent ✨

⭐ Star this repo if you find it useful! ⭐

</div>
