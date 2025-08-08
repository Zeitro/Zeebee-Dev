import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { createRequest } from '@/lib/fetch';
const { request } = createRequest();
export const useUserStore = create(
	persist(
		(set, get, api) => ({
			userInfo: {},
			access: {},
			login: async (values) => {
				const access = await request('/user/login', {
					method: 'POST',
					body: values,
				});
				set({ access });
				return access;
			},
			clearInfo: () => {
				api.persist.clearStorage(); // 清除持久化数据
				set({ access: null, userInfo: null }); // 重置内存状态
			},
			getUserInfo: async () => {
				const userInfo = await request('/user/info', {
					method: 'GET',
				});
				set({ userInfo });
			},
		}),
		{
			name: 'userStore', // 存储唯一标识（必填）
			storage: createJSONStorage(() => localStorage), // 默认 localStorage
			partialize: (state) => ({
				userInfo: state.userInfo,
				access: state.access,
			}),
		},
	),
);
