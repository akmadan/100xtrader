"use client";

import { useCallback } from "react";
import {
  ReactFlow,
  Node,
  addEdge,
  Background,
  Controls,
  MiniMap,
  Connection,
  useNodesState,
  useEdgesState,
  Panel,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import {
  Database,
  TrendingUp,
  Code,
  Target,
  AlertTriangle,
  Play,
} from "lucide-react";

// Custom Node Components
function MarketDataNode({ data }: { data: any }) {
  return (
    <div className="px-4 py-3 bg-primary border border-primary rounded-lg shadow-lg min-w-[200px]">
      <div className="flex items-center gap-2 mb-2">
        <Database className="w-4 h-4 text-accent" />
        <span className="font-helvetica-bold text-primary text-sm">
          Market Data
        </span>
      </div>
      <div className="text-tertiary font-helvetica text-xs">
        {data.label || "Real-time Price Feed"}
      </div>
    </div>
  );
}

function IndicatorNode({ data }: { data: any }) {
  return (
    <div className="px-4 py-3 bg-primary border border-primary rounded-lg shadow-lg min-w-[200px]">
      <div className="flex items-center gap-2 mb-2">
        <TrendingUp className="w-4 h-4 text-accent" />
        <span className="font-helvetica-bold text-primary text-sm">
          {data.label || "Indicator"}
        </span>
      </div>
      <div className="text-tertiary font-helvetica text-xs">
        {data.type || "RSI"}
      </div>
    </div>
  );
}

function CodeBlockNode({ data }: { data: any }) {
  return (
    <div className="px-4 py-3 bg-primary border border-primary rounded-lg shadow-lg min-w-[200px]">
      <div className="flex items-center gap-2 mb-2">
        <Code className="w-4 h-4 text-accent" />
        <span className="font-helvetica-bold text-primary text-sm">
          Code Block
        </span>
      </div>
      <div className="text-tertiary font-helvetica text-xs">
        Custom Python Logic
      </div>
    </div>
  );
}

function ConditionNode({ data }: { data: any }) {
  return (
    <div className="px-4 py-3 bg-primary border border-primary rounded-lg shadow-lg min-w-[200px]">
      <div className="flex items-center gap-2 mb-2">
        <Target className="w-4 h-4 text-accent" />
        <span className="font-helvetica-bold text-primary text-sm">
          Condition
        </span>
      </div>
      <div className="text-tertiary font-helvetica text-xs">
        {data.condition || "IF/THEN/ELSE"}
      </div>
    </div>
  );
}

function RiskManagementNode({ data }: { data: any }) {
  return (
    <div className="px-4 py-3 bg-primary border border-primary rounded-lg shadow-lg min-w-[200px]">
      <div className="flex items-center gap-2 mb-2">
        <AlertTriangle className="w-4 h-4 text-accent" />
        <span className="font-helvetica-bold text-primary text-sm">
          Risk Management
        </span>
      </div>
      <div className="text-tertiary font-helvetica text-xs">
        Stop Loss / Position Size
      </div>
    </div>
  );
}

function BrokerNode({ data }: { data: any }) {
  return (
    <div className="px-4 py-3 bg-primary border-2 border-accent rounded-lg shadow-lg min-w-[200px]">
      <div className="flex items-center gap-2 mb-2">
        <Play className="w-4 h-4 text-accent" />
        <span className="font-helvetica-bold text-primary text-sm">
          {data.broker || "Dhan"} Trade
        </span>
      </div>
      <div className="text-tertiary font-helvetica text-xs">
        Execute Order
      </div>
    </div>
  );
}

// Component Palette
const componentTypes = [
  {
    type: "marketData",
    label: "Market Data",
    icon: Database,
    category: "Data Sources",
  },
  {
    type: "indicator",
    label: "RSI Indicator",
    icon: TrendingUp,
    category: "Indicators",
  },
  {
    type: "indicator",
    label: "MACD Indicator",
    icon: TrendingUp,
    category: "Indicators",
  },
  {
    type: "codeBlock",
    label: "Code Block",
    icon: Code,
    category: "Custom",
  },
  {
    type: "condition",
    label: "IF Condition",
    icon: Target,
    category: "Logic",
  },
  {
    type: "riskManagement",
    label: "Risk Management",
    icon: AlertTriangle,
    category: "Risk",
  },
  {
    type: "broker",
    label: "Dhan Trade",
    icon: Play,
    category: "Execution",
  },
];

// Define node types (after component definitions)
const nodeTypes = {
  marketData: MarketDataNode,
  indicator: IndicatorNode,
  codeBlock: CodeBlockNode,
  condition: ConditionNode,
  riskManagement: RiskManagementNode,
  broker: BrokerNode,
};

export default function CanvasView() {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);

  const onConnect = useCallback(
    (params: Connection) => {
      setEdges((eds) => addEdge(params, eds));
    },
    [setEdges]
  );

  const addNode = useCallback(
    (type: string, label: string) => {
      const newNode: Node = {
        id: `${type}-${Date.now()}`,
        type,
        position: {
          x: Math.random() * 400 + 100,
          y: Math.random() * 400 + 100,
        },
        data: { label, type },
      };
      setNodes((nds) => [...nds, newNode]);
    },
    [setNodes]
  );

  // Custom theme for React Flow
  const flowTheme = {
    background: "transparent",
  };

  return (
    <div className="h-full flex bg-primary">
      {/* Component Palette Sidebar */}
      <div className="w-64 bg-secondary border-r border-primary p-4 overflow-y-auto">
        <h3 className="font-helvetica-bold text-primary mb-4 text-sm">
          Components
        </h3>
        <div className="space-y-2">
          {["Data Sources", "Indicators", "Logic", "Custom", "Risk", "Execution"].map(
            (category) => {
              const categoryComponents = componentTypes.filter((comp) => comp.category === category);
              if (categoryComponents.length === 0) return null;
              
              return (
                <div key={category} className="mb-4">
                  <h4 className="font-helvetica-medium text-tertiary text-xs mb-2 uppercase">
                    {category}
                  </h4>
                  <div className="space-y-1">
                    {categoryComponents.map((comp) => {
                      const Icon = comp.icon;
                      return (
                        <button
                          key={comp.label}
                          onClick={() => addNode(comp.type, comp.label)}
                          className="w-full flex items-center gap-2 px-3 py-2 bg-primary hover:bg-tertiary border border-primary rounded-lg transition-colors duration-200 text-left"
                        >
                          <Icon className="w-4 h-4 text-accent flex-shrink-0" />
                          <span className="font-helvetica text-primary text-sm">
                            {comp.label}
                          </span>
                        </button>
                      );
                    })}
                  </div>
                </div>
              );
            }
          )}
        </div>
      </div>

      {/* Canvas Area */}
      <div className="flex-1 relative">
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          nodeTypes={nodeTypes}
          fitView
          className="bg-secondary/30"
        >
          <Background color="#3a3a3a" gap={16} />
          <Controls className="bg-primary border border-primary" />
          <MiniMap
            className="bg-primary border border-primary"
            nodeColor={(node) => {
              if (node.type === "broker") return "#f59e0b";
              return "#6b7280";
            }}
          />
          <Panel position="top-right" className="bg-primary border border-primary rounded-lg p-2">
            <div className="flex items-center gap-2">
              <button               className="flex items-center gap-2 px-3 py-1.5 bg-accent hover:bg-accent-hover text-inverse font-helvetica-medium text-sm rounded-lg transition-colors">
                <Play className="w-4 h-4 text-inverse" />
                Run
              </button>
              <button className="px-3 py-1.5 bg-secondary hover:bg-tertiary text-primary font-helvetica-medium text-sm rounded-lg transition-colors">
                Save
              </button>
            </div>
          </Panel>
        </ReactFlow>
      </div>
    </div>
  );
}

