'use client';
import { useUserStore } from '@/store/user';
import { setCookie, eventBus } from '@repo/shared-methods/utils';
import { useRouter } from 'next/navigation';
import MainModal from '@/app/setting/organizationProfile/components/Main';
const useContent = () => {
	const router = useRouter();
	const { getUserInfo, getOrganization, setOrganization } = useUserStore();

	// 组织信息完整后登录
	const getInfoGoLogin = async (access, setLoadingFn) => {
		await getUserInfo();
		setCookie('token', access?.access_token, access?.access_token_expire * 1000);
		setLoadingFn?.(false);
		router.replace('/');
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
						await getOrganization();
						getInfoGoLogin(access, setLoadingFn);
					}}
				/>,
			);
			return;
		}
		getInfoGoLogin(access, options?.setLoading);
	};
	return {
		loginByToken,
	};
};
export default useContent;
