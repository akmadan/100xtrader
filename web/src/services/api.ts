// API service for communicating with the backend

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Helper function to handle API responses
async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Unknown error occurred' }));
    throw new Error(error.error || error.message || `HTTP error! status: ${response.status}`);
  }
  return response.json();
}

// Helper function to make API requests
async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  const config: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  const response = await fetch(url, config);
  return handleResponse<T>(response);
}

// Strategy API
export const strategyApi = {
  create: async (userId: number, data: { name: string; description: string }) => {
    const response = await apiRequest<{ message: string; data: { id: string; user_id: number; name: string; description: string; created_at: string; updated_at: string } }>(
      '/strategies',
      {
        method: 'POST',
        body: JSON.stringify({
          user_id: userId,
          name: data.name,
          description: data.description,
        }),
      }
    );
    return response.data;
  },

  getAll: async (userId: number, limit: number = 100, offset: number = 0) => {
    const response = await apiRequest<{ message: string; data: { strategies: Array<{ id: string; user_id: number; name: string; description: string; created_at: string; updated_at: string }>; pagination: { total: number; limit: number; offset: number } } }>(
      `/strategies/user/${userId}?limit=${limit}&offset=${offset}`
    );
    return response.data;
  },

  getById: async (id: string, userId: number) => {
    const response = await apiRequest<{ message: string; data: { id: string; user_id: number; name: string; description: string; created_at: string; updated_at: string } }>(
      `/strategies/${id}?user_id=${userId}`
    );
    return response.data;
  },

  update: async (id: string, userId: number, data: { name: string; description: string }) => {
    const response = await apiRequest<{ message: string; data: { id: string; user_id: number; name: string; description: string; created_at: string; updated_at: string } }>(
      `/strategies/${id}`,
      {
        method: 'PUT',
        body: JSON.stringify({
          id,
          user_id: userId,
          name: data.name,
          description: data.description,
        }),
      }
    );
    return response.data;
  },

  delete: async (id: string, userId: number) => {
    await apiRequest<{ message: string }>(
      `/strategies/${id}?user_id=${userId}`,
      {
        method: 'DELETE',
      }
    );
  },
};

// Rule API
export const ruleApi = {
  create: async (userId: number, data: { name: string; description: string; category: string }) => {
    const response = await apiRequest<{ message: string; data: { id: string; user_id: number; name: string; description: string; category: string; created_at: string; updated_at: string } }>(
      '/rules',
      {
        method: 'POST',
        body: JSON.stringify({
          user_id: userId,
          name: data.name,
          description: data.description,
          category: data.category,
        }),
      }
    );
    return response.data;
  },

  getAll: async (userId: number, limit: number = 100, offset: number = 0) => {
    const response = await apiRequest<{ message: string; data: { rules: Array<{ id: string; user_id: number; name: string; description: string; category: string; created_at: string; updated_at: string }>; pagination: { total: number; limit: number; offset: number } } }>(
      `/rules/user/${userId}?limit=${limit}&offset=${offset}`
    );
    return response.data;
  },

  getById: async (id: string, userId: number) => {
    const response = await apiRequest<{ message: string; data: { id: string; user_id: number; name: string; description: string; category: string; created_at: string; updated_at: string } }>(
      `/rules/${id}?user_id=${userId}`
    );
    return response.data;
  },

  update: async (id: string, userId: number, data: { name: string; description: string; category: string }) => {
    const response = await apiRequest<{ message: string; data: { id: string; user_id: number; name: string; description: string; category: string; created_at: string; updated_at: string } }>(
      `/rules/${id}`,
      {
        method: 'PUT',
        body: JSON.stringify({
          id,
          user_id: userId,
          name: data.name,
          description: data.description,
          category: data.category,
        }),
      }
    );
    return response.data;
  },

  delete: async (id: string, userId: number) => {
    await apiRequest<{ message: string }>(
      `/rules/${id}?user_id=${userId}`,
      {
        method: 'DELETE',
      }
    );
  },
};

// Mistake API
export const mistakeApi = {
  create: async (userId: number, data: { name: string; category: string }) => {
    const response = await apiRequest<{ message: string; data: { id: string; user_id: number; name: string; category: string; created_at: string; updated_at: string } }>(
      '/mistakes',
      {
        method: 'POST',
        body: JSON.stringify({
          user_id: userId,
          name: data.name,
          category: data.category,
        }),
      }
    );
    return response.data;
  },

  getAll: async (userId: number, limit: number = 100, offset: number = 0) => {
    const response = await apiRequest<{ message: string; data: { mistakes: Array<{ id: string; user_id: number; name: string; category: string; created_at: string; updated_at: string }>; pagination: { total: number; limit: number; offset: number } } }>(
      `/mistakes/user/${userId}?limit=${limit}&offset=${offset}`
    );
    return response.data;
  },

  getById: async (id: string, userId: number) => {
    const response = await apiRequest<{ message: string; data: { id: string; user_id: number; name: string; category: string; created_at: string; updated_at: string } }>(
      `/mistakes/${id}?user_id=${userId}`
    );
    return response.data;
  },

  update: async (id: string, userId: number, data: { name: string; category: string }) => {
    const response = await apiRequest<{ message: string; data: { id: string; user_id: number; name: string; category: string; created_at: string; updated_at: string } }>(
      `/mistakes/${id}`,
      {
        method: 'PUT',
        body: JSON.stringify({
          id,
          user_id: userId,
          name: data.name,
          category: data.category,
        }),
      }
    );
    return response.data;
  },

  delete: async (id: string, userId: number) => {
    await apiRequest<{ message: string }>(
      `/mistakes/${id}?user_id=${userId}`,
      {
        method: 'DELETE',
      }
    );
  },
};

