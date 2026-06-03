import { localStorageUtil } from '../lib/local';
import { message } from 'antd';
import Cookie from 'js-cookie';

// 公共提示
function showMessage(msg, type = 'error') {
	message.destroy();
	message.open({
		content: msg,
		type,
	});
}
const str = '/api/v1';
const baseURL = process.env.NEXT_PUBLIC_API_URL + str;

/** 稳定序列化对象，用于生成取消用的默认 key（与属性顺序无关） */
function stableStringify(value) {
	if (value === null || value === undefined) {
		return 'null';
	}
	const t = typeof value;
	if (t !== 'object') {
		return JSON.stringify(value);
	}
	if (Array.isArray(value)) {
		return `[${value.map(stableStringify).join(',')}]`;
	}
	const keys = Object.keys(value).sort();
	return `{${keys.map((k) => `${JSON.stringify(k)}:${stableStringify(value[k])}`).join(',')}}`;
}

function fingerprintBody(body) {
	if (body === undefined || body === null) {
		return '';
	}
	if (typeof body === 'string') {
		try {
			const parsed = JSON.parse(body);
			if (typeof parsed === 'object' && parsed !== null) {
				return stableStringify(parsed);
			}
		} catch {
			return body;
		}
		return body;
	}
	if (typeof body === 'object') {
		if (typeof FormData !== 'undefined' && body instanceof FormData) {
			return '[FormData]';
		}
		if (typeof Blob !== 'undefined' && body instanceof Blob) {
			return '[Blob]';
		}
		return stableStringify(body);
	}
	return String(body);
}

/** 未传 cancelKey 时：method + url + body 指纹，区分同 URL 不同参数与不同方法 */
export function computeDefaultRequestKey(url, options = {}) {
	const method = String(options.method || 'GET').toUpperCase();
	return `${method}\u0000${url}\u0000${fingerprintBody(options.body)}`;
}

