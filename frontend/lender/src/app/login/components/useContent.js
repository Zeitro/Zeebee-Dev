'use client';
import { useUserStore } from '@/store/user';
import { setCookie, eventBus } from '@repo/shared-methods/utils';
import { useRouter } from 'next/navigation';
import MainModal from '@/app/setting/organizationProfile/components/Main';
const useContent = () => {
	const router = useRouter();
	const { getUserInfo, getOrganization, setOrganization } = useUserStore();

	// 组织信息完整后登录
	const getInfoGoLogin = async (access, setLoadingFn, org) => {
		await getUserInfo();
		setCookie('token', access?.access_token, access?.access_token_expire * 1000);
		setLoadingFn?.(false);
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
	const loginByToken = async (access, options) => {
		const org = await getOrganization();
		// 未填写组织信息, 弹出组织相关让其填写
		if (!org?.contact_email || !org?.contact_phone || !org?.contact_name) {
			// 清空本地存储，避免未填写orgProfile信息，刷新页面后就登录进去了
			setOrganization(null);
			options?.setLoading?.(false);
			eventBus.emit(
				'appendPortalComponent',
				<MainModal
					rowData={org}
					setOpen={() => eventBus.emit('removePortalComponent')}
					successCallback={async (setLoadingFn) => {
						let newOrg = await getOrganization();
						getInfoGoLogin(access, setLoadingFn, newOrg);
					}}
				/>,
			);
			return;
		}
		getInfoGoLogin(access, options?.setLoading, org);
	};
	return {
		loginByToken,
	};
};
export default useContent;
