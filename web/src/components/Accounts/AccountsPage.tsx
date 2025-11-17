'use client';

import { useState, useEffect } from 'react';
import { ChevronDown, ChevronUp, CheckCircle2, AlertCircle, RefreshCw } from 'lucide-react';
import { dhanApi } from '@/services/api';

// For now, using a hardcoded user ID. In production, this should come from auth context
const CURRENT_USER_ID = 1;

interface DhanConfig {
  configured: boolean;
  has_credentials?: boolean;
  dhan_client_id?: string;
  dhan_client_name?: string;
  expiry_time?: string;
}

export default function AccountsPage() {
  const [dhanConfig, setDhanConfig] = useState<DhanConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [step, setStep] = useState<'credentials' | 'oauth' | 'complete'>('credentials');
  const [apiKey, setApiKey] = useState('');
  const [apiSecret, setApiSecret] = useState('');
  const [dhanClientId, setDhanClientId] = useState('');
  const [saving, setSaving] = useState(false);
  const [oauthLoading, setOauthLoading] = useState(false);
  const [tokenId, setTokenId] = useState('');
  const [isDhanExpanded, setIsDhanExpanded] = useState(false);
  
  // Zerodha state
  const [isZerodhaExpanded, setIsZerodhaExpanded] = useState(false);
  const [zerodhaAppId, setZerodhaAppId] = useState('');
  const [zerodhaApiSecret, setZerodhaApiSecret] = useState('');
  const [zerodhaSaving, setZerodhaSaving] = useState(false);

  useEffect(() => {
    fetchDhanConfig();
  }, []);

  const fetchDhanConfig = async () => {
    try {
      setLoading(true);
      setError(null);
      const config = await dhanApi.getConfig(CURRENT_USER_ID);
      setDhanConfig(config);
      if (config.configured && config.dhan_client_id) {
        setStep('complete');
      } else if (config.has_credentials && !config.dhan_client_id) {
        // Has API key/secret but no client_id yet - need to provide client_id
        setStep('credentials');
      } else if (!config.has_credentials) {
        // No credentials at all
        setStep('credentials');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch Dhan configuration');
      console.error('Error fetching Dhan config:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSaveCredentials = async () => {
    if (!apiKey.trim() || !apiSecret.trim() || !dhanClientId.trim()) {
      setError('Please enter API Key, API Secret, and Dhan Client ID');
      return;
    }

    try {
      setSaving(true);
      setError(null);

      await dhanApi.saveCredentials(CURRENT_USER_ID, apiKey.trim(), apiSecret.trim(), dhanClientId.trim());
      
      // Move to OAuth step
      setStep('oauth');
      setApiKey('');
      setApiSecret('');
      setDhanClientId('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save credentials');
      console.error('Error saving credentials:', err);
    } finally {
      setSaving(false);
    }
  };

  const handleStartOAuth = async () => {
    try {
      setOauthLoading(true);
      setError(null);

      // Step 1: Generate consent
      const consentResponse = await dhanApi.generateConsent(CURRENT_USER_ID);
      
      // Step 2: Open browser for login
      const loginWindow = window.open(
        consentResponse.login_url,
        'Dhan Login',
        'width=600,height=700,scrollbars=yes,resizable=yes'
      );

      if (!loginWindow) {
        setError('Please allow popups to continue with authentication');
        setOauthLoading(false);
        return;
      }

      // Instructions for user
      alert('After logging in, you will be redirected. Please copy the tokenId from the URL (the part after ?tokenId=) and paste it below.');

      // Wait for user to enter tokenId
      // In a production app, you'd create a callback page that extracts tokenId automatically
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to start OAuth flow');
      console.error('Error starting OAuth:', err);
    } finally {
      setOauthLoading(false);
    }
  };

  const handleConsumeConsent = async () => {
    if (!tokenId.trim()) {
      setError('Please enter the tokenId from the redirect URL');
      return;
    }

    try {
      setOauthLoading(true);
      setError(null);

      // Step 3: Consume consent
      await dhanApi.consumeConsent(CURRENT_USER_ID, tokenId.trim());
      
      // Refresh config
      await fetchDhanConfig();
      setTokenId('');
      setStep('complete');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to complete authentication');
      console.error('Error consuming consent:', err);
    } finally {
      setOauthLoading(false);
    }
  };

  const handleRenewToken = async () => {
    try {
      setOauthLoading(true);
      setError(null);

      await dhanApi.renewToken(CURRENT_USER_ID);
      
      // Refresh config
      await fetchDhanConfig();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to renew token');
      console.error('Error renewing token:', err);
    } finally {
      setOauthLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="p-6">
        <div className="text-primary font-helvetica">Loading...</div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-2xl font-helvetica-bold text-primary mb-2">Accounts</h1>
        <p className="text-tertiary font-helvetica-light">
          Connect and manage your broker accounts
        </p>
      </div>

      {error && (
        <div className="mb-4 p-4 bg-red-500/20 border border-red-500/30 rounded-lg text-red-400">
          <p className="font-helvetica-medium">{error}</p>
        </div>
      )}

      {/* Dhan Broker Card */}
      <div className="bg-primary border border-primary rounded-lg overflow-hidden mb-4 transition-all duration-200 max-w-2xl">
        {/* Card Header - Always Visible */}
        <div 
          className="p-4 cursor-pointer hover:bg-tertiary transition-colors duration-200"
          onClick={() => setIsDhanExpanded(!isDhanExpanded)}
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              {/* Dhan Logo */}
              <div className="w-12 h-12 bg-secondary border border-primary rounded-lg flex items-center justify-center overflow-hidden flex-shrink-0">
                <img
                  src="/brokers/dhan.jpeg"
                  alt="Dhan"
                  className="w-full h-full object-contain"
                />
              </div>
              
              <div className="min-w-0">
                <h2 className="text-lg font-helvetica-bold text-primary">Dhan</h2>
                {dhanConfig?.configured && step === 'complete' ? (
                  <div className="flex items-center gap-2">
                    <CheckCircle2 className="w-3.5 h-3.5 text-green-500 flex-shrink-0" />
                    <span className="text-green-500 font-helvetica-medium text-xs">Configured</span>
                    {dhanConfig.dhan_client_name && (
                      <span className="text-tertiary font-helvetica text-xs truncate">
                        • {dhanConfig.dhan_client_name}
                      </span>
                    )}
                  </div>
                ) : (
                  <div className="flex items-center gap-2">
                    <AlertCircle className="w-3.5 h-3.5 text-yellow-500 flex-shrink-0" />
                    <span className="text-yellow-500 font-helvetica-medium text-xs">Not Configured</span>
                  </div>
                )}
              </div>
            </div>
            
            <div className="flex items-center gap-2 flex-shrink-0">
              {dhanConfig?.configured && step === 'complete' && (
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    handleRenewToken();
                  }}
                  disabled={oauthLoading}
                  className="flex items-center gap-1.5 px-3 py-1.5 bg-accent hover:bg-accent-hover text-inverse font-helvetica-medium text-sm rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <RefreshCw className={`w-3.5 h-3.5 text-inverse ${oauthLoading ? 'animate-spin' : ''}`} />
                  Renew
                </button>
              )}
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  setIsDhanExpanded(!isDhanExpanded);
                }}
                className="text-tertiary hover:text-primary transition-colors p-1"
              >
                {isDhanExpanded ? (
                  <ChevronUp className="w-4 h-4" />
                ) : (
                  <ChevronDown className="w-4 h-4" />
                )}
              </button>
            </div>
          </div>
        </div>

        {/* Expandable Content */}
        {isDhanExpanded && (
          <div className="border-t border-primary p-4 bg-secondary/50">

            {step === 'credentials' && (
              <div className="space-y-3">
                <div className="p-3 bg-primary/30 border border-primary rounded-lg">
                  <p className="text-primary font-helvetica-medium text-xs mb-2">
                    <strong>Step 1: Save API Credentials</strong>
                  </p>
                  <ol className="list-decimal list-inside text-tertiary font-helvetica text-xs space-y-1 ml-1">
                    <li>Log in to web.dhan.co</li>
                    <li>Go to My Profile → Access DhanHQ APIs</li>
                    <li>Toggle to 'API key' and enter your app name</li>
                    <li>Copy your API Key, API Secret, and Dhan Client ID</li>
                    <li>Enter all three below</li>
                  </ol>
                </div>
                <div>
                  <label className="block text-primary font-helvetica-medium text-sm mb-1.5">
                    API Key
                  </label>
                  <input
                    type="text"
                    value={apiKey}
                    onChange={(e) => setApiKey(e.target.value)}
                    placeholder="Enter your API Key"
                    className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary text-sm focus:outline-none focus:ring-2 focus:ring-accent font-helvetica"
                    disabled={saving}
                  />
                </div>
                <div>
                  <label className="block text-primary font-helvetica-medium text-sm mb-1.5">
                    API Secret
                  </label>
                  <input
                    type="password"
                    value={apiSecret}
                    onChange={(e) => setApiSecret(e.target.value)}
                    placeholder="Enter your API Secret"
                    className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary text-sm focus:outline-none focus:ring-2 focus:ring-accent font-helvetica"
                    disabled={saving}
                  />
                </div>
                <div>
                  <label className="block text-primary font-helvetica-medium text-sm mb-1.5">
                    Dhan Client ID <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    value={dhanClientId}
                    onChange={(e) => setDhanClientId(e.target.value)}
                    placeholder="Enter your Dhan Client ID"
                    className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary text-sm focus:outline-none focus:ring-2 focus:ring-accent font-helvetica"
                    disabled={saving}
                  />
                  <p className="text-tertiary font-helvetica text-xs mt-1">
                    Your Dhan Client ID (found in your Dhan account/profile). This is required for authentication.
                  </p>
                </div>
                <button
                  onClick={handleSaveCredentials}
                  disabled={saving || !apiKey.trim() || !apiSecret.trim() || !dhanClientId.trim()}
                  className="w-full bg-accent hover:bg-accent-hover text-inverse font-helvetica-medium py-2.5 px-4 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm"
                >
                  {saving ? 'Saving...' : 'Save Credentials & Continue'}
                </button>
              </div>
            )}

            {step === 'oauth' && (
              <div className="space-y-3">
                <div className="p-3 bg-primary/30 border border-primary rounded-lg">
                  <p className="text-primary font-helvetica-medium text-xs mb-2">
                    <strong>Step 2: Authenticate with Dhan</strong>
                  </p>
                  <ol className="list-decimal list-inside text-tertiary font-helvetica text-xs space-y-1 ml-1">
                    <li>Click "Start Authentication" below</li>
                    <li>Log in with your Dhan credentials in the popup</li>
                    <li>After login, copy the tokenId from the redirect URL</li>
                    <li>Paste it below and click "Complete Authentication"</li>
                  </ol>
                </div>
                <button
                  onClick={handleStartOAuth}
                  disabled={oauthLoading}
                  className="w-full bg-accent hover:bg-accent-hover text-inverse font-helvetica-medium py-2.5 px-4 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm"
                >
                  {oauthLoading ? 'Starting...' : 'Start Authentication'}
                </button>
                <div>
                  <label className="block text-primary font-helvetica-medium text-sm mb-1.5">
                    Token ID (from redirect URL)
                  </label>
                  <input
                    type="text"
                    value={tokenId}
                    onChange={(e) => setTokenId(e.target.value)}
                    placeholder="Paste tokenId from redirect URL"
                    className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary text-sm focus:outline-none focus:ring-2 focus:ring-accent font-helvetica"
                    disabled={oauthLoading}
                  />
                </div>
                <button
                  onClick={handleConsumeConsent}
                  disabled={oauthLoading || !tokenId.trim()}
                  className="w-full bg-accent hover:bg-accent-hover text-inverse font-helvetica-medium py-2.5 px-4 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm"
                >
                  {oauthLoading ? 'Completing...' : 'Complete Authentication'}
                </button>
                <button
                  onClick={() => setStep('credentials')}
                  className="w-full text-tertiary hover:text-primary font-helvetica-medium text-xs py-1.5"
                >
                  ← Back to Credentials
                </button>
              </div>
            )}

            {step === 'complete' && (
              <div className="space-y-3">
                {dhanConfig?.dhan_client_id && (
                  <div className="p-3 bg-primary/30 border border-primary rounded-lg space-y-1.5">
                    <p className="text-tertiary font-helvetica text-xs">
                      <span className="font-medium text-primary">Client ID:</span> {dhanConfig.dhan_client_id}
                    </p>
                    {dhanConfig.expiry_time && (
                      <p className="text-tertiary font-helvetica text-xs">
                        <span className="font-medium text-primary">Token Expires:</span> {new Date(dhanConfig.expiry_time).toLocaleString()}
                      </p>
                    )}
                  </div>
                )}
                <button
                  onClick={() => {
                    setDhanConfig({ configured: false });
                    setStep('credentials');
                    setApiKey('');
                    setApiSecret('');
                    setTokenId('');
                    setIsDhanExpanded(true);
                  }}
                  className="w-full text-red-500 hover:text-red-600 font-helvetica-medium text-xs py-2 border border-red-500/30 hover:bg-red-500/10 rounded-lg transition-colors"
                >
                  Disconnect
                </button>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Zerodha Broker Card */}
      <div className="bg-primary border border-primary rounded-lg overflow-hidden mb-4 transition-all duration-200 max-w-2xl">
        {/* Card Header - Always Visible */}
        <div 
          className="p-4 cursor-pointer hover:bg-tertiary transition-colors duration-200"
          onClick={() => setIsZerodhaExpanded(!isZerodhaExpanded)}
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              {/* Zerodha Logo */}
              <div className="w-12 h-12 bg-secondary border border-primary rounded-lg flex items-center justify-center overflow-hidden flex-shrink-0">
                <img
                  src="/brokers/zerodha.jpeg"
                  alt="Zerodha"
                  className="w-full h-full object-contain"
                />
              </div>
              
              <div className="min-w-0">
                <h2 className="text-lg font-helvetica-bold text-primary">Zerodha</h2>
                <div className="flex items-center gap-2">
                  <AlertCircle className="w-3.5 h-3.5 text-yellow-500 flex-shrink-0" />
                  <span className="text-yellow-500 font-helvetica-medium text-xs">Not Configured</span>
                </div>
              </div>
            </div>
            
            <div className="flex items-center gap-2 flex-shrink-0">
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  setIsZerodhaExpanded(!isZerodhaExpanded);
                }}
                className="text-tertiary hover:text-primary transition-colors p-1"
              >
                {isZerodhaExpanded ? (
                  <ChevronUp className="w-4 h-4" />
                ) : (
                  <ChevronDown className="w-4 h-4" />
                )}
              </button>
            </div>
          </div>
        </div>

        {/* Expandable Content */}
        {isZerodhaExpanded && (
          <div className="border-t border-primary p-4 bg-secondary/50">
            <div className="space-y-3">
              <div className="p-3 bg-primary/30 border border-primary rounded-lg">
                <p className="text-primary font-helvetica-medium text-xs mb-2">
                  <strong>Configure Zerodha</strong>
                </p>
                <p className="text-tertiary font-helvetica text-xs">
                  Enter your Zerodha App ID and API Secret to connect your account.
                </p>
              </div>
              
              <div>
                <label className="block text-primary font-helvetica-medium text-sm mb-1.5">
                  App ID <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  value={zerodhaAppId}
                  onChange={(e) => setZerodhaAppId(e.target.value)}
                  placeholder="Enter your Zerodha App ID"
                  className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary text-sm focus:outline-none focus:ring-2 focus:ring-accent font-helvetica"
                  disabled={zerodhaSaving}
                />
              </div>
              
              <div>
                <label className="block text-primary font-helvetica-medium text-sm mb-1.5">
                  API Secret <span className="text-red-500">*</span>
                </label>
                <input
                  type="password"
                  value={zerodhaApiSecret}
                  onChange={(e) => setZerodhaApiSecret(e.target.value)}
                  placeholder="Enter your Zerodha API Secret"
                  className="w-full px-3 py-2 bg-primary border border-primary rounded-lg text-primary text-sm focus:outline-none focus:ring-2 focus:ring-accent font-helvetica"
                  disabled={zerodhaSaving}
                />
              </div>
              
              <button
                onClick={() => {
                  // TODO: Implement Zerodha credentials saving
                  console.log('Zerodha App ID:', zerodhaAppId);
                  console.log('Zerodha API Secret:', zerodhaApiSecret);
                }}
                disabled={zerodhaSaving || !zerodhaAppId.trim() || !zerodhaApiSecret.trim()}
                className="w-full bg-accent hover:bg-accent-hover text-inverse font-helvetica-medium py-2.5 px-4 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm"
              >
                {zerodhaSaving ? 'Saving...' : 'Save Credentials'}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
