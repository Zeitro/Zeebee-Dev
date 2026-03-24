'use client';
import { usePathname, useRouter, useParams } from 'next/navigation';
import { Button, Dropdown } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { matchRoute } from '@repo/shared-methods/utils';
import { useMemo, useState } from 'react';
import { BtnPopUpModal } from '@repo/shared-methods';
import ImportLoansModal from '@/app/loans/components/ImportLoanModal';
import { PERMISSIONS } from '@config/permissions';
import { useAuth } from '@repo/shared-methods/hooks';

// 自定义全局操作或展示（用于header面包屑平级处展示）
const HeaderCustomAction = () => {
	const { hasPermission, hasAnyPermission } = useAuth();
	const pathname = usePathname();
	const { id } = useParams();
	const router = useRouter();
	const propertySearch = new URLSearchParams({ propertyId: id });
	const [selectedKey, setSelectedKey] = useState('ZeeBee');
	const items = useMemo(() => {
		const arr = [];
		if (hasPermission(PERMISSIONS.LOAN.ADD_ZEEBEE)) {
			arr.push({
				key: 'ZeeBee',
				label: 'Apply ZeeBee loan',
				link: '/loans/createZeeBeeLoan',
			});
		}
		if (hasPermission(PERMISSIONS.LOAN.ADD_3RD_PARTY)) {
			arr.push({
				key: '3rdParty',
				label: 'Add 3rd party loan',
				link: '/loans/createLoan',
			});
		}
		if (hasPermission(PERMISSIONS.LOAN.IMPORT)) {
			arr.push({
				key: 'import',
				label: (
					<BtnPopUpModal
						type="text"
						isButton={false}
						getContent={(setOpen) => <ImportLoansModal setOpen={setOpen} />}
					>
						Import loans
					</BtnPopUpModal>
				),
			});
		}
		return arr;
	}, [hasPermission]);

	const handleClick = (linkUrl, search) => {
		if (linkUrl) {
			router.push(linkUrl + `${search ? `?${search}` : ''}`);
		}
	};

	const createLoanBtn = (search) =>
		hasAnyPermission([
			PERMISSIONS.LOAN.ADD_ZEEBEE,
			PERMISSIONS.LOAN.ADD_3RD_PARTY,
			PERMISSIONS.LOAN.IMPORT,
		]) ? (
			<div className="flex items-center">
				<Dropdown
					menu={{
						items,
						selectable: true,
						onClick: ({ key }) => {
							setSelectedKey(key);
							const selectItem = items.find((i) => i.key === key);
							handleClick(selectItem?.link, search);
						},
						selectedKeys: [selectedKey],
					}}
					placement="bottom"
					arrow
				>
					<Button
						type="primary"
						icon={<PlusOutlined />}
						onClick={() => {
							const selectItem = items.find((i) => i.key === selectedKey);
							handleClick(selectItem?.link, search);
						}}
					>
						{items.find((i) => i.key === selectedKey)?.label}
					</Button>
				</Dropdown>
			</div>
		) : (
			''
		);

	const headerContent = {
		'/loans': createLoanBtn(),
		'/properties/:id': createLoanBtn(propertySearch), // 房产处跳转新增loan，需传propertyId
	};

	const getContent = () => {
		const path = Object.keys(headerContent).find((i) => matchRoute(i, pathname));
		return headerContent[path] || '';
	};
	return getContent();
};

export default HeaderCustomAction;
