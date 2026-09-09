'use client';

import { useUserStore } from '@repo/shared-methods/store';

/**
 * 登录后落地页（非 Dashboard）：问候 + 当前组织 + 简短说明。
 * 导航交给侧边栏；套餐不在此展示。
 */
const Home = () => {
	const { userInfo, organization } = useUserStore();
	const firstName = userInfo?.first_name?.trim();
	const orgName = organization?.name?.trim() || userInfo?.organization_name?.trim();
	const hasMultipleOrgs = (userInfo?.organizations?.length ?? 0) > 1;

	return (
		<div className="flex h-full w-full items-center justify-center overflow-y-auto bg-base-text2 px-6 py-12">
			<div className="w-full max-w-md text-center">
				<h1 className="m-0 text-[32px] font-medium leading-10 tracking-tight text-base-text">
					{firstName ? `Hi, ${firstName}` : 'Welcome back'}
				</h1>

				{orgName ? (
					<p className="m-0 mt-3 text-base leading-6 text-base-gray">
						You&apos;re working in{' '}
						<span className="font-medium text-base-text">{orgName}</span>
					</p>
				) : null}

				<p className="m-0 mx-auto mt-5 max-w-sm text-sm leading-6 text-base-label">
					Your workspace for loans, properties, and transactions. Pick a module from the
					sidebar to get started.
				</p>

				{hasMultipleOrgs ? (
					<p className="m-0 mx-auto mt-10 max-w-sm border-t border-base-border pt-6 text-xs leading-5 text-base-label">
						To switch organizations, click your name at the bottom of the sidebar, then
						choose Change organization.
					</p>
				) : null}
			</div>
		</div>
	);
};

export default Home;
