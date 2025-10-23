"use client";

import { useState } from "react";
import { Plus, Target, BookOpen, AlertTriangle } from "lucide-react";
import { IStrategy, IStrategyFormData, IRule, IRuleCreateRequest, IMistake, IMistakeCreateRequest } from "@/types";

// Strategy Modal Component
function AddStrategyModal({ isOpen, onClose, onSubmit }: { isOpen: boolean; onClose: () => void; onSubmit: (data: IStrategyFormData) => void }) {
  const [formData, setFormData] = useState<IStrategyFormData>({ name: "", description: "" });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim() || !formData.description.trim()) return;

    setIsSubmitting(true);
    try {
      await onSubmit(formData);
      setFormData({ name: "", description: "" });
      onClose();
    } catch (error) {
      console.error("Error creating strategy:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    setFormData({ name: "", description: "" });
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div 
      className="fixed inset-0 bg-black bg-opacity-60 modal-backdrop flex items-center justify-center z-50"
      onClick={handleClose}
    >
      <div 
        className="bg-primary border border-primary rounded-lg shadow-theme-lg w-full max-w-md mx-4 modal-content"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between p-6 border-b border-primary">
          <h2 className="text-xl font-helvetica-bold text-primary">Add New Strategy</h2>
          <button onClick={handleClose} className="p-2 hover:bg-secondary rounded-lg transition-colors">
            <Plus className="w-5 h-5 text-secondary rotate-45" />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-helvetica-medium text-primary mb-2">Strategy Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="Enter strategy name"
              className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary placeholder-tertiary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica transition-colors"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-helvetica-medium text-primary mb-2">Description</label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              placeholder="Enter strategy description"
              rows={4}
              className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary placeholder-tertiary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica resize-none transition-colors"
              required
            />
          </div>
          <div className="pt-4">
            <button
              type="submit"
              disabled={isSubmitting || !formData.name.trim() || !formData.description.trim()}
              className="w-full bg-accent hover:bg-secondary disabled:bg-tertiary disabled:cursor-not-allowed text-primary font-helvetica-medium py-3 px-4 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-accent"
            >
              {isSubmitting ? "Adding Strategy..." : "Add Strategy"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

// Rule Modal Component
function AddRuleModal({ isOpen, onClose, onSubmit }: { isOpen: boolean; onClose: () => void; onSubmit: (data: IRuleCreateRequest) => void }) {
  const [formData, setFormData] = useState<IRuleCreateRequest>({
    userId: "current-user",
    name: "",
    description: "",
    category: "ENTRY" as any,
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim() || !formData.description.trim()) return;

    setIsSubmitting(true);
    try {
      await onSubmit(formData);
      setFormData({ userId: "current-user", name: "", description: "", category: "ENTRY" as any });
      onClose();
    } catch (error) {
      console.error("Error creating rule:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    setFormData({ userId: "current-user", name: "", description: "", category: "ENTRY" as any });
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div 
      className="fixed inset-0 bg-black bg-opacity-60 modal-backdrop flex items-center justify-center z-50"
      onClick={handleClose}
    >
      <div 
        className="bg-primary border border-primary rounded-lg shadow-theme-lg w-full max-w-md mx-4 modal-content"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between p-6 border-b border-primary">
          <h2 className="text-xl font-helvetica-bold text-primary">Add New Rule</h2>
          <button onClick={handleClose} className="p-2 hover:bg-secondary rounded-lg transition-colors">
            <Plus className="w-5 h-5 text-secondary rotate-45" />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-helvetica-medium text-primary mb-2">Rule Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="Enter rule name"
              className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary placeholder-tertiary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica transition-colors"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-helvetica-medium text-primary mb-2">Category</label>
            <select
              value={formData.category}
              onChange={(e) => setFormData({ ...formData, category: e.target.value as any })}
              className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica transition-colors"
            >
              <option value="ENTRY">Entry</option>
              <option value="EXIT">Exit</option>
              <option value="STOP_LOSS">Stop Loss</option>
              <option value="TAKE_PROFIT">Take Profit</option>
              <option value="RISK_MANAGEMENT">Risk Management</option>
              <option value="PSYCHOLOGY">Psychology</option>
              <option value="OTHER">Other</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-helvetica-medium text-primary mb-2">Description</label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              placeholder="Enter rule description"
              rows={4}
              className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary placeholder-tertiary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica resize-none transition-colors"
              required
            />
          </div>
          <div className="pt-4">
            <button
              type="submit"
              disabled={isSubmitting || !formData.name.trim() || !formData.description.trim()}
              className="w-full bg-accent hover:bg-secondary disabled:bg-tertiary disabled:cursor-not-allowed text-primary font-helvetica-medium py-3 px-4 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-accent"
            >
              {isSubmitting ? "Adding Rule..." : "Add Rule"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

// Mistake Modal Component
function AddMistakeModal({ isOpen, onClose, onSubmit }: { isOpen: boolean; onClose: () => void; onSubmit: (data: IMistakeCreateRequest) => void }) {
  const [formData, setFormData] = useState<IMistakeCreateRequest>({
    userId: "current-user",
    name: "",
    category: "PSYCHOLOGICAL" as any,
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim()) return;

    setIsSubmitting(true);
    try {
      await onSubmit(formData);
      setFormData({ userId: "current-user", name: "", category: "PSYCHOLOGICAL" as any });
      onClose();
    } catch (error) {
      console.error("Error creating mistake:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    setFormData({ userId: "current-user", name: "", category: "PSYCHOLOGICAL" as any });
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div 
      className="fixed inset-0 bg-black bg-opacity-60 modal-backdrop flex items-center justify-center z-50"
      onClick={handleClose}
    >
      <div 
        className="bg-primary border border-primary rounded-lg shadow-theme-lg w-full max-w-md mx-4 modal-content"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between p-6 border-b border-primary">
          <h2 className="text-xl font-helvetica-bold text-primary">Add New Mistake</h2>
          <button onClick={handleClose} className="p-2 hover:bg-secondary rounded-lg transition-colors">
            <Plus className="w-5 h-5 text-secondary rotate-45" />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-helvetica-medium text-primary mb-2">Mistake Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="Enter mistake name"
              className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary placeholder-tertiary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica transition-colors"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-helvetica-medium text-primary mb-2">Category</label>
            <select
              value={formData.category}
              onChange={(e) => setFormData({ ...formData, category: e.target.value as any })}
              className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica transition-colors"
            >
              <option value="PSYCHOLOGICAL">Psychological</option>
              <option value="BEHAVIORAL">Behavioral</option>
            </select>
          </div>
          <div className="pt-4">
            <button
              type="submit"
              disabled={isSubmitting || !formData.name.trim()}
              className="w-full bg-accent hover:bg-secondary disabled:bg-tertiary disabled:cursor-not-allowed text-primary font-helvetica-medium py-3 px-4 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-accent"
            >
              {isSubmitting ? "Adding Mistake..." : "Add Mistake"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default function ConfigurePage() {
  const [activeTab, setActiveTab] = useState<"strategies" | "rules" | "mistakes">("strategies");
  const [isStrategyModalOpen, setIsStrategyModalOpen] = useState(false);
  const [isRuleModalOpen, setIsRuleModalOpen] = useState(false);
  const [isMistakeModalOpen, setIsMistakeModalOpen] = useState(false);
  
  const [strategies, setStrategies] = useState<IStrategy[]>([]);
  const [rules, setRules] = useState<IRule[]>([]);
  const [mistakes, setMistakes] = useState<IMistake[]>([]);

  const handleAddStrategy = async (data: IStrategyFormData) => {
    const newStrategy: IStrategy = {
      id: Date.now().toString(),
      name: data.name,
      description: data.description,
      createdAt: new Date(),
      updatedAt: new Date(),
      userId: "current-user",
    };
    setStrategies([...strategies, newStrategy]);
  };

  const handleAddRule = async (data: IRuleCreateRequest) => {
    const newRule: IRule = {
      id: Date.now().toString(),
      name: data.name,
      description: data.description,
      category: data.category,
      createdAt: new Date(),
      updatedAt: new Date(),
      userId: "current-user",
    };
    setRules([...rules, newRule]);
  };

  const handleAddMistake = async (data: IMistakeCreateRequest) => {
    const newMistake: IMistake = {
      id: Date.now().toString(),
      name: data.name,
      category: data.category,
      createdAt: new Date(),
      updatedAt: new Date(),
      userId: "current-user",
    };
    setMistakes([...mistakes, newMistake]);
  };

  const tabs = [
    { id: "strategies", label: "Strategies", icon: Target },
    { id: "rules", label: "Rules", icon: BookOpen },
    { id: "mistakes", label: "Mistakes", icon: AlertTriangle },
  ];

  const renderContent = () => {
    switch (activeTab) {
      case "strategies":
        return (
          <div>
            {strategies.length === 0 ? (
              <div className="text-center py-12">
                <div className="w-16 h-16 bg-secondary rounded-full flex items-center justify-center mx-auto mb-4">
                  <Target className="w-8 h-8 text-tertiary" />
                </div>
                <h3 className="text-lg font-helvetica-medium text-primary mb-2">No strategies yet</h3>
                <p className="text-tertiary font-helvetica mb-4">Create your first trading strategy to get started</p>
                <button
                  onClick={() => setIsStrategyModalOpen(true)}
                  className="bg-accent hover:bg-secondary text-primary font-helvetica-medium py-2 px-4 rounded-lg transition-colors duration-200"
                >
                  Add Your First Strategy
                </button>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {strategies.map((strategy) => (
                  <div key={strategy.id} className="bg-secondary border border-primary rounded-lg p-6 hover:bg-tertiary transition-colors duration-200">
                    <h3 className="text-lg font-helvetica-bold text-primary mb-2">{strategy.name}</h3>
                    <p className="text-tertiary font-helvetica text-sm mb-4 line-clamp-3">{strategy.description}</p>
                    <div className="flex items-center justify-between text-xs text-muted">
                      <span className="font-helvetica">Created {strategy.createdAt.toLocaleDateString()}</span>
                      <span className="font-helvetica">Updated {strategy.updatedAt.toLocaleDateString()}</span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        );
      case "rules":
        return (
          <div>
            {rules.length === 0 ? (
              <div className="text-center py-12">
                <div className="w-16 h-16 bg-secondary rounded-full flex items-center justify-center mx-auto mb-4">
                  <BookOpen className="w-8 h-8 text-tertiary" />
                </div>
                <h3 className="text-lg font-helvetica-medium text-primary mb-2">No rules yet</h3>
                <p className="text-tertiary font-helvetica mb-4">Create your first trading rule to get started</p>
                <button
                  onClick={() => setIsRuleModalOpen(true)}
                  className="bg-accent hover:bg-secondary text-primary font-helvetica-medium py-2 px-4 rounded-lg transition-colors duration-200"
                >
                  Add Your First Rule
                </button>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {rules.map((rule) => (
                  <div key={rule.id} className="bg-secondary border border-primary rounded-lg p-6 hover:bg-tertiary transition-colors duration-200">
                    <h3 className="text-lg font-helvetica-bold text-primary mb-2">{rule.name}</h3>
                    <p className="text-tertiary font-helvetica text-sm mb-2">{rule.description}</p>
                    <span className="inline-block bg-accent text-primary text-xs px-2 py-1 rounded mb-4">{rule.category}</span>
                    <div className="flex items-center justify-between text-xs text-muted">
                      <span className="font-helvetica">Created {rule.createdAt.toLocaleDateString()}</span>
                      <span className="font-helvetica">Updated {rule.updatedAt.toLocaleDateString()}</span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        );
      case "mistakes":
        return (
          <div>
            {mistakes.length === 0 ? (
              <div className="text-center py-12">
                <div className="w-16 h-16 bg-secondary rounded-full flex items-center justify-center mx-auto mb-4">
                  <AlertTriangle className="w-8 h-8 text-tertiary" />
                </div>
                <h3 className="text-lg font-helvetica-medium text-primary mb-2">No mistakes yet</h3>
                <p className="text-tertiary font-helvetica mb-4">Create your first trading mistake to get started</p>
                <button
                  onClick={() => setIsMistakeModalOpen(true)}
                  className="bg-accent hover:bg-secondary text-primary font-helvetica-medium py-2 px-4 rounded-lg transition-colors duration-200"
                >
                  Add Your First Mistake
                </button>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {mistakes.map((mistake) => (
                  <div key={mistake.id} className="bg-secondary border border-primary rounded-lg p-6 hover:bg-tertiary transition-colors duration-200">
                    <h3 className="text-lg font-helvetica-bold text-primary mb-2">{mistake.name}</h3>
                    <span className="inline-block bg-accent text-primary text-xs px-2 py-1 rounded mb-4">{mistake.category}</span>
                    <div className="flex items-center justify-between text-xs text-muted">
                      <span className="font-helvetica">Created {mistake.createdAt.toLocaleDateString()}</span>
                      <span className="font-helvetica">Updated {mistake.updatedAt.toLocaleDateString()}</span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        );
      default:
        return null;
    }
  };

  return (
    <div className="p-6">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-helvetica-bold text-primary">Configure</h1>
          <p className="text-secondary font-helvetica mt-1">Manage your trading strategies, rules, and mistakes</p>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex space-x-1 mb-6 bg-secondary p-1 rounded-lg">
        {tabs.map((tab) => {
          const IconComponent = tab.icon;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as any)}
              className={`flex items-center space-x-2 px-4 py-2 rounded-md font-helvetica-medium transition-colors duration-200 ${
                activeTab === tab.id
                  ? "bg-primary text-secondary"
                  : "text-tertiary hover:text-primary hover:bg-tertiary"
              }`}
            >
              <IconComponent className="w-4 h-4" />
              <span>{tab.label}</span>
            </button>
          );
        })}
      </div>

      {/* Content */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-helvetica-bold text-primary capitalize">{activeTab}</h2>
          <button
            onClick={() => {
              if (activeTab === "strategies") setIsStrategyModalOpen(true);
              else if (activeTab === "rules") setIsRuleModalOpen(true);
              else if (activeTab === "mistakes") setIsMistakeModalOpen(true);
            }}
            className="flex items-center space-x-2 bg-accent hover:bg-secondary text-primary font-helvetica-medium py-2 px-4 rounded-lg transition-colors duration-200"
          >
            <Plus className="w-4 h-4" />
            <span>Add {activeTab.slice(0, -1)}</span>
          </button>
        </div>
        {renderContent()}
      </div>

      {/* Modals */}
      <AddStrategyModal
        isOpen={isStrategyModalOpen}
        onClose={() => setIsStrategyModalOpen(false)}
        onSubmit={handleAddStrategy}
      />
      <AddRuleModal
        isOpen={isRuleModalOpen}
        onClose={() => setIsRuleModalOpen(false)}
        onSubmit={handleAddRule}
      />
      <AddMistakeModal
        isOpen={isMistakeModalOpen}
        onClose={() => setIsMistakeModalOpen(false)}
        onSubmit={handleAddMistake}
      />
    </div>
  );
}
