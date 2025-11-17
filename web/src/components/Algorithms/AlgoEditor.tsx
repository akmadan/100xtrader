"use client";

import { useState, useRef, useEffect } from "react";
import { X } from "lucide-react";
// import CanvasView from "./CanvasView"; // Commented out for now
import CodeView from "./CodeView";

interface AlgoEditorProps {
  onClose: () => void;
}

export default function AlgoEditor({ onClose }: AlgoEditorProps) {
  const [algorithmName, setAlgorithmName] = useState("New Algorithm");
  const [isEditingName, setIsEditingName] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isEditingName && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [isEditingName]);

  const handleNameBlur = () => {
    setIsEditingName(false);
    if (algorithmName.trim() === "") {
      setAlgorithmName("New Algorithm");
    }
  };

  const handleNameKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      setIsEditingName(false);
      if (algorithmName.trim() === "") {
        setAlgorithmName("New Algorithm");
      }
    } else if (e.key === "Escape") {
      setAlgorithmName("New Algorithm");
      setIsEditingName(false);
    }
  };

  return (
    <div className="h-full flex flex-col bg-primary">
      {/* Editor Header */}
      <div className="bg-primary border-b border-primary px-6 py-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            {isEditingName ? (
              <input
                ref={inputRef}
                type="text"
                value={algorithmName}
                onChange={(e) => setAlgorithmName(e.target.value)}
                onBlur={handleNameBlur}
                onKeyDown={handleNameKeyDown}
                className="text-lg font-helvetica-bold text-primary bg-transparent border-b-2 border-accent outline-none px-1 min-w-[200px]"
              />
            ) : (
              <h2
                onClick={() => setIsEditingName(true)}
                className="text-lg font-helvetica-bold text-primary cursor-text hover:text-accent transition-colors"
              >
                {algorithmName}
              </h2>
            )}
            
            {/* View Mode Toggle - Commented out for now */}
            {/* <div className="flex items-center gap-2 bg-secondary border border-primary rounded-lg p-1">
              <button
                onClick={() => setViewMode("canvas")}
                className={`flex items-center gap-2 px-3 py-1.5 rounded-md transition-colors duration-200 ${
                  viewMode === "canvas"
                    ? "bg-accent text-inverse"
                    : "text-tertiary hover:text-primary"
                }`}
              >
                <Workflow className={`w-4 h-4 ${viewMode === "canvas" ? "text-primary" : ""}`} />
                <span className="font-helvetica-medium text-sm">Canvas</span>
              </button>
              <button
                onClick={() => setViewMode("code")}
                className={`flex items-center gap-2 px-3 py-1.5 rounded-md transition-colors duration-200 ${
                  viewMode === "code"
                    ? "bg-accent text-inverse"
                    : "text-tertiary hover:text-primary"
                }`}
              >
                <Code className={`w-4 h-4 ${viewMode === "code" ? "text-primary" : ""}`} />
                <span className="font-helvetica-medium text-sm">Code</span>
              </button>
            </div> */}
          </div>

          <button
            onClick={onClose}
            className="p-2 rounded-lg hover:bg-secondary text-tertiary hover:text-primary transition-colors duration-200"
          >
            <X className="w-5 h-5" />
          </button>
        </div>
      </div>

      {/* Editor Content - Only Code View for now */}
      <div className="flex-1 overflow-hidden">
        <CodeView />
        {/* Canvas view commented out for now */}
        {/* {viewMode === "canvas" ? (
          <CanvasView />
        ) : (
          <CodeView />
        )} */}
      </div>
    </div>
  );
}

