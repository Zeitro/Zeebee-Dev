'use client';
import { FilterTableCard } from '@repo/shared-methods';
import useContent from './components/useContent';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';

const SETTING_USER_API = '/organization_user/page';

/**
 * 设置查询user
 * @param {*} getNewMainModal 获取添加/编辑弹框（与封装的组件差距大，则直接引用自定义的组件）
 * @returns
 */
const SettingUser = ({ getNewMainModal = () => {}, PERMISSIONS }) => {
	const { columns, getRolesOptions, getUserStatusOptions, filterList, tableFilterRef } =
		useContent({
			getNewMainModal,
			PERMISSIONS,
			SETTING_USER_API,
		});

	useDeferredSyncEffect({
		deferred: () => {
			getRolesOptions();
			getUserStatusOptions();
		},
		deps: [],
		isUseMicrotask: true,
	});
	return (
		<FilterTableCard
			filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden"
			ref={tableFilterRef}
			columns={columns}
			scroll={{ x: 'max-content' }}
			filterList={filterList}
			hasFilterField={true}
			api={SETTING_USER_API}
		/>
	);
};
export default SettingUser;
