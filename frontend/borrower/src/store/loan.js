import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
export const useLoanStore = create(
	persist(
		(set, get, api) => ({
			estimateLoan: {}, //官网提交的loan数据，用于计算rate

			clearEstimateLoan: () => {
				api.persist.removeStorage('estimateLoan'); // 清除持久化数据
				set({ estimateLoan: null }); // 重置内存状态
			},
			changeEstimateLoan: async (estimateLoan) => {
				set({ estimateLoan });
				return estimateLoan;
			},
		}),
		{
			name: 'loanStore', // 存储唯一标识（必填）
			storage: createJSONStorage(() => localStorage), // 默认 localStorage
			partialize: (state) => ({
				estimateLoan: state.estimateLoan,
			}),
		},
	),
);
