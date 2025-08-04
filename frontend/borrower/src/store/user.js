import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { createRequest } from '@/lib/fetch';
const { request } = createRequest();
export const useUserStore = create(
	persist(
		(set) => ({
			userInfo: {},
			login: (userData) => set({ userInfo: userData }),
			getUserInfo: async () => {
				const res = await request('/user/info', {
					method: 'GET',
				});
				set({ useInfo: res });
			},
		}),
		{
			name: 'userInfo', // 存储唯一标识（必填）
			storage: createJSONStorage(() => localStorage), // 默认 localStorage
			partialize: (state) => (state.userInfo) 
		},
	),
);
