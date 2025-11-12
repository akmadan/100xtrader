'use client';

import { useState, useEffect } from 'react';
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
      <h1 className="text-2xl font-helvetica-bold text-primary mb-6">Accounts</h1>

      {error && (
        <div className="mb-4 p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
          <p className="text-red-500 font-helvetica">{error}</p>
        </div>
      )}

      {/* Dhan Broker Card */}
      <div className="bg-secondary border border-primary rounded-lg p-6 mb-4">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="text-xl font-helvetica-bold text-primary mb-2">Dhan</h2>
            {dhanConfig?.configured && step === 'complete' ? (
              <div className="space-y-2">
                <p className="text-tertiary font-helvetica text-sm">
                  Status: <span className="text-green-500">Configured</span>
                </p>
                {dhanConfig.dhan_client_name && (
                  <p className="text-tertiary font-helvetica text-sm">
                    Client: {dhanConfig.dhan_client_name}
                  </p>
                )}
                {dhanConfig.dhan_client_id && (
                  <p className="text-tertiary font-helvetica text-sm">
                    Client ID: {dhanConfig.dhan_client_id}
                  </p>
                )}
                {dhanConfig.expiry_time && (
                  <p className="text-tertiary font-helvetica text-sm">
                    Token Expires: {new Date(dhanConfig.expiry_time).toLocaleString()}
                  </p>
                )}
              </div>
            ) : (
              <p className="text-tertiary font-helvetica text-sm">
                Status: <span className="text-yellow-500">Not Configured</span>
              </p>
            )}
          </div>
        </div>

        {step === 'credentials' && (
          <div className="space-y-4">
            <div className="mb-4 p-3 bg-primary/50 border border-primary rounded-lg">
              <p className="text-tertiary font-helvetica text-sm mb-2">
                <strong>Step 1: Save API Credentials</strong>
              </p>
              <ol className="list-decimal list-inside text-tertiary font-helvetica text-sm space-y-1">
                <li>Log in to web.dhan.co</li>
                <li>Go to My Profile → Access DhanHQ APIs</li>
                <li>Toggle to 'API key' and enter your app name</li>
                <li>Copy your API Key, API Secret, and Dhan Client ID</li>
                <li>Enter all three below</li>
              </ol>
            </div>
            <div>
              <label className="block text-primary font-helvetica-medium mb-2">
                API Key
              </label>
              <input
                type="text"
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
                placeholder="Enter your API Key"
                className="w-full px-4 py-2 bg-primary border border-primary rounded-lg text-primary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica"
                disabled={saving}
              />
            </div>
            <div>
              <label className="block text-primary font-helvetica-medium mb-2">
                API Secret
              </label>
              <input
                type="password"
                value={apiSecret}
                onChange={(e) => setApiSecret(e.target.value)}
                placeholder="Enter your API Secret"
                className="w-full px-4 py-2 bg-primary border border-primary rounded-lg text-primary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica"
                disabled={saving}
              />
            </div>
            <div>
              <label className="block text-primary font-helvetica-medium mb-2">
                Dhan Client ID <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                value={dhanClientId}
                onChange={(e) => setDhanClientId(e.target.value)}
                placeholder="Enter your Dhan Client ID"
                className="w-full px-4 py-2 bg-primary border border-primary rounded-lg text-primary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica"
                disabled={saving}
              />
              <p className="text-tertiary font-helvetica text-xs mt-1">
                Your Dhan Client ID (found in your Dhan account/profile). This is required for authentication.
              </p>
            </div>
            <button
              onClick={handleSaveCredentials}
              disabled={saving || !apiKey.trim() || !apiSecret.trim() || !dhanClientId.trim()}
              className="w-full bg-accent hover:bg-accent/80 text-primary font-helvetica-medium py-2 px-4 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {saving ? 'Saving...' : 'Save Credentials & Continue'}
            </button>
          </div>
        )}

        {step === 'oauth' && (
          <div className="space-y-4">
            <div className="mb-4 p-3 bg-primary/50 border border-primary rounded-lg">
              <p className="text-tertiary font-helvetica text-sm mb-2">
                <strong>Step 2: Authenticate with Dhan</strong>
              </p>
              <ol className="list-decimal list-inside text-tertiary font-helvetica text-sm space-y-1">
                <li>Click "Start Authentication" below</li>
                <li>Log in with your Dhan credentials in the popup</li>
                <li>After login, copy the tokenId from the redirect URL</li>
                <li>Paste it below and click "Complete Authentication"</li>
              </ol>
            </div>
            <button
              onClick={handleStartOAuth}
              disabled={oauthLoading}
              className="w-full bg-accent hover:bg-accent/80 text-primary font-helvetica-medium py-2 px-4 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {oauthLoading ? 'Starting...' : 'Start Authentication'}
            </button>
            <div>
              <label className="block text-primary font-helvetica-medium mb-2">
                Token ID (from redirect URL)
              </label>
              <input
                type="text"
                value={tokenId}
                onChange={(e) => setTokenId(e.target.value)}
                placeholder="Paste tokenId from redirect URL"
                className="w-full px-4 py-2 bg-primary border border-primary rounded-lg text-primary focus:outline-none focus:ring-2 focus:ring-accent font-helvetica"
                disabled={oauthLoading}
              />
            </div>
            <button
              onClick={handleConsumeConsent}
              disabled={oauthLoading || !tokenId.trim()}
              className="w-full bg-accent hover:bg-accent/80 text-primary font-helvetica-medium py-2 px-4 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {oauthLoading ? 'Completing...' : 'Complete Authentication'}
            </button>
            <button
              onClick={() => setStep('credentials')}
              className="w-full text-tertiary hover:text-primary font-helvetica-medium text-sm py-2"
            >
              ← Back to Credentials
            </button>
          </div>
        )}

        {step === 'complete' && (
          <div className="space-y-4">
            <button
              onClick={handleRenewToken}
              disabled={oauthLoading}
              className="w-full bg-accent hover:bg-accent/80 text-primary font-helvetica-medium py-2 px-4 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {oauthLoading ? 'Renewing...' : 'Renew Token'}
            </button>
            <button
              onClick={() => {
                setDhanConfig({ configured: false });
                setStep('credentials');
                setApiKey('');
                setApiSecret('');
                setTokenId('');
              }}
              className="w-full text-red-500 hover:text-red-600 font-helvetica-medium text-sm py-2"
            >
              Disconnect
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
