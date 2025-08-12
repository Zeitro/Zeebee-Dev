'use client';
import { usePathname, useRouter, useParams } from 'next/navigation';
import { Button } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { matchRoute } from '@config/common';
export const CustomContent = () => {
	const pathname = usePathname();
	const { id } = useParams();
	const router = useRouter();
	const propertySearch = new URLSearchParams({ propertyId: id });

	const createLoanBtn = (search) => (
		<div className="flex items-center">
			<Button
				type="primary"
				icon={<PlusOutlined />}
				onClick={() => router.push('/loans/createLoan' + `${search ? `?${search}` : ''}`)}
			>
				Create new loan
			</Button>
		</div>
	);
	const headerContent = {
		'/loans': createLoanBtn(),
		'/properties/:id': createLoanBtn(propertySearch),
	};
	const getContent = () => {
		const path = Object.keys(headerContent).find((i) => matchRoute(i, pathname));
		return headerContent[path] || '';
	};
	return getContent();
};
