export const getLocalStorage = (key) => {
    if (!key) return '';
    let value = localStorage.getItem(key);
    if (!value) return '';
    if (value.indexOf('[') !== -1 || value.indexOf('}') !== -1) {
        return JSON.parse(value);
    }
    return value;
};

export const setLocalStorage = (key, value) => {
    if (!key && !value) return '';
    let str = '';
    if (typeof value !== 'string') str = JSON.stringify(value);
    else str = value;
    localStorage.setItem(key, str);
};

export const removeLocalStorage = (key) => {
    if (!key) return '';
    localStorage.removeItem(key);
};

export const clearAll = () => {
    localStorage.clear();
};

// 清空用户信息
export const clearUserInfo = () => {
    removeLocalStorage('userInfo');
    removeLocalStorage('token');
};
