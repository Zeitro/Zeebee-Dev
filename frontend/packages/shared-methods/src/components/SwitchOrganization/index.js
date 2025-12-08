'use client';
import { Loading3QuartersOutlined } from '@ant-design/icons';
import { eventBus } from '@repo/shared-methods/utils';
import { useUserStore } from '@/store/user';

// 切换组织中的页面
export const SwitchOrganization = ({ organization }) => {
	const { getUserInfo } = useUserStore();
	const handleSwitch = async () => {
		await getUserInfo();
		eventBus.emit('removePortalComponent');
	};
	return (
		<div
			className="flex items-center justify-center absolute top-0 left-0 right-0 bottom-0"
			style={{ background: '#7a7a7a', zIndex: 9999 }}
		>
			<div
				style={{ minWidth: 280, maxWidth: 500 }}
				className="flex items-center justify-center flex-col bg-white rounded-lg px-6 py-6 shadow-[0_6px_16px_0_rgba(0,0,0,0.08), 0_3px_6px_-4px_rgba(0,0,0,0.12), 0_9px_28px_8px_rgba(0,0,0,0.05)]"
			>
				<Loading3QuartersOutlined spin style={{ color: '#8C8C8C', fontSize: '24px' }} />
				<div className="text-base font-medium text-black-dark pt-6">
					Switch to {organization}
				</div>
			</div>
		</div>
	);
};
