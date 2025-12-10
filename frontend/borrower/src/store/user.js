import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { createRequest } from '@repo/shared-methods/utils';
const { request } = createRequest();
export const useUserStore = create(
	persist(
		(set, get, api) => ({
			userInfo: {},
			access: {},
			login: async (values, config = {}) => {
				const access = await request('/user/login', {
					method: 'POST',
					body: values,
					...config,
				});
				set({ access });
				return access;
			},
			updateAccess: (access) => {
				set({ access });
			},
			clearUserInfo: () => {
				api.persist.clearStorage(); // 清除持久化数据
				set({ access: null, userInfo: null }); // 重置内存状态
			},
			getUserInfo: async () => {
				const userInfo = await request('/user/info', {
					method: 'GET',
				});
				set({ userInfo });
				return userInfo;
			},
			sendResetPasswordEmail: async (email, config = {}) => {
				// 发送重置密码邮件
				return await request('/user/send_reset_password', {
					method: 'POST',
					body: { email },
					...config,
				});
			},
			resetPassword: async (body) => {
				// 重置密码,存储token
				const access = await request('/user/reset_password', {
					method: 'POST',
					body,
				});
				set({ access });
				return access;
			},
			registerAndSendEmail: async (values, config = {}) => {
				return await request('/user/register', {
					method: 'POST',
					body: values,
					...config,
				});
			},
			resendVerifyEmail: async (email, config = {}) => {
				// 重新发送验证邮件（例如：注册后，需验证邮箱，此时没验证就登录。触发此接口会再给用户发送邮件，提醒用户验证）
				return await request('/user/send_verify_email', {
					method: 'POST',
					body: { email },
					...config,
				});
			},
			verifyEmail: async (body) => {
				// 验证邮箱,存储token
				const access = await request('/user/verify_email', {
					method: 'POST',
					body,
				});
				set({ access });
				return access;
			},
			changePassword: async (body) => {
				// 修改密码
				await request('/user/change_password', {
					method: 'POST',
					body,
				});
				return true;
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
