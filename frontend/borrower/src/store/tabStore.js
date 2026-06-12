import { create } from 'zustand';
import { getTabPolicy, TAB_MODE_SINGLE } from '@/components/TabWorkspace/tabPolicy';

/**
 * borrower 多标签页状态（P0：仅导航与 URL 同步，内容区仍由 Next children 渲染）。
 *
 * @typedef {object} TabItem
 * @property {string} key - fullPath（pathname + search）
 * @property {string} pathname
 * @property {string} search
 * @property {string} title
 * @property {boolean} [closable]
 */

export const useTabStore = create((set, get) => ({
	tabs: [],
	activeKey: '',

	/**
	 * URL 变化时同步 Tab 列表。
	 * 规则：已存在的 Tab 只激活不删除；仅「新打开 + 单 Tab 策略」时才整页替换。
	 */
	syncTabFromUrl(tab) {
		const { tabs } = get();
		const existing = tabs.find((item) => item.key === tab.key);

		// 切换到已打开的 Tab：只更新标题与 activeKey，保留其它 Tab
		if (existing) {
			set({
				tabs: tabs.map((item) =>
					item.key === tab.key ? { ...item, title: tab.title || item.title } : item,
				),
				activeKey: tab.key,
			});
			return;
		}

		const policy = getTabPolicy(tab.pathname);

		// 单 Tab 策略：新路径首次打开时替换全部（侧边栏进入「独占」页）
		if (policy.mode === TAB_MODE_SINGLE) {
			set({
				tabs: [{ ...tab, closable: false }],
				activeKey: tab.key,
			});
			return;
		}

		// 多 Tab：追加新 Tab，超出上限时移除最早可关闭的 Tab
		let nextTabs = [...tabs, { ...tab, closable: true }];
		if (nextTabs.length > policy.maxTabs) {
			const removableIndex = nextTabs.findIndex((item) => item.closable !== false);
			if (removableIndex > -1) {
				nextTabs = nextTabs.filter((_, index) => index !== removableIndex);
			} else {
				nextTabs = nextTabs.slice(-policy.maxTabs);
			}
		}

		set({ tabs: nextTabs, activeKey: tab.key });
	},

	/** eventBus.changeCrumbs 动态更新当前 Tab 标题 */
	updateTabTitle(key, title) {
		if (!key || !title) return;
		set({
			tabs: get().tabs.map((item) => (item.key === key ? { ...item, title } : item)),
		});
	},

	/** 关闭 Tab，返回应跳转的 key（由调用方 router.push） */
	closeTab(key) {
		const { tabs, activeKey } = get();
		const index = tabs.findIndex((item) => item.key === key);
		if (index === -1) return activeKey;

		const nextTabs = tabs.filter((item) => item.key !== key);
		if (!nextTabs.length) {
			set({ tabs: [], activeKey: '' });
			return '/dashboard';
		}

		let nextActiveKey = activeKey;
		if (activeKey === key) {
			const neighbor = nextTabs[Math.min(index, nextTabs.length - 1)];
			nextActiveKey = neighbor.key;
		}

		set({ tabs: nextTabs, activeKey: nextActiveKey });
		return nextActiveKey;
	},

	resetTabs() {
		set({ tabs: [], activeKey: '' });
	},
}));