// Trade API
export const tradeApi = {
  create: async (userId: number, data: {
    symbol: string;
    market_type: string;
    entry_date: string;
    entry_price: number;
    quantity: number;
    total_amount: number;
    exit_price?: number;
    direction: string;
    stop_loss?: number;
    target?: number;
    strategy: string;
    outcome_summary: string;
    trade_analysis?: string;
    rules_followed?: string[];
    screenshots?: string[];
    psychology?: {
      entry_confidence: number;
      satisfaction_rating: number;
      emotional_state: string;
      mistakes_made?: string[];
      lessons_learned?: string;
    };
    // Broker-specific fields (optional, for imported trades)
    trading_broker?: string;
    trader_broker_id?: string;
    exchange_order_id?: string;
    order_id?: string;
    product_type?: string;
    transaction_type?: string;
  }) => {
    const response = await apiRequest<{
      message: string;
      data: {
        id: string;
        user_id: number;
        symbol: string;
        market_type: string;
        entry_date: string;
        entry_price: number;
        quantity: number;
        total_amount: number;
        exit_price?: number;
        direction: string;
        stop_loss?: number;
        target?: number;
        strategy: string;
        outcome_summary: string;
        trade_analysis?: string;
        rules_followed?: string[];
        screenshots?: string[];
        psychology?: {
          entry_confidence: number;
          satisfaction_rating: number;
          emotional_state: string;
          mistakes_made?: string[];
          lessons_learned?: string;
        };
        // Broker-specific fields (optional, for imported trades)
        trading_broker?: string;
        trader_broker_id?: string;
        exchange_order_id?: string;
        order_id?: string;
        product_type?: string;
        transaction_type?: string;
        created_at: string;
        updated_at: string;
      };
    }>(
      '/trades',
      {
        method: 'POST',
        body: JSON.stringify({
          user_id: userId,
          ...data,
        }),
      }
    );
    return response.data;
  },

  getAll: async (userId: number, limit: number = 100, offset: number = 0) => {
    const response = await apiRequest<{
      message: string;
      data: {
        trades: Array<{
          id: string;
          user_id: number;
          symbol: string;
          market_type: string;
          entry_date: string;
          entry_price: number;
          quantity: number;
          total_amount: number;
          exit_price?: number;
          direction: string;
          stop_loss?: number;
          target?: number;
          strategy: string;
          outcome_summary: string;
          trade_analysis?: string;
          rules_followed?: string[];
          screenshots?: string[];
          psychology?: {
            entry_confidence: number;
            satisfaction_rating: number;
            emotional_state: string;
            mistakes_made?: string[];
            lessons_learned?: string;
          };
          // Broker-specific fields (optional, for imported trades)
          trading_broker?: string;
          trader_broker_id?: string;
          exchange_order_id?: string;
          order_id?: string;
          product_type?: string;
          transaction_type?: string;
          created_at: string;
          updated_at: string;
        }>;
        pagination: { total: number; limit: number; offset: number };
      };
    }>(
      `/trades/user/${userId}?limit=${limit}&offset=${offset}`
    );
    return response.data;
  },

  getById: async (id: string, userId: number) => {
    const response = await apiRequest<{
      message: string;
      data: {
        id: string;
        user_id: number;
        symbol: string;
        market_type: string;
        entry_date: string;
        entry_price: number;
        quantity: number;
        total_amount: number;
        exit_price?: number;
        direction: string;
        stop_loss?: number;
        target?: number;
        strategy: string;
        outcome_summary: string;
        trade_analysis?: string;
        rules_followed?: string[];
        screenshots?: string[];
        psychology?: {
          entry_confidence: number;
          satisfaction_rating: number;
          emotional_state: string;
          mistakes_made?: string[];
          lessons_learned?: string;
        };
        // Broker-specific fields (optional, for imported trades)
        trading_broker?: string;
        trader_broker_id?: string;
        exchange_order_id?: string;
        order_id?: string;
        product_type?: string;
        transaction_type?: string;
        created_at: string;
        updated_at: string;
      };
    }>(
      `/trades/${id}?user_id=${userId}`
    );
    return response.data;
  },

  update: async (id: string, userId: number, data: {
    symbol: string;
    market_type: string;
    entry_date: string;
    entry_price: number;
    quantity: number;
    total_amount: number;
    exit_price?: number;
    direction: string;
    stop_loss?: number;
    target?: number;
    strategy: string;
    outcome_summary: string;
    trade_analysis?: string;
    rules_followed?: string[];
    screenshots?: string[];
    psychology?: {
      entry_confidence: number;
      satisfaction_rating: number;
      emotional_state: string;
      mistakes_made?: string[];
      lessons_learned?: string;
    };
    // Broker-specific fields (optional, for imported trades)
    trading_broker?: string;
    trader_broker_id?: string;
    exchange_order_id?: string;
    order_id?: string;
    product_type?: string;
    transaction_type?: string;
  }) => {
    const response = await apiRequest<{
      message: string;
      data: {
        id: string;
        user_id: number;
        symbol: string;
        market_type: string;
        entry_date: string;
        entry_price: number;
        quantity: number;
        total_amount: number;
        exit_price?: number;
        direction: string;
        stop_loss?: number;
        target?: number;
        strategy: string;
        outcome_summary: string;
        trade_analysis?: string;
        rules_followed?: string[];
        screenshots?: string[];
        psychology?: {
          entry_confidence: number;
          satisfaction_rating: number;
          emotional_state: string;
          mistakes_made?: string[];
          lessons_learned?: string;
        };
        // Broker-specific fields (optional, for imported trades)
        trading_broker?: string;
        trader_broker_id?: string;
        exchange_order_id?: string;
        order_id?: string;
        product_type?: string;
        transaction_type?: string;
        created_at: string;
        updated_at: string;
      };
    }>(
      `/trades/${id}`,
      {
        method: 'PUT',
        body: JSON.stringify({
          id,
          user_id: userId,
          ...data,
        }),
      }
    );
    return response.data;
  },

  delete: async (id: string, userId: number) => {
    await apiRequest<{ message: string }>(
      `/trades/${id}?user_id=${userId}`,
      {
        method: 'DELETE',
      }
    );
  },

  syncDhan: async (userId: number) => {
    const response = await apiRequest<{
      message: string;
      data: {
        saved_count: number;
        skipped_count: number;
        error_count: number;
        total_fetched: number;
        from_date: string;
        to_date: string;
        date_range: string;
      };
    }>(`/users/${userId}/trades/sync-dhan`, {
      method: 'POST',
    });
    return response.data;
  },
};

