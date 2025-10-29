import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
const STORAGE_KEY = 'loanStore';
// const EXPIRE_MS = 24 * 60 * 60 * 1000; // 24 hours
const EXPIRE_MS = 2 * 60 * 1000; // 2 minutes
export const useLoanStore = create(
	persist(
		(set, get, api) => ({
			estimateLoan: {}, //官网提交的loan数据，用于计算rate
			preApprovedResult: {},
			clearEstimateLoan: () => {
				api.persist.removeStorage('estimateLoan'); // 清除持久化数据
				set({ estimateLoan: null }); // 重置内存状态
			},
			changeEstimateLoan: async (estimateLoan) => {
				set({ estimateLoan });
				return estimateLoan;
			},
			savePreApprovedResult: async (primary) => {
				set({ preApprovedResult: { primary, updatedAt: Date.now() } });
			},
			clearPreApprovedResult: () => {
				set({ preApprovedResult: {} });
			},
			clearExpiredPreApprovedInStorage: async () => {
				const result = await get().preApprovedResult;
				if (result?.updatedAt && Date.now() - result?.updatedAt > EXPIRE_MS) {
					get()?.clearPreApprovedResult?.(); // 清除持久化数据
				}
			},
		}),
		{
			name: STORAGE_KEY, // 存储唯一标识（必填）
			storage: createJSONStorage(() => localStorage), // 默认 localStorage
			partialize: (state) => ({
				estimateLoan: state.estimateLoan,
				preApprovedResult: state.preApprovedResult,
			}),
		},
	),
);
