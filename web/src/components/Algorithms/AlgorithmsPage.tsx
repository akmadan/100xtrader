"use client";

import { useState } from "react";
import { Plus, Code, Workflow } from "lucide-react";
import AlgoEditor from "./AlgoEditor";

export default function AlgorithmsPage() {
  const [isEditorOpen, setIsEditorOpen] = useState(false);

  return (
    <div className="h-full flex flex-col bg-primary">
      {/* Header */}
      <div className="bg-primary border-b border-primary px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-helvetica-bold text-primary mb-1">
              Algorithms
            </h1>
            <p className="text-tertiary font-helvetica-light text-sm">
              Create and deploy automated trading algorithms
            </p>
          </div>
          {!isEditorOpen && (
            <button
              onClick={() => setIsEditorOpen(true)}
              className="flex items-center gap-2 bg-accent hover:bg-accent-hover text-inverse font-helvetica-medium px-4 py-2 rounded-lg transition-colors duration-200"
            >
              <Plus className="w-4 h-4 text-inverse" />
              Create Algo
            </button>
          )}
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-hidden">
        {isEditorOpen ? (
          <AlgoEditor onClose={() => setIsEditorOpen(false)} />
        ) : (
          <div className="h-full flex items-center justify-center">
            <div className="text-center">
              <div className="w-16 h-16 bg-secondary border border-primary rounded-lg flex items-center justify-center mx-auto mb-4">
                <Code className="w-8 h-8 text-tertiary" />
              </div>
              <h3 className="text-lg font-helvetica-medium text-primary mb-2">
                No algorithms yet
              </h3>
              <p className="text-tertiary font-helvetica-light mb-6 max-w-md">
                Create your first trading algorithm to automate your trading strategy
              </p>
              <button
                onClick={() => setIsEditorOpen(true)}
                className="flex items-center gap-2 bg-accent hover:bg-accent-hover text-inverse font-helvetica-medium px-6 py-3 rounded-lg transition-colors duration-200 mx-auto"
              >
                <Plus className="w-5 h-5 text-inverse" />
                Create Your First Algorithm
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

