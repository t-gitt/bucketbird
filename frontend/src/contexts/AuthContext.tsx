import { useState, useEffect, useCallback } from 'react';
import type { ReactNode } from 'react';

import { resolveApiUrl } from '../lib/runtimeConfig';
import { readStoredJSON, removeStoredItem, writeStoredJSON } from '../lib/storage';
import { AuthContext } from './AuthContextValue';

interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  isReadonly: boolean;
  createdAt: string;
  updatedAt: string;
}

interface AuthTokens {
	accessToken: string;
	accessExpiry: number;
}

interface AuthState {
  user: User | null;
  tokens: AuthTokens | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

export interface AuthContextType extends AuthState {
  login: (email: string, password: string) => Promise<void>;
  demoLogin: () => Promise<void>;
  register: (email: string, password: string, firstName: string, lastName: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshToken: () => Promise<void>;
  updateUser: (updates: Partial<User>) => void;
}

const TOKEN_STORAGE_KEY = 'bucketbird_tokens';
const USER_STORAGE_KEY = 'bucketbird_user';

export function AuthProvider({ children }: { children: ReactNode }) {
	const [state, setState] = useState<AuthState>({
		user: null,
		tokens: null,
		isAuthenticated: false,
		isLoading: true,
	});

	const saveAuth = useCallback((user: User, tokens: AuthTokens) => {
		writeStoredJSON(USER_STORAGE_KEY, user);
		writeStoredJSON(TOKEN_STORAGE_KEY, tokens);
		setState({
			user,
			tokens,
			isAuthenticated: true,
			isLoading: false,
		});
	}, []);

	const clearAuth = useCallback(() => {
		removeStoredItem(USER_STORAGE_KEY);
		removeStoredItem(TOKEN_STORAGE_KEY);
		setState({
			user: null,
			tokens: null,
			isAuthenticated: false,
			isLoading: false,
		});
	}, []);

	const updateUser = (updates: Partial<User>) => {
		setState((prev) => {
			if (!prev.user) {
				return prev;
			}

			const nextUser = { ...prev.user, ...updates };
			writeStoredJSON(USER_STORAGE_KEY, nextUser);

			return {
				...prev,
				user: nextUser,
			};
		});
	};

	const requestRefresh = useCallback(async () => {
		const apiUrl = resolveApiUrl();
		const response = await fetch(`${apiUrl}/api/v1/auth/refresh`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({}),
		});

		if (!response.ok) {
			throw new Error('Refresh failed');
		}

		const data = await response.json();
		saveAuth(data.user, data.auth);
	}, [saveAuth]);

	const login = async (email: string, password: string) => {
		const apiUrl = resolveApiUrl();

		const response = await fetch(`${apiUrl}/api/v1/auth/login`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ email, password }),
		});

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Login failed');
    }

    const data = await response.json();
    saveAuth(data.user, data.auth);
  };

	const demoLogin = async () => {
		const apiUrl = resolveApiUrl();

		const response = await fetch(`${apiUrl}/api/v1/auth/demo`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
		});

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Demo login failed');
    }

    const data = await response.json();
    saveAuth(data.user, data.auth);
  };

	const register = async (email: string, password: string, firstName: string, lastName: string) => {
		const apiUrl = resolveApiUrl();

		const response = await fetch(`${apiUrl}/api/v1/auth/register`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ email, password, firstName, lastName }),
		});

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Registration failed');
    }

    const data = await response.json();
    saveAuth(data.user, data.auth);
  };

	const logout = async () => {
		const apiUrl = resolveApiUrl();

		try {
			await fetch(`${apiUrl}/api/v1/auth/logout`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'include',
				body: JSON.stringify({}),
			});
		} catch (error) {
			console.error('Logout request failed:', error);
		}

		clearAuth();
	};

	const refreshToken = useCallback(async () => {
		try {
			await requestRefresh();
		} catch (error) {
			console.error('Token refresh failed:', error);
			clearAuth();
		}
	}, [requestRefresh, clearAuth]);

	// Load session from storage or refresh via cookie on mount
	useEffect(() => {
		let cancelled = false;

		const restoreFromStorage = (): boolean => {
			const tokens = readStoredJSON<AuthTokens>(TOKEN_STORAGE_KEY);
			const user = readStoredJSON<User>(USER_STORAGE_KEY);
			if (tokens && user) {
				const accessExpiry = new Date(tokens.accessExpiry * 1000);
				if (accessExpiry > new Date()) {
					setState({
						user,
						tokens,
						isAuthenticated: true,
						isLoading: false,
					});
					return true;
				}
			}
			return false;
		};

		const hydrate = async () => {
			const restored = restoreFromStorage();
			if (restored || cancelled) {
				return;
			}
			try {
				await requestRefresh();
			} catch (error) {
				console.error('Failed to refresh session from cookie:', error);
				if (!cancelled) {
					setState({
						user: null,
						tokens: null,
						isAuthenticated: false,
						isLoading: false,
					});
				}
			}
		};

		hydrate();

		return () => {
			cancelled = true;
		};
	}, [requestRefresh]);

	// Auto-refresh token before expiry
	useEffect(() => {
		if (!state.tokens) return;

		const accessExpiry = new Date(state.tokens.accessExpiry * 1000);
		const now = new Date();
		const timeUntilExpiry = accessExpiry.getTime() - now.getTime();

		// Refresh 1 minute before expiry
		const refreshTime = timeUntilExpiry - 60000;

		if (refreshTime > 0) {
			const timeout = setTimeout(() => {
				refreshToken();
			}, refreshTime);

			return () => clearTimeout(timeout);
		}
	}, [state.tokens, refreshToken]);

  return (
    <AuthContext.Provider
      value={{
        ...state,
        login,
        demoLogin,
        register,
        logout,
        refreshToken,
        updateUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
