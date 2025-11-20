'use client';
import { usePathname, useRouter, useParams } from 'next/navigation';
import { Button, Dropdown } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { matchRoute } from '@config/common';
import { useState } from 'react';
import { BtnPopUpModal } from '@repo/shared-methods';
import ImportLoansModal from '@/app/loans/components/ImportLoanModal';
export const CustomContent = () => {
	const pathname = usePathname();
	const { id } = useParams();
	const router = useRouter();
	const propertySearch = new URLSearchParams({ propertyId: id });
	const [selectedKey, setSelectedKey] = useState('ZeeBee');
	const items = [
		{
			key: 'ZeeBee',
			label: 'Apply ZeeBee loan',
			link: '/loans/createZeeBeeLoan',
		},
		{
			key: '3rdParty',
			label: 'Add 3rd party loan',
			link: '/loans/createLoan',
		},
		{
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
		},
	];

	const handleClick = (linkUrl, search) => {
		if (linkUrl) {
			router.push(linkUrl + `${search ? `?${search}` : ''}`);
		}
	};

	const createLoanBtn = (search) => (
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
						console.log('999');

						const selectItem = items.find((i) => i.key === selectedKey);
						handleClick(selectItem?.link, search);
					}}
				>
					{items.find((i) => i.key === selectedKey)?.label}
				</Button>
			</Dropdown>
		</div>
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