export function createRequest() {
	const controllers = new Map();

	function addController(key, controller) {
		let set = controllers.get(key);
		if (!set) {
			set = new Set();
			controllers.set(key, set);
		}
		set.add(controller);
	}

	function removeController(key, controller) {
		const set = controllers.get(key);
		if (!set) return;
		set.delete(controller);
		if (set.size === 0) {
			controllers.delete(key);
		}
	}

	const commonJudge = async ({ ok, isExternalLink, fetchRes, xhr }) => {
		if (!ok) {
			return {
				message: xhr?.statusText || 'System server failure!',
			};
		}
		let data;
		let contentType;
		if (fetchRes) {
			contentType = fetchRes?.headers?.get?.('Content-Type');
			if (contentType?.indexOf('application/json') === -1) {
				return {
					res: fetchRes,
				};
			}
			data = await fetchRes?.json?.();
		}
		if (xhr?.responseText) {
			data = JSON.parse(xhr?.responseText || '{}');
			contentType = data?.headers?.get?.('Content-type');
		}

		if (!data && isExternalLink) {
			return true;
		}
		if (contentType?.indexOf('application/json') === -1) {
			return {
				res: fetchRes || data,
			};
		}

		if (data?.code && data?.code === 401) {
			window.location.href = '/login';
			showMessage('Please log in again!', 'warning');
			return false;
		}
		if (data?.code && data?.code !== 200) {
			console.log('err400', data);

			return {
				message: data.message || `fail:${data.code}`,
				code: data.code,
			};
		}
		return {
			res: data.data,
		};
	};

	/**
	 * @param {string} url
	 * @param {object} [options]
	 * @param {string} [cancelKey] 同一 key 共用一个「槽位」：新请求前 cancel(cancelKey) 可打断同槽旧请求；不传则按 method+url+body 自动区分
	 */
	async function request(url, options = {}, cancelKey) {
		const controller = new AbortController();
		const key =
			cancelKey !== undefined && cancelKey !== null
				? String(cancelKey)
				: computeDefaultRequestKey(url, options);
		addController(key, controller);
		const {
			body,
			method = 'GET',
			isExternalLink = false,
			headers,
			isHideTips,
			...restConfig
		} = options;
		const permissionStore = JSON.parse(
			Cookie.get(`permissionStore_${process.env.NEXT_PUBLIC_SYSTEM_TYPE}`) || '{}',
		);
		const token = permissionStore?.state?.access?.access_token;
		const headersConfig = {
			'Content-Type': 'application/json',
			Authorization: '',
			...(headers || {}),
		};
		if (token && !isExternalLink) {
			headersConfig.Authorization = `Bearer ${token}`;
		}
		const fetchOptions = {
			...restConfig,
			method: method,
			signal: controller.signal,
			headers: headersConfig,
		};
		if (body instanceof Object) {
			fetchOptions.body = JSON.stringify(body);
		}
		let response;
		const finalURL =
			process.env.NEXT_PUBLIC_ENV !== 'local' ? `${baseURL}${url}` : `${str}${url}`;
		try {
			response = await fetch(isExternalLink ? url : finalURL, fetchOptions);
		} catch (error) {
			if (error.name !== 'AbortError' && !isHideTips) {
				console.log(error, '----error----');
				showMessage(error.message || 'Network anomaly');
				return Promise.reject(error);
			}
			return;
		} finally {
			removeController(key, controller);
		}
		const data = await commonJudge({ ok: response?.ok, fetchRes: response, isExternalLink });

		if (!data) return Promise.reject({ message: '401' });
		if (data?.message) {
			if (!isHideTips) {
				showMessage(data.message);
			}
			return Promise.reject(data);
		}
		return data === true ? data : data.res;
	}

	// 上传需要得知进度时，只能用XMLHttpRequest，fetch现不支持
	function uploadWithProgress(url, options, cancelKey) {
		const controller = new AbortController();
		const key =
			cancelKey !== undefined && cancelKey !== null
				? String(cancelKey)
				: computeDefaultRequestKey(url, options);
		addController(key, controller);
		return new Promise((resolve, reject) => {
			const {
				onProgress,
				body,
				method = 'GET',
				isExternalLink = false,
				isHideTips,
			} = options;
			const xhr = new XMLHttpRequest();
			const finalURL =
				process.env.NEXT_PUBLIC_ENV !== 'local' ? `${baseURL}${url}` : `${str}${url}`;

			xhr.open(method, isExternalLink ? url : finalURL, true);

			// ✅ 3. 设置认证 Token Header
			const token = localStorageUtil.get('userStore')?.state?.access?.access_token;
			if (token && !isExternalLink) {
				xhr.setRequestHeader('Authorization', `Bearer ${token}`);
			}
			// ❗️ 注意：此处不能设置 'Content-Type'，必须由浏览器自动处理 FormData

			// 绑定 AbortController 到 XHR：controller.abort() 时触发 xhr.abort()
			const abortHandler = () => {
				xhr.abort();
			};
			if (controller.signal) {
				controller.signal.addEventListener('abort', abortHandler);
			}

			// 清理 helpers
			const cleanup = () => {
				removeController(key, controller);
				controller.signal.removeEventListener('abort', abortHandler);
			};

			// 当中止时，触发 reject（与 fetch 的行为一致）
			xhr.onabort = () => {
				cleanup();
				reject({ name: 'AbortError', message: 'The upload has been aborted.' });
			};

			// 监听上传进度
			xhr.upload.onprogress = (event) => {
				if (event.lengthComputable) {
					const percent = Math.round((event.loaded / event.total) * 100);
					onProgress(percent);
				}
			};

			// 监听请求完成
			xhr.onload = async () => {
				const data = await commonJudge({ ok: xhr.status === 200, xhr, isExternalLink });
				if (!data) {
					cleanup();
					return reject({ message: '401' });
				}
				if (data.message) {
					cleanup();
					if (!isHideTips) {
						showMessage(data.message, 'warning');
					}
					return reject(data);
				}
				cleanup();
				resolve(data === true ? data : data.res);
			};

			// 监听网络等其他错误
			xhr.onerror = () => {
				cleanup();
				showMessage('Network anomaly');
				reject(new Error('Upload request failed due to a network error'));
			};

			// 发起请求
			xhr.send(body);
		});
	}

	/** 精确匹配 cancelKey，取消该槽位上所有进行中的请求 */
	function cancel(key) {
		const set = controllers.get(String(key));
		if (!set) return;
		for (const c of set) {
			c.abort();
		}
		controllers.delete(String(key));
	}

	const cancelAll = () => {
		for (const set of controllers.values()) {
			for (const c of set) {
				c.abort();
			}
		}
		controllers.clear();
	};

	/** 取消所有「已注册 key」包含 partial 的请求（含精确等于） */
	function cancelMatch(partial) {
		if (partial == null || partial === '') {
			cancelAll();
			return;
		}
		const needle = String(partial);
		const keysToCancel = [];
		for (const k of controllers.keys()) {
			if (k === needle || k.includes(needle)) {
				keysToCancel.push(k);
			}
		}
		for (const k of keysToCancel) {
			cancel(k);
		}
	}

	return { request, cancel, cancelMatch, cancelAll, uploadWithProgress };
}
