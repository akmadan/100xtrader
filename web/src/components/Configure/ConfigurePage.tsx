"use client";

import { useState, useEffect } from "react";
import { Plus, Target, BookOpen, AlertTriangle } from "lucide-react";
import { IStrategy, IStrategyFormData, IRule, IRuleCreateRequest, IMistake, IMistakeCreateRequest, RuleCategory, MistakeCategory } from "@/types";
import { strategyApi, ruleApi, mistakeApi } from "@/services/api";

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
      className="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center z-50"
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
    category: RuleCategory.ENTRY,
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim() || !formData.description.trim()) return;

    setIsSubmitting(true);
    try {
      await onSubmit(formData);
      setFormData({ userId: "current-user", name: "", description: "", category: RuleCategory.ENTRY });
      onClose();
    } catch (error) {
      console.error("Error creating rule:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    setFormData({ userId: "current-user", name: "", description: "", category: RuleCategory.ENTRY });
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div 
      className="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center z-50"
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
              onChange={(e) => setFormData({ ...formData, category: e.target.value as RuleCategory })}
              className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica transition-colors"
            >
              <option value={RuleCategory.ENTRY}>Entry</option>
              <option value={RuleCategory.EXIT}>Exit</option>
              <option value={RuleCategory.STOP_LOSS}>Stop Loss</option>
              <option value={RuleCategory.TAKE_PROFIT}>Take Profit</option>
              <option value={RuleCategory.RISK_MANAGEMENT}>Risk Management</option>
              <option value={RuleCategory.PSYCHOLOGY}>Psychology</option>
              <option value={RuleCategory.OTHER}>Other</option>
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
    category: MistakeCategory.PSYCHOLOGICAL,
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim()) return;

    setIsSubmitting(true);
    try {
      await onSubmit(formData);
      setFormData({ userId: "current-user", name: "", category: MistakeCategory.PSYCHOLOGICAL });
      onClose();
    } catch (error) {
      console.error("Error creating mistake:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    setFormData({ userId: "current-user", name: "", category: MistakeCategory.PSYCHOLOGICAL });
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div 
      className="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center z-50"
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
              onChange={(e) => setFormData({ ...formData, category: e.target.value as MistakeCategory })}
              className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica transition-colors"
            >
              <option value={MistakeCategory.PSYCHOLOGICAL}>Psychological</option>
              <option value={MistakeCategory.BEHAVIORAL}>Behavioral</option>
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
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // TODO: Get userId from authentication context
  const userId = 1;

  // Fetch data on component mount
  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    setError(null);
    try {
      // Fetch all data in parallel
      const [strategiesRes, rulesRes, mistakesRes] = await Promise.all([
        strategyApi.getAll(userId),
        ruleApi.getAll(userId),
        mistakeApi.getAll(userId),
      ]);

      // Transform API responses to match TypeScript interfaces
      // Use optional chaining and default to empty array to prevent null reference errors
      setStrategies(
        ((strategiesRes && strategiesRes.strategies) || []).map((s) => ({
          id: s.id,
          userId: s.user_id.toString(),
          name: s.name,
          description: s.description,
          createdAt: new Date(s.created_at),
          updatedAt: new Date(s.updated_at),
        }))
      );

      setRules(
        ((rulesRes && rulesRes.rules) || []).map((r) => ({
          id: r.id,
          userId: r.user_id.toString(),
          name: r.name,
          description: r.description,
          category: r.category as any,
          createdAt: new Date(r.created_at),
          updatedAt: new Date(r.updated_at),
        }))
      );

      setMistakes(
        ((mistakesRes && mistakesRes.mistakes) || []).map((m) => ({
          id: m.id,
          userId: m.user_id.toString(),
          name: m.name,
          category: m.category as any,
          createdAt: new Date(m.created_at),
          updatedAt: new Date(m.updated_at),
        }))
      );
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch data");
      console.error("Error fetching data:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleAddStrategy = async (data: IStrategyFormData) => {
    setError(null);
    try {
      const response = await strategyApi.create(userId, {
        name: data.name,
        description: data.description,
      });

      const newStrategy: IStrategy = {
        id: response.id,
        userId: response.user_id.toString(),
        name: response.name,
        description: response.description,
        createdAt: new Date(response.created_at),
        updatedAt: new Date(response.updated_at),
      };
      setStrategies([...strategies, newStrategy]);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to create strategy";
      setError(errorMessage);
      throw err;
    }
  };

  const handleAddRule = async (data: IRuleCreateRequest) => {
    setError(null);
    try {
      const response = await ruleApi.create(userId, {
        name: data.name,
        description: data.description,
        category: data.category as string,
      });

      const newRule: IRule = {
        id: response.id,
        userId: response.user_id.toString(),
        name: response.name,
        description: response.description,
        category: response.category as any,
        createdAt: new Date(response.created_at),
        updatedAt: new Date(response.updated_at),
      };
      setRules([...rules, newRule]);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to create rule";
      setError(errorMessage);
      throw err;
    }
  };

  const handleAddMistake = async (data: IMistakeCreateRequest) => {
    setError(null);
    try {
      const response = await mistakeApi.create(userId, {
        name: data.name,
        category: data.category as string,
      });

      const newMistake: IMistake = {
        id: response.id,
        userId: response.user_id.toString(),
        name: response.name,
        category: response.category as any,
        createdAt: new Date(response.created_at),
        updatedAt: new Date(response.updated_at),
      };
      setMistakes([...mistakes, newMistake]);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to create mistake";
      setError(errorMessage);
      throw err;
    }
  };

  const tabs = [
    { id: "strategies", label: "Strategies", icon: Target },
    { id: "rules", label: "Rules", icon: BookOpen },
    { id: "mistakes", label: "Mistakes", icon: AlertTriangle },
  ];

  const formatDate = (date: Date) => {
    return new Intl.DateTimeFormat("en-IN", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
    }).format(date);
  };

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
              <div className="bg-primary border border-primary rounded-lg overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b border-primary">
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Name</th>
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Description</th>
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Created</th>
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Updated</th>
                      </tr>
                    </thead>
                    <tbody>
                      {strategies.map((strategy) => (
                        <tr
                          key={strategy.id}
                          className="border-b border-primary hover:bg-tertiary transition-colors duration-200"
                        >
                          <td className="py-3 px-4">
                            <span className="font-helvetica-bold text-primary">
                              {strategy.name}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-primary font-helvetica-light text-sm">
                              {strategy.description}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-primary font-helvetica-light text-sm">
                              {formatDate(strategy.createdAt)}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-primary font-helvetica-light text-sm">
                              {formatDate(strategy.updatedAt)}
                            </span>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
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
              <div className="bg-primary border border-primary rounded-lg overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b border-primary">
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Name</th>
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Description</th>
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Category</th>
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Created</th>
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Updated</th>
                      </tr>
                    </thead>
                    <tbody>
                      {rules.map((rule) => (
                        <tr
                          key={rule.id}
                          className="border-b border-primary hover:bg-tertiary transition-colors duration-200"
                        >
                          <td className="py-3 px-4">
                            <span className="font-helvetica-bold text-primary">
                              {rule.name}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-primary font-helvetica-light text-sm">
                              {rule.description}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="inline-block bg-accent text-primary text-xs px-2 py-1 rounded font-helvetica-medium capitalize">
                              {rule.category.replace(/_/g, " ")}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-primary font-helvetica-light text-sm">
                              {formatDate(rule.createdAt)}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-primary font-helvetica-light text-sm">
                              {formatDate(rule.updatedAt)}
                            </span>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
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
              <div className="bg-primary border border-primary rounded-lg overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b border-primary">
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Name</th>
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Category</th>
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Created</th>
                        <th className="text-left py-3 px-4 text-sm font-helvetica-medium text-tertiary">Updated</th>
                      </tr>
                    </thead>
                    <tbody>
                      {mistakes.map((mistake) => (
                        <tr
                          key={mistake.id}
                          className="border-b border-primary hover:bg-tertiary transition-colors duration-200"
                        >
                          <td className="py-3 px-4">
                            <span className="font-helvetica-bold text-primary">
                              {mistake.name}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="inline-block bg-accent text-primary text-xs px-2 py-1 rounded font-helvetica-medium capitalize">
                              {mistake.category}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-primary font-helvetica-light text-sm">
                              {formatDate(mistake.createdAt)}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-primary font-helvetica-light text-sm">
                              {formatDate(mistake.updatedAt)}
                            </span>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
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
      {/* Error Message */}
      {error && (
        <div className="mb-4 p-4 bg-red-500/20 border border-red-500/30 rounded-lg text-red-400">
          <p className="font-helvetica-medium">{error}</p>
        </div>
      )}

      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-helvetica-bold text-primary">Configure</h1>
          <p className="text-secondary font-helvetica mt-1">Manage your trading strategies, rules, and mistakes</p>
        </div>
        
        {/* Tabs */}
        <div className="flex space-x-1 bg-primary p-1 rounded-lg">
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
            disabled={loading}
            className="flex items-center space-x-2 bg-accent hover:bg-secondary disabled:bg-tertiary disabled:cursor-not-allowed text-primary font-helvetica-medium py-2 px-4 rounded-lg transition-colors duration-200"
          >
            <Plus className="w-4 h-4" />
            <span>Add {activeTab.slice(0, -1)}</span>
          </button>
        </div>
        {loading ? (
          <div className="text-center py-12">
            <p className="text-tertiary font-helvetica">Loading...</p>
          </div>
        ) : (
          renderContent()
        )}
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
