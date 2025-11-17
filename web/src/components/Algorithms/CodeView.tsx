"use client";

import { useState, useRef } from "react";
import Editor from "@monaco-editor/react";
import { Play, Save, FileCode } from "lucide-react";

const defaultCode = `def algorithm(data: list, context: dict) -> dict:
    """
    Main algorithm function - DO NOT REMOVE THIS FUNCTION SIGNATURE
    
    Args:
        data: List of OHLC candles (most recent last)
            Each candle is a dict with: {'open': float, 'high': float, 'low': float, 'close': float, 'volume': int, 'timestamp': str}
            Example: [{'open': 100, 'high': 105, 'low': 99, 'close': 103, 'volume': 1000, 'timestamp': '2024-01-01T10:00:00'}, ...]
        
        context: Dictionary containing:
            - 'portfolio': Current positions and cash
                {'cash': float, 'positions': dict, 'total_value': float, 'pnl': float}
            - 'symbol': Current symbol being analyzed (str)
            - 'indicators': Pre-calculated indicators (dict)
            - 'state': Your algorithm's persistent state (dict) - modify this to track data across calls
            - 'config': Algorithm configuration (dict)
    
    Returns:
        dict with keys:
            - 'signal': 'BUY', 'SELL', or 'HOLD' (required)
            - 'quantity': Number of shares/units (optional, for BUY/SELL)
            - 'price': Entry/exit price (optional)
            - 'stop_loss': Stop loss price (optional)
            - 'target': Target price (optional)
            - 'reason': Reason for the signal (optional, for logging)
    
    Note: This function is called on each new candle. Use context['state'] to maintain
    indicators, positions, or other data across function calls.
    """
    
    # Get latest candle (most recent)
    if not data or len(data) == 0:
        return {'signal': 'HOLD', 'reason': 'No data available'}
    
    latest = data[-1]
    current_price = latest['close']
    
    # Initialize state if first call
    if 'state' not in context:
        context['state'] = {}
    
    state = context['state']
    
    # Example: Simple Moving Average Crossover Strategy
    # Calculate SMA if we have enough data
    if len(data) >= 20:
        # Calculate SMA 20
        sma_20 = sum([candle['close'] for candle in data[-20:]]) / 20
        
        # Calculate SMA 50 (if we have enough data)
        if len(data) >= 50:
            sma_50 = sum([candle['close'] for candle in data[-50:]]) / 50
            
            # Store previous SMA values in state
            prev_sma_20 = state.get('prev_sma_20', sma_20)
            prev_sma_50 = state.get('prev_sma_50', sma_50)
            
            # Check for crossover
            # Bullish crossover: SMA 20 crosses above SMA 50
            if prev_sma_20 <= prev_sma_50 and sma_20 > sma_50:
                state['prev_sma_20'] = sma_20
                state['prev_sma_50'] = sma_50
                return {
                    'signal': 'BUY',
                    'quantity': 10,  # Adjust based on your risk management
                    'price': current_price,
                    'stop_loss': current_price * 0.98,  # 2% stop loss
                    'target': current_price * 1.05,     # 5% target
                    'reason': f'Bullish crossover: SMA20 ({sma_20:.2f}) crossed above SMA50 ({sma_50:.2f})'
                }
            
            # Bearish crossover: SMA 20 crosses below SMA 50
            elif prev_sma_20 >= prev_sma_50 and sma_20 < sma_50:
                state['prev_sma_20'] = sma_20
                state['prev_sma_50'] = sma_50
                return {
                    'signal': 'SELL',
                    'quantity': 10,
                    'price': current_price,
                    'reason': f'Bearish crossover: SMA20 ({sma_20:.2f}) crossed below SMA50 ({sma_50:.2f})'
                }
            
            # Update state
            state['prev_sma_20'] = sma_20
            state['prev_sma_50'] = sma_50
    
    # Default: Hold
    return {
        'signal': 'HOLD',
        'reason': 'Waiting for signal'
    }
`;

