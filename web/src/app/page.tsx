export default function Home() {
  return (
    <div className="p-6">
      <div className="space-y-6">
        <h1 className="text-3xl font-helvetica-bold text-primary">Welcome to 100xTrader</h1>
        <p className="text-secondary font-helvetica">
          Your comprehensive trading journal for tracking trades, setups, and market analysis.
        </p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-secondary p-6 rounded-lg border border-primary">
            <h3 className="text-lg font-helvetica-medium text-primary mb-2">Recent Trades</h3>
            <p className="text-tertiary font-helvetica">No trades yet</p>
          </div>
          <div className="bg-secondary p-6 rounded-lg border border-primary">
            <h3 className="text-lg font-helvetica-medium text-primary mb-2">Performance</h3>
            <p className="text-tertiary font-helvetica">Track your progress</p>
          </div>
          <div className="bg-secondary p-6 rounded-lg border border-primary">
            <h3 className="text-lg font-helvetica-medium text-primary mb-2">Insights</h3>
            <p className="text-tertiary font-helvetica">Analyze your patterns</p>
          </div>
        </div>
      </div>
    </div>
  );
}
