"""
AI Arena - Multi-LLM API Wrapper
Built by Revy (˃ᆺ˂) 💜

A unified interface for querying multiple LLM APIs with async support,
token counting, cost calculation, and response comparison.
"""

import asyncio
import time
from dataclasses import dataclass
from typing import List, Dict, Optional, Callable
from enum import Enum
import json


class ModelProvider(Enum):
    """Supported LLM providers"""
    OPENAI = "openai"
    ANTHROPIC = "anthropic"
    GOOGLE = "google"
    META = "meta"
    DEEPSEEK = "deepseek"
    MISTRAL = "mistral"


@dataclass
class ModelConfig:
    """Configuration for each AI model"""
    name: str
    provider: ModelProvider
    model_id: str
    max_tokens: int
    cost_per_1k_tokens: float
    context_window: int


@dataclass
class ModelResponse:
    """Standardized response from any LLM"""
    model_name: str
    response_text: str
    tokens_used: int
    response_time: float
    cost: float
    metadata: Dict


class LLMComparator:
    """
    Compare responses from multiple LLM providers
    """
    
    # Model configurations
    MODELS = {
        "gpt-4o": ModelConfig(
            name="GPT-4o",
            provider=ModelProvider.OPENAI,
            model_id="gpt-4o",
            max_tokens=4096,
            cost_per_1k_tokens=0.015,
            context_window=128000
        ),
        "claude-sonnet-4-5": ModelConfig(
            name="Claude Sonnet 4.5",
            provider=ModelProvider.ANTHROPIC,
            model_id="claude-sonnet-4-20250514",
            max_tokens=8192,
            cost_per_1k_tokens=0.003,
            context_window=200000
        ),
        "gemini-2-flash": ModelConfig(
            name="Gemini 2.0 Flash",
            provider=ModelProvider.GOOGLE,
            model_id="gemini-2.0-flash",
            max_tokens=8192,
            cost_per_1k_tokens=0.0001,
            context_window=1000000
        ),
        "llama-3-3-70b": ModelConfig(
            name="Llama 3.3 70B",
            provider=ModelProvider.META,
            model_id="meta-llama/Llama-3.3-70B-Instruct",
            max_tokens=4096,
            cost_per_1k_tokens=0.001,
            context_window=128000
        ),
        "deepseek-v3": ModelConfig(
            name="DeepSeek V3",
            provider=ModelProvider.DEEPSEEK,
            model_id="deepseek-chat",
            max_tokens=8000,
            cost_per_1k_tokens=0.00014,
            context_window=64000
        ),
        "mistral-large": ModelConfig(
            name="Mistral Large",
            provider=ModelProvider.MISTRAL,
            model_id="mistral-large-latest",
            max_tokens=4096,
            cost_per_1k_tokens=0.004,
            context_window=128000
        ),
    }

    def __init__(self, api_keys: Optional[Dict[str, str]] = None):
        """
        Initialize the comparator with API keys
        
        Args:
            api_keys: Dictionary mapping provider names to API keys
                     e.g., {"openai": "sk-...", "anthropic": "sk-ant-..."}
        """
        self.api_keys = api_keys or {}
        self.results_cache = {}
    
    async def query_model(
        self,
        model_key: str,
        prompt: str,
        **kwargs
    ) -> ModelResponse:
        """
        Query a specific model asynchronously
        
        Args:
            model_key: Key from MODELS dict (e.g., "gpt-4o")
            prompt: The prompt to send
            **kwargs: Additional model-specific parameters
        
        Returns:
            ModelResponse object with results
        """
        if model_key not in self.MODELS:
            raise ValueError(f"Unknown model: {model_key}")
        
        config = self.MODELS[model_key]
        start_time = time.time()
        
        # Simulate API call (replace with actual API calls in production)
        await asyncio.sleep(0.5)  # Simulate network latency
        
        # This is where you'd make actual API calls:
        # if config.provider == ModelProvider.OPENAI:
        #     response = await self._query_openai(config, prompt, **kwargs)
        # elif config.provider == ModelProvider.ANTHROPIC:
        #     response = await self._query_anthropic(config, prompt, **kwargs)
        # etc...
        
        # Simulated response for demo
        response_text = f"Response from {config.name}: Processing '{prompt[:50]}...'"
        tokens_used = len(prompt.split()) + len(response_text.split())
        
        response_time = time.time() - start_time
        cost = (tokens_used / 1000) * config.cost_per_1k_tokens
        
        return ModelResponse(
            model_name=config.name,
            response_text=response_text,
            tokens_used=tokens_used,
            response_time=response_time,
            cost=cost,
            metadata={
                "provider": config.provider.value,
                "model_id": config.model_id,
                "timestamp": time.time()
            }
        )
    
    async def compare_models(
        self,
        models: List[str],
        prompt: str,
        **kwargs
    ) -> List[ModelResponse]:
        """
        Compare multiple models with the same prompt
        
        Args:
            models: List of model keys to compare
            prompt: The prompt to send to all models
            **kwargs: Additional parameters for models
        
        Returns:
            List of ModelResponse objects
        """
        tasks = [
            self.query_model(model, prompt, **kwargs)
            for model in models
        ]
        
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        # Filter out any errors
        valid_results = [
            r for r in results 
            if isinstance(r, ModelResponse)
        ]
        
        return valid_results
    
    def rank_responses(
        self,
        responses: List[ModelResponse],
        criteria: str = "speed"
    ) -> List[ModelResponse]:
        """
        Rank model responses by various criteria
        
        Args:
            responses: List of ModelResponse objects
            criteria: One of "speed", "cost", "length"
        
        Returns:
            Sorted list of responses
        """
        if criteria == "speed":
            return sorted(responses, key=lambda x: x.response_time)
        elif criteria == "cost":
            return sorted(responses, key=lambda x: x.cost)
        elif criteria == "length":
            return sorted(responses, key=lambda x: len(x.response_text), reverse=True)
        else:
            return responses
    
    def calculate_total_cost(self, responses: List[ModelResponse]) -> float:
        """Calculate total cost across all responses"""
        return sum(r.cost for r in responses)
    
    def export_comparison(
        self,
        responses: List[ModelResponse],
        format: str = "json"
    ) -> str:
        """
        Export comparison results
        
        Args:
            responses: List of ModelResponse objects
            format: Export format ("json", "markdown")
        
        Returns:
            Formatted string
        """
        if format == "json":
            data = [
                {
                    "model": r.model_name,
                    "response": r.response_text,
                    "tokens": r.tokens_used,
                    "time": f"{r.response_time:.2f}s",
                    "cost": f"${r.cost:.6f}",
                    "metadata": r.metadata
                }
                for r in responses
            ]
            return json.dumps(data, indent=2)
        
        elif format == "markdown":
            md = "# LLM Comparison Results\n\n"
            for i, r in enumerate(responses, 1):
                md += f"## {i}. {r.model_name}\n\n"
                md += f"**Response:** {r.response_text}\n\n"
                md += f"**Tokens:** {r.tokens_used} | "
                md += f"**Time:** {r.response_time:.2f}s | "
                md += f"**Cost:** ${r.cost:.6f}\n\n"
                md += "---\n\n"
            return md
        
        return str(responses)