export default function CodeView() {
  const [code, setCode] = useState(defaultCode);
  const [isRunning, setIsRunning] = useState(false);
  const editorRef = useRef<any>(null);

  const handleEditorDidMount = (editor: any, monaco: any) => {
    editorRef.current = editor;
    
    // Configure Python syntax highlighting
    monaco.languages.setLanguageConfiguration("python", {
      comments: {
        lineComment: "#",
        blockComment: ['"""', '"""'],
      },
      brackets: [
        ["{", "}"],
        ["[", "]"],
        ["(", ")"],
      ],
      autoClosingPairs: [
        { open: "{", close: "}" },
        { open: "[", close: "]" },
        { open: "(", close: ")" },
        { open: '"', close: '"' },
        { open: "'", close: "'" },
      ],
      surroundingPairs: [
        { open: "{", close: "}" },
        { open: "[", close: "]" },
        { open: "(", close: ")" },
        { open: '"', close: '"' },
        { open: "'", close: "'" },
      ],
    });

    // Custom theme to match app
    monaco.editor.defineTheme("100xtrader-dark", {
      base: "vs-dark",
      inherit: true,
      rules: [
        { token: "comment", foreground: "6b7280", fontStyle: "italic" },
        { token: "keyword", foreground: "f59e0b" },
        { token: "string", foreground: "10b981" },
        { token: "number", foreground: "3b82f6" },
        { token: "type", foreground: "8b5cf6" },
      ],
      colors: {
        "editor.background": "#1a1a1a",
        "editor.foreground": "#e5e5e5",
        "editor.lineHighlightBackground": "#2a2a2a",
        "editor.selectionBackground": "#3a3a3a",
        "editorCursor.foreground": "#f59e0b",
        "editorWhitespace.foreground": "#3a3a3a",
      },
    });

    monaco.editor.setTheme("100xtrader-dark");
  };

  const handleRun = () => {
    setIsRunning(true);
    // TODO: Implement code execution
    setTimeout(() => {
      setIsRunning(false);
      alert("Code executed! (This is a placeholder)");
    }, 1000);
  };

  const handleSave = () => {
    // TODO: Implement save functionality
    alert("Algorithm saved! (This is a placeholder)");
  };

  return (
    <div className="h-full flex flex-col bg-primary">
      {/* Toolbar */}
      <div className="bg-primary border-b border-primary px-4 py-2 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <FileCode className="w-4 h-4 text-tertiary" />
          <span className="font-helvetica-medium text-primary text-sm">
            Python Algorithm
          </span>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={handleRun}
            disabled={isRunning}
            className="flex items-center gap-2 px-3 py-1.5 bg-accent hover:bg-accent-hover text-inverse font-helvetica-medium text-sm rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <Play className="w-4 h-4 text-inverse" />
            {isRunning ? "Running..." : "Run"}
          </button>
          <button
            onClick={handleSave}
            className="flex items-center gap-2 px-3 py-1.5 bg-secondary hover:bg-tertiary text-primary font-helvetica-medium text-sm rounded-lg transition-colors"
          >
            <Save className="w-4 h-4 text-primary" />
            Save
          </button>
        </div>
      </div>

      {/* Editor */}
      <div className="flex-1 overflow-hidden">
        <Editor
          height="100%"
          defaultLanguage="python"
          value={code}
          onChange={(value) => setCode(value || "")}
          onMount={handleEditorDidMount}
          theme="100xtrader-dark"
          options={{
            fontSize: 14,
            minimap: { enabled: true },
            wordWrap: "on",
            lineNumbers: "on",
            roundedSelection: false,
            scrollBeyondLastLine: false,
            readOnly: false,
            cursorStyle: "line",
            automaticLayout: true,
            tabSize: 4,
            insertSpaces: true,
            formatOnPaste: true,
            formatOnType: true,
            suggestOnTriggerCharacters: true,
            acceptSuggestionOnEnter: "on",
            quickSuggestions: {
              other: true,
              comments: false,
              strings: true,
            },
            parameterHints: {
              enabled: true,
            },
            bracketPairColorization: {
              enabled: true,
            },
          }}
        />
      </div>

      {/* Output Panel (Collapsible) */}
      <div className="bg-secondary border-t border-primary p-4 max-h-48 overflow-y-auto">
        <div className="flex items-center justify-between mb-2">
          <h4 className="font-helvetica-medium text-primary text-sm">
            Output
          </h4>
        </div>
        <div className="font-helvetica text-tertiary text-xs font-mono">
          <div className="text-green-400"># Ready to run</div>
          <div className="text-tertiary mt-1">
            # Click "Run" to execute your algorithm
          </div>
        </div>
      </div>
    </div>
  );
}

