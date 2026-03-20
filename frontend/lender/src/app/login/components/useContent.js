'use client';
import { useUserStore, usePermissionStore } from '@repo/shared-methods/store';
import { eventBus } from '@repo/shared-methods/utils';
import { useRouter } from 'next/navigation';
import MainModal from '@/app/setting/organizationProfile/components/Main';
import { menuSetting } from '@config/routes';
const useContent = () => {
	const router = useRouter();
	const { getUserInfo, getOrganization, setOrganization, getUserMenus } = useUserStore();
	const { setPermissionTags, setAccessPaths, updateAccess } = usePermissionStore();

	// 组织信息完整后登录
	const getUserPermissions = async (access, successCallback, org) => {
		await getUserInfo();
		const result = await getUserMenus(menuSetting);
		setPermissionTags(result?.permissionTags);
		setAccessPaths(result?.accessPaths);
		if (access) {
			updateAccess(access);
			// 	setCookie('token', access?.access_token, access?.access_token_expire * 1000);
		}
		successCallback?.(false);
		if (!org) return;
		const list = org?.organization_settings;
		const preSetting = list?.find((item) => item.purpose === 'preApprovedLetter');
		// 预授权信件未配置,跳转配置页面
		if (!preSetting?.has_configured) {
			const params = new URLSearchParams({ isSettingPreApprovedLetter: true });
			router.replace(`/setting/system?${params}`);
		} else {
			router.replace('/');
		}
	};
	// 拿到用户token后登录
	const loginByToken = async (access, successCallback) => {
		const org = await getOrganization();
		// 未填写组织信息, 弹出组织相关让其填写
		if (!org?.contact_email || !org?.contact_phone || !org?.contact_name) {
			// 清空本地存储，避免未填写orgProfile信息，刷新页面后就登录进去了
			setOrganization(null);
			successCallback?.(false);
			eventBus.emit(
				'appendPortalComponent',
				<MainModal
					rowData={org}
					setOpen={() => eventBus.emit('removePortalComponent')}
					successCallback={async (setLoadingFn) => {
						let newOrg = await getOrganization();
						getUserPermissions(access, setLoadingFn, newOrg);
					}}
				/>,
			);
			return;
		}
		getUserPermissions(access, successCallback, org);
	};
	return {
		loginByToken,
	};
};
export default useContent;