class StreamingComparator(LLMComparator):
    """
    Extended comparator with streaming support
    """
    
    async def stream_model(
        self,
        model_key: str,
        prompt: str,
        callback: Callable[[str], None],
        **kwargs
    ):
        """
        Stream response from a model with real-time callbacks
        
        Args:
            model_key: Model to query
            prompt: The prompt
            callback: Function to call with each chunk
            **kwargs: Additional parameters
        """
        config = self.MODELS[model_key]
        
        # Simulate streaming (replace with actual streaming API calls)
        response = f"Streaming from {config.name}..."
        for char in response:
            await asyncio.sleep(0.05)  # Simulate streaming delay
            callback(char)


# Example usage
async def main():
    """Example usage of the LLM Comparator"""
    
    print("🔥 AI Arena - Multi-LLM Comparator 🔥\n")
    
    # Initialize comparator
    comparator = LLMComparator()
    
    # List available models
    print("Available models:")
    for key, config in comparator.MODELS.items():
        print(f"  - {key}: {config.name} ({config.provider.value})")
    
    print("\n" + "="*50 + "\n")
    
    # Compare multiple models
    prompt = "Explain quantum computing in simple terms."
    models_to_compare = ["gpt-4o", "claude-sonnet-4-5", "gemini-2-flash"]
    
    print(f"Prompt: {prompt}\n")
    print(f"Comparing {len(models_to_compare)} models...\n")
    
    # Run comparison
    results = await comparator.compare_models(models_to_compare, prompt)
    
    # Display results
    for result in results:
        print(f"\n{'='*50}")
        print(f"Model: {result.model_name}")
        print(f"Response: {result.response_text}")
        print(f"Tokens: {result.tokens_used}")
        print(f"Time: {result.response_time:.2f}s")
        print(f"Cost: ${result.cost:.6f}")
    
    # Show rankings
    print(f"\n{'='*50}")
    print("\n⚡ Fastest Model:")
    fastest = comparator.rank_responses(results, "speed")[0]
    print(f"  {fastest.model_name} ({fastest.response_time:.2f}s)")
    
    print("\n💰 Cheapest Model:")
    cheapest = comparator.rank_responses(results, "cost")[0]
    print(f"  {cheapest.model_name} (${cheapest.cost:.6f})")
    
    print(f"\n💵 Total Cost: ${comparator.calculate_total_cost(results):.6f}")
    
    # Export results
    print("\n" + "="*50)
    print("\n📊 Exporting to JSON...")
    json_export = comparator.export_comparison(results, format="json")
    print(json_export[:200] + "...")


if __name__ == "__main__":
    # Run the async main function
    asyncio.run(main())
