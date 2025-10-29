const isBrowser = typeof window !== 'undefined';
// localStorage
export const localStorageUtil = {
	set(key, value) {
		if (!isBrowser || !key) return;
		localStorage.setItem(key, JSON.stringify(value));
	},
	get(key) {
		if (!isBrowser || !key) return;
		const value = localStorage.getItem(key);
		if (value === 'undefined' || value === 'null') return null;
		return value ? JSON.parse(value) : null;
	},
	remove(key) {
		if (!isBrowser || !key) return;
		localStorage.removeItem(key);
	},
	clear() {
		if (!isBrowser) return;
		localStorage.clear();
	},
};

// sessionStorage
export const sessionStorageUtil = {
	set(key, value) {
		if (!isBrowser || !key) return;
		sessionStorage.setItem(key, JSON.stringify(value));
	},
	get(key) {
		if (!isBrowser || !key) return null;
		const value = sessionStorage?.getItem(key);
		if (value === 'undefined' || value === 'null') return null;
		return value ? JSON.parse(value) : null;
	},
	remove(key) {
		if (!isBrowser || !key) return;
		sessionStorage.removeItem(key);
	},
	clear() {
		if (!isBrowser) return;
		sessionStorage.clear();
	},
};
