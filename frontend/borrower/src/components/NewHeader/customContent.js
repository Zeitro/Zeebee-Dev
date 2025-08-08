'use client';
import { usePathname,useRouter } from 'next/navigation';
import { Button } from 'antd';
import {PlusOutlined} from '@ant-design/icons';

export const CustomContent = () => {
	const pathname = usePathname();
    const router = useRouter();
	const createLoanBtn = (
			<div className="flex items-center">
				<Button type="primary" icon={<PlusOutlined />} onClick={() => router.push('/loans/createLoan')}>
					Create new loan
				</Button>
			</div>
		)
	const headerContent = {
		'/loans': createLoanBtn,
	};
	return headerContent[pathname];
};
