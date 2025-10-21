export default function ColorDemo() {
  return (
    <div className="p-6 space-y-4">
      <h2 className="text-2xl font-helvetica-bold text-primary">Color System Demo</h2>
      <p className="text-secondary font-helvetica">Using Helvetica Neue font family across the application</p>
      
      {/* Background Colors */}
      <div className="space-y-2">
        <h3 className="text-lg font-semibold text-primary">Background Colors</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="bg-primary p-4 rounded-lg border border-primary">
            <p className="text-primary font-medium">Primary</p>
          </div>
          <div className="bg-secondary p-4 rounded-lg border border-primary">
            <p className="text-primary font-medium">Secondary</p>
          </div>
          <div className="bg-tertiary p-4 rounded-lg border border-primary">
            <p className="text-primary font-medium">Tertiary</p>
          </div>
          <div className="bg-accent p-4 rounded-lg border border-primary">
            <p className="text-primary font-medium">Accent</p>
          </div>
        </div>
      </div>

      {/* Text Colors */}
      <div className="space-y-2">
        <h3 className="text-lg font-semibold text-primary">Text Colors</h3>
        <div className="bg-primary p-4 rounded-lg border border-primary space-y-2">
          <p className="text-primary">Primary Text</p>
          <p className="text-secondary">Secondary Text</p>
          <p className="text-tertiary">Tertiary Text</p>
          <p className="text-muted">Muted Text</p>
        </div>
      </div>

      {/* Interactive Elements */}
      <div className="space-y-2">
        <h3 className="text-lg font-semibold text-primary">Interactive Elements</h3>
        <div className="flex gap-4">
          <button className="px-4 py-2 bg-secondary hover:bg-tertiary text-primary rounded-lg border border-primary transition-colors">
            Hover Me
          </button>
          <button className="px-4 py-2 bg-primary text-inverse rounded-lg border border-primary hover:bg-secondary transition-colors">
            Primary Button
          </button>
        </div>
      </div>

      {/* Status Colors */}
      <div className="space-y-2">
        <h3 className="text-lg font-semibold text-primary">Status Colors</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="p-3 rounded-lg" style={{ backgroundColor: 'var(--success-bg)', color: 'var(--success)' }}>
            <p className="font-medium">Success</p>
          </div>
          <div className="p-3 rounded-lg" style={{ backgroundColor: 'var(--warning-bg)', color: 'var(--warning)' }}>
            <p className="font-medium">Warning</p>
          </div>
          <div className="p-3 rounded-lg" style={{ backgroundColor: 'var(--error-bg)', color: 'var(--error)' }}>
            <p className="font-medium">Error</p>
          </div>
          <div className="p-3 rounded-lg" style={{ backgroundColor: 'var(--info-bg)', color: 'var(--info)' }}>
            <p className="font-medium">Info</p>
          </div>
        </div>
      </div>

      {/* Typography Demo */}
      <div className="space-y-2">
        <h3 className="text-lg font-semibold text-primary">Helvetica Typography</h3>
        <div className="bg-primary p-4 rounded-lg border border-primary space-y-2">
          <h1 className="text-4xl font-helvetica-bold text-primary">Heading 1 - Bold</h1>
          <h2 className="text-3xl font-helvetica-medium text-primary">Heading 2 - Medium</h2>
          <h3 className="text-2xl font-helvetica text-primary">Heading 3 - Regular</h3>
          <h4 className="text-xl font-helvetica-light text-primary">Heading 4 - Light</h4>
          <p className="text-base font-helvetica text-secondary">Body text - Regular weight</p>
          <p className="text-sm font-helvetica-light text-tertiary">Small text - Light weight</p>
        </div>
      </div>
    </div>
  );
}
