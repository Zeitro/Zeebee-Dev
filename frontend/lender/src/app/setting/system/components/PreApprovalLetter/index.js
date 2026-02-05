import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import useContent from './useContent';
import { Button, Tour } from 'antd';
import { useRouter, useSearchParams } from 'next/navigation';
import { useRef, useState } from 'react';

const PreApprovedLetterSetting = () => {
	const router = useRouter();
	const searchParams = useSearchParams();
	const settingRef = useRef();
	const { switchBtn, switchRef, preApprovalLetterSetting } = useContent();
	const [open, setOpen] = useState(false);
	const steps = [
		{
			title: 'Setting page',
			description: 'Click to jump to the configuration page.',
			target: () => settingRef.current,
		},
		{
			title: 'Action',
			description: 'Turn on the switch.',
			target: () => switchRef.current,
		},
	];

	const handleSetting = () => {
		router.push('/setting/preApprovedLetter');
	};

	useDeferredSyncEffect({
		deferred: async () => {
			if (
				!preApprovalLetterSetting?.has_configured &&
				searchParams.get('isSettingPreApprovedLetter') // 只有第一次登录跳转到/setting/system页面让用户配置preApprovedLetter时才显示引导页面
			) {
				setOpen(true);
			}
		},
		deps: [preApprovalLetterSetting, searchParams],
		isUseMicrotask: false,
		delay: 500,
	});

	return (
		<div className="flex items-center justify-between w-full xl:w-2/3 bg-[#f5f5f5] hover:bg-[#f2f2f2] px-4 py-4 rounded-lg">
			<div>
				<div className="text-base font-bold pb-2">Send Pre-Approval Letter</div>
				<div className="text-sm text-base-label">
					When activated, your borrower will automatically receive a pre-approval letter.
					<Button
						className="underline"
						onClick={handleSetting}
						size="small"
						type="link"
						ref={settingRef}
					>
						Configure detailed settings
					</Button>
				</div>
			</div>
			{switchBtn}
			<Tour open={open} onClose={() => setOpen(false)} steps={steps} />
		</div>
	);
};

export default PreApprovedLetterSetting;
