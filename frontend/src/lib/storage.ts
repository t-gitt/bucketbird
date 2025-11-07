const getSessionStorage = (): Storage | null => {
	if (typeof window === 'undefined') {
		return null;
	}
	try {
		return window.sessionStorage;
	} catch (error) {
		console.error('Session storage unavailable:', error);
		return null;
	}
};

const readStoredJSON = <T,>(key: string): T | null => {
	const storage = getSessionStorage();
	if (!storage) {
		return null;
	}
	const raw = storage.getItem(key);
	if (!raw) {
		return null;
	}
	try {
		return JSON.parse(raw) as T;
	} catch (error) {
		console.error(`Failed to parse stored value for ${key}:`, error);
		return null;
	}
};

const writeStoredJSON = (key: string, value: unknown) => {
	const storage = getSessionStorage();
	if (!storage) {
		return;
	}
	try {
		storage.setItem(key, JSON.stringify(value));
	} catch (error) {
		console.error(`Failed to store value for ${key}:`, error);
	}
};

const removeStoredItem = (key: string) => {
	const storage = getSessionStorage();
	if (!storage) {
		return;
	}
	try {
		storage.removeItem(key);
	} catch (error) {
		console.error(`Failed to remove stored value for ${key}:`, error);
	}
};

export { getSessionStorage, readStoredJSON, writeStoredJSON, removeStoredItem };