// Dhan API
export const dhanApi = {
  getConfig: async (userId: number) => {
    const response = await apiRequest<{
      message: string;
      data: {
        configured: boolean;
        has_credentials?: boolean;
        dhan_client_id?: string;
        dhan_client_name?: string;
        expiry_time?: string;
      };
    }>(`/users/${userId}/dhan/config`, {
      method: 'GET',
    });
    return response.data;
  },

  saveCredentials: async (userId: number, apiKey: string, apiSecret: string, dhanClientId?: string) => {
    const body: { api_key: string; api_secret: string; dhan_client_id?: string } = {
      api_key: apiKey,
      api_secret: apiSecret,
    };
    if (dhanClientId) {
      body.dhan_client_id = dhanClientId;
    }
    const response = await apiRequest<{
      message: string;
      data: {
        configured: boolean;
      };
    }>(`/users/${userId}/dhan/save-credentials`, {
      method: 'POST',
      body: JSON.stringify(body),
    });
    return response.data;
  },

  generateConsent: async (userId: number) => {
    const response = await apiRequest<{
      message: string;
      data: {
        consent_app_id: string;
        consent_app_status: string;
        status: string;
        login_url: string;
      };
    }>(`/users/${userId}/dhan/generate-consent`, {
      method: 'POST',
    });
    return response.data;
  },

  consumeConsent: async (userId: number, tokenId: string) => {
    const response = await apiRequest<{
      message: string;
      data: {
        dhan_client_id: string;
        dhan_client_name: string;
        dhan_client_ucc: string;
        given_power_of_attorney: boolean;
        access_token: string;
        expiry_time: string;
      };
    }>(`/users/${userId}/dhan/consume-consent`, {
      method: 'POST',
      body: JSON.stringify({ token_id: tokenId }),
    });
    return response.data;
  },

  renewToken: async (userId: number) => {
    const response = await apiRequest<{
      message: string;
      data: {
        status: string;
        access_token: string;
        expiry_time: string;
      };
    }>(`/users/${userId}/dhan/renew-token`, {
      method: 'POST',
      body: JSON.stringify({
        access_token: '',
        dhan_client_id: '',
      }),
    });
    return response.data;
  },


  consumeConsent: async (userId: number, tokenId: string) => {
    const response = await apiRequest<{
      message: string;
      data: {
        dhan_client_id: string;
        dhan_client_name: string;
        dhan_client_ucc: string;
        given_power_of_attorney: boolean;
        access_token: string;
        expiry_time: string;
      };
    }>(`/users/${userId}/dhan/consume-consent`, {
      method: 'POST',
      body: JSON.stringify({ token_id: tokenId }),
    });
    return response.data;
  },
};

