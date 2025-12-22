'use client';
import { useEffect, useMemo, useState } from 'react';
import { FieldTips } from '@repo/shared-methods';
import { Button } from 'antd';
import { useUserStore } from '@/store/user';
import PreviewQuoteModal from './PreviewQuoteModal';
// 公共询价按钮
const QuoteBtn = ({
	disabled,
	rowData,
	getRowData,
	isGetOrganization,
	isBtn = true,
	children,
	tooltipProps = {},
	...props
}) => {
	const [open, setOpen] = useState(false);
	const [data, setData] = useState({});
	const { userInfo, organization, getOrganization } = useUserStore();
	const disabledQuote = useMemo(() => {
		const contactInfo = organization?.organization_emails?.find(
			(item) => item.purpose === 'Contact',
		);
		return disabled || !contactInfo?.is_authorized;
	}, [disabled, organization]);

	const tips = useMemo(() => {
		if (disabled) {
			return 'Only for pending stage & quoting status loan';
		}
		const contactInfo = organization?.organization_emails?.find(
			(item) => item.purpose === 'Contact',
		);
		// 邮箱未授权
		if (!contactInfo?.is_authorized) {
			// 管理员
			if (userInfo?.is_admin || userInfo?.is_owner) {
				return (
					<div>
						<p className="text-sm">Email is not authorized for this organization.</p>
						<Button
							size="small"
							onClick={() => window.open('/setting/system', '_blank')}
						>
							Setting
						</Button>
					</div>
				);
			}
			// 非管理员
			return (
				<div className="text-sm">
					Email is not authorized for this organization. <br />
					Please contact your admin.
				</div>
			);
		}
		return '';
	}, [disabled, organization?.organization_emails, userInfo?.is_admin, userInfo?.is_owner]);
	const handleClick = async () => {
		if (disabledQuote) return;
		// 必填字段或需验证字段全对才可进行询价
		const v = (await getRowData?.()) || rowData;
		setData(v);
		setOpen(true);
	};

	useEffect(() => {
		if (isGetOrganization) {
			getOrganization();
		}
	}, [getOrganization, isGetOrganization]);
	return (
		<FieldTips title={tips} {...tooltipProps}>
			{isBtn ? (
				<Button
					type="primary"
					onClick={handleClick}
					ghost
					disabled={disabledQuote}
					{...props}
				>
					{children || 'Request quote'}
				</Button>
			) : (
				<div
					onClick={handleClick}
					style={{
						color: disabledQuote ? 'var(--base-disabled)' : 'var(--foreground)',
						cursor: disabledQuote ? 'not-allowed' : 'pointer',
					}}
				>
					{children || 'Request quote'}
				</div>
			)}
			{open && <PreviewQuoteModal setOpen={setOpen} rowData={data} />}
		</FieldTips>
	);
};

export default QuoteBtn;
