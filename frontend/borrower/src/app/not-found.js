'use client';
import { Button, Result } from 'antd';
import { useRouter } from 'next/navigation';

const NotFound = () => {
	const router = useRouter();
	return (
		<div className="flex items-center justify-center h-full w-full">
			<Result
				status="404"
				title="404"
				subTitle="Sorry, the page you visited does not exist."
				extra={
					<Button type="primary" onClick={() => router.back()}>
						Back to the previous page
					</Button>
				}
				className="shadow-[0_2px_8px_#eee] rounded-lg"
			/>
		</div>
	);
};

export default NotFound;
