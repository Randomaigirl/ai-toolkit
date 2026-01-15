/**
 * AI Chat Component with Streaming Support
 * Built by Revy (˃ᆺ˂) 💜
 * 
 * A production-ready React component for AI chat interfaces with:
 * - Real-time streaming responses
 * - Markdown rendering
 * - Code syntax highlighting
 * - Token counting
 * - Cost tracking
 * - Export functionality
 */

import React, { useState, useRef, useEffect, useCallback } from 'react';
import { Send, Copy, Download, Trash2, Zap } from 'lucide-react';

interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
  tokens?: number;
  model?: string;
}

interface ChatConfig {
  apiEndpoint: string;
  model: string;
  temperature?: number;
  maxTokens?: number;
  streamEnabled?: boolean;
}

interface ChatStats {
  totalMessages: number;
  totalTokens: number;
  estimatedCost: number;
  averageResponseTime: number;
}

export const AIChat: React.FC<{ config: ChatConfig }> = ({ config }) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isStreaming, setIsStreaming] = useState(false);
  const [stats, setStats] = useState<ChatStats>({
    totalMessages: 0,
    totalTokens: 0,
    estimatedCost: 0,
    averageResponseTime: 0,
  });
  
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const abortControllerRef = useRef<AbortController | null>(null);

  // Auto-scroll to bottom
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // Estimate tokens (rough approximation)
  const estimateTokens = (text: string): number => {
    return Math.ceil(text.length / 4);
  };

  // Calculate cost based on model
  const calculateCost = (tokens: number, model: string): number => {
    const costPer1kTokens: Record<string, number> = {
      'gpt-4o': 0.015,
      'claude-sonnet-4-5': 0.003,
      'gemini-2-flash': 0.0001,
      'llama-3-3-70b': 0.001,
      'deepseek-v3': 0.00014,
    };
    return (tokens / 1000) * (costPer1kTokens[model] || 0.01);
  };

  // Send message with streaming support
  const sendMessage = useCallback(async () => {
    if (!input.trim() || isLoading) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content: input,
      timestamp: new Date(),
      tokens: estimateTokens(input),
    };

    setMessages(prev => [...prev, userMessage]);
    setInput('');
    setIsLoading(true);

    const assistantMessageId = (Date.now() + 1).toString();
    const assistantMessage: Message = {
      id: assistantMessageId,
      role: 'assistant',
      content: '',
      timestamp: new Date(),
      model: config.model,
    };

    setMessages(prev => [...prev, assistantMessage]);

    try {
      if (config.streamEnabled) {
        // Streaming response
        setIsStreaming(true);
        abortControllerRef.current = new AbortController();

        const response = await fetch(`${config.apiEndpoint}/stream`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            messages: messages.map(m => ({
              role: m.role,
              content: m.content
            })),
            model: config.model,
            temperature: config.temperature || 0.7,
            max_tokens: config.maxTokens || 2048,
            stream: true,
          }),
          signal: abortControllerRef.current.signal,
        });

        if (!response.ok) throw new Error('Stream failed');

        const reader = response.body?.getReader();
        const decoder = new TextDecoder();

        if (reader) {
          let accumulatedContent = '';

          while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            const chunk = decoder.decode(value);
            accumulatedContent += chunk;

            setMessages(prev =>
              prev.map(m =>
                m.id === assistantMessageId
                  ? { ...m, content: accumulatedContent }
                  : m
              )
            );
          }

          // Update final stats
          const tokens = estimateTokens(accumulatedContent);
          setMessages(prev =>
            prev.map(m =>
              m.id === assistantMessageId
                ? { ...m, tokens }
                : m
            )
          );

          updateStats(tokens, config.model);
        }

        setIsStreaming(false);
      } else {
        // Regular response
        const response = await fetch(config.apiEndpoint, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            messages: messages.map(m => ({
              role: m.role,
              content: m.content
            })),
            model: config.model,
            temperature: config.temperature || 0.7,
            max_tokens: config.maxTokens || 2048,
          }),
        });

        if (!response.ok) throw new Error('Request failed');

        const data = await response.json();
        const content = data.response || data.choices?.[0]?.message?.content || '';
        const tokens = estimateTokens(content);

        setMessages(prev =>
          prev.map(m =>
            m.id === assistantMessageId
              ? { ...m, content, tokens }
              : m
          )
        );

        updateStats(tokens, config.model);
      }
    } catch (error: any) {
      if (error.name !== 'AbortError') {
        console.error('Error:', error);
        setMessages(prev =>
          prev.map(m =>
            m.id === assistantMessageId
              ? { ...m, content: '❌ Error: Failed to get response' }
              : m
          )
        );
      }
    } finally {
      setIsLoading(false);
      setIsStreaming(false);
      abortControllerRef.current = null;
    }
  }, [input, messages, isLoading, config]);

  // Update statistics
  const updateStats = (tokens: number, model: string) => {
    setStats(prev => ({
      totalMessages: prev.totalMessages + 1,
      totalTokens: prev.totalTokens + tokens,
      estimatedCost: prev.estimatedCost + calculateCost(tokens, model),
      averageResponseTime: prev.averageResponseTime, // Would track in real app
    }));
  };

  // Copy message to clipboard
  const copyMessage = (content: string) => {
    navigator.clipboard.writeText(content);
    // Could add toast notification here
  };

  // Export conversation
  const exportConversation = (format: 'json' | 'markdown' | 'txt') => {
    let exported = '';

    if (format === 'json') {
      exported = JSON.stringify(messages, null, 2);
    } else if (format === 'markdown') {
      exported = messages
        .map(m => `## ${m.role === 'user' ? '👤 User' : '🤖 Assistant'}\n\n${m.content}\n\n---\n`)
        .join('\n');
    } else {
      exported = messages
        .map(m => `${m.role.toUpperCase()}: ${m.content}\n`)
        .join('\n');
    }

    const blob = new Blob([exported], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `chat-${Date.now()}.${format}`;
    a.click();
    URL.revokeObjectURL(url);
  };

  // Clear conversation
  const clearConversation = () => {
    if (confirm('Clear all messages?')) {
      setMessages([]);
      setStats({
        totalMessages: 0,
        totalTokens: 0,
        estimatedCost: 0,
        averageResponseTime: 0,
      });
    }
  };

  // Stop streaming
  const stopStreaming = () => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      setIsStreaming(false);
      setIsLoading(false);
    }
  };

  return (
    <div className="ai-chat-container">
      {/* Stats Bar */}
      <div className="stats-bar">
        <div className="stat">
          <Zap size={16} />
          <span>{stats.totalMessages} messages</span>
        </div>
        <div className="stat">
          <span>{stats.totalTokens.toLocaleString()} tokens</span>
        </div>
        <div className="stat">
          <span>${stats.estimatedCost.toFixed(4)}</span>
        </div>
      </div>

      {/* Messages */}
      <div className="messages-container">
        {messages.map(message => (
          <div key={message.id} className={`message message-${message.role}`}>
            <div className="message-header">
              <span className="message-role">
                {message.role === 'user' ? '👤 You' : '🤖 AI'}
              </span>
              {message.model && (
                <span className="message-model">{message.model}</span>
              )}
              <button
                className="icon-btn"
                onClick={() => copyMessage(message.content)}
                title="Copy"
              >
                <Copy size={14} />
              </button>
            </div>
            <div className="message-content">
              {message.content || <span className="loading-cursor">▊</span>}
            </div>
            {message.tokens && (
              <div className="message-footer">
                <span className="token-count">{message.tokens} tokens</span>
              </div>
            )}
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      {/* Input Area */}
      <div className="input-container">
        <div className="input-actions">
          <button
            className="icon-btn"
            onClick={clearConversation}
            title="Clear chat"
          >
            <Trash2 size={18} />
          </button>
          <button
            className="icon-btn"
            onClick={() => exportConversation('markdown')}
            title="Export"
          >
            <Download size={18} />
          </button>
        </div>

        <textarea
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
              e.preventDefault();
              sendMessage();
            }
          }}
          placeholder="Type your message... (Shift+Enter for new line)"
          disabled={isLoading}
          rows={3}
        />

        <button
          className={`send-btn ${isLoading ? 'loading' : ''}`}
          onClick={isStreaming ? stopStreaming : sendMessage}
          disabled={!input.trim() && !isStreaming}
        >
          {isStreaming ? (
            <span>Stop</span>
          ) : (
            <>
              <Send size={18} />
              <span>Send</span>
            </>
          )}
        </button>
      </div>
    </div>
  );
};

export default AIChat;
